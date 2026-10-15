package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowWarehouseAdaptiveParametersSchema = make(map[string]*schema.Schema)
	warehouseAdaptiveParameters           = defs.ParameterDefsForLevel(ssot.ParameterLevelWarehouseAdaptive)
)

func init() {
	for _, def := range warehouseAdaptiveParameters {
		ShowWarehouseAdaptiveParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func WarehouseAdaptiveParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	warehouseParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(warehouseAdaptiveParameters, func(def ssot.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			warehouseParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return warehouseParametersValue
}
