package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowWarehouseParametersSchema            = make(map[string]*schema.Schema)
	ShowWarehouseParametersSchemaInteractive = make(map[string]*schema.Schema)

	warehouseParameters            = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouse)
	warehouseInteractiveParameters = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseInteractive)
)

func init() {
	for _, def := range warehouseParameters {
		ShowWarehouseParametersSchema[def.FieldName()] = ParameterListSchema
	}
	for _, def := range warehouseInteractiveParameters {
		ShowWarehouseParametersSchemaInteractive[def.FieldName()] = ParameterListSchema
	}
}

// ShowAllWarehouseParametersSchema returns a schema containing all warehouse parameters for every warehouse type.
// Used in the warehouses data source to cover all warehouse types in a single schema.
func ShowAllWarehouseParametersSchema() map[string]*schema.Schema {
	return ShowWarehouseParametersSchemaInteractive
}

func WarehouseParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return warehouseParametersToSchema(parameters, providerCtx, warehouseParameters)
}

// WarehouseInteractiveParametersToSchema maps every parameter in ParameterLevelWarehouseInteractive
// (the common warehouse parameters plus the interactive-only FALLBACK_WAREHOUSE).
func WarehouseInteractiveParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return warehouseParametersToSchema(parameters, providerCtx, warehouseInteractiveParameters)
}

// AllWarehouseParametersToSchema maps every warehouse parameter across all warehouse types into the
// ShowAllWarehouseParametersSchema fields.
func AllWarehouseParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	return WarehouseInteractiveParametersToSchema(parameters, providerCtx)
}

func warehouseParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context, defsForLevel []parameterdefs.ParameterDef) map[string]any {
	warehouseParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(defsForLevel, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			warehouseParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return warehouseParametersValue
}
