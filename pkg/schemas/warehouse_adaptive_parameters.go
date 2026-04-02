package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowWarehouseAdaptiveParametersSchema contains all Snowflake parameters for the adaptive warehouses.
// TODO [SNOW-1473425]: descriptions (take from .Description; tool to validate changes later)
// TODO [SNOW-1473425]: should be generated later based on sdk.WarehouseParameters
var ShowWarehouseAdaptiveParametersSchema = map[string]*schema.Schema{
	"statement_queued_timeout_in_seconds": ParameterListSchema,
	"statement_timeout_in_seconds":        ParameterListSchema,
}

// TODO [SNOW-1473425]: validate all present?
func WarehouseAdaptiveParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	warehouseParameters := make(map[string]any)
	for _, param := range parameters {
		parameterSchema := ParameterToSchemaReducedOutput(param, providerCtx)
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds),
			string(sdk.ObjectParameterStatementTimeoutInSeconds):
			warehouseParameters[strings.ToLower(key)] = []map[string]any{parameterSchema}
		}
	}
	return warehouseParameters
}
