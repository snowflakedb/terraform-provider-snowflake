package resourceparametersassert

func (w *WarehouseInteractiveResourceParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseInteractiveResourceParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}

func (w *WarehouseInteractiveResourceParametersAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseInteractiveResourceParametersAssert {
	return w.
		HasStatementQueuedTimeoutInSeconds(0).
		HasStatementQueuedTimeoutInSecondsLevel("")
}

func (w *WarehouseInteractiveResourceParametersAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseInteractiveResourceParametersAssert {
	return w.
		HasStatementTimeoutInSeconds(172800).
		HasStatementTimeoutInSecondsLevel("")
}

func (w *WarehouseInteractiveResourceParametersAssert) HasDefaultFallbackWarehouse() *WarehouseInteractiveResourceParametersAssert {
	return w.
		HasFallbackWarehouse("").
		HasFallbackWarehouseLevel("")
}
