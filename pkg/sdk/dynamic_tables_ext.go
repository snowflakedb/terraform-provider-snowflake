package sdk

func (opts *DynamicTableSet) additionalValidations() error {
	if opts.Warehouse != nil && !ValidObjectIdentifier(*opts.Warehouse) {
		return errInvalidIdentifier("DynamicTableSet", "Warehouse")
	}
	return nil
}

func (opts *AlterDynamicTableOptions) additionalValidations() error {
	var errs []error
	if addSLP := opts.AddStorageLifecyclePolicy; valueSet(addSLP) {
		if !ValidObjectIdentifier(addSLP.StorageLifecyclePolicy) {
			errs = append(errs, errInvalidIdentifier("DynamicTableAddStorageLifecyclePolicy", "StorageLifecyclePolicy"))
		}
		if len(addSLP.On) == 0 {
			errs = append(errs, errNotSet("DynamicTableAddStorageLifecyclePolicy", "On"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowDynamicTableOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("ShowDynamicTableOptions.In", "Account", "Database", "Schema"))
	}
	return JoinErrors(errs...)
}

// TypeString returns the parsed column type when available, otherwise the raw DESC output.
func (d DynamicTableDetails) TypeString() string {
	if d.Type != nil {
		return d.Type.ToSql()
	}
	return d.DataTypeRaw
}
