package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowSchemaParametersSchema = make(map[string]*schema.Schema)
	schemaParameters           = []sdk.AccountParameter{
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
		sdk.AccountParameterPipeExecutionPaused,
	}
)

func init() {
	for _, param := range schemaParameters {
		ShowSchemaParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func SchemaParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	schemaParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(schemaParameters, sdk.AccountParameter(param.Key)) {
			schemaParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return schemaParametersValue
}
