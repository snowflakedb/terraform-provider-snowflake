package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (w *WarehouseShowOutputAssert) HasResourceConstraintEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValueSet("resource_constraint", ""))
	return w
}

func (w *WarehouseShowOutputAssert) HasGenerationEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValueSet("generation", ""))
	return w
}

func (w *WarehouseShowOutputAssert) HasStateNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("state"))
	return w
}

func (w *WarehouseShowOutputAssert) HasStartedClustersNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("started_clusters"))
	return w
}

func (w *WarehouseShowOutputAssert) HasRunningNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("running"))
	return w
}

func (w *WarehouseShowOutputAssert) HasQueuedNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("queued"))
	return w
}

func (w *WarehouseShowOutputAssert) HasAvailableNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("available"))
	return w
}

func (w *WarehouseShowOutputAssert) HasProvisioningNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("provisioning"))
	return w
}

func (w *WarehouseShowOutputAssert) HasQuiescingNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("quiescing"))
	return w
}

func (w *WarehouseShowOutputAssert) HasOtherNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("other"))
	return w
}

func (w *WarehouseShowOutputAssert) HasCreatedOnNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return w
}

func (w *WarehouseShowOutputAssert) HasResumedOnNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("resumed_on"))
	return w
}

func (w *WarehouseShowOutputAssert) HasUpdatedOnNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("updated_on"))
	return w
}

func (w *WarehouseShowOutputAssert) HasOwnerNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return w
}

func (w *WarehouseShowOutputAssert) HasOwnerRoleTypeNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("owner_role_type"))
	return w
}

func (w *WarehouseShowOutputAssert) HasResourceMonitorEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValueSet("resource_monitor", ""))
	return w
}

func (w *WarehouseShowOutputAssert) HasGenerationNotEmpty() *WarehouseShowOutputAssert {
	w.AddAssertion(assert.ResourceShowOutputValuePresent("generation"))
	return w
}
