// Code generated by assertions generator; DO NOT EDIT.

package resourceparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehouseResourceParametersAssert struct {
	*assert.ResourceAssert
}

func WarehouseResourceParameters(t *testing.T, name string) *WarehouseResourceParametersAssert {
	t.Helper()

	resourceParameterAssert := WarehouseResourceParametersAssert{
		ResourceAssert: assert.NewResourceAssert(name, "parameters"),
	}
	resourceParameterAssert.AddAssertion(assert.ValueSet("parameters.#", "1"))
	return &resourceParameterAssert
}

func ImportedWarehouseResourceParameters(t *testing.T, id string) *WarehouseResourceParametersAssert {
	t.Helper()

	resourceParameterAssert := WarehouseResourceParametersAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported parameters"),
	}
	resourceParameterAssert.AddAssertion(assert.ValueSet("parameters.#", "1"))
	return &resourceParameterAssert
}

////////////////////////////
// Parameter value checks //
////////////////////////////

func (w *WarehouseResourceParametersAssert) HasMaxConcurrencyLevel(expected int) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterIntValueSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseResourceParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterIntValueSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseResourceParametersAssert) HasStatementTimeoutInSeconds(expected int) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterIntValueSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}

////////////////////////////
// Parameter level checks //
////////////////////////////

func (w *WarehouseResourceParametersAssert) HasMaxConcurrencyLevelLevel(expected sdk.ParameterType) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterLevelSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseResourceParametersAssert) HasStatementQueuedTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterLevelSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseResourceParametersAssert) HasStatementTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseResourceParametersAssert {
	w.AddAssertion(assert.ResourceParameterLevelSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}
