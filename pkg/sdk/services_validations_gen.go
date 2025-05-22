package sdk

var (
	_ validatable = new(CreateServiceOptions)
	_ validatable = new(AlterServiceOptions)
	_ validatable = new(DropServiceOptions)
	_ validatable = new(ShowServiceOptions)
	_ validatable = new(DescribeServiceOptions)
)

// TODO: Fill validations and add unit tests for them.
func (opts *CreateServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.QueryWarehouse != nil && !ValidObjectIdentifier(opts.QueryWarehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Resume, opts.Suspend, opts.ServiceFromSpecification, opts.ServiceFromSpecificationTemplate, opts.Restore, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterServiceOptions", "Resume", "Suspend", "ServiceFromSpecification", "ServiceFromSpecificationTemplate", "Restore", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if opts.Set.QueryWarehouse != nil && !ValidObjectIdentifier(opts.Set.QueryWarehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if !anyValueSet(opts.Set.MinInstances, opts.Set.MaxInstances, opts.Set.AutoSuspendSecs, opts.Set.MinReadyInstances, opts.Set.QueryWarehouse, opts.Set.AutoResume, opts.Set.ServiceExternalAccessIntegrations, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterServiceOptions.Set", "MinInstances", "MaxInstances", "AutoSuspendSecs", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ServiceExternalAccessIntegrations", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.MinInstances, opts.Unset.AutoSuspendSecs, opts.Unset.MaxInstances, opts.Unset.MinReadyInstances, opts.Unset.QueryWarehouse, opts.Unset.AutoResume, opts.Unset.ExternalAccessIntegrations, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterServiceOptions.Unset", "MinInstances", "AutoSuspendSecs", "MaxInstances", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.Job, opts.ExcludeJobs) {
		errs = append(errs, errOneOf("ShowServiceOptions", "Job", "ExcludeJobs"))
	}
	return JoinErrors(errs...)
}

func (opts *DescribeServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
