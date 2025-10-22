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
