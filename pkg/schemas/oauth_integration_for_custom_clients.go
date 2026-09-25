package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowOauthForCustomClientsParametersSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList)): ParameterListSchema,
}

func OauthForCustomClientsParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	schemaMap := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains([]sdk.AccountParameter{sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList}, sdk.AccountParameter(param.Key)) {
			schemaMap[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return schemaMap
}
