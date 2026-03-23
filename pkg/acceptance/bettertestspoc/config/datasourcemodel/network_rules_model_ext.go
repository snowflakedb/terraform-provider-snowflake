package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (n *NetworkRulesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *NetworkRulesModel {
	return n.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (n *NetworkRulesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *NetworkRulesModel {
	return n.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (n *NetworkRulesModel) WithInAccount() *NetworkRulesModel {
	return n.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (n *NetworkRulesModel) WithRows(rows int) *NetworkRulesModel {
	return n.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (n *NetworkRulesModel) WithRowsAndFrom(rows int, from string) *NetworkRulesModel {
	return n.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
