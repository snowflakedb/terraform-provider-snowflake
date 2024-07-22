// Code generated by sdk-to-schema generator; DO NOT EDIT.

package objectparametersassert

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehouseParametersAssert struct {
	*assert.SnowflakeParametersAssert[sdk.AccountObjectIdentifier]
}

func WarehouseParameters(t *testing.T, id sdk.AccountObjectIdentifier) *WarehouseParametersAssert {
	t.Helper()
	return &WarehouseParametersAssert{
		assert.NewSnowflakeParametersAssertWithProvider(id, sdk.ObjectTypeWarehouse, acc.TestClient().Parameter.ShowWarehouseParameters),
	}
}

func WarehouseParametersPrefetched(t *testing.T, id sdk.AccountObjectIdentifier, parameters []*sdk.Parameter) *WarehouseParametersAssert {
	t.Helper()
	return &WarehouseParametersAssert{
		assert.NewSnowflakeParametersAssertWithParameters(id, sdk.ObjectTypeWarehouse, parameters),
	}
}

//////////////////////////////
// Generic parameter checks //
//////////////////////////////

func (w *WarehouseParametersAssert) HasBoolParameterValue(parameterName sdk.WarehouseParameter, expected bool) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterBoolValueSet(parameterName, expected))
	return w
}

func (w *WarehouseParametersAssert) HasIntParameterValue(parameterName sdk.WarehouseParameter, expected int) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterIntValueSet(parameterName, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStringParameterValue(parameterName sdk.WarehouseParameter, expected string) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterValueSet(parameterName, expected))
	return w
}

func (w *WarehouseParametersAssert) HasDefaultParameterValue(parameterName sdk.WarehouseParameter) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterDefaultValueSet(parameterName))
	return w
}

func (w *WarehouseParametersAssert) HasDefaultParameterValueOnLevel(parameterName sdk.WarehouseParameter, parameterType sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterDefaultValueOnLevelSet(parameterName, parameterType))
	return w
}

///////////////////////////////
// Aggregated generic checks //
///////////////////////////////

// HasAllDefaults checks if all the parameters:
// - have a default value by comparing current value of the sdk.Parameter with its default
// - have an expected level
func (w *WarehouseParametersAssert) HasAllDefaults() *WarehouseParametersAssert {
	return w.
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterMaxConcurrencyLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault)
}

func (w *WarehouseParametersAssert) HasAllDefaultsExplicit() *WarehouseParametersAssert {
	return w.
		HasDefaultMaxConcurrencyLevelValueExplicit().
		HasDefaultStatementQueuedTimeoutInSecondsValueExplicit().
		HasDefaultStatementTimeoutInSecondsValueExplicit()
}

////////////////////////////
// Parameter value checks //
////////////////////////////

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevel(expected int) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterIntValueSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterIntValueSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterIntValueSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}

////////////////////////////
// Parameter level checks //
////////////////////////////

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevelLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}

////////////////////////////////////
// Parameter default value checks //
////////////////////////////////////

func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevelValue() *WarehouseParametersAssert {
	return w.HasDefaultParameterValue(sdk.WarehouseParameterMaxConcurrencyLevel)
}

func (w *WarehouseParametersAssert) HasDefaultStatementQueuedTimeoutInSecondsValue() *WarehouseParametersAssert {
	return w.HasDefaultParameterValue(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)
}

func (w *WarehouseParametersAssert) HasDefaultStatementTimeoutInSecondsValue() *WarehouseParametersAssert {
	return w.HasDefaultParameterValue(sdk.WarehouseParameterStatementTimeoutInSeconds)
}

/////////////////////////////////////////////
// Parameter explicit default value checks //
/////////////////////////////////////////////

func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevelValueExplicit() *WarehouseParametersAssert {
	return w.HasMaxConcurrencyLevel(8)
}

func (w *WarehouseParametersAssert) HasDefaultStatementQueuedTimeoutInSecondsValueExplicit() *WarehouseParametersAssert {
	return w.HasStatementQueuedTimeoutInSeconds(0)
}

func (w *WarehouseParametersAssert) HasDefaultStatementTimeoutInSecondsValueExplicit() *WarehouseParametersAssert {
	return w.HasStatementTimeoutInSeconds(172800)
}
