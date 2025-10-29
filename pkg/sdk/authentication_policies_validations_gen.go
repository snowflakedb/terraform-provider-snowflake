package sdk

var (
	_ validatable = new(CreateAuthenticationPolicyOptions)
	_ validatable = new(AlterAuthenticationPolicyOptions)
	_ validatable = new(DropAuthenticationPolicyOptions)
	_ validatable = new(ShowAuthenticationPolicyOptions)
	_ validatable = new(DescribeAuthenticationPolicyOptions)
)

func (opts *CreateAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateAuthenticationPolicyOptions", "IfNotExists", "OrReplace"))
	}
	if valueSet(opts.MfaPolicy) {
		if !anyValueSet(opts.MfaPolicy.EnforceMfaOnExternalAuthentication, opts.MfaPolicy.AllowedMethods) {
			errs = append(errs, errAtLeastOneOf("CreateAuthenticationPolicyOptions.MfaPolicy", "EnforceMfaOnExternalAuthentication", "AllowedMethods"))
		}
	}
	if valueSet(opts.SecurityIntegrations) {
		if !exactlyOneValueSet(opts.SecurityIntegrations.All, opts.SecurityIntegrations.SecurityIntegrations) {
			errs = append(errs, errExactlyOneOf("CreateAuthenticationPolicyOptions.SecurityIntegrations", "All", "SecurityIntegrations"))
		}
	}
	if valueSet(opts.PatPolicy) {
		if !anyValueSet(opts.PatPolicy.DefaultExpiryInDays, opts.PatPolicy.MaxExpiryInDays, opts.PatPolicy.NetworkPolicyEvaluation) {
			errs = append(errs, errAtLeastOneOf("CreateAuthenticationPolicyOptions.PatPolicy", "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation"))
		}
	}
	if valueSet(opts.WorkloadIdentityPolicy) {
		if !anyValueSet(opts.WorkloadIdentityPolicy.AllowedProviders, opts.WorkloadIdentityPolicy.AllowedAwsAccounts, opts.WorkloadIdentityPolicy.AllowedAzureIssuers, opts.WorkloadIdentityPolicy.AllowedOidcIssuers) {
			errs = append(errs, errAtLeastOneOf("CreateAuthenticationPolicyOptions.WorkloadIdentityPolicy", "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.RenameTo) {
		errs = append(errs, errExactlyOneOf("AlterAuthenticationPolicyOptions", "Set", "Unset", "RenameTo"))
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.AuthenticationMethods, opts.Set.MfaAuthenticationMethods, opts.Set.MfaEnrollment, opts.Set.ClientTypes, opts.Set.SecurityIntegrations, opts.Set.Comment, opts.Set.MfaPolicy, opts.Set.PatPolicy, opts.Set.WorkloadIdentityPolicy) {
			errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set", "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"))
		}
		if valueSet(opts.Set.MfaPolicy) {
			if !anyValueSet(opts.Set.MfaPolicy.EnforceMfaOnExternalAuthentication, opts.Set.MfaPolicy.AllowedMethods) {
				errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.MfaPolicy", "EnforceMfaOnExternalAuthentication", "AllowedMethods"))
			}
		}
		if valueSet(opts.Set.SecurityIntegrations) {
			if !exactlyOneValueSet(opts.Set.SecurityIntegrations.All, opts.Set.SecurityIntegrations.SecurityIntegrations) {
				errs = append(errs, errExactlyOneOf("AlterAuthenticationPolicyOptions.Set.SecurityIntegrations", "All", "SecurityIntegrations"))
			}
		}
		if valueSet(opts.Set.PatPolicy) {
			if !anyValueSet(opts.Set.PatPolicy.DefaultExpiryInDays, opts.Set.PatPolicy.MaxExpiryInDays, opts.Set.PatPolicy.NetworkPolicyEvaluation) {
				errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.PatPolicy", "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation"))
			}
		}
		if valueSet(opts.Set.WorkloadIdentityPolicy) {
			if !anyValueSet(opts.Set.WorkloadIdentityPolicy.AllowedProviders, opts.Set.WorkloadIdentityPolicy.AllowedAwsAccounts, opts.Set.WorkloadIdentityPolicy.AllowedAzureIssuers, opts.Set.WorkloadIdentityPolicy.AllowedOidcIssuers) {
				errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.WorkloadIdentityPolicy", "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers"))
			}
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.ClientTypes, opts.Unset.AuthenticationMethods, opts.Unset.Comment, opts.Unset.SecurityIntegrations, opts.Unset.MfaAuthenticationMethods, opts.Unset.MfaEnrollment, opts.Unset.MfaPolicy, opts.Unset.PatPolicy, opts.Unset.WorkloadIdentityPolicy) {
			errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Unset", "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
