package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Schema maps
// =============================================================================

var warehouseInteractiveParametersAttributesSchemaExt = func() map[string]*schema.Schema {
	m := make(map[string]*schema.Schema)
	for _, p := range defs.ParameterDefsForLevel(ssot.ParameterLevelWarehouseInteractive) {
		m[p.FieldName()] = parameterSchema(p)
	}
	return m
}()

var warehouseInteractiveParametersOutputSchemaExt = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchemaInteractive,
		},
	},
}

var warehouseInteractiveParametersFieldNamesExt = []string{
	strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel)),
	strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)),
	strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds)),
	"fallback_warehouse",
}

var warehouseInteractiveParametersCustomDiffExt = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseInteractiveParameterDiffFunctions,
)

// =============================================================================
// Create / Read / Update / CustomizeDiff handlers
// =============================================================================

func warehouseInteractiveApplyParametersCreateExt(d *schema.ResourceData, req *sdk.CreateInteractiveWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.MaxConcurrencyLevel.FieldName(), &req.MaxConcurrencyLevel),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &req.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &req.StatementTimeoutInSeconds),
		handleParameterCreateWithMapping(d, defs.FallbackWarehouse.FieldName(), &req.FallbackWarehouse, stringToAccountObjectIdentifier),
	)
}

func warehouseInteractiveApplyParametersChangesExt(d *schema.ResourceData, setUnset *warehouseInteractiveAlterRequestsExt) diag.Diagnostics {
	set, unset := setUnset.set, setUnset.unset
	return JoinDiags(
		handleParameterUpdate(d, defs.MaxConcurrencyLevel.FieldName(), &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
		handleParameterUpdateWithMapping(d, defs.FallbackWarehouse.FieldName(), &set.FallbackWarehouse, &unset.FallbackWarehouse, stringToAccountObjectIdentifier),
	)
}

func warehouseInteractiveSetParametersFieldsExt(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
		setResourceData(d, defs.FallbackWarehouse.FieldName(), parameters.FallbackWarehouse.Value.Name()),
	)
}

func warehouseInteractiveParameterDiffFunctions(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.MaxConcurrencyLevel.FieldName(), parameters.MaxConcurrencyLevel, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IdentifierTypedParameterValueComputedIf(defs.FallbackWarehouse.FieldName(), parameters.FallbackWarehouse, sdk.ParameterTypeWarehouse),
	}
}

// =============================================================================
// SDK incision points
// =============================================================================

func warehouseInteractiveShowParametersInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	return client.Warehouses.ShowParameters(ctx, id)
}
