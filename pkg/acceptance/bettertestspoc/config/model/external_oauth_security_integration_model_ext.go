package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// WithExternalOauthTokenUserMappingClaim was added to satisfy the default builders. The method itself is not generated because its type is not yet supported.
// This method will conflict the generated one when the type is supported.
func (e *ExternalOauthSecurityIntegrationModel) WithExternalOauthTokenUserMappingClaim(values []string) *ExternalOauthSecurityIntegrationModel {
	e.ExternalOauthTokenUserMappingClaim = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return e
}

func (e *ExternalOauthSecurityIntegrationModel) WithExternalOauthAllowedRoles(roles ...sdk.AccountObjectIdentifier) *ExternalOauthSecurityIntegrationModel {
	e.ExternalOauthAllowedRolesList = tfconfig.SetVariable(
		collections.Map(roles, func(role sdk.AccountObjectIdentifier) tfconfig.Variable {
			return tfconfig.StringVariable(role.Name())
		})...,
	)
	return e
}

func (e *ExternalOauthSecurityIntegrationModel) WithExternalOauthAudiences(values ...string) *ExternalOauthSecurityIntegrationModel {
	e.ExternalOauthAudienceList = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return e
}

func (e *ExternalOauthSecurityIntegrationModel) WithExternalOauthJwsKeysUrls(urls ...string) *ExternalOauthSecurityIntegrationModel {
	e.ExternalOauthJwsKeysUrl = tfconfig.SetVariable(
		collections.Map(urls, func(url string) tfconfig.Variable {
			return tfconfig.StringVariable(url)
		})...,
	)
	return e
}
