package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowProcedureParametersSchema = make(map[string]*schema.Schema)
	ProcedureParameters           = []sdk.ProcedureParameter{
		sdk.ProcedureParameterEnableConsoleOutput,
		sdk.ProcedureParameterLogLevel,
		sdk.ProcedureParameterMetricLevel,
		sdk.ProcedureParameterTraceLevel,
	}
)

func init() {
	for _, param := range ProcedureParameters {
		ShowProcedureParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func ProcedureParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	ProcedureParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(ProcedureParameters, sdk.ProcedureParameter(param.Key)) {
			ProcedureParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return ProcedureParametersValue
}
