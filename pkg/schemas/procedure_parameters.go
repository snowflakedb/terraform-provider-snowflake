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
	ShowProcedureParametersSchema = make(map[string]*schema.Schema)
	procedureParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelProcedure)
)

func init() {
	for _, def := range procedureParameters {
		ShowProcedureParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func ProcedureParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	procedureParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(procedureParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			procedureParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return procedureParametersValue
}
