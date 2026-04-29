package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (h *HybridTableResourceAssert) HasColumnCount(expected int) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet("column.#", strconv.Itoa(expected)))
	return h
}

func (h *HybridTableResourceAssert) HasColumnName(index int, expected string) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.name", index), expected))
	return h
}

func (h *HybridTableResourceAssert) HasColumnType(index int, expected string) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.type", index), expected))
	return h
}

func (h *HybridTableResourceAssert) HasColumnNullable(index int, expected bool) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.nullable", index), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableResourceAssert) HasColumnComment(index int, expected string) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.comment", index), expected))
	return h
}

func (h *HybridTableResourceAssert) HasPrimaryKeyKeys(expected ...string) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet("primary_key.0.keys.#", strconv.Itoa(len(expected))))
	for i, k := range expected {
		h.AddAssertion(assert.ValueSet(fmt.Sprintf("primary_key.0.keys.%d", i), k))
	}
	return h
}

func (h *HybridTableResourceAssert) HasUniqueConstraintCount(expected int) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet("unique_constraint.#", strconv.Itoa(expected)))
	return h
}

func (h *HybridTableResourceAssert) HasForeignKeyCount(expected int) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet("foreign_key.#", strconv.Itoa(expected)))
	return h
}

func (h *HybridTableResourceAssert) HasColumnDefaultConstant(index int, expected string) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.default.#", index), "1"))
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.default.0.constant", index), expected))
	return h
}

func (h *HybridTableResourceAssert) HasColumnNoDefault(index int) *HybridTableResourceAssert {
	h.AddAssertion(assert.ValueSet(fmt.Sprintf("column.%d.default.#", index), "0"))
	return h
}
