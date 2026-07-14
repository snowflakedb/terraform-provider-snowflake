package sdk

import "errors"

var (
	_ validatable = new(CreateAccountOptions)
	_ validatable = new(AlterAccountOptions)
	_ validatable = new(ShowAccountOptions)
	_ validatable = new(DropAccountOptions)
	_ validatable = new(UndropAccountOptions)
)

func (opts *CreateAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.AdminName == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "AdminName"))
	}
	if !anyValueSet(opts.AdminPassword, opts.AdminRSAPublicKey) {
		errs = append(errs, errAtLeastOneOf("CreateAccountOptions", "AdminPassword", "AdminRSAPublicKey"))
	}
	if opts.Email == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Email"))
	}
	if opts.Edition == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Edition"))
	}
	return errors.Join(errs...)
}

func (opts *AlterAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.Drop, opts.Rename) {
		errs = append(errs, errExactlyOneOf("AlterAccountOptions", "Set", "Unset", "SetTag", "UnsetTag", "Drop", "Rename"))
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Drop) {
		if err := opts.Drop.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Rename) {
		if err := opts.Rename.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if err := opts.additionalValidations(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (opts *AccountSet) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.Parameters, opts.LegacyParameters, opts.ResourceMonitor, opts.PackagesPolicy, opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.FeaturePolicySet, opts.OrgAdmin, opts.ConsumptionBillingEntity) {
		errs = append(errs, errExactlyOneOf("AccountSet", "Parameters", "LegacyParameters", "ResourceMonitor", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "FeaturePolicySet", "OrgAdmin", "ConsumptionBillingEntity"))
	}
	if err := opts.additionalValidations(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (opts *AccountUnset) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.LegacyParameters, opts.Parameters, opts.PackagesPolicy, opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.ResourceMonitor, opts.FeaturePolicyUnset, opts.ConsumptionBillingEntity) {
		errs = append(errs, errExactlyOneOf("AccountUnset", "Parameters", "LegacyParameters", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ResourceMonitor", "FeaturePolicyUnset", "ConsumptionBillingEntity"))
	}
	if err := opts.additionalValidations(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (opts *ShowAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *DropAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *UndropAccountOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *AccountLevelParameters) validate() error {
	var errs []error
	if valueSet(opts.AccountParameters) {
		if err := opts.AccountParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.ObjectParameters) {
		if err := opts.ObjectParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.UserParameters) {
		if err := opts.UserParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *AccountLevelParametersUnset) validate() error {
	if !anyValueSet(opts.AccountParameters, opts.SessionParameters, opts.ObjectParameters, opts.UserParameters) {
		return errAtLeastOneOf("AccountLevelParametersUnset", "LegacyAccountParameters", "SessionParameters", "ObjectParameters", "UserParameters")
	}
	return nil
}

func (opts *AccountRename) validate() error {
	var errs []error
	if !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AccountRename", "NewName"))
	}
	return errors.Join(errs...)
}

func (opts *AccountDrop) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.OldUrl, opts.OldOrganizationUrl) {
		errs = append(errs, errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"))
	}
	return errors.Join(errs...)
}
