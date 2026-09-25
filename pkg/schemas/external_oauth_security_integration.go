package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowExternalOauthParametersSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.AccountParameterExternalOauthAddPrivilegedRolesToBlockedList)): ParameterListSchema,
}

func ExternalOauthParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	schemaMap := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains([]sdk.AccountParameter{sdk.AccountParameterExternalOauthAddPrivilegedRolesToBlockedList}, sdk.AccountParameter(param.Key)) {
			schemaMap[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return schemaMap
}
