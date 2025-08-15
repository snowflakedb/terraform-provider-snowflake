package sdk

var (
	_ validatable = new(CreateSemanticViewOptions)
	_ validatable = new(DropSemanticViewOptions)
	_ validatable = new(DescribeSemanticViewOptions)
	_ validatable = new(ShowSemanticViewOptions)
)

func (opts *CreateSemanticViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateSemanticViewOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *DropSemanticViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeSemanticViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowSemanticViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
