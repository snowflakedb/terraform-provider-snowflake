package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func StorageLifecyclePoliciesDatasourceDescribeOutput(t *testing.T, name string) *StorageLifecyclePolicyDescribeOutputAssert {
	t.Helper()

	s := StorageLifecyclePolicyDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "describe_output", "storage_lifecycle_policies.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &s
}
