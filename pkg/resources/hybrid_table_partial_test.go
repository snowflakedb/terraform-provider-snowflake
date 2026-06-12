package resources

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_hybridTablePartialOnError asserts that every `return diag.FromErr(...)`
// in CreateHybridTable / UpdateHybridTable that sits in a position where a
// remote Snowflake ALTER may have already mutated state is preceded by
// `d.Partial(true)`. This mirrors the convention used in pkg/resources/schema.go
// and was requested in the PR #4689 review.
//
// The rule applied here matches the structural sites the plan touches:
//
//  1. The return is the immediate error handler of an `if err := alter(); err
//     != nil { return diag.FromErr(...) }` block whose init clause invokes
//     `client.HybridTables.Alter`.
//
//  2. The return is inside the body of a `for`/`range` loop that contains a
//     `client.HybridTables.Alter` call — prior loop iterations may already
//     have committed changes.
//
// Notably, the `client.HybridTables.Create` error site in CreateHybridTable is
// NOT flagged: a failed CREATE leaves the resource absent and `d.SetId` has not
// been called yet, so there is no partial state to record.
func Test_hybridTablePartialOnError(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hybrid_table.go", nil, parser.ParseComments)
	require.NoError(t, err)

	bs, err := os.ReadFile("hybrid_table.go")
	require.NoError(t, err)
	src := strings.Split(string(bs), "\n")

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name != "CreateHybridTable" && fn.Name.Name != "UpdateHybridTable" {
			continue
		}

		var violationLines []int
		collectPartialViolations(fn.Body, false, fset, &violationLines)

		for _, line := range violationLines {
			if !precedingLineHas(src, line, "d.Partial(true)") {
				t.Errorf("hybrid_table.go:%d in %s: `return diag.FromErr(...)` after remote ALTER must be preceded by `d.Partial(true)`",
					line, fn.Name.Name)
			}
		}
	}
}

// collectPartialViolations walks `node` and records every `return
// diag.FromErr(...)` that lives inside a flagged context. `inFlaggedContext`
// is true when the return is either:
//
//   - the body of an if-statement whose init clause invokes
//     `client.HybridTables.Alter`, or
//   - inside a for/range loop body that contains a
//     `client.HybridTables.Alter` call.
//
// The flag is propagated INTO descendants so that a return nested several
// levels below a flagged trigger is still reported.
func collectPartialViolations(node ast.Node, inFlaggedContext bool, fset *token.FileSet, violations *[]int) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.IfStmt:
		bodyFlagged := inFlaggedContext || initInvokesAlter(n.Init)
		collectInBlock(n.Body, bodyFlagged, fset, violations)
		if n.Else != nil {
			collectPartialViolations(n.Else, inFlaggedContext, fset, violations)
		}
		return

	case *ast.ForStmt:
		bodyFlagged := inFlaggedContext || blockContainsAlter(n.Body)
		collectInBlock(n.Body, bodyFlagged, fset, violations)
		return

	case *ast.RangeStmt:
		bodyFlagged := inFlaggedContext || blockContainsAlter(n.Body)
		collectInBlock(n.Body, bodyFlagged, fset, violations)
		return

	case *ast.ReturnStmt:
		if inFlaggedContext && isDiagFromErrReturn(n) {
			*violations = append(*violations, fset.Position(n.Pos()).Line)
		}
		return

	case *ast.BlockStmt:
		collectInBlock(n, inFlaggedContext, fset, violations)
		return
	}

	// Generic fallthrough: walk children that contain control-flow nodes.
	ast.Inspect(node, func(child ast.Node) bool {
		if child == node || child == nil {
			return true
		}
		switch child.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.ReturnStmt, *ast.BlockStmt:
			collectPartialViolations(child, inFlaggedContext, fset, violations)
			return false
		}
		return true
	})
}

func collectInBlock(b *ast.BlockStmt, inFlaggedContext bool, fset *token.FileSet, violations *[]int) {
	if b == nil {
		return
	}
	for _, stmt := range b.List {
		collectPartialViolations(stmt, inFlaggedContext, fset, violations)
	}
}

// initInvokesAlter reports whether the init statement of an if-statement
// invokes `client.HybridTables.Alter` (i.e. matches `if err := alter(...);
// err != nil`).
func initInvokesAlter(init ast.Stmt) bool {
	if init == nil {
		return false
	}
	return nodeContainsCall(init, "client.HybridTables.Alter")
}

// blockContainsAlter reports whether the block contains a call to
// `client.HybridTables.Alter` at any depth.
func blockContainsAlter(b *ast.BlockStmt) bool {
	if b == nil {
		return false
	}
	return nodeContainsCall(b, "client.HybridTables.Alter")
}

func nodeContainsCall(n ast.Node, qualifiedName string) bool {
	found := false
	ast.Inspect(n, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if calls(call, qualifiedName) {
			found = true
			return false
		}
		return true
	})
	return found
}

func calls(call *ast.CallExpr, qualifiedName string) bool {
	parts := strings.Split(qualifiedName, ".")
	cursor := call.Fun
	for i := len(parts) - 1; i >= 0; i-- {
		switch e := cursor.(type) {
		case *ast.SelectorExpr:
			if e.Sel.Name != parts[i] {
				return false
			}
			cursor = e.X
		case *ast.Ident:
			return i == 0 && e.Name == parts[i]
		default:
			return false
		}
	}
	return false
}

func isDiagFromErrReturn(r *ast.ReturnStmt) bool {
	if len(r.Results) != 1 {
		return false
	}
	call, ok := r.Results[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	return calls(call, "diag.FromErr")
}

func precedingLineHas(src []string, line int, needle string) bool {
	if line < 2 || line > len(src) {
		return false
	}
	// Walk back to the nearest non-empty, non-comment line.
	for i := line - 2; i >= 0; i-- {
		s := strings.TrimSpace(src[i])
		if s == "" || strings.HasPrefix(s, "//") {
			continue
		}
		return strings.Contains(s, needle)
	}
	return false
}
