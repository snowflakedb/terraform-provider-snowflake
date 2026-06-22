package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *ApiIntegrationGitRepositoryOauth2Model) WithApiAllowedPrefixes(apiAllowedPrefixes []string) *ApiIntegrationGitRepositoryOauth2Model {
	prefixVars := collections.Map(apiAllowedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiAllowedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}

func (a *ApiIntegrationGitRepositoryOauth2Model) WithApiBlockedPrefixes(apiBlockedPrefixes []string) *ApiIntegrationGitRepositoryOauth2Model {
	prefixVars := collections.Map(apiBlockedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiBlockedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}

func (a *ApiIntegrationGitRepositoryOauth2Model) WithOauthAllowedScopes(scopes []string) *ApiIntegrationGitRepositoryOauth2Model {
	scopeVars := collections.Map(scopes, func(s string) tfconfig.Variable { return tfconfig.StringVariable(s) })
	a.WithOauthAllowedScopesValue(tfconfig.ListVariable(scopeVars...))
	return a
}
