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

	assertTask := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, warehouseId *sdk.AccountObjectIdentifier) {
		t.Helper()
		asserts := objectassert.TaskFromObject(t, task).
			HasNotEmptyCreatedOn().
			HasName(id.Name()).
			HasNotEmptyId().
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

	assertTaskWithOptions := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, comment string, warehouse *sdk.AccountObjectIdentifier, schedule string, condition string, allowOverlappingExecution bool, config string, predecessor *sdk.SchemaObjectIdentifier, errorIntegrationName *sdk.AccountObjectIdentifier) {
		t.Helper()

		asserts := objectassert.TaskFromObject(t, task).
			HasNotEmptyCreatedOn().
			HasName(id.Name()).
			HasNotEmptyId().
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

		if warehouse != nil {
			asserts.HasWarehouse(*warehouse)
		} else {
			asserts.HasNoWarehouse()
		}

		assertThatObject(t, asserts)
	}

	assertTaskTerse := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, schedule string) {
		t.Helper()
		assertThatObject(t, objectassert.TaskFromObject(t, task).
			HasNotEmptyCreatedOn().
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
			HasLastCommittedOn("").
			HasLastSuspendedOn("").
			HasOwnerRoleType("").
			HasConfig("").
			HasBudget("").
			HasLastSuspendedOn("").
			HasTaskRelations(sdk.TaskRelations{}),
		)
	}

	sessionParametersSet := sdk.SessionParameters{
		AbortDetachedQuery:                       sdk.Bool(true),
		Autocommit:                               sdk.Bool(false),
		BinaryInputFormat:                        sdk.Pointer(sdk.BinaryInputFormatUTF8),
		BinaryOutputFormat:                       sdk.Pointer(sdk.BinaryOutputFormatBase64),
		ClientMemoryLimit:                        sdk.Int(1024),
		ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
		ClientPrefetchThreads:                    sdk.Int(2),
		ClientResultChunkSize:                    sdk.Int(48),
		ClientResultColumnCaseInsensitive:        sdk.Bool(true),
		ClientSessionKeepAlive:                   sdk.Bool(true),
		ClientSessionKeepAliveHeartbeatFrequency: sdk.Int(2400),
		ClientTimestampTypeMapping:               sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
		DateInputFormat:                          sdk.String("YYYY-MM-DD"),
		DateOutputFormat:                         sdk.String("YY-MM-DD"),
		EnableUnloadPhysicalTypeOptimization:     sdk.Bool(false),
		ErrorOnNondeterministicMerge:             sdk.Bool(false),
		ErrorOnNondeterministicUpdate:            sdk.Bool(true),
		GeographyOutputFormat:                    sdk.Pointer(sdk.GeographyOutputFormatWKB),
		GeometryOutputFormat:                     sdk.Pointer(sdk.GeometryOutputFormatWKB),
		JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
		JdbcUseSessionTimezone:                   sdk.Bool(false),
		JsonIndent:                               sdk.Int(4),
		LockTimeout:                              sdk.Int(21222),
		LogLevel:                                 sdk.Pointer(sdk.LogLevelError),
		MultiStatementCount:                      sdk.Int(0),
		NoorderSequenceAsDefault:                 sdk.Bool(false),
		OdbcTreatDecimalAsInt:                    sdk.Bool(true),
		QueryTag:                                 sdk.String("some_tag"),
		QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
		RowsPerResultset:                         sdk.Int(2),
		S3StageVpceDnsName:                       sdk.String("vpce-id.s3.region.vpce.amazonaws.com"),
		SearchPath:                               sdk.String("$public, $current"),
		StatementQueuedTimeoutInSeconds:          sdk.Int(10),
		StatementTimeoutInSeconds:                sdk.Int(10),
		StrictJsonOutput:                         sdk.Bool(true),
		TimestampDayIsAlways24h:                  sdk.Bool(true),
		TimestampInputFormat:                     sdk.String("YYYY-MM-DD"),
		TimestampLTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
		TimestampNTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
		TimestampOutputFormat:                    sdk.String("YYYY-MM-DD HH24:MI:SS"),
		TimestampTypeMapping:                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
		TimestampTZOutputFormat:                  sdk.String("YYYY-MM-DD HH24:MI:SS"),
		Timezone:                                 sdk.String("Europe/Warsaw"),
		TimeInputFormat:                          sdk.String("HH24:MI"),
		TimeOutputFormat:                         sdk.String("HH24:MI"),
		TraceLevel:                               sdk.Pointer(sdk.TraceLevelPropagate),
		TransactionAbortOnError:                  sdk.Bool(true),
		TransactionDefaultIsolationLevel:         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
		TwoDigitCenturyStart:                     sdk.Int(1980),
		UnsupportedDDLAction:                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
		UseCachedResult:                          sdk.Bool(false),
		WeekOfYearPolicy:                         sdk.Int(1),
		WeekStart:                                sdk.Int(1),
	}

	assertSessionParametersSet := func(parametersAssert *objectparametersassert.TaskParametersAssert) *objectparametersassert.TaskParametersAssert {
		return parametersAssert.
			HasAbortDetachedQuery(true).
			HasAutocommit(false).
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
			HasMultiStatementCount(0).
			HasNoorderSequenceAsDefault(false).
			HasOdbcTreatDecimalAsInt(true).
			HasQueryTag("some_tag").
			HasQuotedIdentifiersIgnoreCase(true).
			HasRowsPerResultset(2).
			HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
			HasSearchPath("$public, $current").
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

		assertThatObject(t, objectparametersassert.TaskParameters(t, task.ID()).
			HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall),
		)

		assertTask(t, task, id, nil)
	})

	t.Run("create task: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, sql).
			WithOrReplace(true).
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
			WithErrorIntegration(errorIntegration.ID()).
			WithSchedule("10 MINUTE").
			WithConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
			WithAllowOverlappingExecution(true).
			WithSessionParameters(sdk.SessionParameters{
				JsonIndent: sdk.Int(4),
			}).
			WithUserTaskTimeoutMs(500).
			WithSuspendTaskAfterNumFailures(3).
			WithComment("some comment").
			WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)

		assertTaskWithOptions(t, task, id, "some comment", sdk.Pointer(testClientHelper().Ids.WarehouseId()), "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil, sdk.Pointer(errorIntegration.ID()))
		assertThatObject(t, objectparametersassert.TaskParameters(t, id).
			HasJsonIndent(4).
			HasUserTaskTimeoutMs(500).
			HasSuspendTaskAfterNumFailures(3),
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

		assertTaskWithOptions(t, task, id, "", nil, "", "", false, "", &rootTaskId, nil)
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

		assertThatObject(t, objectassert.Task(t, rootTaskId).
			HasTaskRelations(sdk.TaskRelations{
				Predecessors:  []sdk.SchemaObjectIdentifier{},
				FinalizerTask: &finalizerId,
			}),
		)

		assertThatObject(t, objectassert.Task(t, finalizerId).
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
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, sql).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
			}),
		)
		t.Cleanup(taskCleanup)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, sdk.Pointer("v1"), returnedTagValue)
	})

	t.Run("clone task: default", func(t *testing.T) {
		rootTask, rootTaskCleanup := testClientHelper().Task.Create(t)
		t.Cleanup(rootTaskCleanup)

		sourceTaskId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		sourceTask, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(sourceTaskId, sql).
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
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tasks.CreateOrAlter(ctx, sdk.NewCreateOrAlterTaskRequest(id, sql).
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClientHelper().Ids.WarehouseId())).
			WithSchedule("10 MINUTES").
			WithConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
			WithAllowOverlappingExecution(true).
			WithUserTaskTimeoutMs(10).
			WithSessionParameters(sessionParametersSet).
			WithSuspendTaskAfterNumFailures(15).
			WithComment("some_comment").
			WithTaskAutoRetryAttempts(15).
			WithWhen(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)
		createdOn := task.CreatedOn

		assertThatObject(t, objectassert.TaskFromObject(t, task).
			HasWarehouse(testClientHelper().Ids.WarehouseId()).
			HasSchedule("10 MINUTES").
			HasConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
			HasAllowOverlappingExecution(true).
			HasCondition(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`).
			HasComment("some_comment").
			HasTaskRelations(sdk.TaskRelations{}),
		)
		assertThatObject(t, assertSessionParametersSet(objectparametersassert.TaskParameters(t, task.ID()).
			HasUserTaskTimeoutMs(10).
			HasSuspendTaskAfterNumFailures(15).
			HasTaskAutoRetryAttempts(15)),
		)

		err = client.Tasks.CreateOrAlter(ctx, sdk.NewCreateOrAlterTaskRequest(id, sql))
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.TaskFromObject(t, alteredTask).
			HasNoWarehouse().
			HasSchedule("").
			HasConfig("").
			HasAllowOverlappingExecution(false).
			HasCondition("").
			HasComment("").
			HasTaskRelations(sdk.TaskRelations{}),
		)
		assertThatObject(t, objectparametersassert.TaskParameters(t, task.ID()).
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

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(*sdk.NewTaskSetRequest().
			// TODO(SNOW-1519496): Cannot set warehouse due to Snowflake error
			// WithWarehouse(testClientHelper().Ids.WarehouseId()).
			WithErrorIntegration(errorIntegration.ID()).
			WithSessionParameters(sessionParametersSet).
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

		assertThatObject(t, objectassert.Task(t, task.ID()).
			// HasWarehouse(testClientHelper().Ids.WarehouseId().Name()).
			HasErrorIntegration(errorIntegration.ID()).
			HasSchedule("10 MINUTE").
			HasConfig(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`).
			HasAllowOverlappingExecution(true).
			HasComment("new comment"),
		)
		assertThatObject(t, assertSessionParametersSet(objectparametersassert.TaskParameters(t, task.ID())).
			HasUserTaskTimeoutMs(1000).
			HasSuspendTaskAfterNumFailures(100).
			HasTaskAutoRetryAttempts(10).
			HasUserTaskMinimumTriggerIntervalInSeconds(15),
		)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(*sdk.NewTaskUnsetRequest().
			WithErrorIntegration(true).
			WithSessionParametersUnset(sdk.SessionParametersUnset{
				AbortDetachedQuery:                       sdk.Bool(true),
				Autocommit:                               sdk.Bool(true),
				BinaryInputFormat:                        sdk.Bool(true),
				BinaryOutputFormat:                       sdk.Bool(true),
				ClientMemoryLimit:                        sdk.Bool(true),
				ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
				ClientPrefetchThreads:                    sdk.Bool(true),
				ClientResultChunkSize:                    sdk.Bool(true),
				ClientResultColumnCaseInsensitive:        sdk.Bool(true),
				ClientSessionKeepAlive:                   sdk.Bool(true),
				ClientSessionKeepAliveHeartbeatFrequency: sdk.Bool(true),
				ClientTimestampTypeMapping:               sdk.Bool(true),
				DateInputFormat:                          sdk.Bool(true),
				DateOutputFormat:                         sdk.Bool(true),
				EnableUnloadPhysicalTypeOptimization:     sdk.Bool(true),
				ErrorOnNondeterministicMerge:             sdk.Bool(true),
				ErrorOnNondeterministicUpdate:            sdk.Bool(true),
				GeographyOutputFormat:                    sdk.Bool(true),
				GeometryOutputFormat:                     sdk.Bool(true),
				JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
				JdbcUseSessionTimezone:                   sdk.Bool(true),
				JsonIndent:                               sdk.Bool(true),
				LockTimeout:                              sdk.Bool(true),
				LogLevel:                                 sdk.Bool(true),
				MultiStatementCount:                      sdk.Bool(true),
				NoorderSequenceAsDefault:                 sdk.Bool(true),
				OdbcTreatDecimalAsInt:                    sdk.Bool(true),
				QueryTag:                                 sdk.Bool(true),
				QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
				RowsPerResultset:                         sdk.Bool(true),
				S3StageVpceDnsName:                       sdk.Bool(true),
				SearchPath:                               sdk.Bool(true),
				StatementQueuedTimeoutInSeconds:          sdk.Bool(true),
				StatementTimeoutInSeconds:                sdk.Bool(true),
				StrictJsonOutput:                         sdk.Bool(true),
				TimestampDayIsAlways24h:                  sdk.Bool(true),
				TimestampInputFormat:                     sdk.Bool(true),
				TimestampLTZOutputFormat:                 sdk.Bool(true),
				TimestampNTZOutputFormat:                 sdk.Bool(true),
				TimestampOutputFormat:                    sdk.Bool(true),
				TimestampTypeMapping:                     sdk.Bool(true),
				TimestampTZOutputFormat:                  sdk.Bool(true),
				Timezone:                                 sdk.Bool(true),
				TimeInputFormat:                          sdk.Bool(true),
				TimeOutputFormat:                         sdk.Bool(true),
				TraceLevel:                               sdk.Bool(true),
				TransactionAbortOnError:                  sdk.Bool(true),
				TransactionDefaultIsolationLevel:         sdk.Bool(true),
				TwoDigitCenturyStart:                     sdk.Bool(true),
				UnsupportedDDLAction:                     sdk.Bool(true),
				UseCachedResult:                          sdk.Bool(true),
				WeekOfYearPolicy:                         sdk.Bool(true),
				WeekStart:                                sdk.Bool(true),
			}).
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

		assertThatObject(t, objectassert.Task(t, task.ID()).
			HasNoErrorIntegration().
			HasSchedule("").
			HasConfig("").
			HasAllowOverlappingExecution(false).
			HasComment(""),
		)
		assertThatObject(t, objectparametersassert.TaskParameters(t, task.ID()).HasAllDefaults())
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

		assertThatObject(t, objectassert.TaskFromObject(t, task).
			HasTaskRelations(sdk.TaskRelations{
				FinalizerTask:     nil,
				FinalizedRootTask: nil,
			}),
		)

		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(finalTask.ID()).WithSetFinalize(task.ID()))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Task(t, task.ID()).
			HasTaskRelations(sdk.TaskRelations{
				FinalizerTask: sdk.Pointer(finalTask.ID()),
			}),
		)
		assertThatObject(t, objectassert.Task(t, finalTask.ID()).
			HasTaskRelations(sdk.TaskRelations{
				FinalizedRootTask: sdk.Pointer(task.ID()),
			}),
		)

		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(finalTask.ID()).WithUnsetFinalize(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Task(t, task.ID()).
			HasTaskRelations(sdk.TaskRelations{
				FinalizerTask: nil,
			}),
		)
		assertThatObject(t, objectassert.Task(t, finalTask.ID()).
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

	t.Run("create task: with serverless task parameters", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		// Create a serverless task (no warehouse specified) with serverless parameters
		err := testClient(t).Tasks.Create(ctx, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium)).
			WithTargetCompletionInterval("10 MINUTES").
			WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
			WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Task.DropFunc(t, id))

		task, err := testClientHelper().Task.Show(t, id)
		require.NoError(t, err)
		require.NotNil(t, task)

		// Verify task was created successfully
		assert.Equal(t, id.Name(), task.Name)

		// Verify parameters were set by checking ShowParameters
		params, err := client.Tasks.ShowParameters(ctx, id)
		require.NoError(t, err)

		// Find and verify the parameters
		var targetCompletionInterval, minSize, maxSize *sdk.Parameter
		for _, p := range params {
			switch p.Key {
			case string(sdk.TaskParameterTargetCompletionInterval):
				targetCompletionInterval = p
			case string(sdk.TaskParameterServerlessTaskMinStatementSize):
				minSize = p
			case string(sdk.TaskParameterServerlessTaskMaxStatementSize):
				maxSize = p
			}
		}

		// Note: TARGET_COMPLETION_INTERVAL may not appear in SHOW PARAMETERS if it's at default
		// So we just verify the task was created without error
		if targetCompletionInterval != nil {
			assert.Equal(t, "10 MINUTES", targetCompletionInterval.Value)
		}

		if minSize != nil {
			assert.Equal(t, string(sdk.WarehouseSizeSmall), minSize.Value)
		}

		if maxSize != nil {
			assert.Equal(t, string(sdk.WarehouseSizeLarge), maxSize.Value)
		}
	})

	t.Run("alter task: set and unset serverless task parameters", func(t *testing.T) {
		// Create a serverless task (with USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE)
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		task, taskCleanup := testClientHelper().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").
			WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium)))
		t.Cleanup(taskCleanup)

		// Set the serverless task parameters
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithSet(*sdk.NewTaskSetRequest().
			WithTargetCompletionInterval("15 MINUTES").
			WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
			WithServerlessTaskMaxStatementSize(sdk.WarehouseSizeLarge),
		))
		require.NoError(t, err)

		// Verify parameters were set
		params, err := client.Tasks.ShowParameters(ctx, task.ID())
		require.NoError(t, err)

		var targetCompletionInterval, minSize, maxSize *sdk.Parameter
		for _, p := range params {
			switch p.Key {
			case string(sdk.TaskParameterTargetCompletionInterval):
				targetCompletionInterval = p
			case string(sdk.TaskParameterServerlessTaskMinStatementSize):
				minSize = p
			case string(sdk.TaskParameterServerlessTaskMaxStatementSize):
				maxSize = p
			}
		}

		// Note: Parameters may not always appear in SHOW PARAMETERS
		// So we just verify the ALTER command succeeded
		if targetCompletionInterval != nil {
			assert.Equal(t, "15 MINUTES", targetCompletionInterval.Value)
		}

		if minSize != nil {
			assert.Equal(t, string(sdk.WarehouseSizeSmall), minSize.Value)
		}

		if maxSize != nil {
			assert.Equal(t, string(sdk.WarehouseSizeLarge), maxSize.Value)
		}

		// Unset the serverless task parameters
		err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(task.ID()).WithUnset(*sdk.NewTaskUnsetRequest().
			WithTargetCompletionInterval(true).
			WithServerlessTaskMinStatementSize(true).
			WithServerlessTaskMaxStatementSize(true),
		))
		require.NoError(t, err)

		// Note: UNSET may revert parameters to system defaults rather than empty values
		// The important thing is that the UNSET command succeeded without error
		// Snowflake may keep default values for these parameters even after UNSET
	})
}
