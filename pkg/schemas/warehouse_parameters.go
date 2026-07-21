package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// showWarehouseParametersSchemaCommon contains the warehouse parameters present for all warehouse types.
// TODO [SNOW-1473425]: descriptions (take from .Description; tool to validate changes later)
// TODO [SNOW-1473425]: should be generated later based on sdk.WarehouseParameters
var showWarehouseParametersSchemaCommon = map[string]*schema.Schema{
	"max_concurrency_level":               ParameterListSchema,
	"statement_queued_timeout_in_seconds": ParameterListSchema,
	"statement_timeout_in_seconds":        ParameterListSchema,
}

// ShowWarehouseParametersSchema contains all Snowflake parameters for the (standard/adaptive) warehouses.
var ShowWarehouseParametersSchema = collections.MergeMaps(showWarehouseParametersSchemaCommon)

// ShowWarehouseParametersSchemaInteractive contains common and interactive-only warehouse parameters
// (used by the interactive warehouse resource).
var ShowWarehouseParametersSchemaInteractive = collections.MergeMaps(showWarehouseParametersSchemaCommon, map[string]*schema.Schema{
	"fallback_warehouse": ParameterListSchema,
})

// commonWarehouseParametersToSchema maps the warehouse parameters present for all warehouse types
// (showWarehouseParametersSchemaCommon) into the given map.
func commonWarehouseParametersToSchema(warehouseParameters map[string]any, parameters []*sdk.Parameter, providerCtx *provider.Context) {
	for _, param := range parameters {
		parameterSchema := ParameterToSchemaReducedOutput(param, providerCtx)
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.WarehouseParameterMaxConcurrencyLevel),
			string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds),
			string(sdk.WarehouseParameterStatementTimeoutInSeconds):
			warehouseParameters[strings.ToLower(key)] = []map[string]any{parameterSchema}
		}
	}
}

// TODO [SNOW-1473425]: validate all present?
func WarehouseParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	warehouseParameters := make(map[string]any)
	commonWarehouseParametersToSchema(warehouseParameters, parameters, providerCtx)
	return warehouseParameters
}

// WarehouseInteractiveParametersToSchema maps the common warehouse parameters plus the interactive-only
// FALLBACK_WAREHOUSE parameter (an account-object identifier).
func WarehouseInteractiveParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	warehouseParameters := make(map[string]any)
	commonWarehouseParametersToSchema(warehouseParameters, parameters, providerCtx)
	for _, param := range parameters {
		if strings.EqualFold(param.Key, string(sdk.WarehouseParameterFallbackWarehouse)) {
			warehouseParameters["fallback_warehouse"] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return warehouseParameters
}
