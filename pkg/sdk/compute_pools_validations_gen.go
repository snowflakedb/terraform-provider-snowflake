package sdk

var (
	_ validatable = new(CreateComputePoolOptions)
	_ validatable = new(AlterComputePoolOptions)
	_ validatable = new(DropComputePoolOptions)
	_ validatable = new(ShowComputePoolOptions)
	_ validatable = new(DescribeComputePoolOptions)
)

func (opts *CreateComputePoolOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// Validation added manually.
	if !validateIntGreaterThan(opts.MinNodes, 0) {
		errs = append(errs, errIntValue("CreateComputePoolOptions", "MinNodes", IntErrGreater, 0))
	}
	// Validation added manually.
	if !validateIntGreaterThanOrEqual(opts.MaxNodes, opts.MinNodes) {
		errs = append(errs, errIntValue("CreateComputePoolOptions", "MaxNodes", IntErrGreaterOrEqual, opts.MinNodes))
	}
	return JoinErrors(errs...)
}

func (opts *AlterComputePoolOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterComputePoolOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		// Validation added manually.
		if valueSet(opts.Set.MinNodes) && !validateIntGreaterThan(*opts.Set.MinNodes, 0) {
			errs = append(errs, errIntValue("AlterComputePoolOptions", "Set.MinNodes", IntErrGreater, 0))
		}
		// Validation added manually.
		if valueSet(opts.Set.MaxNodes) && !validateIntGreaterThan(*opts.Set.MaxNodes, 0) {
			errs = append(errs, errIntValue("AlterComputePoolOptions", "Set.MaxNodes", IntErrGreater, 0))
		}
		if !anyValueSet(opts.Set.MinNodes, opts.Set.MaxNodes, opts.Set.AutoResume, opts.Set.AutoSuspendSecs, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterComputePoolOptions.Set", "MinNodes", "MaxNodes", "AutoResume", "AutoSuspendSecs", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.AutoResume, opts.Unset.AutoSuspendSecs, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterComputePoolOptions.Unset", "AutoResume", "AutoSuspendSecs", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropComputePoolOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowComputePoolOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeComputePoolOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
