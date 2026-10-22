package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehouseAdaptiveParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseAdaptive) {
		warehouseAdaptiveParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

func handleWarehouseAdaptiveParametersCreate(d *schema.ResourceData, opts *sdk.CreateAdaptiveWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &opts.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &opts.StatementTimeoutInSeconds),
	)
}

func handleWarehouseAdaptiveParametersChanges(d *schema.ResourceData, set *sdk.WarehouseSetRequest, unset *sdk.WarehouseUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	)
}

func handleWarehouseAdaptiveParameterRead(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
	)
}

var warehouseAdaptiveParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseAdaptiveParameterDiffFunctions,
)

func warehouseAdaptiveParameterDiffFunctions(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
	}
}
