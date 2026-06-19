package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func StorageLifecyclePoliciesDatasourceShowOutput(t *testing.T, name string) *StorageLifecyclePolicyShowOutputAssert {
	t.Helper()

	s := StorageLifecyclePolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "storage_lifecycle_policies.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
