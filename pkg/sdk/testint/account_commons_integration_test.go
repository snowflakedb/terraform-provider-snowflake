package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// setAndUnsetAccountParametersTest is a common test used for different account kinds.
func setAndUnsetAccountParametersTest(
	setParameters func(ctx context.Context, parameters sdk.AccountParametersRequest) error,
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
		err := setParameters(context.Background(),
			*sdk.NewAccountParametersRequest().
				WithAbortDetachedQuery(true).
				WithAllowBindValuesAccess(true).
				WithAllowClientMfaCaching(true).
				WithAllowedSpcsWorkloadTypes("ALL").
				WithAllowIdToken(true).
				WithAutocommit(false).
				WithBaseLocationPrefix("STORAGE_BASE_URL/").
				WithBinaryInputFormat(sdk.BinaryInputFormatBase64).
				WithBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
				WithCatalog(helpers.TestDatabaseCatalog).
				WithClientEnableLogInfoStatementParameters(true).
				WithClientEncryptionKeySize(256).
				WithClientMemoryLimit(1540).
				WithClientMetadataRequestUseConnectionCtx(true).
				WithClientMetadataUseSessionDatabase(true).
				WithClientPrefetchThreads(5).
				WithClientResultChunkSize(159).
				WithClientResultColumnCaseInsensitive(true).
				WithClientSessionKeepAlive(true).
				WithClientSessionKeepAliveHeartbeatFrequency(3599).
				WithClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
				WithCortexCodeCliDailyEstCreditLimitPerUser(10).
				WithCortexCodeDesktopDailyEstCreditLimitPerUser(20).
				WithCortexCodeSnowsightDailyEstCreditLimitPerUser(30).
				WithCortexEnabledCrossRegion("ANY_REGION").
				WithCortexModelsAllowlist("All").
				WithCsvTimestampFormat("YYYY-MM-DD").
				WithDataMetricSchedule("60 MINUTES").
				WithDataRetentionTimeInDays(2).
				WithDateInputFormat("YYYY-MM-DD").
				WithDateOutputFormat("YYYY-MM-DD").
				WithDefaultDbtVersion("1.9.4").
				WithDefaultDdlCollation(sdk.StringAllowEmpty{Value: "en-cs"}).
				WithDefaultNotebookComputePoolCpu("CPU_X64_S").
				WithDefaultNotebookComputePoolGpu("GPU_NV_S").
				WithDefaultNullOrdering(sdk.DefaultNullOrderingFirst).
				WithDefaultStreamlitComputePool("SYSTEM_COMPUTE_POOL_GPU").
				WithDefaultStreamlitNotebookWarehouse(warehouseId).
				WithDisallowedSpcsWorkloadTypes("").
				WithDisableUiDownloadButton(true).
				WithDisableUserPrivilegeGrants(true).
				WithEnableAutomaticSensitiveDataClassificationLog(false).
				WithEnableBudgetEventLogging(true).
				WithEnableDataCompaction(true).
				WithEnableEgressCostOptimizer(false).
				WithEnableGetDdlUseDataTypeAlias(false).
				WithEnableIcebergMergeOnRead(true).
				WithEnableNotebookCreationInPersonalDb(false).
				WithEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement(false).
				WithEnableTagPropagationEventLogging(false).
				WithEnableIdentifierFirstLogin(false).
				WithEnableTriSecretAndRekeyOptOutForImageRepository(true).
				WithEnableTriSecretAndRekeyOptOutForSpcsBlockStorage(true).
				WithEnableUnhandledExceptionsReporting(false).
				WithEnableUnloadPhysicalTypeOptimization(false).
				WithEnableUnredactedQuerySyntaxError(true).
				WithEnableUnredactedSecureObjectError(true).
				WithEnforceNetworkRulesForInternalStages(true).
				WithErrorOnNondeterministicMerge(false).
				WithErrorOnNondeterministicUpdate(true).
				WithEventTable(eventTable.ID()).
				WithExternalOauthAddPrivilegedRolesToBlockedList(false).
				WithExternalVolume(externalVolumeId).
				WithGeographyOutputFormat(sdk.GeographyOutputFormatWKT).
				WithGeometryOutputFormat(sdk.GeometryOutputFormatWKT).
				WithHybridTableLockTimeout(3599).
				WithIcebergVersionDefault(2).
				WithInitialReplicationSizeLimitInTb(9.9).
				WithJdbcTreatDecimalAsInt(false).
				WithJdbcTreatTimestampNtzAsUtc(true).
				WithJdbcUseSessionTimezone(false).
				WithJsonIndent(4).
				WithJsTreatIntegerAsBigint(true).
				WithListingAutoFulfillmentReplicationRefreshSchedule("2 minutes").
				WithLockTimeout(43201).
				WithLogLevel(sdk.LogLevelInfo).
				WithLogEventLevel(sdk.LogLevelInfo).
				WithMaxConcurrencyLevel(7).
				WithMaxDataExtensionTimeInDays(13).
				WithMetricLevel(sdk.MetricLevelAll).
				WithMinDataRetentionTimeInDays(1).
				WithMultiStatementCount(0).
				WithNetworkPolicy(networkPolicy.ID()).
				WithNoorderSequenceAsDefault(false).
				WithOauthAddPrivilegedRolesToBlockedList(false).
				WithOdbcTreatDecimalAsInt(true).
				WithPeriodicDataRekeying(false).
				WithPipeExecutionPaused(true).
				WithPreventUnloadToInlineUrl(true).
				WithPreventUnloadToInternalStages(true).
				WithPythonProfilerTargetStage(stage.ID()).
				WithQueryTag("test-query-tag").
				WithQuotedIdentifiersIgnoreCase(true).
				WithReadConsistencyMode("SESSION").
				WithReplaceInvalidCharacters(true).
				WithRequireStorageIntegrationForStageCreation(true).
				WithRequireStorageIntegrationForStageOperation(true).
				WithRowTimestampDefault(false).
				WithRowsPerResultset(1000).
				WithSearchPath("$current, $public").
				WithServerlessTaskMaxStatementSize(sdk.WarehouseSize("6X-LARGE")).
				WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
				WithSsoLoginPage(true).
				WithSqlTraceQueryText("OFF").
				WithStatementQueuedTimeoutInSeconds(1).
				WithStatementTimeoutInSeconds(1).
				WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
				WithStrictJsonOutput(true).
				WithSuspendTaskAfterNumFailures(3).
				WithTaskAutoRetryAttempts(3).
				WithTimestampDayIsAlways24H(true).
				WithTimestampInputFormat("YYYY-MM-DD").
				WithTimestampLtzOutputFormat("YYYY-MM-DD").
				WithTimestampNtzOutputFormat("YYYY-MM-DD").
				WithTimestampOutputFormat("YYYY-MM-DD").
				WithTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
				WithTimestampTzOutputFormat("YYYY-MM-DD").
				WithTimezone("Europe/London").
				WithTimeInputFormat("YYYY-MM-DD").
				WithTimeOutputFormat("YYYY-MM-DD").
				WithTraceLevel(sdk.TraceLevelPropagate).
				WithTransactionAbortOnError(true).
				WithTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
				WithTwoDigitCenturyStart(1971).
				WithUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
				WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeX6Large).
				WithUserTaskMinimumTriggerIntervalInSeconds(10).
				WithUserTaskTimeoutMs(10).
				WithUseCachedResult(false).
				WithUseWorkspacesForSql("unset").
				WithWeekOfYearPolicy(1).
				WithWeekStart(1),
		)
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
			HasCortexCodeCliDailyEstCreditLimitPerUser(10).
			HasCortexCodeDesktopDailyEstCreditLimitPerUser(20).
			HasCortexCodeSnowsightDailyEstCreditLimitPerUser(30).
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
			HasDefaultStreamlitComputePool("SYSTEM_COMPUTE_POOL_GPU").
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

// expectedAccountPolicy describes a single policy expected to be attached to the account.
type expectedAccountPolicy struct {
	Kind sdk.PolicyKind
	ID   sdk.SchemaObjectIdentifier
}

func assertThatPolicyIsSetOnAccount(t *testing.T, expectedPolicies ...expectedAccountPolicy) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.NoError(t, err)
	for _, expected := range expectedPolicies {
		_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
			return reference.PolicyName == expected.ID.Name() && reference.PolicyKind == expected.Kind
		})
		require.NoError(t, err, "expected a policy reference of kind %s with name %s to be set on the account", expected.Kind, expected.ID.Name())
	}
}

func assertThatAuthenticationPolicyIsSetOnAccountWithScopes(t *testing.T, id sdk.SchemaObjectIdentifier, expectedScopes ...sdk.AuthenticationPolicyTargetScope) {
	t.Helper()

	policies, err := testClient(t).AuthenticationPolicies.Show(context.Background(), sdk.NewShowAuthenticationPolicyRequest().WithOn(sdk.On{Account: new(true)}))
	require.NoError(t, err)

	attachedPolicy, err := collections.FindFirst(policies, func(p sdk.AuthenticationPolicy) bool {
		return p.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
	require.NoError(t, err, "expected authentication policy with name %s to be returned by SHOW AUTHENTICATION POLICIES ON ACCOUNT", id.FullyQualifiedName())

	assertThatObject(
		t, objectassert.AuthenticationPolicyFromObject(t, attachedPolicy).
			HasTargetScopes(expectedScopes...),
	)
}

func assertThatSessionPolicyIsSetOnAccountWithScopes(t *testing.T, id sdk.SchemaObjectIdentifier, expectedScopes ...sdk.SessionPolicyTargetScope) {
	t.Helper()

	policies, err := testClient(t).SessionPolicies.Show(context.Background(), sdk.NewShowSessionPolicyRequest().WithOn(sdk.On{Account: new(true)}))
	require.NoError(t, err)

	attachedPolicy, err := collections.FindFirst(policies, func(p sdk.SessionPolicy) bool {
		return p.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
	require.NoError(t, err, "expected session policy with name %s to be returned by SHOW SESSION POLICIES ON ACCOUNT", id.FullyQualifiedName())

	assertThatObject(
		t, objectassert.SessionPolicyFromObject(t, attachedPolicy).
			HasTargetScopes(expectedScopes...),
	)
}

func assertThatNoPolicyIsSetOnAccount(t *testing.T) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.Empty(t, policies)
	require.NoError(t, err)
}
