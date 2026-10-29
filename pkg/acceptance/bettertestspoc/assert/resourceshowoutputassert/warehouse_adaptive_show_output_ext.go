package resourceshowoutputassert

func (w *WarehouseAdaptiveShowOutputAssert) HasStateNotEmpty() *WarehouseAdaptiveShowOutputAssert {
	w.ValuePresent("state")
	return w
}

func (w *WarehouseAdaptiveShowOutputAssert) HasCommentEmpty() *WarehouseAdaptiveShowOutputAssert {
	w.StringValueSet("comment", "")
	return w
}

func (w *WarehouseAdaptiveShowOutputAssert) HasResourceMonitorEmpty() *WarehouseAdaptiveShowOutputAssert {
	w.StringValueSet("resource_monitor", "")
	return w
}

func (w *WarehouseAdaptiveShowOutputAssert) HasOwnerNotEmpty() *WarehouseAdaptiveShowOutputAssert {
	w.ValuePresent("owner")
	return w
}

func (w *WarehouseAdaptiveShowOutputAssert) HasOwnerRoleTypeNotEmpty() *WarehouseAdaptiveShowOutputAssert {
	w.ValuePresent("owner_role_type")
	return w
}
