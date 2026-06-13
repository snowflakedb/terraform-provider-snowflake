package resourceshowoutputassert

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// HybridTableDescribeOutputAssert asserts fields of a single row in the
// describe_output list. Because hybrid table describe output has one row per
// column, callers must supply the 0-based row index (= column position in the
// physical column order returned by DESCRIBE TABLE).
type HybridTableDescribeOutputAssert struct {
	*assert.ResourceAssert
	rowIndex int
}

func HybridTableDescribeOutput(t *testing.T, name string, rowIndex int) *HybridTableDescribeOutputAssert {
	t.Helper()
	return &HybridTableDescribeOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "describe_output"),
		rowIndex:       rowIndex,
	}
}

func (h *HybridTableDescribeOutputAssert) field(name string) string {
	return fmt.Sprintf("describe_output.%d.%s", h.rowIndex, name)
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (h *HybridTableDescribeOutputAssert) HasName(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("name"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasType(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("type"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasCollation(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("collation"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasKind(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("kind"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasIsNullable(expected bool) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("is_nullable"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasDefault(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("default"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasPrimaryKey(expected bool) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("primary_key"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasUniqueKey(expected bool) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("unique_key"), strconv.FormatBool(expected)))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasCheck(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("check"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasExpression(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("expression"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasComment(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("comment"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasPolicyName(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("policy_name"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasPrivacyDomain(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("privacy_domain"), expected))
	return h
}

func (h *HybridTableDescribeOutputAssert) HasSchemaEvolutionRecord(expected string) *HybridTableDescribeOutputAssert {
	h.AddAssertion(assert.ValueSet(h.field("schema_evolution_record"), expected))
	return h
}
