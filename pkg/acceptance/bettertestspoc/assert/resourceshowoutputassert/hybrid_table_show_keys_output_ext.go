package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

type HybridTableShowKeysOutputRowAssert struct {
	*assert.ResourceAssert
}

func HybridTableShowKeysOutputRow(t *testing.T, name string, rowIndex int) *HybridTableShowKeysOutputRowAssert {
	t.Helper()
	return &HybridTableShowKeysOutputRowAssert{
		ResourceAssert: assert.NewResourceShowOutputAssertAtRowWithPath(name, resources.ShowKeysOutputAttributeName, rowIndex),
	}
}

func HybridTablesDatasourceShowKeysOutputRow(t *testing.T, name string, rowIndex int) *HybridTableShowKeysOutputRowAssert {
	t.Helper()
	return &HybridTableShowKeysOutputRowAssert{
		ResourceAssert: assert.NewResourceShowOutputAssertAtRowWithPath(name, "hybrid_tables.0.show_keys_output", rowIndex),
	}
}

func (h *HybridTableShowKeysOutputRowAssert) HasName(expected string) *HybridTableShowKeysOutputRowAssert {
	h.StringValueSet("name", expected)
	return h
}

func (h *HybridTableShowKeysOutputRowAssert) HasNameNotEmpty() *HybridTableShowKeysOutputRowAssert {
	h.ValuePresent("name")
	return h
}

func (h *HybridTableShowKeysOutputRowAssert) HasKind(expected string) *HybridTableShowKeysOutputRowAssert {
	h.StringValueSet("kind", expected)
	return h
}

func (h *HybridTableShowKeysOutputRowAssert) HasColumns(expected ...string) *HybridTableShowKeysOutputRowAssert {
	h.ListContainsExactlyStringValuesInOrder("columns", expected...)
	return h
}

func (h *HybridTableShowKeysOutputRowAssert) HasReferencedTable(expected string) *HybridTableShowKeysOutputRowAssert {
	h.StringValueSet("referenced_table", expected)
	return h
}

func (h *HybridTableShowKeysOutputRowAssert) HasReferencedColumns(expected ...string) *HybridTableShowKeysOutputRowAssert {
	h.ListContainsExactlyStringValuesInOrder("referenced_columns", expected...)
	return h
}
