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
	ShowTableParametersSchema = make(map[string]*schema.Schema)
	tableParameters           = defs.ParameterDefsForLevel(ssot.ParameterLevelTable)
)

func init() {
	for _, def := range tableParameters {
		ShowTableParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func TableParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	tableParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(tableParameters, func(def ssot.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			tableParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return tableParametersValue
}
