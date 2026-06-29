package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (h *HybridTableResourceAssert) HasColumnNullable(index int, expected bool) *HybridTableResourceAssert {
	h.ValueSet(fmt.Sprintf("column.%d.nullable", index), strconv.FormatBool(expected))
	return h
}

func (h *HybridTableResourceAssert) HasColumnComment(index int, expected string) *HybridTableResourceAssert {
	h.ValueSet(fmt.Sprintf("column.%d.comment", index), expected)
	return h
}

func (h *HybridTableResourceAssert) HasPrimaryKeyKeys(expected ...string) *HybridTableResourceAssert {
	h.ValueSet("primary_key.0.keys.#", strconv.Itoa(len(expected)))
	for i, k := range expected {
		h.ValueSet(fmt.Sprintf("primary_key.0.keys.%d", i), k)
	}
	return h
}

func (h *HybridTableResourceAssert) HasUniqueConstraintCount(expected int) *HybridTableResourceAssert {
	h.ValueSet("unique_constraint.#", strconv.Itoa(expected))
	return h
}

func (h *HybridTableResourceAssert) HasForeignKeyCount(expected int) *HybridTableResourceAssert {
	h.ValueSet("foreign_key.#", strconv.Itoa(expected))
	return h
}

func (h *HybridTableResourceAssert) HasColumnDefaultConstant(index int, expected string) *HybridTableResourceAssert {
	h.ValueSet(fmt.Sprintf("column.%d.default.#", index), "1")
	h.ValueSet(fmt.Sprintf("column.%d.default.0.constant", index), expected)
	return h
}

func (h *HybridTableResourceAssert) HasColumnNoDefault(index int) *HybridTableResourceAssert {
	h.ValueSet(fmt.Sprintf("column.%d.default.#", index), "0")
	return h
}
