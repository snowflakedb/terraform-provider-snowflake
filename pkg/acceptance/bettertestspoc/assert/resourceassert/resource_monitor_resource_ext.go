package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (r *ResourceMonitorResourceAssert) HasStartTimestampNotEmpty() *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValuePresent("start_timestamp"))
	return r
}

func (r *ResourceMonitorResourceAssert) HasEndTimestampNotEmpty() *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValuePresent("end_timestamp"))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyUsersLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_users.#", strconv.FormatInt(int64(len), 10)))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyUser(index int, userName string) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet(fmt.Sprintf("notify_users.%d", index), userName))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyTriggersLen(len int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_triggers.#", strconv.FormatInt(int64(len), 10)))
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyTrigger(index int, threshold int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet(fmt.Sprintf("notify_triggers.%d", index), strconv.Itoa(threshold)))
	return r
}
