package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func ProcedureParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	ProcedureParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(ProcedureParameters, sdk.ProcedureParameter(param.Key)) {
			ProcedureParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return ProcedureParametersValue
}
