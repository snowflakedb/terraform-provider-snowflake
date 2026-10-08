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
	ShowOpenflowDeploymentParametersSchema = make(map[string]*schema.Schema)
	openflowDeploymentParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelOpenflowDeployment)
)

func init() {
	for _, def := range openflowDeploymentParameters {
		ShowOpenflowDeploymentParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func OpenflowDeploymentParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(openflowDeploymentParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			result[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return result
}
