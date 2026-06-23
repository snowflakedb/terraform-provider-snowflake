package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *ApiIntegrationExternalMcpOAuth2Model) WithApiAllowedPrefixes(apiAllowedPrefixes []string) *ApiIntegrationExternalMcpOAuth2Model {
	prefixVars := collections.Map(apiAllowedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiAllowedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}

func (a *ApiIntegrationExternalMcpOAuth2Model) WithApiBlockedPrefixes(apiBlockedPrefixes []string) *ApiIntegrationExternalMcpOAuth2Model {
	prefixVars := collections.Map(apiBlockedPrefixes, func(p string) tfconfig.Variable { return tfconfig.StringVariable(p) })
	a.WithApiBlockedPrefixesValue(tfconfig.ListVariable(prefixVars...))
	return a
}
