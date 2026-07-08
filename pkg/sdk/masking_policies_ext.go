package sdk

type MaskingPolicyOptions struct {
	ExemptOtherPolicies bool `json:"EXEMPT_OTHER_POLICIES"`
}

// additionalConvert extracts ExemptOtherPolicies from the JSON-parsed Options field.
func (r maskingPolicyDBRow) additionalConvert(result *MaskingPolicy) error {
	result.ExemptOtherPolicies = result.Options.ExemptOtherPolicies
	return nil
}

// additionalValidations checks that NewName stays within the same database and schema.
func (opts *AlterMaskingPolicyOptions) additionalValidations() error {
	if opts.NewName != nil {
		var errs []error
		if !ValidObjectIdentifier(opts.NewName) {
			errs = append(errs, errInvalidIdentifier("AlterMaskingPolicyOptions", "NewName"))
		}
		if opts.name.DatabaseName() != opts.NewName.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
		if opts.name.SchemaName() != opts.NewName.SchemaName() {
			errs = append(errs, ErrDifferentSchema)
		}
		return JoinErrors(errs...)
	}
	return nil
}
