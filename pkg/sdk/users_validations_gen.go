package sdk

import "errors"

var (
	_ validatable = new(CreateUserOptions)
	_ validatable = new(AlterUserOptions)
	_ validatable = new(DropUserOptions)
	_ validatable = new(describeUserOptions)
	_ validatable = new(ShowUserOptions)
)

func (opts *CreateUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.ObjectProperties) {
		if valueSet(opts.ObjectProperties.WorkloadIdentity) {
			if err := opts.ObjectProperties.WorkloadIdentity.validate(); err != nil {
				return err
			}
		}
		if valueSet(opts.ObjectProperties.DefaultSecondaryRoles) {
			if err := opts.ObjectProperties.DefaultSecondaryRoles.validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (opts *AlterUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.NewName, opts.ResetPassword, opts.AbortAllQueries, opts.AddDelegatedAuthorization, opts.RemoveDelegatedAuthorization, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterUserOptions", "NewName", "ResetPassword", "AbortAllQueries", "AddDelegatedAuthorization", "RemoveDelegatedAuthorization", "Set", "Unset", "SetTag", "UnsetTag"))
	}
	if valueSet(opts.RemoveDelegatedAuthorization) {
		if err := opts.RemoveDelegatedAuthorization.validate(); err != nil {
			errs = append(errs, err)
		}
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
	return errors.Join(errs...)
}

func (opts *DropUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *describeUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *ShowUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *showUserAuthenticationMethodOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (opts *UserSet) validate() error {
	if !anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return errAtLeastOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters")
	}
	if moreThanOneValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) {
		return errOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy")
	}
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be set with user properties or parameters at the same time")
	}
	if valueSet(opts.ObjectProperties) {
		if valueSet(opts.ObjectProperties.WorkloadIdentity) {
			if err := opts.ObjectProperties.WorkloadIdentity.validate(); err != nil {
				return err
			}
		}
		if valueSet(opts.ObjectProperties.DefaultSecondaryRoles) {
			if err := opts.ObjectProperties.DefaultSecondaryRoles.validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (opts *UserUnset) validate() error {
	// TODO [SNOW-1645875]: change validations with policies
	if !anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters, opts.AuthenticationPolicy) {
		return errAtLeastOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters")
	}
	if moreThanOneValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) {
		return errOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy")
	}
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be unset with user properties or parameters at the same time")
	}
	return nil
}

func (opts *SecondaryRoles) validate() error {
	if !exactlyOneValueSet(opts.All, opts.None) {
		return errExactlyOneOf("SecondaryRoles", "All", "None")
	}
	return nil
}

func (opts *RemoveDelegatedAuthorization) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.Role, opts.Authorizations) {
		errs = append(errs, errExactlyOneOf("RemoveDelegatedAuthorization", "Role", "Authorization"))
	}
	if !valueSet(opts.Integration) {
		errs = append(errs, errNotSet("RemoveDelegatedAuthorization", "Integration"))
	}
	return errors.Join(errs...)
}

func (workloadIdentity *UserObjectWorkloadIdentityProperties) validate() error {
	var errs []error
	if !exactlyOneValueSet(workloadIdentity.AwsType, workloadIdentity.AzureType, workloadIdentity.GcpType, workloadIdentity.OidcType) {
		errs = append(errs, errExactlyOneOf("UserObjectWorkloadIdentityProperties", "AwsType", "AzureType", "GcpType", "OidcType"))
	}
	return errors.Join(errs...)
}
