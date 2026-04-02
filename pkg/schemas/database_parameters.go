package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowDatabaseParametersSchema = make(map[string]*schema.Schema)
	databaseParameters           = []sdk.AccountParameter{
		sdk.AccountParameterDataRetentionTimeInDays,
		sdk.AccountParameterMaxDataExtensionTimeInDays,
		sdk.AccountParameterExternalVolume,
		sdk.AccountParameterCatalog,
		sdk.AccountParameterReplaceInvalidCharacters,
		sdk.AccountParameterDefaultDDLCollation,
		sdk.AccountParameterStorageSerializationPolicy,
		sdk.AccountParameterLogLevel,
		sdk.AccountParameterTraceLevel,
		sdk.AccountParameterSuspendTaskAfterNumFailures,
		sdk.AccountParameterTaskAutoRetryAttempts,
		sdk.AccountParameterUserTaskManagedInitialWarehouseSize,
		sdk.AccountParameterUserTaskTimeoutMs,
		sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds,
		sdk.AccountParameterQuotedIdentifiersIgnoreCase,
		sdk.AccountParameterEnableConsoleOutput,
	}
)

func init() {
	for _, param := range databaseParameters {
		ShowDatabaseParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func DatabaseParametersToSchema(parameters []*sdk.Parameter, providerCtx *provider.Context) map[string]any {
	databaseParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(databaseParameters, sdk.AccountParameter(param.Key)) {
			databaseParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchemaReducedOutput(param, providerCtx)}
		}
	}
	return databaseParametersValue
}
