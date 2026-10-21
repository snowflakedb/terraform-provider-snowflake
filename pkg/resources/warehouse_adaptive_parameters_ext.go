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

var warehouseAdaptiveParametersAttributesSchemaExt = func() map[string]*schema.Schema {
	m := make(map[string]*schema.Schema)
	for _, p := range defs.ParameterDefsForLevel(ssot.ParameterLevelWarehouseAdaptive) {
		m[p.FieldName()] = parameterSchema(p)
	}
	return m
}()

var warehouseAdaptiveParametersOutputSchemaExt = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseAdaptiveParametersSchema,
		},
	},
}

var warehouseAdaptiveParametersFieldNamesExt = []string{
	strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)),
	strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds)),
}

var warehouseAdaptiveParametersCustomDiffExt = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseAdaptiveParameterDiffFunctions,
)

// =============================================================================
// Create / Read / Update / CustomizeDiff handlers
// =============================================================================

func warehouseAdaptiveApplyParametersCreateExt(d *schema.ResourceData, opts *sdk.CreateAdaptiveWarehouseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &opts.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &opts.StatementTimeoutInSeconds),
	)
}

func warehouseAdaptiveApplyParametersChangesExt(d *schema.ResourceData, setUnset *warehouseAdaptiveAlterRequestsExt) diag.Diagnostics {
	set, unset := setUnset.set, setUnset.unset
	return JoinDiags(
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	)
}

func warehouseAdaptiveSetParametersFieldsExt(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
	)
}

func warehouseAdaptiveParameterDiffFunctions(parameters *sdk.WarehouseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse),
	}
}

// =============================================================================
// SDK incision points
// =============================================================================

func warehouseAdaptiveShowParametersInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	return client.Warehouses.ShowParameters(ctx, id)
}
