package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (t *TagResourceAssert) HasMaskingPoliciesLength(len int) *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("masking_policies.#", fmt.Sprintf("%d", len)))
	return t
}

func (t *TagResourceAssert) HasAllowedValuesLength(len int) *TagResourceAssert {
	t.AddAssertion(assert.ValueSet("allowed_values.#", fmt.Sprintf("%d", len)))
	return t
}
