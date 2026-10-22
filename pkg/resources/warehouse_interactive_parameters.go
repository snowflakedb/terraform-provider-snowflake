package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehouseInteractiveParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseInteractive) {
		warehouseInteractiveParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

func handleWarehouseInteractiveParametersCreate(d *schema.ResourceData, req *sdk.CreateInteractiveWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.MaxConcurrencyLevel.FieldName(), &req.MaxConcurrencyLevel),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &req.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &req.StatementTimeoutInSeconds),
		handleParameterCreateWithMapping(d, defs.FallbackWarehouse.FieldName(), &req.FallbackWarehouse, stringToAccountObjectIdentifier),
	)
}

func handleWarehouseInteractiveParametersChanges(d *schema.ResourceData, set *sdk.WarehouseSetRequest, unset *sdk.WarehouseUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.MaxConcurrencyLevel.FieldName(), &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
		handleParameterUpdateWithMapping(d, defs.FallbackWarehouse.FieldName(), &set.FallbackWarehouse, &unset.FallbackWarehouse, stringToAccountObjectIdentifier),
	)
}

func handleWarehouseInteractiveParameterRead(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
		setResourceData(d, defs.FallbackWarehouse.FieldName(), parameters.FallbackWarehouse.Value.Name()),
	)
}

var warehouseInteractiveParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseInteractiveParameterDiffFunctions,
)

func warehouseInteractiveParameterDiffFunctions(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IdentifierTypedParameterValueComputedIf(defs.FallbackWarehouse.FieldName(), parameters.FallbackWarehouse, sdk.ParameterTypeWarehouse),
	}
}
