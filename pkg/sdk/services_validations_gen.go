package sdk

var (
	_ validatable = new(CreateServiceOptions)
	_ validatable = new(AlterServiceOptions)
	_ validatable = new(DropServiceOptions)
	_ validatable = new(ShowServiceOptions)
	_ validatable = new(DescribeServiceOptions)
	_ validatable = new(ExecuteJobServiceOptions)
)

func (opts *CreateServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.FromSpecification, opts.FromSpecificationTemplate) {
		errs = append(errs, errExactlyOneOf("CreateServiceOptions", "FromSpecification", "FromSpecificationTemplate"))
	}
	if opts.QueryWarehouse != nil && !ValidObjectIdentifier(opts.QueryWarehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.FromSpecification) {
		if !exactlyOneValueSet(opts.FromSpecification.SpecificationFile, opts.FromSpecification.Specification) {
			errs = append(errs, errExactlyOneOf("CreateServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
		}
		if everyValueSet(opts.FromSpecification.Location, opts.FromSpecification.Specification) {
			errs = append(errs, errOneOf("CreateServiceOptions.FromSpecification", "Location", "Specification"))
		}
	}
	if valueSet(opts.FromSpecificationTemplate) {
		if !exactlyOneValueSet(opts.FromSpecificationTemplate.SpecificationTemplateFile, opts.FromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errExactlyOneOf("CreateServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
		}
		if everyValueSet(opts.FromSpecificationTemplate.Location, opts.FromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errOneOf("CreateServiceOptions.FromSpecificationTemplate", "Location", "SpecificationTemplate"))
		}
	}
	// Validation added manually.
	if valueSet(opts.MinReadyInstances) {
		if !validateIntGreaterThan(*opts.MinReadyInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinReadyInstances", IntErrGreater, 0))
		}
		if valueSet(opts.MinInstances) && !validateIntGreaterThanOrEqual(*opts.MinInstances, *opts.MinReadyInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreaterOrEqual, *opts.MinReadyInstances))
		}
		if valueSet(opts.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.MaxInstances, *opts.MinReadyInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, *opts.MinReadyInstances))
		}
	}
	// Validation added manually.
	if valueSet(opts.MinInstances) {
		if !validateIntGreaterThan(*opts.MinInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreater, 0))
		}
		if valueSet(opts.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.MaxInstances, *opts.MinInstances) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, *opts.MinInstances))
		}
	}
	// Validation added manually.
	if valueSet(opts.MaxInstances) {
		if !validateIntGreaterThan(*opts.MaxInstances, 0) {
			errs = append(errs, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreater, 0))
		}
	}
	// Validation added manually.
	if valueSet(opts.AutoSuspendSecs) && !validateIntGreaterThanOrEqual(*opts.AutoSuspendSecs, 0) {
		errs = append(errs, errIntValue("CreateServiceOptions", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
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
	if !exactlyOneValueSet(opts.Resume, opts.Suspend, opts.FromSpecification, opts.FromSpecificationTemplate, opts.Restore, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterServiceOptions", "Resume", "Suspend", "FromSpecification", "FromSpecificationTemplate", "Restore", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.FromSpecification) {
		if !exactlyOneValueSet(opts.FromSpecification.SpecificationFile, opts.FromSpecification.Specification) {
			errs = append(errs, errExactlyOneOf("AlterServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
		}
		if everyValueSet(opts.FromSpecification.Location, opts.FromSpecification.Specification) {
			errs = append(errs, errOneOf("AlterServiceOptions.FromSpecification", "Location", "Specification"))
		}
	}
	if valueSet(opts.FromSpecificationTemplate) {
		if !exactlyOneValueSet(opts.FromSpecificationTemplate.SpecificationTemplateFile, opts.FromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errExactlyOneOf("AlterServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
		}
		if everyValueSet(opts.FromSpecificationTemplate.Location, opts.FromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errOneOf("AlterServiceOptions.FromSpecificationTemplate", "Location", "SpecificationTemplate"))
		}
	}
	if valueSet(opts.Restore) {
		if !ValidObjectIdentifier(opts.Restore.FromSnapshot) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if valueSet(opts.Set) {
		if opts.Set.QueryWarehouse != nil && !ValidObjectIdentifier(opts.Set.QueryWarehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if !anyValueSet(opts.Set.MinInstances, opts.Set.MaxInstances, opts.Set.AutoSuspendSecs, opts.Set.MinReadyInstances, opts.Set.QueryWarehouse, opts.Set.AutoResume, opts.Set.ExternalAccessIntegrations, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterServiceOptions.Set", "MinInstances", "MaxInstances", "AutoSuspendSecs", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"))
		}
		// Validation added manually.
		if valueSet(opts.Set.MinReadyInstances) {
			if !validateIntGreaterThan(*opts.Set.MinReadyInstances, 0) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinReadyInstances", IntErrGreater, 0))
			}
			if valueSet(opts.Set.MinInstances) && !validateIntGreaterThanOrEqual(*opts.Set.MinInstances, *opts.Set.MinReadyInstances) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreaterOrEqual, *opts.Set.MinReadyInstances))
			}
			if valueSet(opts.Set.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.Set.MaxInstances, *opts.Set.MinReadyInstances) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, *opts.Set.MinReadyInstances))
			}
		}
		// Validation added manually.
		if valueSet(opts.Set.MinInstances) {
			if !validateIntGreaterThan(*opts.Set.MinInstances, 0) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreater, 0))
			}
			if valueSet(opts.Set.MaxInstances) && !validateIntGreaterThanOrEqual(*opts.Set.MaxInstances, *opts.Set.MinInstances) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, *opts.Set.MinInstances))
			}
		}
		// Validation added manually.
		if valueSet(opts.Set.MaxInstances) {
			if !validateIntGreaterThan(*opts.Set.MaxInstances, 0) {
				errs = append(errs, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreater, 0))
			}
		}
		// Validation added manually.
		if valueSet(opts.Set.AutoSuspendSecs) && !validateIntGreaterThanOrEqual(*opts.Set.AutoSuspendSecs, 0) {
			errs = append(errs, errIntValue("AlterServiceOptions.Set", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
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

func (opts *ExecuteJobServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.Name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.JobServiceFromSpecification, opts.JobServiceFromSpecificationTemplate) {
		errs = append(errs, errExactlyOneOf("ExecuteJobServiceOptions", "JobServiceFromSpecification", "JobServiceFromSpecificationTemplate"))
	}
	if !ValidObjectIdentifier(opts.InComputePool) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.QueryWarehouse != nil && !ValidObjectIdentifier(opts.QueryWarehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.JobServiceFromSpecification) {
		if !exactlyOneValueSet(opts.JobServiceFromSpecification.SpecificationFile, opts.JobServiceFromSpecification.Specification) {
			errs = append(errs, errExactlyOneOf("ExecuteJobServiceOptions.JobServiceFromSpecification", "SpecificationFile", "Specification"))
		}
		if !exactlyOneValueSet(opts.JobServiceFromSpecification.Location, opts.JobServiceFromSpecification.Specification) {
			errs = append(errs, errExactlyOneOf("ExecuteJobServiceOptions.JobServiceFromSpecification", "Location", "Specification"))
		}
	}
	if valueSet(opts.JobServiceFromSpecificationTemplate) {
		if !exactlyOneValueSet(opts.JobServiceFromSpecificationTemplate.SpecificationTemplateFile, opts.JobServiceFromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errExactlyOneOf("ExecuteJobServiceOptions.JobServiceFromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
		}
		if !exactlyOneValueSet(opts.JobServiceFromSpecificationTemplate.Location, opts.JobServiceFromSpecificationTemplate.SpecificationTemplate) {
			errs = append(errs, errExactlyOneOf("ExecuteJobServiceOptions.JobServiceFromSpecificationTemplate", "Location", "SpecificationTemplate"))
		}
	}
	return JoinErrors(errs...)
}
