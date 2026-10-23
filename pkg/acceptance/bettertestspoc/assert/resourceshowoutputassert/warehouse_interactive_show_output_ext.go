package resourceshowoutputassert

func (w *WarehouseInteractiveShowOutputAssert) HasStateNotEmpty() *WarehouseInteractiveShowOutputAssert {
	w.ValuePresent("state")
	return w
}

func (w *WarehouseInteractiveShowOutputAssert) HasTables(expected ...string) *WarehouseInteractiveShowOutputAssert {
	w.SetContainsExactlyStringValues("tables", expected...)
	return w
}
