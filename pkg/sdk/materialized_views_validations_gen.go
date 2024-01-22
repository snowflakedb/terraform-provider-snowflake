package sdk

var (
	_ validatable = new(CreateMaterializedViewOptions)
	_ validatable = new(AlterMaterializedViewOptions)
	_ validatable = new(DropMaterializedViewOptions)
	_ validatable = new(ShowMaterializedViewOptions)
	_ validatable = new(DescribeMaterializedViewOptions)
)

func (opts *CreateMaterializedViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateMaterializedViewOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.RowAccessPolicy) {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.RowAccessPolicy) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterMaterializedViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.ClusterBy, opts.DropClusteringKey, opts.SuspendRecluster, opts.ResumeRecluster, opts.Suspend, opts.Resume, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterMaterializedViewOptions", "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Secure, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterMaterializedViewOptions.Set", "Secure", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Secure, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterMaterializedViewOptions.Unset", "Secure", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropMaterializedViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowMaterializedViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeMaterializedViewOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
