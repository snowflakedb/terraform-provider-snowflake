package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func NetworkRuleFromId(
	id sdk.SchemaObjectIdentifier,
	mode sdk.NetworkRuleMode,
	type_ sdk.NetworkRuleType,
	valueList []string,
) *NetworkRuleModel {
	n := &NetworkRuleModel{ResourceModelMeta: config.Meta("test", resources.NetworkRule)}
	n.WithDatabase(id.DatabaseName())
	n.WithSchema(id.SchemaName())
	n.WithName(id.Name())
	n.WithMode(string(mode))
	n.WithType(string(type_))
	n.WithValueList(valueList)
	return n
}

func (n *NetworkRuleModel) WithValueList(valueList []string) *NetworkRuleModel {
	if len(valueList) == 0 {
		return n.WithValueListValue(config.EmptyListVariable())
	}
	return n.WithValueListValue(
		tfconfig.SetVariable(
			collections.Map(valueList, func(v string) tfconfig.Variable { return tfconfig.StringVariable(v) })...,
		),
	)
}
