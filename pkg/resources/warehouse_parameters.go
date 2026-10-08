package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehouseParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouse) {
		warehouseParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

func handleWarehouseParametersCreate(d *schema.ResourceData, request *sdk.CreateWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.MaxConcurrencyLevel.FieldName(), &request.MaxConcurrencyLevel),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &request.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &request.StatementTimeoutInSeconds),
	)
}

func handleWarehouseParametersChanges(d *schema.ResourceData, set *sdk.WarehouseSetRequest, unset *sdk.WarehouseUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.MaxConcurrencyLevel.FieldName(), &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	)
}

func handleWarehouseParameterRead(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
	)
}

var warehouseParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseParameterDiffFunctions,
)

func warehouseParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.WarehouseParametersDetails, error) {
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Warehouses.ShowParametersDetails(ctx, id)
}

func warehouseParameterDiffFunctions(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
	}
}
