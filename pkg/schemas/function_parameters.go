package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowFunctionParametersSchema = make(map[string]*schema.Schema)
	functionParameters           = []sdk.FunctionParameter{
		sdk.FunctionParameterEnableConsoleOutput,
		sdk.FunctionParameterLogLevel,
		sdk.FunctionParameterMetricLevel,
		sdk.FunctionParameterTraceLevel,
	}
)

func init() {
	for _, param := range functionParameters {
		ShowFunctionParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func FunctionParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	functionParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(functionParameters, sdk.FunctionParameter(param.Key)) {
			functionParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return functionParametersValue
}
