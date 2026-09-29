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
	ShowSchemaParametersSchema = make(map[string]*schema.Schema)
	schemaParameters           = defs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema)
)

func init() {
	for _, def := range schemaParameters {
		ShowSchemaParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func SchemaParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	schemaParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(schemaParameters, func(def parameterdefs.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			schemaParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return schemaParametersValue
}
