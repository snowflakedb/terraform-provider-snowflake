package sdk

var (
	_ validatable = new(createDynamicTableOptions)
	_ validatable = new(alterDynamicTableOptions)
	_ validatable = new(dropDynamicTableOptions)
	_ validatable = new(showDynamicTableOptions)
	_ validatable = new(describeDynamicTableOptions)
	_ validatable = new(DynamicTableSet)
)

func (tl *TargetLag) validate() error {
	if tl == nil {
		return ErrNilOptions
	}
	if everyValueSet(tl.MaximumDuration, tl.Downstream) {
		return errOneOf("TargetLag", "MaximumDuration", "Downstream")
	}
	return nil
}

func (opts *createDynamicTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.warehouse) {
		errs = append(errs, errInvalidIdentifier("createDynamicTableOptions", "warehouse"))
	}
	return JoinErrors(errs...)
}

func (dts *DynamicTableSet) validate() error {
	var errs []error
	if valueSet(dts.TargetLag) {
		errs = append(errs, dts.TargetLag.validate())
	}
	if dts.Warehouse != nil && !ValidObjectIdentifier(*dts.Warehouse) {
		errs = append(errs, errInvalidIdentifier("DynamicTableSet", "Warehouse"))
	}
	return JoinErrors(errs...)
}

func (opts *alterDynamicTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Suspend, opts.Resume, opts.Refresh, opts.Set, opts.AddStorageLifecyclePolicy, opts.DropStorageLifecyclePolicy); !ok {
		errs = append(errs, errExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy"))
	}
	if valueSet(opts.Set) && valueSet(opts.Set.TargetLag) {
		errs = append(errs, opts.Set.TargetLag.validate())
	}
	if addStorageLifecyclePolicy := opts.AddStorageLifecyclePolicy; valueSet(addStorageLifecyclePolicy) {
		if !ValidObjectIdentifier(addStorageLifecyclePolicy.StorageLifecyclePolicy) {
			errs = append(errs, errInvalidIdentifier("DynamicTableAddStorageLifecyclePolicy", "StorageLifecyclePolicy"))
		}
		if len(addStorageLifecyclePolicy.On) == 0 {
			errs = append(errs, errNotSet("DynamicTableAddStorageLifecyclePolicy", "On"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *showDynamicTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("showDynamicTableOptions.In", "Account", "Database", "Schema"))
	}
	return JoinErrors(errs...)
}

func (opts *dropDynamicTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !ValidObjectIdentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *describeDynamicTableOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !ValidObjectIdentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}
