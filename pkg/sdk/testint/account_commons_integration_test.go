package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// setAndUnsetAccountParametersTest is a common test used for different account kinds.
func setAndUnsetAccountParametersTest(
	setParameters func(ctx context.Context, parameters sdk.AccountParameters) error,
	unsetAllParameters func(ctx context.Context) error,
	showParameters func(ctx context.Context) ([]*sdk.Parameter, error),
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		id := testClientHelper().Context.CurrentAccountId(t)

		warehouseId := testClientHelper().Ids.WarehouseId()

		eventTable, eventTableCleanup := testClientHelper().EventTable.Create(t)
		t.Cleanup(eventTableCleanup)

		externalVolumeId, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		createNetworkPolicyRequest := sdk.NewCreateNetworkPolicyRequest(testClientHelper().Ids.RandomAccountObjectIdentifier()).WithAllowedIpList([]sdk.IPRequest{*sdk.NewIPRequest("0.0.0.0/0")})
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, createNetworkPolicyRequest)
		t.Cleanup(networkPolicyCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		// TODO(SNOW-2138715): Test all parameters, the following parameters were not tested due to more complex setup:
		// - ActivePythonProfiler
		// - CatalogSync
		// - EnableInternalStagesPrivatelink
		// - PythonProfilerModules
		// - S3StageVpceDnsName
		// - SimulatedDataSharingConsumer
		err := setParameters(context.Background(), sdk.AccountParameters{
			AbortDetachedQuery:                            new(true),
			AllowBindValuesAccess:                         new(true),
			AllowClientMFACaching:                         new(true),
			AllowedSpcsWorkloadTypes:                      new("ALL"),
			AllowIDToken:                                  new(true),
			Autocommit:                                    new(false),
			BaseLocationPrefix:                            new("STORAGE_BASE_URL/"),
			BinaryInputFormat:                             sdk.Pointer(sdk.BinaryInputFormatBase64),
			BinaryOutputFormat:                            sdk.Pointer(sdk.BinaryOutputFormatBase64),
			Catalog:                                       new(helpers.TestDatabaseCatalog.Name()),
			ClientEnableLogInfoStatementParameters:        new(true),
			ClientEncryptionKeySize:                       new(256),
			ClientMemoryLimit:                             new(1540),
			ClientMetadataRequestUseConnectionCtx:         new(true),
			ClientMetadataUseSessionDatabase:              new(true),
			ClientPrefetchThreads:                         new(5),
			ClientResultChunkSize:                         new(159),
			ClientResultColumnCaseInsensitive:             new(true),
			ClientSessionKeepAlive:                        new(true),
			ClientSessionKeepAliveHeartbeatFrequency:      new(3599),
			ClientTimestampTypeMapping:                    sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
			CortexEnabledCrossRegion:                      new("ANY_REGION"),
			CortexModelsAllowlist:                         new("All"),
			CsvTimestampFormat:                            new("YYYY-MM-DD"),
			DataMetricSchedule:                            new("60 MINUTES"),
			DataRetentionTimeInDays:                       new(2),
			DateInputFormat:                               new("YYYY-MM-DD"),
			DateOutputFormat:                              new("YYYY-MM-DD"),
			DefaultDbtVersion:                             new("1.9.4"),
			DefaultDDLCollation:                           new("en-cs"),
			DefaultNotebookComputePoolCpu:                 new("CPU_X64_S"),
			DefaultNotebookComputePoolGpu:                 new("GPU_NV_S"),
			DefaultNullOrdering:                           sdk.Pointer(sdk.DefaultNullOrderingFirst),
			DefaultStreamlitNotebookWarehouse:             new(warehouseId),
			DisallowedSpcsWorkloadTypes:                   new(""),
			DisableUiDownloadButton:                       new(true),
			DisableUserPrivilegeGrants:                    new(true),
			EnableAutomaticSensitiveDataClassificationLog: new(false),
			EnableBudgetEventLogging:                      new(true),
			EnableDataCompaction:                          new(true),
			EnableEgressCostOptimizer:                     new(false),
			EnableGetDdlUseDataTypeAlias:                  new(false),
			EnableIcebergMergeOnRead:                      new(true),
			EnableNotebookCreationInPersonalDb:            new(false),
			EnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement: new(false),
			EnableTagPropagationEventLogging:                         new(false),
			EnableIdentifierFirstLogin:                               new(false),
			EnableTriSecretAndRekeyOptOutForImageRepository:          new(true),
			EnableTriSecretAndRekeyOptOutForSpcsBlockStorage:         new(true),
			EnableUnhandledExceptionsReporting:                       new(false),
			EnableUnloadPhysicalTypeOptimization:                     new(false),
			EnableUnredactedQuerySyntaxError:                         new(true),
			EnableUnredactedSecureObjectError:                        new(true),
			EnforceNetworkRulesForInternalStages:                     new(true),
			ErrorOnNondeterministicMerge:                             new(false),
			ErrorOnNondeterministicUpdate:                            new(true),
			EventTable:                                               new(eventTable.ID()),
			ExternalOAuthAddPrivilegedRolesToBlockedList:             new(false),
			ExternalVolume:                                           new(externalVolumeId),
			GeographyOutputFormat:                                    sdk.Pointer(sdk.GeographyOutputFormatWKT),
			GeometryOutputFormat:                                     sdk.Pointer(sdk.GeometryOutputFormatWKT),
			HybridTableLockTimeout:                                   new(3599),
			IcebergVersionDefault:                                    new(2),
			InitialReplicationSizeLimitInTB:                          new("9.9"),
			JdbcTreatDecimalAsInt:                                    new(false),
			JdbcTreatTimestampNtzAsUtc:                               new(true),
			JdbcUseSessionTimezone:                                   new(false),
			JsonIndent:                                               new(4),
			JsTreatIntegerAsBigInt:                                   new(true),
			ListingAutoFulfillmentReplicationRefreshSchedule:         new("2 minutes"),
			LockTimeout:                                              new(43201),
			LogLevel:                                                 sdk.Pointer(sdk.LogLevelInfo),
			LogEventLevel:                                            sdk.Pointer(sdk.LogLevelInfo),
			MaxConcurrencyLevel:                                      new(7),
			MaxDataExtensionTimeInDays:                               new(13),
			MetricLevel:                                              sdk.Pointer(sdk.MetricLevelAll),
			MinDataRetentionTimeInDays:                               new(1),
			MultiStatementCount:                                      new(0),
			NetworkPolicy:                                            new(networkPolicy.ID()),
			NoorderSequenceAsDefault:                                 new(false),
			OAuthAddPrivilegedRolesToBlockedList:                     new(false),
			OdbcTreatDecimalAsInt:                                    new(true),
			PeriodicDataRekeying:                                     new(false),
			PipeExecutionPaused:                                      new(true),
			PreventUnloadToInlineURL:                                 new(true),
			PreventUnloadToInternalStages:                            new(true),
			PythonProfilerTargetStage:                                new(stage.ID()),
			QueryTag:                                                 new("test-query-tag"),
			QuotedIdentifiersIgnoreCase:                              new(true),
			ReadConsistencyMode:                                      new("SESSION"),
			ReplaceInvalidCharacters:                                 new(true),
			RequireStorageIntegrationForStageCreation:                new(true),
			RequireStorageIntegrationForStageOperation:               new(true),
			RowTimestampDefault:                                      new(false),
			RowsPerResultset:                                         new(1000),
			SearchPath:                                               new("$current, $public"),
			ServerlessTaskMaxStatementSize:                           sdk.Pointer(sdk.WarehouseSize("6X-LARGE")),
			ServerlessTaskMinStatementSize:                           sdk.Pointer(sdk.WarehouseSizeSmall),
			SsoLoginPage:                                             new(true),
			SqlTraceQueryText:                                        new("OFF"),
			StatementQueuedTimeoutInSeconds:                          new(1),
			StatementTimeoutInSeconds:                                new(1),
			StorageSerializationPolicy:                               sdk.Pointer(sdk.StorageSerializationPolicyOptimized),
			StrictJsonOutput:                                         new(true),
			SuspendTaskAfterNumFailures:                              new(3),
			TaskAutoRetryAttempts:                                    new(3),
			TimestampDayIsAlways24h:                                  new(true),
			TimestampInputFormat:                                     new("YYYY-MM-DD"),
			TimestampLtzOutputFormat:                                 new("YYYY-MM-DD"),
			TimestampNtzOutputFormat:                                 new("YYYY-MM-DD"),
			TimestampOutputFormat:                                    new("YYYY-MM-DD"),
			TimestampTypeMapping:                                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
			TimestampTzOutputFormat:                                  new("YYYY-MM-DD"),
			Timezone:                                                 new("Europe/London"),
			TimeInputFormat:                                          new("YYYY-MM-DD"),
			TimeOutputFormat:                                         new("YYYY-MM-DD"),
			TraceLevel:                                               sdk.Pointer(sdk.TraceLevelPropagate),
			TransactionAbortOnError:                                  new(true),
			TransactionDefaultIsolationLevel:                         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
			TwoDigitCenturyStart:                                     new(1971),
			UnsupportedDdlAction:                                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
			UserTaskManagedInitialWarehouseSize:                      sdk.Pointer(sdk.WarehouseSizeX6Large),
			UserTaskMinimumTriggerIntervalInSeconds:                  new(10),
			UserTaskTimeoutMs:                                        new(10),
			UseCachedResult:                                          new(false),
			UseWorkspacesForSql:                                      new("unset"),
			WeekOfYearPolicy:                                         new(1),
			WeekStart:                                                new(1),
		})
		require.NoError(t, err)

		parameters, err := showParameters(context.Background())
		require.NoError(t, err)

		objectparametersassert.AccountParametersPrefetched(t, id, parameters).
			HasAbortDetachedQuery(true).
			HasAllowClientMfaCaching(true).
			HasAllowIdToken(true).
			HasAutocommit(false).
			HasBaseLocationPrefix("STORAGE_BASE_URL/").
			HasBinaryInputFormat(sdk.BinaryInputFormatBase64).
			HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
			HasCatalog(helpers.TestDatabaseCatalog.Name()).
			HasClientEnableLogInfoStatementParameters(true).
			HasClientEncryptionKeySize(256).
			HasClientMemoryLimit(1540).
			HasClientMetadataRequestUseConnectionCtx(true).
			HasClientMetadataUseSessionDatabase(true).
			HasClientPrefetchThreads(5).
			HasClientResultChunkSize(159).
			HasClientResultColumnCaseInsensitive(true).
			HasClientSessionKeepAlive(true).
			HasClientSessionKeepAliveHeartbeatFrequency(3599).
			HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
			HasCortexEnabledCrossRegion("ANY_REGION").
			HasCortexModelsAllowlist("All").
			HasCsvTimestampFormat("YYYY-MM-DD").
			HasDataRetentionTimeInDays(2).
			HasDateInputFormat("YYYY-MM-DD").
			HasDateOutputFormat("YYYY-MM-DD").
			HasDefaultDdlCollation("en-cs").
			HasDefaultNotebookComputePoolCpu("CPU_X64_S").
			HasDefaultNotebookComputePoolGpu("GPU_NV_S").
			HasDefaultNullOrdering(sdk.DefaultNullOrderingFirst).
			HasDefaultStreamlitNotebookWarehouse(warehouseId.Name()).
			HasDisableUiDownloadButton(true).
			HasDisableUserPrivilegeGrants(true).
			HasEnableAutomaticSensitiveDataClassificationLog(false).
			HasEnableEgressCostOptimizer(false).
			HasEnableIdentifierFirstLogin(false).
			HasEnableTriSecretAndRekeyOptOutForImageRepository(true).
			HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorage(true).
			HasEnableUnhandledExceptionsReporting(false).
			HasEnableUnloadPhysicalTypeOptimization(false).
			HasEnableUnredactedQuerySyntaxError(true).
			HasEnableUnredactedSecureObjectError(true).
			HasEnforceNetworkRulesForInternalStages(true).
			HasErrorOnNondeterministicMerge(false).
			HasErrorOnNondeterministicUpdate(true).
			HasEventTable(eventTable.ID().FullyQualifiedName()).
			HasExternalOauthAddPrivilegedRolesToBlockedList(false).
			HasExternalVolume(externalVolumeId.Name()).
			HasGeographyOutputFormat(sdk.GeographyOutputFormatWKT).
			HasGeometryOutputFormat(sdk.GeometryOutputFormatWKT).
			HasHybridTableLockTimeout(3599).
			HasInitialReplicationSizeLimitInTb("9.9").
			HasJdbcTreatDecimalAsInt(false).
			HasJdbcTreatTimestampNtzAsUtc(true).
			HasJdbcUseSessionTimezone(false).
			HasJsonIndent(4).
			HasJsTreatIntegerAsBigint(true).
			HasListingAutoFulfillmentReplicationRefreshSchedule("2 minutes").
			HasLockTimeout(43201).
			HasLogLevel(sdk.LogLevelInfo).
			HasLogEventLevel(sdk.LogLevelInfo).
			HasMaxConcurrencyLevel(7).
			HasMaxDataExtensionTimeInDays(13).
			HasMetricLevel(sdk.MetricLevelAll).
			HasMinDataRetentionTimeInDays(1).
			HasMultiStatementCount(0).
			HasNetworkPolicy(networkPolicy.ID().Name()).
			HasNoorderSequenceAsDefault(false).
			HasOauthAddPrivilegedRolesToBlockedList(false).
			HasOdbcTreatDecimalAsInt(true).
			HasPeriodicDataRekeying(false).
			HasPipeExecutionPaused(true).
			HasPreventUnloadToInlineUrl(true).
			HasPreventUnloadToInternalStages(true).
			HasQueryTag("test-query-tag").
			HasQuotedIdentifiersIgnoreCase(true).
			HasReplaceInvalidCharacters(true).
			HasRequireStorageIntegrationForStageCreation(true).
			HasRequireStorageIntegrationForStageOperation(true).
			HasRowsPerResultset(1000).
			HasSearchPath("$current, $public").
			HasServerlessTaskMaxStatementSize("6X-LARGE").
			HasServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
			HasSsoLoginPage(true).
			HasStatementQueuedTimeoutInSeconds(1).
			HasStatementTimeoutInSeconds(1).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStrictJsonOutput(true).
			HasSuspendTaskAfterNumFailures(3).
			HasTaskAutoRetryAttempts(3).
			HasTimestampDayIsAlways24h(true).
			HasTimestampInputFormat("YYYY-MM-DD").
			HasTimestampLtzOutputFormat("YYYY-MM-DD").
			HasTimestampNtzOutputFormat("YYYY-MM-DD").
			HasTimestampOutputFormat("YYYY-MM-DD").
			HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
			HasTimestampTzOutputFormat("YYYY-MM-DD").
			HasTimezone("Europe/London").
			HasTimeInputFormat("YYYY-MM-DD").
			HasTimeOutputFormat("YYYY-MM-DD").
			HasTraceLevel(sdk.TraceLevelPropagate).
			HasTransactionAbortOnError(true).
			HasTransactionDefaultIsolationLevel(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
			HasTwoDigitCenturyStart(1971).
			HasUnsupportedDdlAction(string(sdk.UnsupportedDDLActionFail)).
			HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeX6Large).
			HasUserTaskMinimumTriggerIntervalInSeconds(10).
			HasUserTaskTimeoutMs(10).
			HasUseCachedResult(false).
			HasWeekOfYearPolicy(1).
			HasWeekStart(1).
			HasAllowBindValuesAccess(true).
			HasAllowedSpcsWorkloadTypes("ALL").
			HasDataMetricSchedule("60 MINUTES").
			HasDefaultDbtVersion("1.9.4").
			HasDisallowedSpcsWorkloadTypes("").
			HasEnableBudgetEventLogging(true).
			HasEnableDataCompaction(true).
			HasEnableGetDdlUseDataTypeAlias(false).
			HasEnableIcebergMergeOnRead(true).
			HasEnableNotebookCreationInPersonalDb(false).
			HasEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement(false).
			HasEnableTagPropagationEventLogging(false).
			HasIcebergVersionDefault(2).
			HasReadConsistencyMode("SESSION").
			HasRowTimestampDefault(false).
			HasSqlTraceQueryText("OFF").
			HasUseWorkspacesForSql("unset")

		err = unsetAllParameters(context.Background())
		require.NoError(t, err)

		parameters, err = showParameters(context.Background())
		require.NoError(t, err)

		objectparametersassert.AccountParametersPrefetched(t, id, parameters).
			HasAllDefaults()
	}
}

func assertThatPolicyIsSetOnAccount(t *testing.T, kind sdk.PolicyKind, id sdk.SchemaObjectIdentifier) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.NoError(t, err)
	_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
		return reference.PolicyName == id.Name() && reference.PolicyKind == kind
	})
	require.NoError(t, err)
}

func assertThatNoPolicyIsSetOnAccount(t *testing.T) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.Empty(t, policies)
	require.NoError(t, err)
}
