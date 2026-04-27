package gen

// MethodParameter represents a single named parameter in a custom interface method signature.
// The ctx context.Context parameter is always the first parameter and should not be included here.
type MethodParameter struct {
	// Name is the parameter name, e.g. "id" or "kind"
	Name string
	// Kind is the Go type string, e.g. "AccountObjectIdentifier" or "[]*Parameter"
	Kind string
}

// NewMethodParameter creates a new MethodParameter with the given name and kind.
func NewMethodParameter(name string, kind string) *MethodParameter {
	return &MethodParameter{Name: name, Kind: kind}
}

// CustomInterfaceMethod declares a method in the generated interface without a generated implementation.
// The implementation is expected to be provided manually by the user in a separate (non-generated) file (typically *_ext.go).
type CustomInterfaceMethod struct {
	// Name is the method name, e.g. "ShowParameters"
	Name string
	// Doc is a documentation comment emitted directly before the method in the interface.
	Doc string
	// Parameters is the list of method parameters excluding ctx context.Context which is always first.
	Parameters []*MethodParameter
	// ReturnTypes is the list of return types, e.g. ["[]*Parameter", "error"] or ["error"].
	ReturnTypes []string
}

// WithCustomInterfaceMethod adds a custom method declaration to the generated interface.
// The method will appear in the generated interface but no implementation will be generated for it.
// ctx context.Context is always the first parameter and must not be listed in params.
func (i *Interface) WithCustomInterfaceMethod(name, doc string, params []*MethodParameter, returnTypes ...string) *Interface {
	i.CustomMethods = append(i.CustomMethods, &CustomInterfaceMethod{
		Name:        name,
		Doc:         doc,
		Parameters:  params,
		ReturnTypes: returnTypes,
	})
	return i
}
