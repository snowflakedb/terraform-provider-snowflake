package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (r *ResourceMonitorResourceAssert) HasNotifyUsers(expected ...string) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_users.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		r.AddAssertion(assert.SetElem("notify_users.*", v))
	}
	return r
}

func (r *ResourceMonitorResourceAssert) HasNotifyTriggers(expected ...int) *ResourceMonitorResourceAssert {
	r.AddAssertion(assert.ValueSet("notify_triggers.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		r.AddAssertion(assert.SetElem("notify_triggers.*", fmt.Sprintf("%d", v)))
	}
	return r
}
