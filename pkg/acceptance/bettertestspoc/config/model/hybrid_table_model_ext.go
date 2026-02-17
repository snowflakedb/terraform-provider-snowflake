// Manually written extensions for complex hybrid table configurations

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// ColumnOpts defines options for configuring a column
type ColumnOpts struct {
	Nullable   *bool
	PrimaryKey bool
	Unique     bool
	Comment    string
	Collate    string
	Default    *ColumnDefaultOpts
	Identity   *ColumnIdentityOpts
	ForeignKey *InlineForeignKeyOpts
}

// ColumnDefaultOpts defines options for column default values
type ColumnDefaultOpts struct {
	Expression string
	Sequence   string
}

// ColumnIdentityOpts defines options for identity columns
type ColumnIdentityOpts struct {
	StartNum int
	StepNum  int
}

// InlineForeignKeyOpts defines options for inline foreign key constraints
type InlineForeignKeyOpts struct {
	TableName  string
	ColumnName string
}

// WithColumnDescs adds multiple columns at once
func (h *HybridTableModel) WithColumnDescs(columns []ColumnDesc) *HybridTableModel {
	vars := make([]tfconfig.Variable, len(columns))
	for i, col := range columns {
		vars[i] = col.toVariable()
	}
	h.Column = tfconfig.ListVariable(vars...)
	return h
}

// ColumnDesc describes a column configuration
type ColumnDesc struct {
	Name       string
	DataType   string
	Nullable   *bool
	PrimaryKey bool
	Unique     bool
	Comment    string
	Collate    string
	Default    *ColumnDefaultOpts
	Identity   *ColumnIdentityOpts
	ForeignKey *InlineForeignKeyOpts
}

// toVariable converts a ColumnDesc to a tfconfig.Variable
func (c ColumnDesc) toVariable() tfconfig.Variable {
	m := map[string]tfconfig.Variable{
		"name": tfconfig.StringVariable(c.Name),
		"type": tfconfig.StringVariable(c.DataType),
	}

	// Handle nullable - default is true if not specified
	if c.Nullable != nil {
		m["nullable"] = tfconfig.BoolVariable(*c.Nullable)
	} else {
		m["nullable"] = tfconfig.BoolVariable(true)
	}

	if c.PrimaryKey {
		m["primary_key"] = tfconfig.BoolVariable(true)
	}

	if c.Unique {
		m["unique"] = tfconfig.BoolVariable(true)
	}

	if c.Comment != "" {
		m["comment"] = tfconfig.StringVariable(c.Comment)
	}

	if c.Collate != "" {
		m["collate"] = tfconfig.StringVariable(c.Collate)
	}

	if c.Default != nil {
		defaultMap := map[string]tfconfig.Variable{}
		if c.Default.Expression != "" {
			defaultMap["expression"] = tfconfig.StringVariable(c.Default.Expression)
		}
		if c.Default.Sequence != "" {
			defaultMap["sequence"] = tfconfig.StringVariable(c.Default.Sequence)
		}
		m["default"] = tfconfig.ListVariable(tfconfig.MapVariable(defaultMap))
	}

	if c.Identity != nil {
		identityMap := map[string]tfconfig.Variable{
			"start_num": tfconfig.IntegerVariable(c.Identity.StartNum),
			"step_num":  tfconfig.IntegerVariable(c.Identity.StepNum),
		}
		m["identity"] = tfconfig.ListVariable(tfconfig.MapVariable(identityMap))
	}

	if c.ForeignKey != nil {
		fkMap := map[string]tfconfig.Variable{
			"table_name":  tfconfig.StringVariable(c.ForeignKey.TableName),
			"column_name": tfconfig.StringVariable(c.ForeignKey.ColumnName),
		}
		m["foreign_key"] = tfconfig.ListVariable(tfconfig.MapVariable(fkMap))
	}

	return tfconfig.MapVariable(m)
}

// WithPrimaryKeyColumns adds an out-of-line primary key constraint
func (h *HybridTableModel) WithPrimaryKeyColumns(columns ...string) *HybridTableModel {
	return h.WithPrimaryKeyNamed("", columns...)
}

// WithPrimaryKeyNamed adds an out-of-line primary key constraint with a name
func (h *HybridTableModel) WithPrimaryKeyNamed(name string, columns ...string) *HybridTableModel {
	columnVars := make([]tfconfig.Variable, len(columns))
	for i, col := range columns {
		columnVars[i] = tfconfig.StringVariable(col)
	}

	m := map[string]tfconfig.Variable{
		"columns": tfconfig.ListVariable(columnVars...),
	}

	if name != "" {
		m["name"] = tfconfig.StringVariable(name)
	}

	h.PrimaryKey = tfconfig.ListVariable(tfconfig.MapVariable(m))
	return h
}

// WithIndexes adds indexes to the hybrid table
func (h *HybridTableModel) WithIndexes(indexes []IndexDesc) *HybridTableModel {
	vars := make([]tfconfig.Variable, len(indexes))
	for i, idx := range indexes {
		columnVars := make([]tfconfig.Variable, len(idx.Columns))
		for j, col := range idx.Columns {
			columnVars[j] = tfconfig.StringVariable(col)
		}

		m := map[string]tfconfig.Variable{
			"name":    tfconfig.StringVariable(idx.Name),
			"columns": tfconfig.ListVariable(columnVars...),
		}

		vars[i] = tfconfig.MapVariable(m)
	}
	h.Index = tfconfig.SetVariable(vars...)
	return h
}

// IndexDesc describes an index configuration
type IndexDesc struct {
	Name    string
	Columns []string
}

// WithUniqueConstraints adds unique constraints to the hybrid table
func (h *HybridTableModel) WithUniqueConstraints(constraints []UniqueConstraintDesc) *HybridTableModel {
	vars := make([]tfconfig.Variable, len(constraints))
	for i, uc := range constraints {
		columnVars := make([]tfconfig.Variable, len(uc.Columns))
		for j, col := range uc.Columns {
			columnVars[j] = tfconfig.StringVariable(col)
		}

		m := map[string]tfconfig.Variable{
			"columns": tfconfig.ListVariable(columnVars...),
		}

		if uc.Name != "" {
			m["name"] = tfconfig.StringVariable(uc.Name)
		}

		vars[i] = tfconfig.MapVariable(m)
	}
	h.UniqueConstraint = tfconfig.SetVariable(vars...)
	return h
}

// UniqueConstraintDesc describes a unique constraint configuration
type UniqueConstraintDesc struct {
	Name    string
	Columns []string
}

// WithForeignKeys adds foreign key constraints to the hybrid table
func (h *HybridTableModel) WithForeignKeys(foreignKeys []ForeignKeyDesc) *HybridTableModel {
	vars := make([]tfconfig.Variable, len(foreignKeys))
	for i, fk := range foreignKeys {
		columnVars := make([]tfconfig.Variable, len(fk.Columns))
		for j, col := range fk.Columns {
			columnVars[j] = tfconfig.StringVariable(col)
		}

		refColumnVars := make([]tfconfig.Variable, len(fk.ReferencesColumns))
		for j, col := range fk.ReferencesColumns {
			refColumnVars[j] = tfconfig.StringVariable(col)
		}

		m := map[string]tfconfig.Variable{
			"columns":            tfconfig.ListVariable(columnVars...),
			"references_table":   tfconfig.StringVariable(fk.ReferencesTable),
			"references_columns": tfconfig.ListVariable(refColumnVars...),
		}

		if fk.Name != "" {
			m["name"] = tfconfig.StringVariable(fk.Name)
		}

		vars[i] = tfconfig.MapVariable(m)
	}
	h.ForeignKey = tfconfig.SetVariable(vars...)
	return h
}

// ForeignKeyDesc describes a foreign key constraint configuration
type ForeignKeyDesc struct {
	Name              string
	Columns           []string
	ReferencesTable   string
	ReferencesColumns []string
}

// Helper function to create a pointer to bool (for nullable)
func Bool(b bool) *bool {
	return &b
}
