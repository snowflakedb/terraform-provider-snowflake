package gen

import (
	"fmt"
	"strings"
)

// objectIdentifierFallbackKind stands in for the generic ObjectIdentifier interface kind: any
// concrete kind works, since ValidObjectIdentifier only checks Name() != "".
const objectIdentifierFallbackKind = "AccountObjectIdentifier"

// emptyIdentifierVar returns the pkg/sdk test-only variable name for an empty (invalid) identifier of the given kind,
// following the convention empty<Kind> (e.g. "emptySchemaObjectIdentifier").
// The corresponding vars are declared in pkg/sdk/random_test.go.
func emptyIdentifierVar(kind string) string {
	if kind == "ObjectIdentifier" {
		kind = objectIdentifierFallbackKind
	}
	return "empty" + kind
}

// randomIdentifierCall returns the pkg/sdk test-only call expression that produces a valid random identifier of the given kind,
// following the convention random<Kind>() (e.g. "randomSchemaObjectIdentifier()").
// The corresponding functions are declared in pkg/sdk/random_test.go.
func randomIdentifierCall(kind string) string {
	return "random" + kind + "()"
}

// zeroValuesForPredefinedTypes lists zero values for predefined types not defined directly in the generator definition.
var zeroValuesForPredefinedTypes = map[string]string{
	"Location": "nil",
}

// zeroValueFor returns the Go expression that makes f unset in the valueSet() sense.
func zeroValueFor(f *Field) string {
	predefinedTypeZeroValue, isPredefinedType := zeroValuesForPredefinedTypes[f.Kind]
	switch {
	case f.IsPointer():
		return "nil"
	case f.IsSlice():
		return "nil"
	case f.IsIdentifier():
		return emptyIdentifierVar(f.KindNoPtr())
	case f.IsStruct():
		return f.KindNoPtr() + "{}"
	case f.Kind == "bool":
		return "false"
	case f.Kind == "int":
		return "0"
	case f.Kind == "any":
		return "nil"
	case isPredefinedType:
		return predefinedTypeZeroValue
	case strings.Contains(f.Kind, "."):
		return "nil"
	default:
		return `""`
	}
}

// nonZeroValueFor returns the Go expression that makes f set in the valueSet() sense, and true if one could be derived.
func nonZeroValueFor(f *Field) (string, bool) {
	switch {
	case f.IsSlice():
		return "", false
	case f.IsStruct():
		if f.IsPointer() {
			return fmt.Sprintf("&%s{}", f.KindNoPtr()), true
		}
		return f.KindNoPtr() + "{}", true
	case f.IsPointer():
		inner, ok := nonZeroScalarOrIdentifier(f.KindNoPtr(), f.IsIdentifier())
		if !ok {
			return "", false
		}
		return fmt.Sprintf("new(%s)", inner), true
	default:
		return nonZeroScalarOrIdentifier(f.Kind, f.IsIdentifier())
	}
}

func nonZeroScalarOrIdentifier(kind string, isIdentifier bool) (string, bool) {
	if isIdentifier {
		return randomIdentifierCall(kind), true
	}
	switch kind {
	case "bool":
		return "true", true
	case "int":
		return "1", true
	case "string":
		return `"foo"`, true
	default:
		return "", false
	}
}

// primeAncestors returns the statements needed so that target itself becomes reachable,
// i.e. so that every pointer or slice ancestor along the path from the root to target's parent is non-nil/non-empty.
func primeAncestors(target *Field) []string {
	stmts := make([]string, 0)
	for _, ancestor := range target.AncestorsFromRoot() {
		switch {
		case ancestor.IsPointer():
			stmts = append(stmts, fmt.Sprintf("opts%s = &%s{}", ancestor.IndexedPath(), ancestor.KindNoPtr()))
		case ancestor.IsSlice():
			// Best-effort single empty element; deeper nesting through indexed path.
			stmts = append(stmts, fmt.Sprintf("opts%s = []%s{{}}", ancestor.IndexedPath(), ancestor.KindNoPtr()))
		}
	}
	return stmts
}
