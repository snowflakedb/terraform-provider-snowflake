package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowIcebergTableParametersSchema contains the Snowflake parameters surfaced for iceberg tables.
var ShowIcebergTableParametersSchema = map[string]*schema.Schema{
	"external_volume":            ParameterListSchema,
	"catalog":                    ParameterListSchema,
	"replace_invalid_characters": ParameterListSchema,
}

func IcebergTableParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	result := make(map[string]any)
	for _, param := range parameters {
		parameterSchema := ParameterToSchemaReducedOutput(param, providerCtx)
		switch key := strings.ToUpper(param.Key); key {
		case string(sdk.IcebergTableParameterExternalVolume),
			string(sdk.IcebergTableParameterCatalog),
			string(sdk.IcebergTableParameterReplaceInvalidCharacters):
			result[strings.ToLower(key)] = []map[string]any{parameterSchema}
		}
	}
	return result
}
