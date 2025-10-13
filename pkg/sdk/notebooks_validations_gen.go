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
	if opts.IdleAutoShutdownTimeSeconds != nil && !validateIntInRangeInclusive(*opts.IdleAutoShutdownTimeSeconds, 60, 259200) {
		errs = append(errs, errIntBetween("CreateNotebookOptions", "IdleAutoShutdownTimeSeconds", 60, 259200))
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
	if !exactlyOneValueSet(opts.Set) {
		errs = append(errs, errExactlyOneOf("AlterNotebookOptions", "Set"))
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
		// Validation added manually.
		if opts.Set.IdleAutoShutdownTimeSeconds != nil && !validateIntInRangeInclusive(*opts.Set.IdleAutoShutdownTimeSeconds, 60, 259200) {
			errs = append(errs, errIntBetween("CreateNotebookOptions", "IdleAutoShutdownTimeSeconds", 60, 259200))
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
