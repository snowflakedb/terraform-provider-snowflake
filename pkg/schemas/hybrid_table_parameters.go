package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowHybridTableParametersSchema = map[string]*schema.Schema{
	"data_retention_time_in_days":     ParameterListSchema,
	"max_data_extension_time_in_days": ParameterListSchema,
}

func HybridTableParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, param := range parameters {
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.ObjectParameterDataRetentionTimeInDays),
			string(sdk.ObjectParameterMaxDataExtensionTimeInDays):
			result[strings.ToLower(key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return result
}
