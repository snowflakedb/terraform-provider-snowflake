package resourceshowoutputassert

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// HybridTableDescribeOutputRowAssert asserts fields of a single row in the
// describe_output list. Because hybrid table describe output has one row per
// column, callers must supply the 0-based row index (= column position in the
// physical column order returned by DESCRIBE TABLE).
type HybridTableDescribeOutputRowAssert struct {
	*assert.ResourceAssert
}

func HybridTableDescribeOutputRow(t *testing.T, name string, rowIndex int) *HybridTableDescribeOutputRowAssert {
	t.Helper()
	return &HybridTableDescribeOutputRowAssert{
		ResourceAssert: assert.NewResourceDescribeOutputAssertAtRow(name, rowIndex),
	}
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (h *HybridTableDescribeOutputRowAssert) HasName(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("name", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasType(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("type", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasCollation(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("collation", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasKind(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("kind", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasIsNullable(expected bool) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("is_nullable", strconv.FormatBool(expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasDefault(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("default", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPrimaryKey(expected bool) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("primary_key", strconv.FormatBool(expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasUniqueKey(expected bool) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("unique_key", strconv.FormatBool(expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasCheck(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("check", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasExpression(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("expression", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasComment(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("comment", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPolicyName(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("policy_name", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPrivacyDomain(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("privacy_domain", expected)
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasSchemaEvolutionRecord(expected string) *HybridTableDescribeOutputRowAssert {
	h.StringValueSet("schema_evolution_record", expected)
	return h
}
