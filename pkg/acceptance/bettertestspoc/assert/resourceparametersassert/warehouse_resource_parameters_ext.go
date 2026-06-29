package resourceparametersassert

func (w *WarehouseResourceParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}

func (w *WarehouseResourceParametersAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseResourceParametersAssert {
	return w.
		HasStatementQueuedTimeoutInSeconds(0).
		HasStatementQueuedTimeoutInSecondsLevel("")
}

func (w *WarehouseResourceParametersAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseResourceParametersAssert {
	return w.
		HasStatementTimeoutInSeconds(172800).
		HasStatementTimeoutInSecondsLevel("")
}
