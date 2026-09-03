package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowOpenflowDeploymentParametersSchema = map[string]*schema.Schema{
	"event_table": ParameterListSchema,
}

func OpenflowDeploymentParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, param := range parameters {
		if key := strings.ToUpper(param.Key); key == string(sdk.OpenflowDeploymentParameterEventTable) {
			result[strings.ToLower(key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return result
}
