package resourceassert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (h *HybridTableResourceAssert) HasColumns(columns []sdk.TableColumnSignature) *HybridTableResourceAssert {
	h.ValueSet("column.#", strconv.Itoa(len(columns)))
	for i, col := range columns {
		h.ValueSet(fmt.Sprintf("column.%d.name", i), col.Name)
		// Read substitutes the user's config spelling when the DESCRIBE value is
		// canonically equivalent (see buildHybridColumnStateFromDescribe), so
		// state holds the same form the model writes to HCL — Type.ToSql().
		h.ValueSet(fmt.Sprintf("column.%d.type", i), col.Type.ToSql())
	}
	return h
}

// HasColumnConfigs asserts all per-column fields (name, type, not_null, comment,
// collate, default) in one call. Use this when the model was built with
// WithColumnConfigs. A nil NotNull is treated as the schema default (false).
func (h *HybridTableResourceAssert) HasColumnConfigs(columns []model.HybridTableColumnConfig) *HybridTableResourceAssert {
	h.ValueSet("column.#", strconv.Itoa(len(columns)))
	for i, col := range columns {
		h.ValueSet(fmt.Sprintf("column.%d.name", i), col.Name)
		h.ValueSet(fmt.Sprintf("column.%d.type", i), col.Type)
		notNull := false
		if col.NotNull != nil {
			notNull = *col.NotNull
		}
		h.ValueSet(fmt.Sprintf("column.%d.not_null", i), strconv.FormatBool(notNull))
		h.ValueSet(fmt.Sprintf("column.%d.comment", i), col.Comment)
		h.ValueSet(fmt.Sprintf("column.%d.collate", i), col.Collate)
		if col.Default != nil {
			h.ValueSet(fmt.Sprintf("column.%d.default.#", i), "1")
			constant, expression, sequence := "", "", ""
			if col.Default.Constant != nil {
				constant = *col.Default.Constant
			}
			if col.Default.Expression != nil {
				expression = *col.Default.Expression
			}
			if col.Default.Sequence != nil {
				sequence = *col.Default.Sequence
			}
			h.ValueSet(fmt.Sprintf("column.%d.default.0.constant", i), constant)
			h.ValueSet(fmt.Sprintf("column.%d.default.0.expression", i), expression)
			h.ValueSet(fmt.Sprintf("column.%d.default.0.sequence", i), sequence)
		} else {
			h.ValueSet(fmt.Sprintf("column.%d.default.#", i), "0")
		}
	}
	return h
}

func (h *HybridTableResourceAssert) HasPrimaryKeyColumns(expected ...string) *HybridTableResourceAssert {
	h.ValueSet("primary_key_constraint.0.columns.#", strconv.Itoa(len(expected)))
	for i, k := range expected {
		h.ValueSet(fmt.Sprintf("primary_key_constraint.0.columns.%d", i), k)
	}
	return h
}

func (h *HybridTableResourceAssert) HasUniqueConstraints(constraints ...model.HybridTableUniqueConstraintConfig) *HybridTableResourceAssert {
	h.ValueSet("unique_constraint.#", strconv.Itoa(len(constraints)))
	for _, uc := range constraints {
		attrs := map[string]string{
			"columns.#": strconv.Itoa(len(uc.Columns)),
		}
		if uc.Name != nil {
			attrs["name"] = *uc.Name
		}
		for i, col := range uc.Columns {
			attrs[fmt.Sprintf("columns.%d", i)] = col
		}
		h.SetContainsElemNested("unique_constraint", attrs)
	}
	return h
}

func (h *HybridTableResourceAssert) HasForeignKeyConstraints(constraints ...model.HybridTableForeignKeyConstraintConfig) *HybridTableResourceAssert {
	h.ValueSet("foreign_key_constraint.#", strconv.Itoa(len(constraints)))
	for _, fk := range constraints {
		attrs := map[string]string{
			"columns.#":     strconv.Itoa(len(fk.Columns)),
			"table_name":    fk.TableName,
			"ref_columns.#": strconv.Itoa(len(fk.RefColumns)),
		}
		if fk.Name != nil {
			attrs["name"] = *fk.Name
		}
		for i, col := range fk.Columns {
			attrs[fmt.Sprintf("columns.%d", i)] = col
		}
		for i, col := range fk.RefColumns {
			attrs[fmt.Sprintf("ref_columns.%d", i)] = col
		}
		h.SetContainsElemNested("foreign_key_constraint", attrs)
	}
	return h
}

func (h *HybridTableResourceAssert) HasIndexes(indexes ...model.HybridTableIndexConfig) *HybridTableResourceAssert {
	h.ValueSet("index.#", strconv.Itoa(len(indexes)))
	for _, idx := range indexes {
		attrs := map[string]string{
			"name":              idx.Name,
			"columns.#":         strconv.Itoa(len(idx.Columns)),
			"include_columns.#": strconv.Itoa(len(idx.IncludeColumns)),
		}
		for i, col := range idx.Columns {
			attrs[fmt.Sprintf("columns.%d", i)] = col
		}
		for _, col := range idx.IncludeColumns {
			// Nested TypeSet keys are hashes; match the resource's include_columns Set func.
			attrs[fmt.Sprintf("include_columns.%d", schema.HashString(strings.ToUpper(col)))] = col
		}
		h.SetContainsElemNested("index", attrs)
	}
	return h
}
