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
	ShowHybridTableParametersSchema = make(map[string]*schema.Schema)
	hybridTableParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelHybridTable)
)

func init() {
	for _, def := range hybridTableParameters {
		ShowHybridTableParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func HybridTableParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(hybridTableParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			result[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return result
}
