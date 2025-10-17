package sdk

var (
	_ validatable = new(CreateNotebookOptions)
	_ validatable = new(AlterNotebookOptions)
	_ validatable = new(DropNotebookOptions)
	_ validatable = new(DescribeNotebookOptions)
	_ validatable = new(ShowNotebookOptions)
)

func (opts *CreateNotebookOptions) validate() error {
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
	if opts.Warehouse != nil && !ValidObjectIdentifier(opts.Warehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateNotebookOptions", "IfNotExists", "OrReplace"))
	}
	if opts.ComputePool != nil && !ValidObjectIdentifier(opts.ComputePool) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	// Validation added manually.
	if opts.IdleAutoShutdownTimeSeconds != nil && !validateIntGreaterThan(*opts.IdleAutoShutdownTimeSeconds, 0) {
		errs = append(errs, errIntValue("CreateNotebookOptions", "IdleAutoShutdownTimeSeconds", IntErrGreater, 0))
	}
	// Validation added manually.
	if opts.ExternalAccessIntegrations != nil {
		for _, integration := range opts.ExternalAccessIntegrations {
			if !ValidObjectIdentifier(integration) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterNotebookOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterNotebookOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if opts.Set.QueryWarehouse != nil && !ValidObjectIdentifier(opts.Set.QueryWarehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if opts.Set.Warehouse != nil && !ValidObjectIdentifier(opts.Set.Warehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if opts.Set.ComputePool != nil && !ValidObjectIdentifier(opts.Set.ComputePool) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if !anyValueSet(opts.Set.Comment, opts.Set.QueryWarehouse, opts.Set.IdleAutoShutdownTimeSeconds, opts.Set.Secrets, opts.Set.MainFile, opts.Set.Warehouse, opts.Set.RuntimeName, opts.Set.ComputePool, opts.Set.ExternalAccessIntegrations, opts.Set.RuntimeEnvironmentVersion) {
			errs = append(errs, errAtLeastOneOf("AlterNotebookOptions.Set", "Comment", "QueryWarehouse", "IdleAutoShutdownTimeSeconds", "Secrets", "MainFile", "Warehouse", "RuntimeName", "ComputePool", "ExternalAccessIntegrations", "RuntimeEnvironmentVersion"))
		}
		// Validation added manually.
		if opts.Set.IdleAutoShutdownTimeSeconds != nil && !validateIntGreaterThan(*opts.Set.IdleAutoShutdownTimeSeconds, 0) {
			errs = append(errs, errIntValue("AlterNotebookOptions", "IdleAutoShutdownTimeSeconds", IntErrGreater, 0))
		}
		// Validation added manually.
		if opts.Set.ExternalAccessIntegrations != nil {
			for _, integration := range opts.Set.ExternalAccessIntegrations {
				if !ValidObjectIdentifier(integration) {
					errs = append(errs, ErrInvalidObjectIdentifier)
				}
			}
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Comment, opts.Unset.QueryWarehouse, opts.Unset.Secrets, opts.Unset.Warehouse, opts.Unset.RuntimeName, opts.Unset.ComputePool, opts.Unset.ExternalAccessIntegrations, opts.Unset.RuntimeEnvironmentVersion) {
			errs = append(errs, errAtLeastOneOf("AlterNotebookOptions.Unset", "Comment", "QueryWarehouse", "Secrets", "Warehouse", "RuntimeName", "ComputePool", "ExternalAccessIntegrations", "RuntimeEnvironmentVersion"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropNotebookOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeNotebookOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowNotebookOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
