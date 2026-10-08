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
	ShowServiceParametersSchema = make(map[string]*schema.Schema)
	serviceParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelService)
)

func init() {
	for _, def := range serviceParameters {
		ShowServiceParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func ServiceParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	serviceParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(serviceParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			serviceParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return serviceParametersValue
}
