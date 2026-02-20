package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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
