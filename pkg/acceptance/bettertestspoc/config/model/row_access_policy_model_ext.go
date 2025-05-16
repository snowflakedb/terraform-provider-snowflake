package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func RowAccessPolicyDynamicArguments(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	body string,
) *RowAccessPolicyModel {
	m := &RowAccessPolicyModel{ResourceModelMeta: config.Meta(resourceName, resources.RowAccessPolicy)}
	m.WithDatabase(id.DatabaseName())
	m.WithSchema(id.SchemaName())
	m.WithName(id.Name())
	m.WithBody(body)
	return m.WithDynamicBlock(config.NewDynamicBlock("argument", "arguments", []string{"name", "type"}))
}

func RowAccessPolicyFromId(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	body string,
) *RowAccessPolicyModel {
	m := &RowAccessPolicyModel{ResourceModelMeta: config.Meta(resourceName, resources.RowAccessPolicy)}
	m.WithDatabase(id.DatabaseName())
	m.WithSchema(id.SchemaName())
	m.WithName(id.Name())
	m.WithBody(body)
	return m
}

func (r *RowAccessPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *RowAccessPolicyModel {
	maps := make([]tfconfig.Variable, len(argument))
	for i, v := range argument {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	r.Argument = tfconfig.SetVariable(maps...)
	return r
}
