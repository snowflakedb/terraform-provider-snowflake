package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TableWithId(
	resourceName string,
	tableId sdk.SchemaObjectIdentifier,
	column []sdk.TableColumnSignature,
) *TableModel {
	return Table(resourceName, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), column)
}

func (t *TableModel) WithColumn(column []sdk.TableColumnSignature) *TableModel {
	maps := make([]tfconfig.Variable, len(column))
	for i, v := range column {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	t.Column = tfconfig.SetVariable(maps...)
	return t
}
