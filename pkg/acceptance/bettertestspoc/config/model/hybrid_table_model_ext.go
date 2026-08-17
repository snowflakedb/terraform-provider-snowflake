package model

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HybridTableColumnDefaultConfig represents the default value block of a column.
// Tests should set at most one of Constant, Expression, or Sequence; mutual
// exclusivity is enforced by the resource at apply time, so tests deliberately
// constructing invalid combinations remain possible.
type HybridTableColumnDefaultConfig struct {
	Constant   *string
	Expression *string
	Sequence   *string
}

// HybridTableColumnConfig is a richer column definition used in tests that need
// column fields beyond name and type (e.g. comment, not_null, collate, default).
type HybridTableColumnConfig struct {
	Name    string
	Type    string
	Comment string
	NotNull *bool
	Collate string
	Default *HybridTableColumnDefaultConfig
}

// WithColumnConfigs sets the column list from richer column definitions.
// Use instead of WithColumn when tests require comment, not_null, collate, or default.
func (h *HybridTableModel) WithColumnConfigs(columns []HybridTableColumnConfig) *HybridTableModel {
	objs := make([]tfconfig.Variable, len(columns))
	for i, col := range columns {
		m := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(col.Name),
			"type": tfconfig.StringVariable(col.Type),
		}
		if col.Comment != "" {
			m["comment"] = tfconfig.StringVariable(col.Comment)
		}
		if col.NotNull != nil {
			m["not_null"] = tfconfig.BoolVariable(*col.NotNull)
		}
		if col.Collate != "" {
			m["collate"] = tfconfig.StringVariable(col.Collate)
		}
		if col.Default != nil {
			defMap := map[string]tfconfig.Variable{}
			if col.Default.Constant != nil {
				defMap["constant"] = tfconfig.StringVariable(*col.Default.Constant)
			}
			if col.Default.Expression != nil {
				defMap["expression"] = tfconfig.StringVariable(*col.Default.Expression)
			}
			if col.Default.Sequence != nil {
				defMap["sequence"] = tfconfig.StringVariable(*col.Default.Sequence)
			}
			if len(defMap) > 0 {
				m["default"] = tfconfig.ListVariable(tfconfig.ObjectVariable(defMap))
			}
		}
		objs[i] = tfconfig.ObjectVariable(m)
	}
	h.Column = tfconfig.SetVariable(objs...)
	return h
}

func buildColumnConfigs(column, primaryKey []sdk.TableColumnSignature) []HybridTableColumnConfig {
	pkNames := make(map[string]struct{}, len(primaryKey))
	for _, k := range primaryKey {
		pkNames[k.Name] = struct{}{}
	}
	return collections.Map(column, func(v sdk.TableColumnSignature) HybridTableColumnConfig {
		cfg := HybridTableColumnConfig{Name: v.Name, Type: v.Type.ToSql()}
		if _, isPK := pkNames[v.Name]; isPK {
			cfg.NotNull = new(true)
		}
		return cfg
	})
}

func HybridTableFromId(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	column []sdk.TableColumnSignature,
	primaryKey []sdk.TableColumnSignature,
) *HybridTableModel {
	return HybridTable(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), column, primaryKey).
		WithColumnConfigs(buildColumnConfigs(column, primaryKey))
}

func HybridTableWithImplicitDependencies(
	resourceName string,
	name string,
	column []sdk.TableColumnSignature,
	primaryKey []sdk.TableColumnSignature,
	schemaModel *SchemaModel,
	databaseModel *DatabaseModel,
) *HybridTableModel {
	return HybridTable(resourceName, "", "", name, column, primaryKey).
		WithColumnConfigs(buildColumnConfigs(column, primaryKey)).
		WithDatabaseValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", databaseModel.ResourceReference()))).
		WithSchemaValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", schemaModel.ResourceReference())))
}

// WithColumn satisfies the generated constructor's call for the complex list attribute.
func (h *HybridTableModel) WithColumn(column []sdk.TableColumnSignature) *HybridTableModel {
	return h.WithColumnConfigs(collections.Map(column, func(v sdk.TableColumnSignature) HybridTableColumnConfig {
		return HybridTableColumnConfig{Name: v.Name, Type: v.Type.ToSql()}
	}))
}

// WithPrimaryKeyConstraint sets the primary_key_constraint block from column
// signatures. Only Name is used — Type is ignored.
func (h *HybridTableModel) WithPrimaryKeyConstraint(primaryKey []sdk.TableColumnSignature) *HybridTableModel {
	cols := collections.Map(primaryKey, func(v sdk.TableColumnSignature) tfconfig.Variable {
		return tfconfig.StringVariable(v.Name)
	})
	h.PrimaryKeyConstraint = tfconfig.SetVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"columns": tfconfig.ListVariable(cols...),
		}),
	)
	return h
}

type HybridTableUniqueConstraintConfig struct {
	Name    string
	Columns []string
}

// WithUniqueConstraints sets the unique_constraint block from one or more definitions.
// Supports both named and anonymous constraints in a single call.
func (h *HybridTableModel) WithUniqueConstraints(constraints ...HybridTableUniqueConstraintConfig) *HybridTableModel {
	objs := make([]tfconfig.Variable, len(constraints))
	for i, uc := range constraints {
		colVars := collections.Map(uc.Columns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
		m := map[string]tfconfig.Variable{
			"columns": tfconfig.ListVariable(colVars...),
		}
		if uc.Name != "" {
			m["name"] = tfconfig.StringVariable(uc.Name)
		}
		objs[i] = tfconfig.ObjectVariable(m)
	}
	h.UniqueConstraint = tfconfig.SetVariable(objs...)
	return h
}

type HybridTableIndexConfig struct {
	Name           string
	Columns        []string // required; the schema enforces MinItems:1
	IncludeColumns []string // optional (the INCLUDE payload)
}

// WithIndex sets the index block from one or more index definitions.
func (h *HybridTableModel) WithIndex(indexes ...HybridTableIndexConfig) *HybridTableModel {
	objs := make([]tfconfig.Variable, len(indexes))
	for i, idx := range indexes {
		colVars := collections.Map(idx.Columns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
		m := map[string]tfconfig.Variable{
			"name":    tfconfig.StringVariable(idx.Name),
			"columns": tfconfig.ListVariable(colVars...),
		}
		if len(idx.IncludeColumns) > 0 {
			incVars := collections.Map(idx.IncludeColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
			m["include_columns"] = tfconfig.SetVariable(incVars...)
		}
		objs[i] = tfconfig.ObjectVariable(m)
	}
	h.Index = tfconfig.SetVariable(objs...)
	return h
}

type HybridTableForeignKeyConstraintConfig struct {
	Name       string
	Columns    []string
	TableName  string
	RefColumns []string
}

// WithForeignKeyConstraints sets the foreign_key_constraint block from one or more definitions.
func (h *HybridTableModel) WithForeignKeyConstraints(constraints ...HybridTableForeignKeyConstraintConfig) *HybridTableModel {
	objs := make([]tfconfig.Variable, len(constraints))
	for i, fk := range constraints {
		lcVars := collections.Map(fk.Columns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
		rcVars := collections.Map(fk.RefColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
		m := map[string]tfconfig.Variable{
			"columns":     tfconfig.ListVariable(lcVars...),
			"table_name":  tfconfig.StringVariable(fk.TableName),
			"ref_columns": tfconfig.ListVariable(rcVars...),
		}
		if fk.Name != "" {
			m["name"] = tfconfig.StringVariable(fk.Name)
		}
		objs[i] = tfconfig.ObjectVariable(m)
	}
	h.ForeignKeyConstraint = tfconfig.SetVariable(objs...)
	return h
}
