//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Task_ProveSessionParameterBehavior(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	providerModel := providermodel.SnowflakeProvider().WithParamsValue(
		configvariable.ObjectVariable(
			map[string]configvariable.Variable{
				"statement_timeout_in_seconds": configvariable.IntegerVariable(12345),
			},
		),
	)
	taskModel := model.TaskWithId("test", id, false, statement)
	executeCheckOnSession := model.ExecuteWithNoOpActions("t1").
		WithQuery("SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION")
	executeCheckOnTask := model.ExecuteWithNoOpActions("t2").
		WithQuery(fmt.Sprintf(`SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN TASK %s`, id.FullyQualifiedName())).
		WithDependsOn(taskModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Task_ProveSessionParameterBehavior"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, taskModel, executeCheckOnSession, executeCheckOnTask),
				Check: assertThat(
					t,
					resourceassert.ExecuteResource(t, executeCheckOnSession.ResourceReference()).
						HasQueryResultsLength(1).
						HasKeyValueOnIdx("value", "12345", 0).
						HasKeyValueOnIdx("level", string(sdk.ParameterTypeSession), 0),

					// the parameter set on session is not used in object creation
					resourceassert.ExecuteResource(t, executeCheckOnTask.ResourceReference()).
						HasQueryResultsLength(1).
						HasKeyValueOnIdx("value", "172800", 0).
						HasKeyValueOnIdx("level", "", 0),
				),
			},
		},
	})
}

func TestAcc_Task_ProveCurrentDriftBehavior(t *testing.T) {
	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	taskModel := model.TaskWithId("test", id, false, statement)
	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	providerModelWithFeatureEnabled := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ParametersIgnoreValueChangesIfNotOnObjectLevel)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactoryWithoutCache(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, taskModel),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, taskModel.ResourceReference()).
						HasNameString(id.Name()).
						HasStatementTimeoutInSecondsString("172800"),
				),
			},
			{
				PreConfig: func() {
					revertParameter := secondaryTestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43200")
					t.Cleanup(revertParameter)
				},
				Config: config.FromModels(t, providerModel, taskModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(taskModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, taskModel.ResourceReference()).
						HasNameString(id.Name()).
						HasStatementTimeoutInSecondsString("43200"),
					// modifying the account-level parameter one more time as assertion, as this is the moment in-between apply and refresh
					invokeactionassert.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43201", secondaryTestClient()),
				),
				ExpectError: regexp.MustCompile(`(?s)After applying this test step, the non-refresh plan was not empty.*~ update in-place.*# snowflake_task.test will be updated in-place.*~ statement_timeout_in_seconds\s+=\s+43200 -> \(known after apply\)`),
			},
			// one more time but with the experimental feature enabled
			{
				PreConfig: func() {
					revertParameter := secondaryTestClient().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43200")
					t.Cleanup(revertParameter)
				},
				Config: config.FromModels(t, providerModelWithFeatureEnabled, taskModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(taskModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, taskModel.ResourceReference()).
						HasNameString(id.Name()).
						HasStatementTimeoutInSecondsString("43200"),
					// modifying the account-level parameter one more time as assertion, as this is the moment in-between apply and refresh
					invokeactionassert.UpdateAccountParameterTemporarily(t, sdk.AccountParameterStatementTimeoutInSeconds, "43201", secondaryTestClient()),
				),
			},
		},
	})
}

// All below tests in this file are temporarily moved to account level tests due to STATEMENT_TIMEOUT_IN_SECONDS being set on warehouse level and messing with the results.

func TestAcc_Task_BasicUseCase(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	configModel := model.TaskWithId("test", id, false, statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasExecuteAsUserString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasExecuteAsUserEmpty().
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanFalse).
						HasErrorIntegrationString("").
						HasExecuteAsUserString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString("").
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

// execute_as_user is omitted: it needs IMPERSONATE on the user and the owner role granted to that user.
// Covered by TestAcc_Task_ExecuteAsUser so this portable complete smoke stays free of that setup.
func TestAcc_Task_CompleteUseCase(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	errorNotificationIntegration := gcpPubSubNotificationIntegration()

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	taskConfig := `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	configModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(10).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(testClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasExecuteAsUserString("").
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasWarehouse(testClient().Ids.WarehouseId()).
						HasScheduleMinutes(10).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasExecuteAsUserEmpty().
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(taskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName:    configModel.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(testClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasExecuteAsUserString("").
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString(condition).
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

func TestAcc_Task_CompleteUseCase_AllParameters(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	configModel := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5)
	configModelWithAllParametersSet := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(15).
		WithTaskAutoRetryAttempts(15).
		WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall).
		WithUserTaskMinimumTriggerIntervalInSeconds(30).
		WithUserTaskTimeoutMs(1000).
		WithAbortDetachedQuery(true).
		WithAutocommit(true).
		WithBinaryInputFormatEnum(sdk.BinaryInputFormatUTF8).
		WithBinaryOutputFormatEnum(sdk.BinaryOutputFormatBase64).
		WithClientMemoryLimit(1024).
		WithClientMetadataRequestUseConnectionCtx(true).
		WithClientPrefetchThreads(2).
		WithClientResultChunkSize(48).
		WithClientResultColumnCaseInsensitive(true).
		WithClientSessionKeepAlive(true).
		WithClientSessionKeepAliveHeartbeatFrequency(2400).
		WithClientTimestampTypeMappingEnum(sdk.ClientTimestampTypeMappingNtz).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YY-MM-DD").
		WithEnableUnloadPhysicalTypeOptimization(false).
		WithErrorOnNondeterministicMerge(false).
		WithErrorOnNondeterministicUpdate(true).
		WithGeographyOutputFormatEnum(sdk.GeographyOutputFormatWKB).
		WithGeometryOutputFormatEnum(sdk.GeometryOutputFormatWKB).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithLockTimeout(21222).
		WithLogLevelEnum(sdk.LogLevelError).
		WithLogEventLevelEnum(sdk.LogLevelError).
		WithMultiStatementCount(0).
		WithNoorderSequenceAsDefault(false).
		WithOdbcTreatDecimalAsInt(true).
		WithQueryTag("some_tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithRowsPerResultset(2).
		WithS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
		WithStatementQueuedTimeoutInSeconds(10).
		WithStatementTimeoutInSeconds(10).
		WithStrictJsonOutput(true).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampTypeMappingEnum(sdk.TimestampTypeMappingLtz).
		WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimezone("Europe/Warsaw").
		WithTimeInputFormat("HH24:MI").
		WithTimeOutputFormat("HH24:MI").
		WithTraceLevelEnum(sdk.TraceLevelPropagate).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevelEnum(sdk.TransactionDefaultIsolationLevelReadCommitted).
		WithTwoDigitCenturyStart(1980).
		WithUnsupportedDdlActionEnum(sdk.UnsupportedDDLActionFail).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(
					t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaultsForEnvironment(t, testClient().SnowflakeDefaults).
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName:    configModel.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assertThatImport(
					t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				Check: assertThat(
					t,
					objectparametersassert.TaskParameters(t, id).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(true).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasLogEventLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$current, $public").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
					resourceparametersassert.TaskResourceParameters(t, configModelWithAllParametersSet.ResourceReference()).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(true).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasLogEventLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$current, $public").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// import when all parameters set
			{
				ResourceName:    configModelWithAllParametersSet.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				ImportStateCheck: assertThatImport(
					t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(true).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasLogEventLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$current, $public").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// unset all the parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(
					t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaultsForEnvironment(t, testClient().SnowflakeDefaults).
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
		},
	})
}

// Covers CREATE with execute_as_user set, then ALTER SET / UNSET, import, and external drift revert.
func TestAcc_Task_ExecuteAsUser(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	user1, user1Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)
	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	testClient().Grant.GrantPrivilegesOnUserToAccountRole(t, currentRole, user1.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate}, false)
	t.Cleanup(func() {
		testClient().Grant.RevokePrivilegesOnUserFromAccountRole(t, currentRole, user1.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate})
	})
	testClient().Grant.GrantPrivilegesOnUserToAccountRole(t, currentRole, user2.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate}, false)
	t.Cleanup(func() {
		testClient().Grant.RevokePrivilegesOnUserFromAccountRole(t, currentRole, user2.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate})
	})
	testClient().Role.GrantRoleToUser(t, currentRole, user1.ID())
	testClient().Role.GrantRoleToUser(t, currentRole, user2.ID())

	configUser1 := model.TaskWithId("test", id, false, statement).WithExecuteAsUser(user1.ID().Name())
	configUser2 := model.TaskWithId("test", id, false, statement).WithExecuteAsUser(user2.ID().Name())
	configUnset := model.TaskWithId("test", id, false, statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configUser1),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configUser1.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasExecuteAsUserString(user1.ID().Name()).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configUser1.ResourceReference()).
						HasName(id.Name()).
						HasExecuteAsUser(user1.ID()),
					objectassert.Task(t, id).
						HasExecuteAsUser(user1.ID()),
				),
			},
			{
				ResourceName: configUser1.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasExecuteAsUserString(user1.ID().Name()).
						HasSqlStatementString(statement),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configUser2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configUser2),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configUser2.ResourceReference()).
						HasExecuteAsUserString(user2.ID().Name()),
					resourceshowoutputassert.TaskShowOutput(t, configUser2.ResourceReference()).
						HasExecuteAsUser(user2.ID()),
					objectassert.Task(t, id).
						HasExecuteAsUser(user2.ID()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configUnset.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configUnset),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configUnset.ResourceReference()).
						HasExecuteAsUserString(""),
					resourceshowoutputassert.TaskShowOutput(t, configUnset.ResourceReference()).
						HasExecuteAsUserEmpty(),
					objectassert.Task(t, id).
						HasNoExecuteAsUser(),
				),
			},
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSetExecuteAsUser(user1.ID()))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(configUnset.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, configUnset),
				Check: assertThat(
					t,
					resourceassert.TaskResource(t, configUnset.ResourceReference()).
						HasExecuteAsUserString(""),
					resourceshowoutputassert.TaskShowOutput(t, configUnset.ResourceReference()).
						HasExecuteAsUserEmpty(),
					objectassert.Task(t, id).
						HasNoExecuteAsUser(),
				),
			},
		},
	})
}
