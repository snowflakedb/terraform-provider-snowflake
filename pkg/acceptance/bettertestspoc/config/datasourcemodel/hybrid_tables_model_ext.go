package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// WithLike adds a LIKE filter with pattern matching
func (h *HybridTablesModel) WithLike(pattern string) *HybridTablesModel {
	h.Like = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"pattern": tfconfig.StringVariable(pattern),
	})
	return h
}

// WithInAccount adds an IN ACCOUNT filter
func (h *HybridTablesModel) WithInAccount(account bool) *HybridTablesModel {
	h.In = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"account": tfconfig.BoolVariable(account),
	})
	return h
}

// WithInDatabase adds an IN DATABASE filter
func (h *HybridTablesModel) WithInDatabase(database string) *HybridTablesModel {
	h.In = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"database": tfconfig.StringVariable(database),
	})
	return h
}

// WithInSchema adds an IN SCHEMA filter
func (h *HybridTablesModel) WithInSchema(schema string) *HybridTablesModel {
	h.In = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"schema": tfconfig.StringVariable(schema),
	})
	return h
}

// WithLimit adds a LIMIT filter
func (h *HybridTablesModel) WithLimit(rows int) *HybridTablesModel {
	h.Limit = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"rows": tfconfig.IntegerVariable(rows),
	})
	return h
}

// WithLimitFrom adds a LIMIT filter with FROM cursor
func (h *HybridTablesModel) WithLimitFrom(rows int, from string) *HybridTablesModel {
	h.Limit = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"rows": tfconfig.IntegerVariable(rows),
		"from": tfconfig.StringVariable(from),
	})
	return h
}
