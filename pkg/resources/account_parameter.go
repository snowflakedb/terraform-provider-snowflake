package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO(next prs): Deprecate account_parameter in favor of current_account (the list should stay as client.Parameters.SetAccountParameter was not updated to handle newly defined parameters).
var accountParameterSupportedParameters = []sdk.AccountParameter{
	sdk.AccountParameterAllowClientMFACaching,
	sdk.AccountParameterAllowIDToken,
	sdk.AccountParameterClientEncryptionKeySize,
	sdk.AccountParameterCortexEnabledCrossRegion,
	sdk.AccountParameterDisableUserPrivilegeGrants,
	sdk.AccountParameterEnableIdentifierFirstLogin,
	sdk.AccountParameterEnableInternalStagesPrivatelink,
	sdk.AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository,
	sdk.AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage,
	sdk.AccountParameterEnableUnhandledExceptionsReporting,
	sdk.AccountParameterEnforceNetworkRulesForInternalStages,
	sdk.AccountParameterEventTable,
	sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList,
	sdk.AccountParameterInitialReplicationSizeLimitInTB,
	sdk.AccountParameterMinDataRetentionTimeInDays,
	sdk.AccountParameterNetworkPolicy,
	sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList,
	sdk.AccountParameterPeriodicDataRekeying,
	sdk.AccountParameterPreventLoadFromInlineURL,
	sdk.AccountParameterPreventUnloadToInlineURL,
	sdk.AccountParameterRequireStorageIntegrationForStageCreation,
	sdk.AccountParameterRequireStorageIntegrationForStageOperation,
	sdk.AccountParameterSsoLoginPage,

	sdk.AccountParameterAbortDetachedQuery,
	sdk.AccountParameterActivePythonProfiler,
	sdk.AccountParameterAutocommit,
	sdk.AccountParameterBinaryInputFormat,
	sdk.AccountParameterBinaryOutputFormat,
	sdk.AccountParameterClientEnableLogInfoStatementParameters,
	sdk.AccountParameterClientMemoryLimit,
	sdk.AccountParameterClientMetadataRequestUseConnectionCtx,
	sdk.AccountParameterClientMetadataUseSessionDatabase,
	sdk.AccountParameterClientPrefetchThreads,
	sdk.AccountParameterClientResultChunkSize,
	sdk.AccountParameterClientSessionKeepAlive,
	sdk.AccountParameterClientSessionKeepAliveHeartbeatFrequency,
	sdk.AccountParameterClientTimestampTypeMapping,
	sdk.AccountParameterEnableUnloadPhysicalTypeOptimization,
	sdk.AccountParameterClientResultColumnCaseInsensitive,
	sdk.AccountParameterCsvTimestampFormat,
	sdk.AccountParameterDateInputFormat,
	sdk.AccountParameterDateOutputFormat,
	sdk.AccountParameterErrorOnNondeterministicMerge,
	sdk.AccountParameterErrorOnNondeterministicUpdate,
	sdk.AccountParameterGeographyOutputFormat,
	sdk.AccountParameterGeometryOutputFormat,
	sdk.AccountParameterHybridTableLockTimeout,
	sdk.AccountParameterJdbcTreatDecimalAsInt,
	sdk.AccountParameterJdbcTreatTimestampNtzAsUtc,
	sdk.AccountParameterJdbcUseSessionTimezone,
	sdk.AccountParameterJsonIndent,
	sdk.AccountParameterJsTreatIntegerAsBigInt,
	sdk.AccountParameterLockTimeout,
	sdk.AccountParameterMultiStatementCount,
	sdk.AccountParameterNoorderSequenceAsDefault,
	sdk.AccountParameterOdbcTreatDecimalAsInt,
	sdk.AccountParameterPythonProfilerModules,
	sdk.AccountParameterPythonProfilerTargetStage,
	sdk.AccountParameterQueryTag,
	sdk.AccountParameterQuotedIdentifiersIgnoreCase,
	sdk.AccountParameterRowsPerResultset,
	sdk.AccountParameterS3StageVpceDnsName,
	sdk.AccountParameterSearchPath,
	sdk.AccountParameterSimulatedDataSharingConsumer,
	sdk.AccountParameterStatementTimeoutInSeconds,
	sdk.AccountParameterStrictJsonOutput,
	sdk.AccountParameterTimeInputFormat,
	sdk.AccountParameterTimeOutputFormat,
	sdk.AccountParameterTimestampDayIsAlways24h,
	sdk.AccountParameterTimestampInputFormat,
	sdk.AccountParameterTimestampLtzOutputFormat,
	sdk.AccountParameterTimestampNtzOutputFormat,
	sdk.AccountParameterTimestampOutputFormat,
	sdk.AccountParameterTimestampTypeMapping,
	sdk.AccountParameterTimestampTzOutputFormat,
	sdk.AccountParameterTimezone,
	sdk.AccountParameterTransactionAbortOnError,
	sdk.AccountParameterTransactionDefaultIsolationLevel,
	sdk.AccountParameterTwoDigitCenturyStart,
	sdk.AccountParameterUnsupportedDdlAction,
	sdk.AccountParameterUseCachedResult,
	sdk.AccountParameterWeekOfYearPolicy,
	sdk.AccountParameterWeekStart,

	sdk.AccountParameterCatalog,
	sdk.AccountParameterDataRetentionTimeInDays,
	sdk.AccountParameterDefaultDDLCollation,
	sdk.AccountParameterExternalVolume,
	sdk.AccountParameterLogLevel,
	sdk.AccountParameterMaxConcurrencyLevel,
	sdk.AccountParameterMaxDataExtensionTimeInDays,
	sdk.AccountParameterPipeExecutionPaused,
	sdk.AccountParameterPreventUnloadToInternalStages,
	sdk.AccountParameterReplaceInvalidCharacters,
	sdk.AccountParameterStatementQueuedTimeoutInSeconds,
	sdk.AccountParameterStorageSerializationPolicy,
	sdk.AccountParameterShareRestrictions,
	sdk.AccountParameterSuspendTaskAfterNumFailures,
	sdk.AccountParameterTraceLevel,
	sdk.AccountParameterUserTaskManagedInitialWarehouseSize,
	sdk.AccountParameterUserTaskTimeoutMs,
	sdk.AccountParameterTaskAutoRetryAttempts,
	sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds,
	sdk.AccountParameterMetricLevel,
	sdk.AccountParameterEnableConsoleOutput,
	sdk.AccountParameterEnableUnredactedQuerySyntaxError,
	sdk.AccountParameterEnablePersonalDatabase,
}

func ToAccountParameter(s string) (sdk.AccountParameter, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(accountParameterSupportedParameters, sdk.AccountParameter(s)) {
		return "", fmt.Errorf("invalid account parameter: %s", s)
	}
	return sdk.AccountParameter(s), nil
}

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToAccountParameter),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToAccountParameter),
		Description:      fmt.Sprintf("Name of account parameter. Valid values are (case-insensitive): %s. Deprecated parameters are not supported in the provider.", possibleValuesListed(sdk.AsStringList(accountParameterSupportedParameters))),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation. The parameter values are validated in Snowflake.",
	},
}

func AccountParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.AccountParameter, CreateAccountParameter),
		ReadContext:   TrackingReadWrapper(resources.AccountParameter, ReadAccountParameter),
		UpdateContext: TrackingUpdateWrapper(resources.AccountParameter, UpdateAccountParameter),
		DeleteContext: TrackingDeleteWrapper(resources.AccountParameter, DeleteAccountParameter),

		Description: "Resource used to manage current account parameters. For more information, check [parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters).",

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	parameter, err := ToAccountParameter(key)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Parameters.SetAccountParameter(ctx, parameter, value)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(string(parameter)))
	return ReadAccountParameter(ctx, d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parameterNameRaw := d.Id()
	parameterName, err := ToAccountParameter(parameterNameRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	parameter, err := client.Parameters.ShowAccountParameter(ctx, parameterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading account parameter: %w", err))
	}
	errs := errors.Join(
		d.Set("value", parameter.Value),
		d.Set("key", parameter.Key),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateAccountParameter(ctx, d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	parameter := sdk.AccountParameter(key)

	err := client.Parameters.UnsetAccountParameter(ctx, parameter)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", err))
	}

	d.SetId("")
	return nil
}
