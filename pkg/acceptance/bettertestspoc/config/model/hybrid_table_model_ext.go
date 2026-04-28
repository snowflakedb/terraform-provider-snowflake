package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HybridTableFromId(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	column []sdk.TableColumnSignature,
	primaryKey []sdk.TableColumnSignature,
) *HybridTableModel {
	return HybridTable(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), column, primaryKey)
}

// WithColumn satisfies the generated constructor's call for the complex list attribute.
func (h *HybridTableModel) WithColumn(column []sdk.TableColumnSignature) *HybridTableModel {
	maps := make([]tfconfig.Variable, len(column))
	for i, v := range column {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	h.Column = tfconfig.SetVariable(maps...)
	return h
}

// WithPrimaryKey satisfies the generated constructor's call for the complex list attribute.
// When called from the constructor, primaryKey is []sdk.TableColumnSignature where Name = column name.
func (h *HybridTableModel) WithPrimaryKey(primaryKey []sdk.TableColumnSignature) *HybridTableModel {
	keys := make([]tfconfig.Variable, len(primaryKey))
	for i, v := range primaryKey {
		keys[i] = tfconfig.StringVariable(v.Name)
	}
	h.PrimaryKey = tfconfig.SetVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"keys": tfconfig.ListVariable(keys...),
		}),
	)
	return h
}
