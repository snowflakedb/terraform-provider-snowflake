package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func warehouseApplyParametersCreateExt(d *schema.ResourceData, request *sdk.CreateWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.MaxConcurrencyLevel.FieldName(), &request.MaxConcurrencyLevel),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &request.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &request.StatementTimeoutInSeconds),
	)
}

func warehouseApplyParametersChangesExt(d *schema.ResourceData, reqs *warehouseAlterRequestsExt) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.MaxConcurrencyLevel.FieldName(), &reqs.set.MaxConcurrencyLevel, &reqs.unset.MaxConcurrencyLevel),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &reqs.set.StatementQueuedTimeoutInSeconds, &reqs.unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &reqs.set.StatementTimeoutInSeconds, &reqs.unset.StatementTimeoutInSeconds),
	)
}

func warehouseSetParametersFieldsExt(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
	)
}

func warehouseParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.WarehouseParametersDetails, error) {
	id, err := warehouseParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Warehouses.ShowParametersDetails(ctx, id)
}

func warehouseParameterDiffFunctionsExt(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
	}
}

func warehouseShowParametersInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	return client.Warehouses.ShowParameters(ctx, id)
}

func warehouseParametersDetailsFromRawExt(params []*sdk.Parameter) (*sdk.WarehouseParametersDetails, error) {
	return sdk.ToWarehouseParametersDetails(params)
}
