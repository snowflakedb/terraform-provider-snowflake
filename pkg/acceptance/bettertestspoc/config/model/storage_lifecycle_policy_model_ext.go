package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func StorageLifecyclePolicyDynamicArguments(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	body string,
) *StorageLifecyclePolicyModel {
	s := &StorageLifecyclePolicyModel{ResourceModelMeta: config.Meta(resourceName, resources.StorageLifecyclePolicy)}
	s.WithDatabase(id.DatabaseName())
	s.WithSchema(id.SchemaName())
	s.WithName(id.Name())
	s.WithBody(body)
	return s.WithDynamicBlock(config.NewDynamicBlock("argument", "arguments", []string{"name", "type"}))
}

func (s *StorageLifecyclePolicyModel) WithArgument(argument []sdk.TableColumnSignature) *StorageLifecyclePolicyModel {
	maps := make([]tfconfig.Variable, len(argument))
	for i, v := range argument {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	s.Argument = tfconfig.SetVariable(maps...)
	return s
}
