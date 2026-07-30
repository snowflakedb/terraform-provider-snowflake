package resourceshowoutputassert

func (w *WarehouseShowOutputAssert) HasResourceConstraintEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("resource_constraint", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasGenerationEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("generation", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasStateNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("state")
	return w
}

func (w *WarehouseShowOutputAssert) HasStartedClustersNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("started_clusters")
	return w
}

func (w *WarehouseShowOutputAssert) HasRunningNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("running")
	return w
}

func (w *WarehouseShowOutputAssert) HasQueuedNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("queued")
	return w
}

func (w *WarehouseShowOutputAssert) HasAvailableNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("available")
	return w
}

func (w *WarehouseShowOutputAssert) HasProvisioningNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("provisioning")
	return w
}

func (w *WarehouseShowOutputAssert) HasQuiescingNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("quiescing")
	return w
}

func (w *WarehouseShowOutputAssert) HasOtherNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("other")
	return w
}

func (w *WarehouseShowOutputAssert) HasCreatedOnNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("created_on")
	return w
}

func (w *WarehouseShowOutputAssert) HasResumedOnNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("resumed_on")
	return w
}

func (w *WarehouseShowOutputAssert) HasUpdatedOnNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("updated_on")
	return w
}

func (w *WarehouseShowOutputAssert) HasOwnerNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("owner")
	return w
}

func (w *WarehouseShowOutputAssert) HasOwnerRoleTypeNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("owner_role_type")
	return w
}

func (w *WarehouseShowOutputAssert) HasResourceMonitorEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("resource_monitor", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasGenerationNotEmpty() *WarehouseShowOutputAssert {
	w.ValuePresent("generation")
	return w
}

func (w *WarehouseShowOutputAssert) HasCommentEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("comment", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasSizeEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("size", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasScalingPolicyEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("scaling_policy", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasMaxQueryPerformanceLevelEmpty() *WarehouseShowOutputAssert {
	w.StringValueSet("max_query_performance_level", "")
	return w
}

func (w *WarehouseShowOutputAssert) HasTables(expected ...string) *WarehouseShowOutputAssert {
	w.SetContainsExactlyStringValues("tables", expected...)
	return w
}

func (w *WarehouseShowOutputAssert) HasNoTables() *WarehouseShowOutputAssert {
	w.ValueSet("tables.#", "0")
	return w
}
