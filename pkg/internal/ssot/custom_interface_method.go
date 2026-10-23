package ssot

// CustomInterfaceMethodParameter is one parameter of a custom interface method, excluding ctx.
type CustomInterfaceMethodParameter struct {
	Name string
	// Kind is the Go type. Empty uses the object's IdType.
	Kind string
}

// CustomInterfaceMethodSpec declares one interface method. The body stays in the SDK ext file.
// A nil Parameters slice means the method takes id of IdType.
// A non-nil Parameters slice is the parameter list, excluding ctx.
// An empty Returns slice means the method returns error.
// An empty Doc emits no comment.
type CustomInterfaceMethodSpec struct {
	Name       string
	Doc        string
	Parameters []CustomInterfaceMethodParameter
	Returns    []string
}
