package resourceshowoutputassert

import (
	"fmt"
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
	rowIndex int
}

func HybridTableDescribeOutputRow(t *testing.T, name string, rowIndex int) *HybridTableDescribeOutputRowAssert {
	t.Helper()
	return &HybridTableDescribeOutputRowAssert{
		ResourceAssert: assert.NewResourceAssert(name, fmt.Sprintf("describe_output.%d", rowIndex)),
		rowIndex:       rowIndex,
	}
}

func (h *HybridTableDescribeOutputRowAssert) field(name string) string {
	return fmt.Sprintf("describe_output.%d.%s", h.rowIndex, name)
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (h *HybridTableDescribeOutputRowAssert) HasName(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("name"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasType(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("type"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasCollation(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("collation"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasKind(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("kind"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasIsNullable(expected bool) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("is_nullable"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasDefault(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("default"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPrimaryKey(expected bool) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("primary_key"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasUniqueKey(expected bool) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("unique_key"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasCheck(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("check"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasExpression(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("expression"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasComment(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("comment"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPolicyName(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("policy_name"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasPrivacyDomain(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("privacy_domain"), expected))
	return h
}

func (h *HybridTableDescribeOutputRowAssert) HasSchemaEvolutionRecord(expected string) *HybridTableDescribeOutputRowAssert {
	h.AddAssertion(assert.ValueSet(h.field("schema_evolution_record"), expected))
	return h
}
