package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (p *MaskingPolicyShowOutputAssert) HasCreatedOnNotEmpty() *MaskingPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}

// MaskingPoliciesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func MaskingPoliciesDatasourceShowOutput(t *testing.T, name string) *MaskingPolicyShowOutputAssert {
	t.Helper()

	m := MaskingPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "masking_policies.0."),
	}
	m.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &m
}
