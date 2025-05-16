package model

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (n *NetworkPolicyModel) WithAllowedNetworkRules(rules ...sdk.SchemaObjectIdentifier) *NetworkPolicyModel {
	return n.WithAllowedNetworkRuleListValue(
		tfconfig.SetVariable(
			collections.Map(rules, func(rule sdk.SchemaObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(rule.FullyQualifiedName())
			})...,
		),
	)
}

func (n *NetworkPolicyModel) WithAllowedNetworkRulesUnquotedNamePart(rules ...sdk.SchemaObjectIdentifier) *NetworkPolicyModel {
	return n.WithAllowedNetworkRuleListValue(
		tfconfig.SetVariable(
			collections.Map(rules, func(rule sdk.SchemaObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(fmt.Sprintf("\"%s\".\"%s\".%s", rule.DatabaseName(), rule.SchemaName(), rule.Name()))
			})...,
		),
	)
}

func (n *NetworkPolicyModel) WithBlockedNetworkRules(rules ...sdk.SchemaObjectIdentifier) *NetworkPolicyModel {
	return n.WithBlockedNetworkRuleListValue(
		tfconfig.SetVariable(
			collections.Map(rules, func(rule sdk.SchemaObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(rule.FullyQualifiedName())
			})...,
		),
	)
}

func (n *NetworkPolicyModel) WithBlockedNetworkRulesUnquotedNamePart(rules ...sdk.SchemaObjectIdentifier) *NetworkPolicyModel {
	return n.WithBlockedNetworkRuleListValue(
		tfconfig.SetVariable(
			collections.Map(rules, func(rule sdk.SchemaObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(fmt.Sprintf("\"%s\".\"%s\".%s", rule.DatabaseName(), rule.SchemaName(), rule.Name()))
			})...,
		),
	)
}

func (n *NetworkPolicyModel) WithAllowedIps(ips ...string) *NetworkPolicyModel {
	return n.WithAllowedIpListValue(
		tfconfig.SetVariable(
			collections.Map(ips, func(ip string) tfconfig.Variable { return tfconfig.StringVariable(ip) })...,
		),
	)
}

func (n *NetworkPolicyModel) WithBlockedIps(ips ...string) *NetworkPolicyModel {
	return n.WithBlockedIpListValue(
		tfconfig.SetVariable(
			collections.Map(ips, func(ip string) tfconfig.Variable { return tfconfig.StringVariable(ip) })...,
		),
	)
}
