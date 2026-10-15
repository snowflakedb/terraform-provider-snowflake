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
	ShowFunctionParametersSchema = make(map[string]*schema.Schema)
	functionParameters           = defs.ParameterDefsForLevel(ssot.ParameterLevelFunction)
)

func init() {
	for _, def := range functionParameters {
		ShowFunctionParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func FunctionParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	functionParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(functionParameters, func(def ssot.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			functionParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return functionParametersValue
}
