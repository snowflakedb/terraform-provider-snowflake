package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateAuthenticationPolicyOptions]   = new(CreateAuthenticationPolicyRequest)
	_ optionsProvider[AlterAuthenticationPolicyOptions]    = new(AlterAuthenticationPolicyRequest)
	_ optionsProvider[DropAuthenticationPolicyOptions]     = new(DropAuthenticationPolicyRequest)
	_ optionsProvider[ShowAuthenticationPolicyOptions]     = new(ShowAuthenticationPolicyRequest)
	_ optionsProvider[DescribeAuthenticationPolicyOptions] = new(DescribeAuthenticationPolicyRequest)
)

type CreateAuthenticationPolicyRequest struct {
	OrReplace                *bool
	IfNotExists              *bool
	name                     SchemaObjectIdentifier // required
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *MfaEnrollmentOption
	MfaPolicy                *AuthenticationPolicyMfaPolicyRequest
	ClientTypes              []ClientTypes
	SecurityIntegrations     *SecurityIntegrationsOptionRequest
	PatPolicy                *AuthenticationPolicyPatPolicyRequest
	WorkloadIdentityPolicy   *AuthenticationPolicyWorkloadIdentityPolicyRequest
	Comment                  *string
}

type AuthenticationPolicyMfaPolicyRequest struct {
	EnforceMfaOnExternalAuthentication *EnforceMfaOnExternalAuthenticationOption
	AllowedMethods                     []AuthenticationPolicyMfaPolicyListItem
}

type SecurityIntegrationsOptionRequest struct {
	All                  *bool
	SecurityIntegrations []AccountObjectIdentifier
}

type AuthenticationPolicyPatPolicyRequest struct {
	DefaultExpiryInDays     *int
	MaxExpiryInDays         *int
	NetworkPolicyEvaluation *NetworkPolicyEvaluationOption
}

type AuthenticationPolicyWorkloadIdentityPolicyRequest struct {
	AllowedProviders    []AuthenticationPolicyAllowedProviderListItem
	AllowedAwsAccounts  []StringListItemWrapper
	AllowedAzureIssuers []StringListItemWrapper
	AllowedOidcIssuers  []StringListItemWrapper
}

type AlterAuthenticationPolicyRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *AuthenticationPolicySetRequest
	Unset    *AuthenticationPolicyUnsetRequest
	RenameTo *SchemaObjectIdentifier
}

type AuthenticationPolicySetRequest struct {
	AuthenticationMethods    []AuthenticationMethods
	MfaAuthenticationMethods []MfaAuthenticationMethods
	MfaEnrollment            *MfaEnrollmentOption
	MfaPolicy                *AuthenticationPolicyMfaPolicyRequest
	ClientTypes              []ClientTypes
	SecurityIntegrations     *SecurityIntegrationsOptionRequest
	PatPolicy                *AuthenticationPolicyPatPolicyRequest
	WorkloadIdentityPolicy   *AuthenticationPolicyWorkloadIdentityPolicyRequest
	Comment                  *string
}

type AuthenticationPolicyUnsetRequest struct {
	ClientTypes              *bool
	AuthenticationMethods    *bool
	SecurityIntegrations     *bool
	MfaAuthenticationMethods *bool
	MfaEnrollment            *bool
	MfaPolicy                *bool
	PatPolicy                *bool
	WorkloadIdentityPolicy   *bool
	Comment                  *bool
}

type DropAuthenticationPolicyRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowAuthenticationPolicyRequest struct {
	Like       *Like
	In         *ExtendedIn
	On         *On
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeAuthenticationPolicyRequest struct {
	name SchemaObjectIdentifier // required
}
