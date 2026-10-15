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
	ShowAccountParametersSchema = make(map[string]*schema.Schema)
	accountParameters           = defs.ParameterDefsForLevel(ssot.ParameterLevelAccountExt)
)

func init() {
	for _, def := range accountParameters {
		ShowAccountParametersSchema[def.FieldName()] = ParameterListSchema
	}
}

func AccountParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	accountParametersValue := make(map[string]any)
	for _, parameter := range parameters {
		if slices.ContainsFunc(accountParameters, func(def ssot.ParameterDef) bool { return def.SqlName == parameter.Key }) {
			accountParametersValue[strings.ToLower(parameter.Key)] = []map[string]any{ParameterToSchemaReducedOutput(parameter, providerCtx)}
		}
	}
	return accountParametersValue
}
