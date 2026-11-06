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
	ShowPipeParametersSchema = make(map[string]*schema.Schema)
	pipeParameters           = defs.ParameterDefsForLevel(ssot.ParameterLevelPipe)
)

func init() {
	for _, def := range pipeParameters {
		ShowPipeParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func PipeParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	pipeParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(pipeParameters, func(def ssot.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			pipeParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return pipeParametersValue
}
