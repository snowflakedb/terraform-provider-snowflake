package sdk

var (
	_ validatable = new(CreateOrganizationAccountOptions)
	_ validatable = new(AlterOrganizationAccountOptions)
	_ validatable = new(ShowOrganizationAccountOptions)
)

func (opts *CreateOrganizationAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !anyValueSet(opts.AdminPassword, opts.AdminRsaPublicKey) {
		errs = append(errs, errAtLeastOneOf("CreateOrganizationAccountOptions", "AdminPassword", "AdminRsaPublicKey"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterOrganizationAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if opts.Name != nil && !ValidObjectIdentifier(opts.Name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Name, opts.Set) {
		errs = append(errs, errOneOf("AlterOrganizationAccountOptions", "Name", "Set"))
	}
	if everyValueSet(opts.Name, opts.Unset) {
		errs = append(errs, errOneOf("AlterOrganizationAccountOptions", "Name", "Unset"))
	}
	if everyValueSet(opts.Name, opts.SetTags) {
		errs = append(errs, errOneOf("AlterOrganizationAccountOptions", "Name", "SetTags"))
	}
	if everyValueSet(opts.Name, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterOrganizationAccountOptions", "Name", "UnsetTags"))
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.RenameTo, opts.DropOldUrl) {
		errs = append(errs, errExactlyOneOf("AlterOrganizationAccountOptions", "Set", "Unset", "SetTags", "UnsetTags", "RenameTo", "DropOldUrl"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Parameters, opts.Set.ResourceMonitor, opts.Set.PasswordPolicy, opts.Set.SessionPolicy) {
			errs = append(errs, errAtLeastOneOf("AlterOrganizationAccountOptions.Set", "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Parameters, opts.Unset.ResourceMonitor, opts.Unset.PasswordPolicy, opts.Unset.SessionPolicy) {
			errs = append(errs, errAtLeastOneOf("AlterOrganizationAccountOptions.Unset", "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowOrganizationAccountOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
