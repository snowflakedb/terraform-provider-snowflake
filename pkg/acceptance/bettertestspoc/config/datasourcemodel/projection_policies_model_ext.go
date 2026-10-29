package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (p *ProjectionPoliciesModel) WithEmptyIn() *ProjectionPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (p *ProjectionPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *ProjectionPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (p *ProjectionPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *ProjectionPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (p *ProjectionPoliciesModel) WithLimit(rows int) *ProjectionPoliciesModel {
	return p.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}
