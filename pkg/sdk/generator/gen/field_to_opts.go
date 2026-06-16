package gen

import "strings"

// toOpts template variable name constants. Changing a value here updates both the Go methods
// and the template output, keeping them in sync without touching any logic.
const (
	toOptsSourceRootVar = "r"    // source variable name outside slice context
	toOptsSliceElemVar  = "v"    // loop element variable inside slice range (block-scoped, no conflict)
	toOptsOptsVar       = "opts" // opts variable name outside slice context
)

// ToOptsSliceVar returns a field-derived local variable name for the slice being built.
// Unlike the loop element variable (which is block-scoped to the for loop and harmless to shadow),
// the slice build variable is declared in the enclosing block and would shadow an outer declaration
// for the rest of that block. Field-derived names prevent this.
// E.g. field "LogicalTables" → "logicalTables", field "UniqueKeys" → "uniqueKeys".
func (f *Field) ToOptsSliceVar() string {
	return strings.ToLower(f.Name[:1]) + f.Name[1:]
}

// ToOptsSourcePath returns the source-access path for use in the toOpts template.
// When a slice ancestor exists, the path is relative to that slice element (stops at the first slice).
// When a shared-toOpts ancestor exists, the path is relative to that struct (stops there).
// Otherwise it equals Path() — the full path from root.
func (f *Field) ToOptsSourcePath() string {
	if f.IsRoot() {
		return ""
	}
	if f.Parent.IsSlice() {
		return "." + f.Name
	}
	if f.Parent.GenerateSharedToOpts {
		return "." + f.Name
	}
	if f.Parent.IsRoot() {
		return f.Path() // no slice ancestor — full path
	}
	return f.Parent.ToOptsSourcePath() + "." + f.Name
}

// ToOptsSourceRoot returns the source root variable for use in the toOpts template.
// Returns the slice element variable when a slice ancestor exists, otherwise the request variable.
func (f *Field) ToOptsSourceRoot() string {
	if f.hasSliceAncestor() {
		return toOptsSliceElemVar
	}
	return toOptsSourceRootVar
}

// ToOptsOptsRef returns the opts assignment target for use in the toOpts template.
// Returns the indexed slice reference when a slice ancestor exists, otherwise the opts variable.
func (f *Field) ToOptsOptsRef() string {
	if f.hasSliceAncestor() {
		return f.sliceAncestor().ToOptsSliceVar() + "[i]"
	}
	return toOptsOptsVar
}

// hasSliceAncestor reports whether any ancestor in the parent chain is a slice field.
func (f *Field) hasSliceAncestor() bool {
	if f.IsRoot() {
		return false
	}
	if f.Parent.IsSlice() {
		return true
	}
	return f.Parent.hasSliceAncestor()
}

// sliceAncestor returns the nearest slice ancestor field in the parent chain.
func (f *Field) sliceAncestor() *Field {
	if f.Parent.IsSlice() {
		return f.Parent
	}
	return f.Parent.sliceAncestor()
}
