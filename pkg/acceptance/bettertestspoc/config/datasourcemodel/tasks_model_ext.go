package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (t *TasksModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *TasksModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (t *TasksModel) WithInAccount() *TasksModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (t *TasksModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *TasksModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (t *TasksModel) WithLimitRows(rows int) *TasksModel {
	return t.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (t *TasksModel) WithLimitRowsAndFrom(rows int, from string) *TasksModel {
	return t.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
