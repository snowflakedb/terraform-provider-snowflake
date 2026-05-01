package sdk

var (
	_ validatable = new(CreateOpenflowRuntimeOptions)
	_ validatable = new(AlterOpenflowRuntimeOptions)
	_ validatable = new(DropOpenflowRuntimeOptions)
	_ validatable = new(ShowOpenflowRuntimeOptions)
	_ validatable = new(DescribeOpenflowRuntimeOptions)
)

func (opts *CreateOpenflowRuntimeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !validateIntGreaterThan(opts.MinNodes, 0) {
		errs = append(errs, errIntValue("CreateOpenflowRuntimeOptions", "MinNodes", IntErrGreater, 0))
	}
	if !validateIntGreaterThanOrEqual(opts.MaxNodes, opts.MinNodes) {
		errs = append(errs, errIntValue("CreateOpenflowRuntimeOptions", "MaxNodes", IntErrGreaterOrEqual, opts.MinNodes))
	}
	return JoinErrors(errs...)
}

func (opts *AlterOpenflowRuntimeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Suspend, opts.Resume, opts.Terminate, opts.TerminateCascade, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterOpenflowRuntimeOptions", "Suspend", "Resume", "Terminate", "TerminateCascade", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.MinNodes, opts.Set.MaxNodes, opts.Set.ExecuteAsRole, opts.Set.DisplayName, opts.Set.Comment) && len(opts.Set.ExternalAccessIntegrations) == 0 {
			errs = append(errs, errAtLeastOneOf("AlterOpenflowRuntimeOptions.Set", "MinNodes", "MaxNodes", "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.DisplayName, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterOpenflowRuntimeOptions.Unset", "DisplayName", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropOpenflowRuntimeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowOpenflowRuntimeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	return nil
}

func (opts *DescribeOpenflowRuntimeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
