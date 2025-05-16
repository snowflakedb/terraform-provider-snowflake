package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

// ResourceMonitorDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ResourceMonitorDatasourceShowOutput(t *testing.T, name string) *ResourceMonitorShowOutputAssert {
	t.Helper()

	u := ResourceMonitorShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "resource_monitors.0."),
	}
	u.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &u
}

func (r *ResourceMonitorShowOutputAssert) HasStartTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("start_time"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasEndTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("end_time"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasCreatedOnNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasOwnerNotEmpty() *ResourceMonitorShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return r
}
