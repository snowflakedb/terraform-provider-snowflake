package sdk

var (
	_ validatable = new(CreateExternalAccessIntegrationOptions)
	_ validatable = new(AlterExternalAccessIntegrationOptions)
	_ validatable = new(DropExternalAccessIntegrationOptions)
	_ validatable = new(ShowExternalAccessIntegrationOptions)
	_ validatable = new(DescribeExternalAccessIntegrationOptions)
)

func (opts *CreateExternalAccessIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateExternalAccessIntegrationOptions", "IfNotExists", "OrReplace"))
	}
	if len(opts.AllowedNetworkRules) == 0 {
		errs = append(errs, errNotSet("CreateExternalAccessIntegrationOptions", "AllowedNetworkRules"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterExternalAccessIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterExternalAccessIntegrationOptions", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.AllowedNetworkRules, opts.Set.AllowedAuthenticationSecrets, opts.Set.Enabled, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterExternalAccessIntegrationOptions.Set", "AllowedNetworkRules", "AllowedAuthenticationSecrets", "Enabled", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.AllowedAuthenticationSecrets, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterExternalAccessIntegrationOptions.Unset", "AllowedAuthenticationSecrets", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropExternalAccessIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowExternalAccessIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeExternalAccessIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
