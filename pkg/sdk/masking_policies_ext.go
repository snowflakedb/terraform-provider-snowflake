package sdk

type MaskingPolicyOptions struct {
	ExemptOtherPolicies bool `json:"EXEMPT_OTHER_POLICIES"`
}

// additionalConvert extracts ExemptOtherPolicies from the JSON-parsed Options field.
func (r maskingPolicyDBRow) additionalConvert(result *MaskingPolicy) error {
	if result.Options != nil {
		result.ExemptOtherPolicies = result.Options.ExemptOtherPolicies
	}
	return nil
}

// additionalValidations checks that NewName stays within the same database and schema.
func (opts *AlterMaskingPolicyOptions) additionalValidations() error {
	if opts.RenameTo != nil {
		var errs []error
		if !ValidObjectIdentifier(opts.RenameTo) {
			errs = append(errs, errInvalidIdentifier("AlterMaskingPolicyOptions", "RenameTo"))
		}
		if opts.name.DatabaseName() != opts.RenameTo.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
		if opts.name.SchemaName() != opts.RenameTo.SchemaName() {
			errs = append(errs, ErrDifferentSchema)
		}
		return JoinErrors(errs...)
	}
	return nil
}
