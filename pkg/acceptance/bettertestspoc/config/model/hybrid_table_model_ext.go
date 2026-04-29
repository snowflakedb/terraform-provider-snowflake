package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HybridTableColumnConfig is a richer column definition used in tests that need
// column fields beyond name and type (e.g. comment, nullable, collate).
type HybridTableColumnConfig struct {
	Name     string
	Type     string
	Comment  string
	Nullable *bool
	Collate  string
}

// WithColumnConfigs sets the column list from richer column definitions.
// Use instead of WithColumn when tests require comment, nullable, or collate.
func (h *HybridTableModel) WithColumnConfigs(columns []HybridTableColumnConfig) *HybridTableModel {
	maps := make([]tfconfig.Variable, len(columns))
	for i, col := range columns {
		m := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(col.Name),
			"type": tfconfig.StringVariable(col.Type),
		}
		if col.Comment != "" {
			m["comment"] = tfconfig.StringVariable(col.Comment)
		}
		if col.Nullable != nil {
			m["nullable"] = tfconfig.BoolVariable(*col.Nullable)
		}
		if col.Collate != "" {
			m["collate"] = tfconfig.StringVariable(col.Collate)
		}
		maps[i] = tfconfig.MapVariable(m)
	}
	h.Column = tfconfig.SetVariable(maps...)
	return h
}

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
