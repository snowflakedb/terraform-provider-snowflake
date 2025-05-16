package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

func FunctionParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	functionParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(functionParameters, sdk.FunctionParameter(param.Key)) {
			functionParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return functionParametersValue
}
