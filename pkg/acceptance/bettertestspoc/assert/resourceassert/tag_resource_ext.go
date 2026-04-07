package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasMaskingPoliciesLength checks that the masking_policies field has the expected length
func (t *TagResourceAssert) HasMaskingPoliciesLength(expected int) *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("masking_policies.#", strconv.FormatInt(int64(expected), 10)))
	return t
}

// HasAllowedValuesLength checks that the allowed_values field has the expected length
func (t *TagResourceAssert) HasAllowedValuesLength(expected int) *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("allowed_values.#", strconv.FormatInt(int64(expected), 10)))
	return t
}

func (t *TagResourceAssert) HasOnConflictCustomValue(expected string) *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("on_conflict.0.custom_value", expected))
	return t
}

func (t *TagResourceAssert) HasOnConflictAllowedValuesSequence() *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("on_conflict.0.allowed_values_sequence", "true"))
	return t
}

// HasAllowedValuesOrder checks that the allowed_values_order field has the expected values in exact order.
func (t *TagResourceAssert) HasAllowedValuesOrder(expected ...string) *TagResourceAssert {
	t.ListContainsExactlyStringValuesInOrder("allowed_values_order", expected...)
	return t
}

// HasPropagateEnum
func (t *TagResourceAssert) HasPropagateEnum(expected sdk.TagPropagation) *TagResourceAssert {
	return t.HasPropagateString(string(expected))
}
