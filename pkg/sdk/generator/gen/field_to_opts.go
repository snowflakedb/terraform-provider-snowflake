package gen

// toOpts template variable name constants. Changing a value here updates both the Go methods
// and the template output, keeping them in sync without touching any logic.
const (
	toOptsSourceRootVar = "r"    // source variable name outside slice context
	toOptsSliceElemVar  = "v"    // loop element variable inside slice range
	toOptsOptsVar       = "opts" // opts variable name outside slice context
	toOptsSliceVar      = "s"    // local slice build variable
)

// ToOptsSourcePath returns the source-access path for use in the toOpts template.
// When a slice ancestor exists, the path is relative to that slice element (stops at the first slice).
// Otherwise it equals Path() — the full path from root.
func (f *Field) ToOptsSourcePath() string {
	if f.IsRoot() {
		return ""
	}
	if f.Parent.IsSlice() {
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
		return toOptsSliceVar + "[i]"
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
