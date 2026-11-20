package model

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
	columns := make([]tfconfig.Variable, len(column))
	for idx, v := range column {
		columns[idx] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	t.Column = tfconfig.ListVariable(columns...)
	return t
}

// TableColumnConfig is a richer column definition used in tests that need
// column fields beyond name and type (comment, nullable).
type TableColumnConfig struct {
	Name     string
	Type     datatypes.DataType
	Comment  string
	Nullable *bool
}

func (t *TableModel) WithColumnConfigs(columns ...TableColumnConfig) *TableModel {
	vars := make([]tfconfig.Variable, len(columns))
	for idx, c := range columns {
		m := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(c.Name),
			"type": tfconfig.StringVariable(c.Type.ToSql()),
		}
		if c.Comment != "" {
			m["comment"] = tfconfig.StringVariable(c.Comment)
		}
		if c.Nullable != nil {
			m["nullable"] = tfconfig.BoolVariable(*c.Nullable)
		}
		vars[idx] = tfconfig.ObjectVariable(m)
	}
	t.Column = tfconfig.ListVariable(vars...)
	return t
}

func (t *TableModel) WithClusterBy(clusterBy ...string) *TableModel {
	return t.WithClusterByValue(tfconfig.ListVariable(collections.Map(clusterBy, func(s string) tfconfig.Variable {
		return tfconfig.StringVariable(s)
	})...))
}
