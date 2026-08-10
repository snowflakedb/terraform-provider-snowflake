package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalAccessIntegrationFromId(
	id sdk.AccountObjectIdentifier,
	networkRules []sdk.SchemaObjectIdentifier,
	enabled bool,
) *ExternalAccessIntegrationModel {
	e := &ExternalAccessIntegrationModel{ResourceModelMeta: config.Meta("test", resources.ExternalAccessIntegration)}
	e.WithName(id.Name())
	e.WithAllowedNetworkRules(collections.Map(networkRules, sdk.SchemaObjectIdentifier.FullyQualifiedName))
	e.WithEnabled(enabled)
	return e
}

func (e *ExternalAccessIntegrationModel) WithAllowedNetworkRules(networkRules []string) *ExternalAccessIntegrationModel {
	if len(networkRules) == 0 {
		return e.WithAllowedNetworkRulesValue(config.EmptyListVariable())
	}
	return e.WithAllowedNetworkRulesValue(
		tfconfig.SetVariable(
			collections.Map(networkRules, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...,
		),
	)
}

func (e *ExternalAccessIntegrationModel) WithAllowedAuthenticationSecretsSecrets(secrets []string) *ExternalAccessIntegrationModel {
	return e.WithAllowedAuthenticationSecretsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"secrets": tfconfig.SetVariable(
			collections.Map(secrets, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...,
		),
	})))
}

func (e *ExternalAccessIntegrationModel) WithAllowedAuthenticationSecretsAll() *ExternalAccessIntegrationModel {
	return e.WithAllowedAuthenticationSecretsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.BoolVariable(true),
	})))
}

func (e *ExternalAccessIntegrationModel) WithAllowedAuthenticationSecretsNone() *ExternalAccessIntegrationModel {
	return e.WithAllowedAuthenticationSecretsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"none": tfconfig.BoolVariable(true),
	})))
}

func (e *ExternalAccessIntegrationModel) WithAllowedApiAuthenticationIntegrationsIntegrations(integrations []string) *ExternalAccessIntegrationModel {
	return e.WithAllowedApiAuthenticationIntegrationsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"integrations": tfconfig.SetVariable(
			collections.Map(integrations, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...,
		),
	})))
}

func (e *ExternalAccessIntegrationModel) WithAllowedApiAuthenticationIntegrationsNone() *ExternalAccessIntegrationModel {
	return e.WithAllowedApiAuthenticationIntegrationsValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"none": tfconfig.BoolVariable(true),
	})))
}
