package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowServiceParametersSchema = make(map[string]*schema.Schema)
	serviceParameters           = []sdk.ServiceParameter{
		sdk.ServiceParameterServiceCallerTokenValiditySecs,
	}
)

func init() {
	for _, param := range serviceParameters {
		ShowServiceParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func ServiceParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	serviceParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(serviceParameters, sdk.ServiceParameter(param.Key)) {
			serviceParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return serviceParametersValue
}
