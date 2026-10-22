//go:build account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tasks(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	sql := "SELECT CURRENT_TIMESTAMP"

	errorIntegration, ErrorIntegrationCleanup := testClientHelper().NotificationIntegration.CreateWithGcpPubSub(t)
	t.Cleanup(ErrorIntegrationCleanup)

	createUserForTaskImpersonation := func(t *testing.T) *sdk.User {
		t.Helper()
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)
		currentRole := testClientHelper().Context.CurrentRole(t)
		testClientHelper().Grant.GrantPrivilegesOnUserToAccountRole(t, currentRole, user.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate}, false)
		t.Cleanup(func() {
			testClientHelper().Grant.RevokePrivilegesOnUserFromAccountRole(t, currentRole, user.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeImpersonate})
		})
		testClientHelper().Role.GrantRoleToUser(t, currentRole, user.ID())
		return user
	}

	assertTask := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, warehouseId *sdk.AccountObjectIdentifier) {
		t.Helper()
		asserts := objectassert.TaskFromObject(t, task).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasIdNotEmpty().
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasOwner("ACCOUNTADMIN").
			HasComment("").
			HasSchedule("").
			HasPredecessorsInAnyOrder().
			HasState(sdk.TaskStateSuspended).
			HasDefinition(sql).
			HasCondition("").
			HasAllowOverlappingExecution(false).
			HasNoErrorIntegration().
			HasNoExecuteAsUser().
			HasLastCommittedOn("").
			HasLastSuspendedOn("").
			HasOwnerRoleType("ROLE").
			HasConfig("").
			HasBudget("").
			HasLastSuspendedOn("").
			HasTaskRelations(sdk.TaskRelations{})

		if warehouseId != nil {
			asserts.HasWarehouse(*warehouseId)
		} else {
			asserts.HasNoWarehouse()
		}

		assertThatObject(t, asserts)
	}

	assertTaskWithOptions := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, comment string, warehouse *sdk.AccountObjectIdentifier, schedule string, condition string, allowOverlappingExecution bool, config string, predecessor *sdk.SchemaObjectIdentifier, errorIntegrationName *sdk.AccountObjectIdentifier, executeAsUser *sdk.AccountObjectIdentifier) {
		t.Helper()

		asserts := objectassert.TaskFromObject(t, task).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasIdNotEmpty().
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasOwner("ACCOUNTADMIN").
			HasComment(comment).
			HasSchedule(schedule).
			HasState(sdk.TaskStateSuspended).
			HasDefinition(sql).
			HasCondition(condition).
			HasAllowOverlappingExecution(allowOverlappingExecution).
			HasLastCommittedOn("").
			HasLastSuspendedOn("").
			HasOwnerRoleType("ROLE").
			HasConfig(config).
			HasBudget("").
			HasLastSuspendedOn("")

		if predecessor != nil {
			asserts.HasPredecessorsInAnyOrder(*predecessor)
			asserts.HasTaskRelations(sdk.TaskRelations{
				Predecessors: []sdk.SchemaObjectIdentifier{*predecessor},
			})
		} else {
			asserts.HasPredecessorsInAnyOrder()
			asserts.HasTaskRelations(sdk.TaskRelations{})
		}

		if errorIntegrationName != nil {
			asserts.HasErrorIntegration(*errorIntegrationName)
		} else {
			asserts.HasNoErrorIntegration()
		}

		if executeAsUser != nil {
			asserts.HasExecuteAsUser(*executeAsUser)
		} else {
			asserts.HasNoExecuteAsUser()
		}

		if warehouse != nil {
			asserts.HasWarehouse(*warehouse)
		} else {
			asserts.HasNoWarehouse()
		}

		assertThatObject(t, asserts)
	}

	assertTaskTerse := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, schedule string) {
		t.Helper()
		assertThatObject(
			t, objectassert.TaskFromObject(t, task).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
				HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
				HasSchedule(schedule).
				// all below are not contained in the terse response, that's why all of them we expect to be empty
				HasId("").
				HasOwner("").
				HasComment("").
				HasNoWarehouse().
				HasPredecessorsInAnyOrder().
				HasState("").
				HasDefinition("").
				HasCondition("").
				HasAllowOverlappingExecution(false).
				HasNoErrorIntegration().
				HasNoExecuteAsUser().
				HasLastCommittedOn("").
				HasLastSuspendedOn("").
				HasOwnerRoleType("").
				HasConfig("").
				HasBudget("").
				HasLastSuspendedOn("").
				HasTaskRelations(sdk.TaskRelations{}),
		)
	}

	assertSessionParametersSet := func(parametersAssert *objectparametersassert.TaskParametersAssert) *objectparametersassert.TaskParametersAssert {
		return parametersAssert.
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
			HasJdbcTreatTimestampNtzAsUtc(true).
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
			HasWeekStart(1)
	}

	t.Run("create task: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTask(t, task, id, nil)

		assertThatObject(t, objectparametersassert.TaskParameters(t, id).HasAllDefaults())
	})

	t.Run("create task: with initial warehouse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall),
		)

		assertTask(t, task, id, nil)
	})

	t.Run("create task: complete case", func(t *testing.T) {
		user := createUserForTaskImpersonation(t)
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).
			WithOrReplace(true).
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
			WithErrorIntegration(errorIntegration.ID()).
			WithSchedule("10 MINUTE").
			WithConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
			WithAllowOverlappingExecution(true).
			WithAbortDetachedQuery(true).
			WithAutocommit(true).
			WithBinaryInputFormat(sdk.BinaryInputFormatUTF8).
			WithBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
			WithClientMemoryLimit(1024).
			WithClientMetadataRequestUseConnectionCtx(true).
			WithClientPrefetchThreads(2).
			WithClientResultChunkSize(48).
			WithClientResultColumnCaseInsensitive(true).
			WithClientSessionKeepAlive(true).
			WithClientSessionKeepAliveHeartbeatFrequency(2400).
			WithClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
			WithDateInputFormat("YYYY-MM-DD").
			WithDateOutputFormat("YY-MM-DD").
			WithEnableUnloadPhysicalTypeOptimization(false).
			WithErrorOnNondeterministicMerge(false).
			WithErrorOnNondeterministicUpdate(true).
			WithGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
			WithGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
			WithJdbcTreatTimestampNtzAsUtc(true).
			WithJdbcUseSessionTimezone(false).
			WithJsonIndent(4).
			WithLockTimeout(21222).
			WithLogLevel(sdk.LogLevelError).
			WithLogEventLevel(sdk.LogLevelError).
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
			WithTimestampDayIsAlways24H(true).
			WithTimestampInputFormat("YYYY-MM-DD").
			WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
			WithTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
			WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			WithTimezone("Europe/Warsaw").
			WithTimeInputFormat("HH24:MI").
			WithTimeOutputFormat("HH24:MI").
			WithTraceLevel(sdk.TraceLevelPropagate).
			WithTransactionAbortOnError(true).
			WithTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
			WithTwoDigitCenturyStart(1980).
			WithUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
			WithUseCachedResult(false).
			WithWeekOfYearPolicy(1).
			WithWeekStart(1).
			WithUserTaskTimeoutMs(500).
			WithSuspendTaskAfterNumFailures(3).
			WithTaskAutoRetryAttempts(5).
			WithUserTaskMinimumTriggerIntervalInSeconds(20).
			WithComment("some comment").
			WithExecuteAsUser(user.ID()).
			WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTaskWithOptions(t, task, id, "some comment", sdk.Pointer(testClientHelper().Ids.WarehouseId()), "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil, sdk.Pointer(errorIntegration.ID()), sdk.Pointer(user.ID()))
		assertThatObject(
			t, assertSessionParametersSet(objectparametersassert.TaskParameters(t, id)).
				HasUserTaskTimeoutMs(500).
				HasSuspendTaskAfterNumFailures(3).
				HasTaskAutoRetryAttempts(5).
				HasUserTaskMinimumTriggerIntervalInSeconds(20),
		)
	})

	t.Run("create task: with after", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(rootTaskId, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, rootTaskId))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithAfter([]sdk.SchemaObjectIdentifier{rootTaskId}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTaskWithOptions(t, task, id, "", nil, "", "", false, "", &rootTaskId, nil, nil)
	})

	t.Run("create task: with after and finalizer", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		finalizerId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(rootTaskId, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, rootTaskId))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithAfter([]sdk.SchemaObjectIdentifier{rootTaskId}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(finalizerId, sql).WithFinalize(rootTaskId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, finalizerId))

		assertThatObject(
			t, objectassert.Task(t, rootTaskId).
				HasTaskRelations(sdk.TaskRelations{
					Predecessors:  []sdk.SchemaObjectIdentifier{},
					FinalizerTask: &finalizerId,
				}),
		)

		assertThatObject(
			t, objectassert.Task(t, finalizerId).
				HasTaskRelations(sdk.TaskRelations{
					Predecessors:      []sdk.SchemaObjectIdentifier{},
					FinalizedRootTask: &rootTaskId,
				}),
		)
	})

	// Tested graph
	//		 t1
	// 	   /    \
	// root	     t3
	//	   \    /
	//		 t2
	t.Run("create dag of tasks", func(t *testing.T) {
		rootId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root, rootCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootId, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(rootCleanup)
		require.Empty(t, root.TaskRelations.Predecessors)

		t1, t1Cleanup := testClientHelper().Task.CreateWithAfter(t, rootId)
		t.Cleanup(t1Cleanup)
		require.Equal(t, []sdk.SchemaObjectIdentifier{rootId}, t1.TaskRelations.Predecessors)

		t2, t2Cleanup := testClientHelper().Task.CreateWithAfter(t, t1.ID(), rootId)
		t.Cleanup(t2Cleanup)

		require.Contains(t, t2.TaskRelations.Predecessors, rootId)
		require.Contains(t, t2.TaskRelations.Predecessors, t1.ID())
		require.Len(t, t2.TaskRelations.Predecessors, 2)

		t3, t3Cleanup := testClientHelper().Task.CreateWithAfter(t, t2.ID(), t1.ID())
		t.Cleanup(t3Cleanup)

		require.Contains(t, t3.TaskRelations.Predecessors, t2.ID())
		require.Contains(t, t3.TaskRelations.Predecessors, t1.ID())
		require.Len(t, t3.TaskRelations.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, rootId)
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t1.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t2.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t3.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		// cannot set ALLOW_OVERLAPPING_EXECUTION on child task
		alterRequest := sdk.NewAlterTaskRequest(t1.ID()).WithSet(*sdk.NewTaskSetRequest().WithAllowOverlappingExecution(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Cannot set allow_overlapping_execution on non-root task")

		// can set ALLOW_OVERLAPPING_EXECUTION on root task
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithSet(*sdk.NewTaskSetRequest().WithAllowOverlappingExecution(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// can create cycle, because DAG is suspended
		alterRequest = sdk.NewAlterTaskRequest(t1.ID()).WithAddAfter([]sdk.SchemaObjectIdentifier{t3.ID()})
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// can get the root task even with cycle
		rootTasks, err = sdk.GetRootTasks(client.Tasks, ctx, t3.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 1)
		require.Equal(t, rootId, rootTasks[0].ID())

		// we get an error when trying to start
		alterRequest = sdk.NewAlterTaskRequest(rootId).WithResume(true)
		err = client.Tasks.Alter(ctx, alterRequest)
		require.ErrorContains(t, err, "Graph has at least one cycle containing task")
	})

	// Tested graph
	// root1
	//      \
	//       t1
	//      /
	// root2
	t.Run("create dag of tasks - multiple roots", func(t *testing.T) {
		root1Id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root1, root1Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(root1Id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(root1Cleanup)
		require.Empty(t, root1.TaskRelations.Predecessors)

		root2Id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		root2, root2Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(root2Id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(root2Cleanup)
		require.Empty(t, root2.TaskRelations.Predecessors)

		t1, t1Cleanup := testClientHelper().Task.CreateWithAfter(t, root1.ID(), root2.ID())
		t.Cleanup(t1Cleanup)

		require.Contains(t, t1.TaskRelations.Predecessors, root1Id)
		require.Contains(t, t1.TaskRelations.Predecessors, root2Id)
		require.Len(t, t1.TaskRelations.Predecessors, 2)

		rootTasks, err := sdk.GetRootTasks(client.Tasks, ctx, t1.ID())
		require.NoError(t, err)
		require.Len(t, rootTasks, 2)
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[0].ID())
		require.Contains(t, []sdk.SchemaObjectIdentifier{root1Id, root2Id}, rootTasks[1].ID())

		// we get an error when trying to start
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(root1Id).WithResume(true))
		require.ErrorContains(t, err, "The graph has more than one root task (one without predecessors)")
	})

	t.Run("validate: finalizer set on non-root task", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		finalizerId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(rootTaskId, sql))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, rootTaskId))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).WithAfter([]sdk.SchemaObjectIdentifier{rootTaskId}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		err = testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(finalizerId, sql).WithFinalize(id))
		require.ErrorContains(t, err, "cannot finalize a non-root task")
	})

	t.Run("create task: with tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(
			t, sdk.NewCreateTaskRequest(id, sql).
				WithTag([]sdk.TagAssociation{
					{
						Name:  tag.ID(),
						Value: "v1",
					},
				}),
		)
		t.Cleanup(taskCleanup)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, sdk.NewGetTagRequest(*sdk.NewGetTagArgumentsRequest(tag.ID(), task.ID(), sdk.ObjectTypeTask)))
		require.NoError(t, err)

		assert.Equal(t, sdk.Pointer("v1"), returnedTagValue.Tag)
	})

	t.Run("create task: with serverless task parameters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		// Create a serverless task (WithUserTaskManagedInitialWarehouseSize) with serverless parameters
		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium)).
			WithTargetCompletionInterval("10 MINUTES").
			WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
			WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		// Verify parameters were set
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, id).
				HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
				HasServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
				HasServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge),
		)
		// target_completion_interval is returned by SHOW command, not SHOW PARAMETERS
		assertThatObject(
			t, objectassert.Task(t, id).
				HasNoWarehouse().
				HasTargetCompletionInterval(sdk.TaskTargetCompletionInterval{Minutes: sdk.Pointer(10)}),
		)
	})

	t.Run("clone task: default", func(t *testing.T) {
		rootTask, rootTaskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(rootTaskCleanup)

		sourceTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		sourceTask, taskCleanup := testClientHelper().Task.CreateWithRequest(
			t, sdk.NewCreateTaskRequest(sourceTaskId, sql).
				WithAfter([]sdk.SchemaObjectIdentifier{rootTask.ID()}).
				WithAllowOverlappingExecution(false).
				WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
				WithComment(random.Comment()).
				WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`),
		)
		t.Cleanup(taskCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tasks.Clone(ctx, sdk.NewCloneTaskRequest(id, sourceTask.ID()))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sourceTask.Definition, task.Definition)
		assert.Equal(t, sourceTask.Config, task.Config)
		assert.Equal(t, sourceTask.Condition, task.Condition)
		assert.Equal(t, sourceTask.Warehouse, task.Warehouse)
		assert.Equal(t, sourceTask.TaskRelations.Predecessors, task.TaskRelations.Predecessors)
		assert.Equal(t, sourceTask.AllowOverlappingExecution, task.AllowOverlappingExecution)
		assert.Equal(t, sourceTask.Comment, task.Comment)
		assert.Equal(t, sourceTask.ErrorIntegration, task.ErrorIntegration)
		assert.Equal(t, sourceTask.Schedule, task.Schedule)
		assert.Equal(t, sourceTask.TaskRelations, task.TaskRelations)
	})

	t.Run("create or alter: complete", func(t *testing.T) {
		user := createUserForTaskImpersonation(t)
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tasks.CreateOrAlter(
			ctx, sdk.NewCreateOrAlterTaskRequest(id, sql).
				WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
				WithSchedule("10 MINUTES").
				WithConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
				WithAllowOverlappingExecution(true).
				WithUserTaskTimeoutMs(10).
				WithAbortDetachedQuery(true).
				WithAutocommit(true).
				WithBinaryInputFormat(sdk.BinaryInputFormatUTF8).
				WithBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
				WithClientMemoryLimit(1024).
				WithClientMetadataRequestUseConnectionCtx(true).
				WithClientPrefetchThreads(2).
				WithClientResultChunkSize(48).
				WithClientResultColumnCaseInsensitive(true).
				WithClientSessionKeepAlive(true).
				WithClientSessionKeepAliveHeartbeatFrequency(2400).
				WithClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
				WithDateInputFormat("YYYY-MM-DD").
				WithDateOutputFormat("YY-MM-DD").
				WithEnableUnloadPhysicalTypeOptimization(false).
				WithErrorOnNondeterministicMerge(false).
				WithErrorOnNondeterministicUpdate(true).
				WithGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
				WithGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
				WithJdbcTreatTimestampNtzAsUtc(true).
				WithJdbcUseSessionTimezone(false).
				WithJsonIndent(4).
				WithLockTimeout(21222).
				WithLogLevel(sdk.LogLevelError).
				WithLogEventLevel(sdk.LogLevelError).
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
				WithTimestampDayIsAlways24H(true).
				WithTimestampInputFormat("YYYY-MM-DD").
				WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
				WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimezone("Europe/Warsaw").
				WithTimeInputFormat("HH24:MI").
				WithTimeOutputFormat("HH24:MI").
				WithTraceLevel(sdk.TraceLevelPropagate).
				WithTransactionAbortOnError(true).
				WithTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
				WithTwoDigitCenturyStart(1980).
				WithUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
				WithUseCachedResult(false).
				WithWeekOfYearPolicy(1).
				WithWeekStart(1).
				WithSuspendTaskAfterNumFailures(15).
				WithComment("some_comment").
				WithTaskAutoRetryAttempts(15).
				WithExecuteAsUser(user.ID()).
				WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)
		createdOn := task.CreatedOn

		assertThatObject(
			t, objectassert.TaskFromObject(t, task).
				HasWarehouse(testClientHelper().Ids.WarehouseId()).
				HasSchedule("10 MINUTES").
				HasConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
				HasAllowOverlappingExecution(true).
				HasCondition(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`).
				HasComment("some_comment").
				HasExecuteAsUser(user.ID()).
				HasTaskRelations(sdk.TaskRelations{}),
		)
		assertThatObject(
			t, assertSessionParametersSet(objectparametersassert.TaskParameters(t, task.ID()).
				HasUserTaskTimeoutMs(10).
				HasSuspendTaskAfterNumFailures(15).
				HasTaskAutoRetryAttempts(15)),
		)

		err = client.Tasks.CreateOrAlter(ctx, sdk.NewCreateOrAlterTaskRequest(id, sql))
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.TaskFromObject(t, alteredTask).
				HasNoWarehouse().
				HasSchedule("").
				HasConfig("").
				HasAllowOverlappingExecution(false).
				HasCondition("").
				HasComment("").
				HasNoExecuteAsUser().
				HasTaskRelations(sdk.TaskRelations{}),
		)
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasDefaultAutocommitValue().
				HasDefaultAbortDetachedQueryValue().
				HasDefaultUserTaskTimeoutMsValue().
				HasDefaultSuspendTaskAfterNumFailuresValue().
				HasDefaultTaskAutoRetryAttemptsValue(),
		)

		require.Equal(t, createdOn, alteredTask.CreatedOn)
	})

	t.Run("drop task: existing", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(task.ID()))
		require.NoError(t, err)

		_, err = client.Tasks.ShowByID(ctx, task.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				// TODO(SNOW-1519496): Cannot set warehouse due to Snowflake error
				// WithWarehouse(testClientHelper().Ids.WarehouseId()).
				WithErrorIntegration(errorIntegration.ID()).
				WithAbortDetachedQuery(true).
				WithAutocommit(true).
				WithBinaryInputFormat(sdk.BinaryInputFormatUTF8).
				WithBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
				WithClientMemoryLimit(1024).
				WithClientMetadataRequestUseConnectionCtx(true).
				WithClientPrefetchThreads(2).
				WithClientResultChunkSize(48).
				WithClientResultColumnCaseInsensitive(true).
				WithClientSessionKeepAlive(true).
				WithClientSessionKeepAliveHeartbeatFrequency(2400).
				WithClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
				WithDateInputFormat("YYYY-MM-DD").
				WithDateOutputFormat("YY-MM-DD").
				WithEnableUnloadPhysicalTypeOptimization(false).
				WithErrorOnNondeterministicMerge(false).
				WithErrorOnNondeterministicUpdate(true).
				WithGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
				WithGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
				WithJdbcTreatTimestampNtzAsUtc(true).
				WithJdbcUseSessionTimezone(false).
				WithJsonIndent(4).
				WithLockTimeout(21222).
				WithLogLevel(sdk.LogLevelError).
				WithLogEventLevel(sdk.LogLevelError).
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
				WithTimestampDayIsAlways24H(true).
				WithTimestampInputFormat("YYYY-MM-DD").
				WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
				WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
				WithTimezone("Europe/Warsaw").
				WithTimeInputFormat("HH24:MI").
				WithTimeOutputFormat("HH24:MI").
				WithTraceLevel(sdk.TraceLevelPropagate).
				WithTransactionAbortOnError(true).
				WithTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
				WithTwoDigitCenturyStart(1980).
				WithUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
				WithUseCachedResult(false).
				WithWeekOfYearPolicy(1).
				WithWeekStart(1).
				WithSchedule("10 MINUTE").
				WithConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
				WithAllowOverlappingExecution(true).
				WithUserTaskTimeoutMs(1000).
				WithSuspendTaskAfterNumFailures(100).
				WithComment("new comment").
				WithTaskAutoRetryAttempts(10).
				WithUserTaskMinimumTriggerIntervalInSeconds(15),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				// HasWarehouse(testClientHelper().Ids.WarehouseId().Name()).
				HasErrorIntegration(errorIntegration.ID()).
				HasSchedule("10 MINUTE").
				HasConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
				HasAllowOverlappingExecution(true).
				HasComment("new comment"),
		)
		assertThatObject(
			t, assertSessionParametersSet(objectparametersassert.TaskParameters(t, task.ID())).
				HasUserTaskTimeoutMs(1000).
				HasSuspendTaskAfterNumFailures(100).
				HasTaskAutoRetryAttempts(10).
				HasUserTaskMinimumTriggerIntervalInSeconds(15),
		)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(
			*sdk.NewTaskUnsetRequest().
				WithErrorIntegration(true).
				WithAbortDetachedQuery(true).
				WithAutocommit(true).
				WithBinaryInputFormat(true).
				WithBinaryOutputFormat(true).
				WithClientMemoryLimit(true).
				WithClientMetadataRequestUseConnectionCtx(true).
				WithClientPrefetchThreads(true).
				WithClientResultChunkSize(true).
				WithClientResultColumnCaseInsensitive(true).
				WithClientSessionKeepAlive(true).
				WithClientSessionKeepAliveHeartbeatFrequency(true).
				WithClientTimestampTypeMapping(true).
				WithDateInputFormat(true).
				WithDateOutputFormat(true).
				WithEnableUnloadPhysicalTypeOptimization(true).
				WithErrorOnNondeterministicMerge(true).
				WithErrorOnNondeterministicUpdate(true).
				WithGeographyOutputFormat(true).
				WithGeometryOutputFormat(true).
				WithJdbcTreatTimestampNtzAsUtc(true).
				WithJdbcUseSessionTimezone(true).
				WithJsonIndent(true).
				WithLockTimeout(true).
				WithLogLevel(true).
				WithLogEventLevel(true).
				WithMultiStatementCount(true).
				WithNoorderSequenceAsDefault(true).
				WithOdbcTreatDecimalAsInt(true).
				WithQueryTag(true).
				WithQuotedIdentifiersIgnoreCase(true).
				WithRowsPerResultset(true).
				WithS3StageVpceDnsName(true).
				WithStatementQueuedTimeoutInSeconds(true).
				WithStatementTimeoutInSeconds(true).
				WithStrictJsonOutput(true).
				WithTimestampDayIsAlways24H(true).
				WithTimestampInputFormat(true).
				WithTimestampLtzOutputFormat(true).
				WithTimestampNtzOutputFormat(true).
				WithTimestampOutputFormat(true).
				WithTimestampTypeMapping(true).
				WithTimestampTzOutputFormat(true).
				WithTimezone(true).
				WithTimeInputFormat(true).
				WithTimeOutputFormat(true).
				WithTraceLevel(true).
				WithTransactionAbortOnError(true).
				WithTransactionDefaultIsolationLevel(true).
				WithTwoDigitCenturyStart(true).
				WithUnsupportedDdlAction(true).
				WithUseCachedResult(true).
				WithWeekOfYearPolicy(true).
				WithWeekStart(true).
				WithWarehouse(true).
				WithSchedule(true).
				WithConfig(true).
				WithAllowOverlappingExecution(true).
				WithUserTaskTimeoutMs(true).
				WithSuspendTaskAfterNumFailures(true).
				WithComment(true).
				WithTaskAutoRetryAttempts(true).
				WithUserTaskMinimumTriggerIntervalInSeconds(true),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoErrorIntegration().
				HasSchedule("").
				HasConfig("").
				HasAllowOverlappingExecution(false).
				HasComment(""),
		)
		assertThatObject(t, objectparametersassert.TaskParameters(t, task.ID()).HasAllDefaults())
	})

	// AUTOCOMMIT only accepts TRUE on TASK, narrower than Snowflake's generic session-parameter rules.
	// (SEARCH_PATH's rejection at the TASK level is no longer exercisable here: its ParameterDef has no
	// ParameterLevelTask, so the catalog never generates a TaskSetRequest.WithSearchPath in the first
	// place - the invariant is now enforced at compile time instead of by a runtime SQL error.)
	t.Run("alter task: rejected parameters", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().WithAutocommit(false),
		))
		require.Error(t, err)

		assertThatObject(t, objectparametersassert.TaskParameters(t, task.ID()).HasAllDefaults())
	})

	t.Run("alter task: set and unset execute as user", func(t *testing.T) {
		user1 := createUserForTaskImpersonation(t)
		user2 := createUserForTaskImpersonation(t)
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)
		assertThatObject(t, objectassert.Task(t, task.ID()).HasNoExecuteAsUser())

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSetExecuteAsUser(user1.ID()))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Task(t, task.ID()).HasExecuteAsUser(user1.ID()))

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSetExecuteAsUser(user2.ID()))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Task(t, task.ID()).HasExecuteAsUser(user2.ID()))

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnsetExecuteAsUser(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Task(t, task.ID()).HasNoExecuteAsUser())
	})

	t.Run("alter task: resume and suspend", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(taskCleanup)

		assert.Equal(t, sdk.TaskStateSuspended, task.State)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithResume(true))
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateStarted, alteredTask.State)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(true))
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, sdk.TaskStateSuspended, alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(rootTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(taskCleanup)

		assert.Contains(t, task.TaskRelations.Predecessors, rootTask.ID())

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithRemoveAfter([]sdk.SchemaObjectIdentifier{rootTask.ID()}))
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, task.ID())

		require.NoError(t, err)
		assert.Empty(t, task.TaskRelations.Predecessors)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithAddAfter([]sdk.SchemaObjectIdentifier{rootTask.ID()}))
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, task.ID())

		require.NoError(t, err)
		assert.Contains(t, task.TaskRelations.Predecessors, rootTask.ID())
	})

	t.Run("alter task: set and unset final task", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		finalTask, finalTaskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(finalTaskCleanup)

		assertThatObject(
			t, objectassert.TaskFromObject(t, task).
				HasTaskRelations(sdk.TaskRelations{
					FinalizerTask:     nil,
					FinalizedRootTask: nil,
				}),
		)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(finalTask.ID()).WithSetFinalize(task.ID()))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasTaskRelations(sdk.TaskRelations{
					FinalizerTask: sdk.Pointer(finalTask.ID()),
				}),
		)
		assertThatObject(
			t, objectassert.Task(t, finalTask.ID()).
				HasTaskRelations(sdk.TaskRelations{
					FinalizedRootTask: sdk.Pointer(task.ID()),
				}),
		)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(finalTask.ID()).WithUnsetFinalize(true))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasTaskRelations(sdk.TaskRelations{
					FinalizerTask: nil,
				}),
		)
		assertThatObject(
			t, objectassert.Task(t, finalTask.ID()).
				HasTaskRelations(sdk.TaskRelations{
					FinalizedRootTask: nil,
				}),
		)
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		newSql := "SELECT CURRENT_DATE"
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithModifyAs(newSql))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Task(t, task.ID()).HasDefinition(newSql))

		newWhen := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithModifyWhen(newWhen))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Task(t, task.ID()).HasCondition(newWhen))

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithRemoveWhen(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Task(t, task.ID()).HasCondition(""))
	})

	t.Run("show task: default", func(t *testing.T) {
		task1, task1Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task1Cleanup)

		task2, task2Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task2Cleanup)

		returnedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}))
		require.NoError(t, err)

		require.Len(t, returnedTasks, 2)
		assert.Contains(t, returnedTasks, *task1)
		assert.Contains(t, returnedTasks, *task2)
	})

	t.Run("show task: terse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).WithSchedule("10 MINUTE"))
		t.Cleanup(taskCleanup)

		returnedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).WithTerse(true))
		require.NoError(t, err)

		require.Len(t, returnedTasks, 1)
		assertTaskTerse(t, &returnedTasks[0], task.ID(), "10 MINUTE")
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1, task1Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task1Cleanup)

		task2, task2Cleanup := testClientHelper().Task.Create(t)
		t.Cleanup(task2Cleanup)

		returnedTasks, err := client.Tasks.Show(ctx, sdk.NewShowTaskRequest().
			WithLike(sdk.Like{Pattern: &task1.Name}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(5)}))

		require.NoError(t, err)
		assert.Len(t, returnedTasks, 1)
		assert.Contains(t, returnedTasks, *task1)
		assert.NotContains(t, returnedTasks, *task2)
	})

	t.Run("describe task: default", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		returnedTask, err := client.Tasks.Describe(ctx, task.ID())
		require.NoError(t, err)

		assertTask(t, returnedTask, task.ID(), sdk.Pointer(testClientHelper().Ids.WarehouseId()))
	})

	t.Run("execute task: default", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		err := client.Tasks.Execute(ctx, sdk.NewExecuteTaskRequest(task.ID()))
		require.NoError(t, err)
	})

	t.Run("execute task: retry last after successful last task", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		_, subTaskCleanup := testClientHelper().Task.CreateWithAfter(t, task.ID())
		t.Cleanup(subTaskCleanup)

		err := client.Tasks.Execute(ctx, sdk.NewExecuteTaskRequest(task.ID()))
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			err = client.Tasks.Execute(ctx, sdk.NewExecuteTaskRequest(task.ID()).WithRetryLast(true))
			return strings.Contains(err.Error(), fmt.Sprintf("Cannot perform retry: no suitable run of graph with root task %s to retry.", task.ID().Name())) ||
				strings.Contains(err.Error(), fmt.Sprintf("graph with root task %s had no failures.", task.ID().Name()))
		}, 2*time.Second, time.Millisecond*300)
	})

	t.Run("execute task: retry last after failed last task", func(t *testing.T) {
		task, taskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(taskCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, subTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, "select * from not_existing_table"))
		t.Cleanup(subTaskCleanup)

		err := client.Tasks.Execute(ctx, sdk.NewExecuteTaskRequest(task.ID()))
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			err := client.Tasks.Execute(ctx, sdk.NewExecuteTaskRequest(task.ID()).WithRetryLast(true))
			return err != nil
		}, 2*time.Second, time.Millisecond*500)
	})

	t.Run("temporarily suspend root tasks", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(taskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, task.ID(), task.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)
	})

	// Tested graph
	// root1
	//      \
	//       t1
	//      /
	// root2
	// Because graph validation occurs only after resuming the root task, we assume that Snowflake will throw
	// validation error with given graph configuration.
	t.Run("resume root tasks within a graph containing more than one root task", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		secondRootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secondRootTask, secondRootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(secondRootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(secondRootTaskCleanup)

		_, cleanupTask := testClientHelper().Task.CreateWithAfter(t, rootTask.ID(), secondRootTask.ID())
		t.Cleanup(cleanupTask)

		require.ErrorContains(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)), "The graph has more than one root task (one without predecessors)")
		require.ErrorContains(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(secondRootTask.ID()).WithResume(true)), "The graph has more than one root task (one without predecessors)")
	})

	t.Run("suspend root tasks temporarily with three sequentially connected tasks - last in DAG", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		middleTask, middleTaskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(middleTaskCleanup)

		task, taskCleanup := testClientHelper().Task.CreateWithAfter(t, middleTask.ID())
		t.Cleanup(taskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(middleTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(middleTask.ID()).WithSuspend(true)))
		})

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, task.ID(), task.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)
		require.Contains(t, tasksToResume, rootTask.ID())

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		middleTaskStatus, err := client.Tasks.ShowByID(ctx, middleTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, middleTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)

		middleTaskStatus, err = client.Tasks.ShowByID(ctx, middleTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, middleTaskStatus.State)
	})

	t.Run("suspend root tasks temporarily with three sequentially connected tasks - middle in DAG", func(t *testing.T) {
		rootTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		rootTask, rootTaskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, sql).WithSchedule("60 MINUTES"))
		t.Cleanup(rootTaskCleanup)

		middleTask, middleTaskCleanup := testClientHelper().Task.CreateWithAfter(t, rootTask.ID())
		t.Cleanup(middleTaskCleanup)

		childTask, childTaskCleanup := testClientHelper().Task.CreateWithAfter(t, middleTask.ID())
		t.Cleanup(childTaskCleanup)

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(childTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(childTask.ID()).WithSuspend(true)))
		})

		require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true)))
		t.Cleanup(func() {
			require.NoError(t, client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)))
		})

		tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, middleTask.ID(), middleTask.ID())
		require.NoError(t, err)
		require.NotEmpty(t, tasksToResume)

		rootTaskStatus, err := client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateSuspended, rootTaskStatus.State)

		childTaskStatus, err := client.Tasks.ShowByID(ctx, childTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, childTaskStatus.State)

		require.NoError(t, client.Tasks.ResumeTasks(ctx, tasksToResume))

		rootTaskStatus, err = client.Tasks.ShowByID(ctx, rootTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, rootTaskStatus.State)

		childTaskStatus, err = client.Tasks.ShowByID(ctx, childTask.ID())
		require.NoError(t, err)
		require.Equal(t, sdk.TaskStateStarted, childTaskStatus.State)
	})

	t.Run("alter task: set and unset serverless task parameters", func(t *testing.T) {
		// Create a serverless task (with USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE)
		task, taskCleanup := testClientHelper().Task.CreateServerless(t)
		t.Cleanup(taskCleanup)

		// Set the serverless task parameters
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				WithTargetCompletionInterval("15 MINUTES").
				WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
				WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge),
		))
		require.NoError(t, err)

		// Verify parameters were set
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
				HasServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge),
		)
		// target_completion_interval is returned by SHOW command, not SHOW PARAMETERS
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoWarehouse().
				HasTargetCompletionInterval(sdk.TaskTargetCompletionInterval{Minutes: sdk.Pointer(15)}),
		)

		// Unset the serverless task parameters
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(
			*sdk.NewTaskUnsetRequest().
				WithTargetCompletionInterval(true).
				WithServerlessTaskMinStatementSize(true).
				WithServerlessTaskMaxStatementSize(true),
		))
		require.NoError(t, err)

		// Verify parameters were unset (reverted to defaults)
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasDefaultServerlessTaskMinStatementSizeValue().
				HasDefaultServerlessTaskMaxStatementSizeValue(),
		)
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoTargetCompletionInterval(),
		)
	})

	t.Run("serverless task parameters on non-serverless task", func(t *testing.T) {
		warehouseId := testClientHelper().Ids.WarehouseId()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(warehouseId)))
		t.Cleanup(taskCleanup)

		assertThatObject(t, objectassert.Task(t, task.ID()).HasWarehouse(warehouseId))

		// Setting serverless task parameters on a non-serverless task fails, even if we set USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
				WithTargetCompletionInterval("10 MINUTES").
				WithServerlessTaskMinStatementSize(sdk.WarehouseSizeXSmall).
				WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeSmall),
		))
		require.ErrorContains(t, err, "Cannot set USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE on a non-serverless task.")

		// First, unset the warehouse to make the task serverless.
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(
			*sdk.NewTaskUnsetRequest().
				WithWarehouse(true),
		))
		require.NoError(t, err)

		// Assert default attributes.
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasDefaultUserTaskManagedInitialWarehouseSizeValue().
				HasDefaultServerlessTaskMinStatementSizeValue().
				HasDefaultServerlessTaskMaxStatementSizeValue(),
		)
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoWarehouse().
				HasNoTargetCompletionInterval(),
		)

		// Set serverless task parameters.
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				WithTargetCompletionInterval("10 MINUTES").
				WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeSmall).
				WithServerlessTaskMinStatementSize(sdk.WarehouseSizeXSmall).
				WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeMedium),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeSmall).
				HasServerlessTaskMinStatementSize(sdk.WarehouseSizeXSmall).
				HasServerlessTaskMaxStatementSize(sdk.WarehouseSizeMedium),
		)
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoWarehouse().
				HasTargetCompletionInterval(sdk.TaskTargetCompletionInterval{Minutes: sdk.Pointer(10)}),
		)

		// Convert back to warehouse-based task - fails
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				WithWarehouse(warehouseId),
		))
		require.ErrorContains(t, err, "091857 (42601): TARGET_COMPLETION_INTERVAL is only allowed for serverless Tasks.")

		// Unset the serverless parameters on the warehouse-based task
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(
			*sdk.NewTaskUnsetRequest().
				WithTargetCompletionInterval(true).
				WithUserTaskManagedInitialWarehouseSize(true).
				WithServerlessTaskMinStatementSize(true).
				WithServerlessTaskMaxStatementSize(true),
		))
		require.NoError(t, err)

		// Verify parameters were unset
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasDefaultUserTaskManagedInitialWarehouseSizeValue().
				HasDefaultServerlessTaskMinStatementSizeValue().
				HasDefaultServerlessTaskMaxStatementSizeValue(),
		)
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasNoWarehouse().
				HasNoTargetCompletionInterval(),
		)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(
			*sdk.NewTaskSetRequest().
				WithWarehouse(warehouseId),
		))
		require.NoError(t, err)

		// Verify task is back to warehouse-based and parameters are still set
		assertThatObject(
			t, objectparametersassert.TaskParameters(t, task.ID()).
				HasDefaultUserTaskManagedInitialWarehouseSizeValue().
				HasDefaultServerlessTaskMinStatementSizeValue().
				HasDefaultServerlessTaskMaxStatementSizeValue(),
		)
		assertThatObject(
			t, objectassert.Task(t, task.ID()).
				HasWarehouse(warehouseId).
				HasNoTargetCompletionInterval(),
		)
	})

	// TODO(SNOW-1277135): Create more tests with different sets of roots/children and see if the current implementation
	// acts correctly in certain situations/edge cases.
}

func TestInt_TasksShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, t1Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id1, "SELECT CURRENT_TIMESTAMP"))
		_, t2Cleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id2, "SELECT CURRENT_TIMESTAMP"))
		t.Cleanup(t1Cleanup)
		t.Cleanup(t2Cleanup)

		e1, err := client.Tasks.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Tasks.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
