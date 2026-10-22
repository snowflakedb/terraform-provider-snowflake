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
	ShowTaskParametersSchema = make(map[string]*schema.Schema)
	// SEARCH_PATH has no ParameterLevelTask membership in the catalog (it is not settable on a task at
	// all), but SHOW PARAMETERS IN TASK still reports it, and the resource keeps a search_path field for
	// that, so it is kept here too alongside the catalog-derived list.
	taskParameters = append(slices.Clone(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelTask)), defs.SearchPath)
)

func init() {
	for _, def := range taskParameters {
		ShowTaskParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func TaskParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	taskParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(taskParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			taskParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return taskParametersValue
}
