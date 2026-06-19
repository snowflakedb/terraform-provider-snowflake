package model

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TableWithId(
	resourceName string,
	tableId sdk.SchemaObjectIdentifier,
	column []sdk.TableColumnSignature,
) *TableModel {
	return Table(resourceName, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), column)
}

func TableWithImplicitDependencies(
	resourceName string,
	tableName string,
	column []sdk.TableColumnSignature,
	schemaModel *SchemaModel,
	databaseModel *DatabaseModel,
) *TableModel {
	return Table(resourceName, "", "", tableName, column).
		WithDatabaseValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", databaseModel.ResourceReference()))).
		WithSchemaValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", schemaModel.ResourceReference())))
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
