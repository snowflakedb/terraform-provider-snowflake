package resources

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	databaseParametersSchema              = make(map[string]*schema.Schema)
	sharedDatabaseParametersSchema        = make(map[string]*schema.Schema)
	sharedDatabaseNotApplicableParameters = []sdk.ObjectParameter{
		sdk.ObjectParameterDataRetentionTimeInDays,
		sdk.ObjectParameterMaxDataExtensionTimeInDays,
	}
	databaseParametersCustomDiff = ParametersCustomDiff(
		databaseParametersProvider,
		parameter[sdk.AccountParameter]{sdk.AccountParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterExternalVolume, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterCatalog, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterDefaultDDLCollation, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterLogLevel, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterTraceLevel, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeDatabase},
	)
)

func databaseParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), databaseParametersProviderFunc, sdk.ParseAccountObjectIdentifier)
}

func databaseParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.AccountObjectIdentifier] {
	return c.Databases.ShowParameters
}

func init() {
	// TODO [SNOW-1645342]: use parameterDef
	databaseParameterFields := []struct {
		Name         sdk.ObjectParameter
		Type         schema.ValueType
		Description  string
		DiffSuppress schema.SchemaDiffSuppressFunc
		ValidateDiag schema.SchemaValidateDiagFunc
	}{
		{
			Name:        sdk.ObjectParameterDataRetentionTimeInDays,
			Type:        schema.TypeInt,
			Description: "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
			// Choosing higher range (for the standard edition or transient databases, the maximum number is 1)
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 90)),
		},
		{
			Name:        sdk.ObjectParameterDefaultDDLCollation,
			Type:        schema.TypeString,
			Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
		},
		{
			Name:         sdk.ObjectParameterCatalog,
			Type:         schema.TypeString,
			Description:  "The database parameter that specifies the default catalog to use for Iceberg tables. For more information, see [CATALOG](https://docs.snowflake.com/en/sql-reference/parameters#catalog).",
			ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppress: suppressIdentifierQuoting,
		},
		{
			Name:         sdk.ObjectParameterExternalVolume,
			Type:         schema.TypeString,
			Description:  "The database parameter that specifies the default external volume to use for Iceberg tables. For more information, see [EXTERNAL_VOLUME](https://docs.snowflake.com/en/sql-reference/parameters#external-volume).",
			ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppress: suppressIdentifierQuoting,
		},
		{
			Name:         sdk.ObjectParameterLogLevel,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
			ValidateDiag: sdkValidation(sdk.ToLogLevel),
			DiffSuppress: NormalizeAndCompare(sdk.ToLogLevel),
		},
		{
			Name:         sdk.ObjectParameterTraceLevel,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", docs.PossibleValuesListed(sdk.AsStringList(sdk.AllTraceLevels))),
			ValidateDiag: sdkValidation(sdk.ToTraceLevel),
			DiffSuppress: NormalizeAndCompare(sdk.ToTraceLevel),
		},
		{
			Name:         sdk.ObjectParameterMaxDataExtensionTimeInDays,
			Type:         schema.TypeInt,
			Description:  "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale. For a detailed description of this parameter, see [MAX_DATA_EXTENSION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters.html#label-max-data-extension-time-in-days).",
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 90)),
		},
		{
			Name:        sdk.ObjectParameterReplaceInvalidCharacters,
			Type:        schema.TypeBool,
			Description: "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog. For more information, see [REPLACE_INVALID_CHARACTERS](https://docs.snowflake.com/en/sql-reference/parameters#replace-invalid-characters).",
		},
		{
			Name:         sdk.ObjectParameterStorageSerializationPolicy,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("The storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake. For more information, see [STORAGE_SERIALIZATION_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#storage-serialization-policy).", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
			ValidateDiag: sdkValidation(sdk.ToStorageSerializationPolicy),
			DiffSuppress: NormalizeAndCompare(sdk.ToStorageSerializationPolicy),
		},
		{
			Name:         sdk.ObjectParameterSuspendTaskAfterNumFailures,
			Type:         schema.TypeInt,
			Description:  "How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending. For more information, see [SUSPEND_TASK_AFTER_NUM_FAILURES](https://docs.snowflake.com/en/sql-reference/parameters#suspend-task-after-num-failures).",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterTaskAutoRetryAttempts,
			Type:         schema.TypeInt,
			Description:  "Maximum automatic retries allowed for a user task. For more information, see [TASK_AUTO_RETRY_ATTEMPTS](https://docs.snowflake.com/en/sql-reference/parameters#task-auto-retry-attempts).",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterUserTaskManagedInitialWarehouseSize,
			Type:         schema.TypeString,
			Description:  "The initial size of warehouse to use for managed warehouses in the absence of history. For more information, see [USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE](https://docs.snowflake.com/en/sql-reference/parameters#user-task-managed-initial-warehouse-size).",
			ValidateDiag: sdkValidation(sdk.ToWarehouseSize),
			DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize),
		},
		{
			Name:         sdk.ObjectParameterUserTaskTimeoutMs,
			Type:         schema.TypeInt,
			Description:  "User task execution timeout in milliseconds. For more information, see [USER_TASK_TIMEOUT_MS](https://docs.snowflake.com/en/sql-reference/parameters#user-task-timeout-ms).",
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 86400000)),
		},
		{
			Name:        sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds,
			Type:        schema.TypeInt,
			Description: "Minimum amount of time between Triggered Task executions in seconds.",
			// TODO(DOC-2511): ValidateDiag: Not documented
		},
		{
			Name:        sdk.ObjectParameterQuotedIdentifiersIgnoreCase,
			Type:        schema.TypeBool,
			Description: "If true, the case of quoted identifiers is ignored. For more information, see [QUOTED_IDENTIFIERS_IGNORE_CASE](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case).",
		},
		{
			Name:        sdk.ObjectParameterEnableConsoleOutput,
			Type:        schema.TypeBool,
			Description: "If true, enables stdout/stderr fast path logging for anonymous stored procedures.",
		},
	}

	for _, field := range databaseParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		databaseParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
		}

		if !slices.Contains(sharedDatabaseNotApplicableParameters, field.Name) {
			sharedDatabaseParametersSchema[fieldName] = &schema.Schema{
				Type:             field.Type,
				Description:      field.Description,
				ForceNew:         true,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: field.ValidateDiag,
				DiffSuppressFunc: field.DiffSuppress,
			}
		}
	}
}

func handleDatabaseParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateDatabaseOptions) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ObjectParameterDataRetentionTimeInDays, &createOpts.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &createOpts.MaxDataExtensionTimeInDays),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterExternalVolume, &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterCatalog, &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.ObjectParameterReplaceInvalidCharacters, &createOpts.ReplaceInvalidCharacters),
		handleParameterCreate(d, sdk.ObjectParameterDefaultDDLCollation, &createOpts.DefaultDDLCollation),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterLogLevel, &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterTraceLevel, &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.ObjectParameterEnableConsoleOutput, &createOpts.EnableConsoleOutput),
	)
}

func handleSecondaryDatabaseParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateSecondaryDatabaseOptions) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ObjectParameterDataRetentionTimeInDays, &createOpts.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &createOpts.MaxDataExtensionTimeInDays),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterExternalVolume, &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterCatalog, &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.ObjectParameterReplaceInvalidCharacters, &createOpts.ReplaceInvalidCharacters),
		handleParameterCreate(d, sdk.ObjectParameterDefaultDDLCollation, &createOpts.DefaultDDLCollation),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterLogLevel, &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterTraceLevel, &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.ObjectParameterEnableConsoleOutput, &createOpts.EnableConsoleOutput),
	)
}

func handleSharedDatabaseParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateSharedDatabaseOptions) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, sdk.ObjectParameterExternalVolume, &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterCatalog, &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.ObjectParameterReplaceInvalidCharacters, &createOpts.ReplaceInvalidCharacters),
		handleParameterCreate(d, sdk.ObjectParameterDefaultDDLCollation, &createOpts.DefaultDDLCollation),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterLogLevel, &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterTraceLevel, &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.ObjectParameterEnableConsoleOutput, &createOpts.EnableConsoleOutput),
	)
}

func handleDatabaseParametersChanges(d *schema.ResourceData, set *sdk.DatabaseSet, unset *sdk.DatabaseUnset) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ObjectParameterDataRetentionTimeInDays, &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterExternalVolume, &set.ExternalVolume, &unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterCatalog, &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, sdk.ObjectParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdate(d, sdk.ObjectParameterDefaultDDLCollation, &set.DefaultDDLCollation, &unset.DefaultDDLCollation),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterLogLevel, &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterTraceLevel, &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskTimeoutMs, &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, sdk.ObjectParameterEnableConsoleOutput, &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
	)
}

func handleDatabaseParameterRead(d *schema.ResourceData, databaseParameters []*sdk.Parameter) diag.Diagnostics {
	for _, parameter := range databaseParameters {
		switch parameter.Key {
		case
			string(sdk.ObjectParameterDataRetentionTimeInDays),
			string(sdk.ObjectParameterMaxDataExtensionTimeInDays),
			string(sdk.ObjectParameterSuspendTaskAfterNumFailures),
			string(sdk.ObjectParameterTaskAutoRetryAttempts),
			string(sdk.ObjectParameterUserTaskTimeoutMs),
			string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds):
			value, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterDefaultDDLCollation),
			string(sdk.ObjectParameterStorageSerializationPolicy),
			string(sdk.ObjectParameterLogLevel),
			string(sdk.ObjectParameterTraceLevel),
			string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize),
			string(sdk.ObjectParameterExternalVolume),
			string(sdk.ObjectParameterCatalog):
			if err := d.Set(strings.ToLower(parameter.Key), parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterReplaceInvalidCharacters),
			string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase),
			string(sdk.ObjectParameterEnableConsoleOutput):
			value, err := strconv.ParseBool(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
