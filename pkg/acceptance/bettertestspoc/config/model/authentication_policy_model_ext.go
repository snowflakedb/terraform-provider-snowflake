package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AuthenticationPolicyModel) WithAuthenticationMethods(authenticationMethods ...sdk.AuthenticationMethodsOption) *AuthenticationPolicyModel {
	return a.WithAuthenticationMethodsValue(
		tfconfig.SetVariable(
			collections.Map(authenticationMethods, func(authenticationMethod sdk.AuthenticationMethodsOption) tfconfig.Variable {
				return tfconfig.StringVariable(string(authenticationMethod))
			})...,
		),
	)
}

func (a *AuthenticationPolicyModel) WithMfaAuthenticationMethods(mfaAuthenticationMethods ...sdk.MfaAuthenticationMethodsOption) *AuthenticationPolicyModel {
	return a.WithMfaAuthenticationMethodsValue(
		tfconfig.SetVariable(
			collections.Map(mfaAuthenticationMethods, func(mfaAuthenticationMethod sdk.MfaAuthenticationMethodsOption) tfconfig.Variable {
				return tfconfig.StringVariable(string(mfaAuthenticationMethod))
			})...,
		),
	)
}

func (a *AuthenticationPolicyModel) WithClientTypes(clientTypes ...sdk.ClientTypesOption) *AuthenticationPolicyModel {
	return a.WithClientTypesValue(
		tfconfig.SetVariable(
			collections.Map(clientTypes, func(clientType sdk.ClientTypesOption) tfconfig.Variable {
				return tfconfig.StringVariable(string(clientType))
			})...,
		),
	)
}

func (a *AuthenticationPolicyModel) WithSecurityIntegrations(securityIntegrations ...string) *AuthenticationPolicyModel {
	return a.WithSecurityIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(securityIntegrations, func(securityIntegration string) tfconfig.Variable {
				return tfconfig.StringVariable(securityIntegration)
			})...,
		),
	)
}

func (a *AuthenticationPolicyModel) WithMfaEnrollmentEnum(mfaEnrollment sdk.MfaEnrollmentOption) *AuthenticationPolicyModel {
	return a.WithMfaEnrollmentValue(tfconfig.StringVariable(string(mfaEnrollment)))
}

func (a *AuthenticationPolicyModel) WithMfaPolicy(mfaPolicy sdk.AuthenticationPolicyMfaPolicyRequest) *AuthenticationPolicyModel {
	m := make(map[string]tfconfig.Variable)
	if mfaPolicy.EnforceMfaOnExternalAuthentication != nil {
		m["enforce_mfa_on_external_authentication"] = tfconfig.StringVariable(string(*mfaPolicy.EnforceMfaOnExternalAuthentication))
	}
	if mfaPolicy.AllowedMethods != nil {
		m["allowed_methods"] = tfconfig.ListVariable(
			collections.Map(mfaPolicy.AllowedMethods, func(allowedMethod sdk.AuthenticationPolicyMfaPolicyListItem) tfconfig.Variable {
				return tfconfig.StringVariable(string(allowedMethod.Method))
			})...,
		)
	}
	return a.WithMfaPolicyValue(tfconfig.ObjectVariable(m))
}

func (a *AuthenticationPolicyModel) WithPatPolicy(patPolicy sdk.AuthenticationPolicyPatPolicyRequest) *AuthenticationPolicyModel {
	m := make(map[string]tfconfig.Variable)
	if patPolicy.DefaultExpiryInDays != nil {
		m["default_expiry_in_days"] = tfconfig.IntegerVariable(*patPolicy.DefaultExpiryInDays)
	}
	if patPolicy.MaxExpiryInDays != nil {
		m["max_expiry_in_days"] = tfconfig.IntegerVariable(*patPolicy.MaxExpiryInDays)
	}
	if patPolicy.NetworkPolicyEvaluation != nil {
		m["network_policy_evaluation"] = tfconfig.StringVariable(string(*patPolicy.NetworkPolicyEvaluation))
	}
	return a.WithPatPolicyValue(tfconfig.ObjectVariable(m))
}

func (a *AuthenticationPolicyModel) WithWorkloadIdentityPolicy(workloadIdentityPolicy sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest) *AuthenticationPolicyModel {
	m := make(map[string]tfconfig.Variable)
	if workloadIdentityPolicy.AllowedProviders != nil {
		m["allowed_providers"] = tfconfig.SetVariable(
			collections.Map(workloadIdentityPolicy.AllowedProviders, func(allowedProvider sdk.AuthenticationPolicyAllowedProviderListItem) tfconfig.Variable {
				return tfconfig.StringVariable(string(allowedProvider.Provider))
			})...,
		)
	}
	if workloadIdentityPolicy.AllowedAwsAccounts != nil {
		m["allowed_aws_accounts"] = tfconfig.SetVariable(
			collections.Map(workloadIdentityPolicy.AllowedAwsAccounts, func(allowedAwsAccount sdk.StringListItemWrapper) tfconfig.Variable {
				return tfconfig.StringVariable(allowedAwsAccount.Value)
			})...,
		)
	}
	if workloadIdentityPolicy.AllowedAzureIssuers != nil {
		m["allowed_azure_issuers"] = tfconfig.SetVariable(
			collections.Map(workloadIdentityPolicy.AllowedAzureIssuers, func(allowedAzureIssuer sdk.StringListItemWrapper) tfconfig.Variable {
				return tfconfig.StringVariable(allowedAzureIssuer.Value)
			})...,
		)
	}
	if workloadIdentityPolicy.AllowedOidcIssuers != nil {
		m["allowed_oidc_issuers"] = tfconfig.SetVariable(
			collections.Map(workloadIdentityPolicy.AllowedOidcIssuers, func(allowedOidcIssuer sdk.StringListItemWrapper) tfconfig.Variable {
				return tfconfig.StringVariable(allowedOidcIssuer.Value)
			})...,
		)
	}
	return a.WithWorkloadIdentityPolicyValue(tfconfig.ObjectVariable(m))
}
