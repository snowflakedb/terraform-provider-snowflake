package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (p *AggregationPoliciesModel) WithEmptyIn() *AggregationPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (p *AggregationPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *AggregationPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (p *AggregationPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *AggregationPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (p *AggregationPoliciesModel) WithLimit(rows int) *AggregationPoliciesModel {
	return p.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (p *AggregationPoliciesModel) WithInAccount() *AggregationPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (p *AggregationPoliciesModel) WithRowsAndFrom(rows int, from string) *AggregationPoliciesModel {
	return p.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
