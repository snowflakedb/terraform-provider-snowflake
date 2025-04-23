package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func MaskingPolicyDynamicArguments(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	body string,
	returnDataType sdk.DataType,
) *MaskingPolicyModel {
	m := &MaskingPolicyModel{ResourceModelMeta: config.Meta(resourceName, resources.MaskingPolicy)}
	m.WithDatabase(id.DatabaseName())
	m.WithSchema(id.SchemaName())
	m.WithName(id.Name())
	m.WithBody(body)
	m.WithReturnDataType(string(returnDataType))
	return m.WithDynamicBlock(config.NewDynamicBlock("argument", "arguments", []string{"name", "type"}))
}

func (p *MaskingPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *MaskingPolicyModel {
	maps := make([]tfconfig.Variable, len(argument))
	for i, v := range argument {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	p.Argument = tfconfig.SetVariable(maps...)
	return p
}
