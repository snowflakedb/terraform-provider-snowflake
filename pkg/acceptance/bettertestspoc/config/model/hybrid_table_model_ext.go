package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

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
// column fields beyond name and type (e.g. comment, nullable, collate, default).
type HybridTableColumnConfig struct {
	Name     string
	Type     string
	Comment  string
	Nullable *bool
	Collate  string
	Default  *HybridTableColumnDefaultConfig
}

// WithColumnConfigs sets the column list from richer column definitions.
// Use instead of WithColumn when tests require comment, nullable, collate, or default.
//
// Uses ObjectVariable rather than MapVariable for both column maps and the default
// block, because terraform-plugin-testing's MapVariable.MarshalJSON requires every
// value to be the same underlying type, and these blocks mix string + bool + list.
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
		if col.Nullable != nil {
			m["nullable"] = tfconfig.BoolVariable(*col.Nullable)
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
// Only the Name field of each TableColumnSignature is used — the Type field is
// ignored. Tests that build the PK separately should prefer WithPrimaryKeyNames,
// which takes plain column names and avoids the misleading Type-less signature.
func (h *HybridTableModel) WithPrimaryKey(primaryKey []sdk.TableColumnSignature) *HybridTableModel {
	names := make([]string, len(primaryKey))
	for i, v := range primaryKey {
		names[i] = v.Name
	}
	return h.WithPrimaryKeyNames(names...)
}

// WithPrimaryKeyNames sets the primary_key block from a slice of column names.
// This is the preferred form for tests that compose models incrementally: the
// resource's primary_key.keys is just a list of column names, not signatures.
func (h *HybridTableModel) WithPrimaryKeyNames(names ...string) *HybridTableModel {
	keys := collections.Map(names, func(n string) tfconfig.Variable { return tfconfig.StringVariable(n) })
	h.PrimaryKey = tfconfig.SetVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"keys": tfconfig.ListVariable(keys...),
		}),
	)
	return h
}

// WithUniqueConstraint sets a single unnamed unique constraint on the given columns.
func (h *HybridTableModel) WithUniqueConstraint(columns []string) *HybridTableModel {
	colVars := collections.Map(columns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
	h.UniqueConstraint = tfconfig.SetVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"columns": tfconfig.ListVariable(colVars...),
		}),
	)
	return h
}

// HybridTableUniqueConstraintConfig is a unique constraint definition for tests.
// Name is optional; if empty, Snowflake generates a SYS_CONSTRAINT_-prefixed name.
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

// HybridTableIndexConfig is a single secondary-index definition for tests.
// IncludeColumns is optional (the INCLUDE payload).
type HybridTableIndexConfig struct {
	Name           string
	Columns        []string // required; the schema enforces MinItems:1
	IncludeColumns []string // optional (the INCLUDE payload)
}

// WithIndex sets the index block from one or more index definitions.
//
// Uses ObjectVariable rather than MapVariable because each index block mixes a
// string (name) with list values (columns, include_columns), which MapVariable's
// MarshalJSON rejects ("maps must contain the same type").
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

// WithForeignKey sets a single unnamed foreign key constraint. localColumns are the columns
// in this table, refTableId is the fully-qualified name of the referenced table,
// and refColumns are the columns in the referenced table.
//
// Uses ObjectVariable instead of MapVariable for both the outer foreign-key block
// and the inner references block, because each mixes string and list values, which
// MapVariable's MarshalJSON rejects ("maps must contain the same type").
func (h *HybridTableModel) WithForeignKey(localColumns []string, refTableId string, refColumns []string) *HybridTableModel {
	lcVars := collections.Map(localColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
	rcVars := collections.Map(refColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
	h.ForeignKey = tfconfig.SetVariable(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"columns": tfconfig.ListVariable(lcVars...),
			"references": tfconfig.ListVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"table_id": tfconfig.StringVariable(refTableId),
					"columns":  tfconfig.ListVariable(rcVars...),
				}),
			),
		}),
	)
	return h
}

// WithNamedForeignKey is like WithForeignKey but includes an explicit constraint name.
func (h *HybridTableModel) WithNamedForeignKey(name string, localColumns []string, refTableId string, refColumns []string) *HybridTableModel {
	lcVars := collections.Map(localColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
	rcVars := collections.Map(refColumns, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })
	h.ForeignKey = tfconfig.SetVariable(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"name":    tfconfig.StringVariable(name),
			"columns": tfconfig.ListVariable(lcVars...),
			"references": tfconfig.ListVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"table_id": tfconfig.StringVariable(refTableId),
					"columns":  tfconfig.ListVariable(rcVars...),
				}),
			),
		}),
	)
	return h
}
