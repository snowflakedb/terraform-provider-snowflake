package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *ApiIntegrationGitRepositoryPrivateLinkModel) WithApiAllowedPrefixes(apiAllowedPrefixes []string) *ApiIntegrationGitRepositoryPrivateLinkModel {
	prefixVars := collections.Map(apiAllowedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiAllowedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}

func (a *ApiIntegrationGitRepositoryPrivateLinkModel) WithApiBlockedPrefixes(apiBlockedPrefixes []string) *ApiIntegrationGitRepositoryPrivateLinkModel {
	prefixVars := collections.Map(apiBlockedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiBlockedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}

func (a *ApiIntegrationGitRepositoryPrivateLinkModel) WithAllowedAuthenticationSecrets(secrets []string) *ApiIntegrationGitRepositoryPrivateLinkModel {
	secretVars := collections.Map(secrets, func(s string) tfconfig.Variable { return tfconfig.StringVariable(s) })
	a.WithAllowedAuthenticationSecretsValue(tfconfig.SetVariable(secretVars...))
	return a
}
