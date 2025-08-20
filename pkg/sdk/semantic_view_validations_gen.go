package sdk

var (
	_ validatable = new(CreateSemanticViewOptions)
	_ validatable = new(DropSemanticViewOptions)
	_ validatable = new(DescribeSemanticViewOptions)
	_ validatable = new(ShowSemanticViewOptions)
	_ validatable = new(AlterSemanticViewOptions)
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
	if valueSet(opts.semanticViewMetrics) {
		for _, v := range opts.semanticViewMetrics {
			if !exactlyOneValueSet(v.semanticExpression, v.windowFunctionMetricDefinition) {
				errs = append(errs, errExactlyOneOf("CreateSemanticViewOptions.semanticViewMetrics", "semanticExpression", "windowFunctionMetricDefinition"))
			}
		}
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

func (opts *AlterSemanticViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
