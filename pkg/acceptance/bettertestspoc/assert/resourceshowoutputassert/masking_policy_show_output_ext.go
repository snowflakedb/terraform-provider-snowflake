package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// MaskingPoliciesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func MaskingPoliciesDatasourceShowOutput(t *testing.T, datasourceReference string) *MaskingPolicyShowOutputAssert {
	t.Helper()

	m := MaskingPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "masking_policies.0."),
	}
	m.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &m
}

func (p *MaskingPolicyShowOutputAssert) HasCreatedOnNotEmpty() *MaskingPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}

func (p *MaskingPolicyShowOutputAssert) HasOwnerNotEmpty() *MaskingPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return p
}

func (p *MaskingPolicyShowOutputAssert) HasOwnerRoleTypeNotEmpty() *MaskingPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("owner_role_type"))
	return p
}
