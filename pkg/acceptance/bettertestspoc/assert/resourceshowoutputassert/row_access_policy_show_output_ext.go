package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (r *RowAccessPolicyShowOutputAssert) HasCreatedOnNotEmpty() *RowAccessPolicyShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return r
}

// RowAccessPoliciesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func RowAccessPoliciesDatasourceShowOutput(t *testing.T, name string) *RowAccessPolicyShowOutputAssert {
	t.Helper()

	r := RowAccessPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "row_access_policies.0."),
	}
	r.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &r
}
