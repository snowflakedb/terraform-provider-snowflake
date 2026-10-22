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
	ShowUserParametersSchema = make(map[string]*schema.Schema)
	userParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelUser)
)

func init() {
	for _, def := range userParameters {
		ShowUserParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func UserParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	userParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(userParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			userParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return userParametersValue
}
