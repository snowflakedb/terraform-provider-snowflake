//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Pipes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	copyStatement := testClientHelper().Pipe.CopyStatement(t, table, stage)

	assertPipe := func(t *testing.T, pipeDetails *sdk.Pipe, expectedName string, expectedComment string) {
		t.Helper()
		assertThatObject(t, objectassert.PipeFromObject(t, pipeDetails).
			HasCreatedOnNotEmpty().
			HasName(expectedName).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasDefinition(copyStatement).
			HasOwner(testClientHelper().Context.CurrentRole(t).Name()).
			HasNotificationChannel("").
			HasComment(expectedComment).
			HasIntegration("").
			HasPattern("").
			HasErrorIntegration("").
			HasOwnerRoleType("ROLE").
			HasInvalidReason(""))
	}

	t.Run("create: schema name contains special characters", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithName(t, testClientHelper().Ids.RandomDatabaseObjectIdentifierWithPrefix("tcK1>AJ+").Name())
		t.Cleanup(schemaCleanup)
		specialTable, specialTableCleanup := testClientHelper().Table.CreateInSchema(t, schema.ID())
		t.Cleanup(specialTableCleanup)

		t.Run("fully qualified copy statement", func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

			err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, testClientHelper().Pipe.CopyStatement(t, specialTable, stage)))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))
		})

		// The stage is in the test schema. A name-only reference is resolved in the pipe's schema.
		t.Run("stage referenced only by name", func(t *testing.T) {
			statement := fmt.Sprintf("COPY INTO %s\nFROM @\"%s\"", specialTable.ID().FullyQualifiedName(), stage.Name)

			err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()), statement))
			require.ErrorContains(t, err, "object does not exist or not authorized")
		})
	})

	t.Run("create: with comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement).
			WithIfNotExists(true).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		pipe, err := client.Pipes.ShowByID(ctx, id)
		require.NoError(t, err)
		assertPipe(t, pipe, id.Name(), comment)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		pipe, err := client.Pipes.ShowByID(ctx, id)
		require.NoError(t, err)
		assertPipe(t, pipe, id.Name(), "")
	})

	t.Run("create: complete error notifications", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		integration, cleanupIntegration := testClientHelper().NotificationIntegration.CreateWithAmazonPush(t)
		t.Cleanup(cleanupIntegration)

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement).
			WithAutoIngest(false).
			WithComment("complete").
			WithErrorIntegration(integration.Name))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		assertThatObject(t, objectassert.Pipe(t, id).
			HasName(id.Name()).
			HasDefinition(copyStatement).
			HasComment("complete").
			HasErrorIntegration(integration.Name).
			HasNotificationChannel(""))
		assertThatObject(t, objectassert.PipeDetails(t, id).
			HasName(id.Name()).
			HasDefinition(copyStatement).
			HasComment("complete").
			HasErrorIntegration(integration.Name).
			HasNotificationChannel(""))
	})

	t.Run("create: if not exists preserves existing pipe", func(t *testing.T) {
		pipe, cleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanup)

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(pipe.ID(), copyStatement).
			WithIfNotExists(true).
			WithComment("ignored"))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasDefinition(copyStatement).
			HasComment(""))
	})

	t.Run("create: or replace changes existing pipe", func(t *testing.T) {
		pipe, cleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanup)
		changed := copyStatement + "\nON_ERROR = CONTINUE"

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(pipe.ID(), changed).
			WithOrReplace(true).
			WithComment("replaced"))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasDefinition(changed).
			HasComment("replaced"))
		assertThatObject(t, objectassert.PipeDetails(t, pipe.ID()).
			HasDefinition(changed).
			HasComment("replaced"))
	})

	t.Run("create: duplicate without modifier fails", func(t *testing.T) {
		pipe, cleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanup)

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(pipe.ID(), copyStatement))
		require.Error(t, err)
	})

	t.Run("create: auto ingest external S3 stage", func(t *testing.T) {
		externalStage, cleanupStage := testClientHelper().Stage.CreateStageWithURL(t)
		t.Cleanup(cleanupStage)
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		statement := testClientHelper().Pipe.CopyStatement(t, table, externalStage)

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, statement).
			WithAutoIngest(true))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		assertThatObject(t, objectassert.Pipe(t, id).
			HasNotificationChannelNotEmpty().
			HasDefinition(statement))
		assertThatObject(t, objectassert.PipeDetails(t, id).
			HasNotificationChannelNotEmpty())
	})

	t.Run("create: auto ingest GCS requires an accessible subscription", func(t *testing.T) {
		location := "gcs://test-bucket/test-path/"
		storageIntegrationId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, cleanupStorageIntegration := testClientHelper().StorageIntegration.CreateWithRequest(t, storageIntegrationId,
			sdk.NewCreateStorageIntegrationRequest(storageIntegrationId, true, []sdk.StorageLocation{{Path: location}}).
				WithGCSStorageProviderParams(*sdk.NewGCSStorageParamsRequest()))
		t.Cleanup(cleanupStorageIntegration)

		stageId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Stages.CreateOnGCS(ctx, sdk.NewCreateOnGCSStageRequest(stageId, *sdk.NewExternalGCSStageParamsRequest(location).
			WithStorageIntegration(storageIntegrationId)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, stageId))

		externalStage, err := client.Stages.ShowByID(ctx, stageId)
		require.NoError(t, err)

		integration, cleanupIntegration := testClientHelper().NotificationIntegration.CreateWithGcpPubSub(t)
		t.Cleanup(cleanupIntegration)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		// The notification helper uses an unprovisioned Pub/Sub subscription.
		err = client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, testClientHelper().Pipe.CopyStatement(t, table, externalStage)).
			WithAutoIngest(true).
			WithIntegration(integration.Name))
		require.ErrorContains(t, err, "Subscription not found")

		_, err = client.Pipes.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("create: SNS topic requires an accessible topic", func(t *testing.T) {
		externalStage, cleanupStage := testClientHelper().Stage.CreateStageWithURL(t)
		t.Cleanup(cleanupStage)
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		// A syntactically valid but unprovisioned SNS topic must surface the Snowflake error.
		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, testClientHelper().Pipe.CopyStatement(t, table, externalStage)).
			WithAutoIngest(true).
			WithAwsSnsTopic("arn:aws:sns:us-east-2:123456789012:MissingTopic"))
		require.ErrorContains(t, err, "Pipe Notifications bind failure")
	})

	t.Run("create: missing ingestion integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))

		err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement).
			WithAutoIngest(true).
			WithIntegration(NonExistingAccountObjectIdentifier.Name()))
		require.Error(t, err)
	})

	t.Run("show and describe", func(t *testing.T) {
		table2, table2Cleanup := testClientHelper().Table.Create(t)
		t.Cleanup(table2Cleanup)

		pipe1CopyStatement := testClientHelper().Pipe.CopyStatement(t, table, stage)
		pipe1, pipe1Cleanup := testClientHelper().Pipe.CreatePipe(t, pipe1CopyStatement)
		t.Cleanup(pipe1Cleanup)

		pipe2CopyStatement := testClientHelper().Pipe.CopyStatement(t, table2, stage)
		pipe2, pipe2Cleanup := testClientHelper().Pipe.CreatePipe(t, pipe2CopyStatement)
		t.Cleanup(pipe2Cleanup)

		t.Run("show: without options", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest())
			require.NoError(t, err)
			assert.Len(t, pipes, 2)
			assert.Contains(t, pipes, *pipe1)
			assert.Contains(t, pipes, *pipe2)
		})

		t.Run("show: in schema", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest().
				WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
			require.NoError(t, err)
			assert.Len(t, pipes, 2)
			assert.Contains(t, pipes, *pipe1)
			assert.Contains(t, pipes, *pipe2)
		})

		t.Run("show: like", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest().
				WithLike(sdk.Like{Pattern: sdk.String(pipe1.Name)}))
			require.NoError(t, err)
			assert.Len(t, pipes, 1)
			assert.Contains(t, pipes, *pipe1)
		})

		t.Run("show: non-existent pipe", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest().
				WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
			require.NoError(t, err)
			assert.Empty(t, pipes)
		})

		t.Run("show: in account", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest().
				WithIn(sdk.In{Account: sdk.Bool(true)}).
				WithLike(sdk.Like{Pattern: sdk.String(pipe1.Name)}))
			require.NoError(t, err)
			require.Len(t, pipes, 1)

			assertThatObject(t, objectassert.PipeFromObject(t, &pipes[0]).
				HasName(pipe1.Name).
				HasDatabaseName(pipe1.DatabaseName).
				HasSchemaName(pipe1.SchemaName).
				HasDefinition(pipe1.Definition))
		})

		t.Run("show: in database", func(t *testing.T) {
			pipes, err := client.Pipes.Show(ctx, sdk.NewShowPipeRequest().
				WithIn(sdk.In{Database: testClientHelper().Ids.DatabaseId()}).
				WithLike(sdk.Like{Pattern: sdk.String(pipe1.Name)}))
			require.NoError(t, err)
			require.Len(t, pipes, 1)

			assertThatObject(t, objectassert.PipeFromObject(t, &pipes[0]).
				HasName(pipe1.Name).
				HasDatabaseName(pipe1.DatabaseName).
				HasSchemaName(pipe1.SchemaName).
				HasDefinition(pipe1.Definition))
		})

		t.Run("describe: existing pipe", func(t *testing.T) {
			pipeDetails, err := client.Pipes.Describe(ctx, pipe1.ID())
			require.NoError(t, err)

			assertThatObject(t, objectassert.PipeDetailsFromObject(t, pipeDetails).
				HasCreatedOnNotEmpty().
				HasName(pipe1.Name).
				HasDatabaseName(pipe1.DatabaseName).
				HasSchemaName(pipe1.SchemaName).
				HasDefinition(pipe1CopyStatement).
				HasOwner(pipe1.Owner).
				HasComment("").
				HasNotificationChannel("").
				HasIntegration("").
				HasPattern("").
				HasErrorIntegration("").
				HasOwnerRoleType("ROLE").
				HasInvalidReason(""))
		})

		t.Run("describe: non-existing pipe", func(t *testing.T) {
			_, err := client.Pipes.Describe(ctx, NonExistingSchemaObjectIdentifier)
			assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	})

	t.Run("drop: existing", func(t *testing.T) {
		pipe, cleanupPipe := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanupPipe)

		err := client.Pipes.Drop(ctx, sdk.NewDropPipeRequest(pipe.ID()))
		require.NoError(t, err)

		_, err = client.Pipes.Describe(ctx, pipe.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: missing", func(t *testing.T) {
		err := client.Pipes.Drop(ctx, sdk.NewDropPipeRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on a missing pipe", func(t *testing.T) {
		err := client.Pipes.Drop(ctx, sdk.NewDropPipeRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop: if exists on an existing pipe", func(t *testing.T) {
		pipe, cleanupPipe := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanupPipe)

		err := client.Pipes.Drop(ctx, sdk.NewDropPipeRequest(pipe.ID()).
			WithIfExists(true))
		require.NoError(t, err)

		_, err = client.Pipes.ShowByID(ctx, pipe.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("alter: set and unset comment and execution paused", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithSet(*sdk.NewPipeSetRequest().
				WithComment("new comment").
				WithPipeExecutionPaused(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasComment("new comment"))

		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithUnset(*sdk.NewPipeUnsetRequest().
				WithComment(true).
				WithPipeExecutionPaused(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasComment(""))
	})

	t.Run("alter: refresh with prefix and modified after", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithRefresh(*sdk.NewPipeRefreshRequest().
				WithPrefix("/d1").
				WithModifiedAfter("2018-07-30T13:56:46-07:00")))
		require.NoError(t, err)
	})

	t.Run("alter: set, replace, and unset error integration", func(t *testing.T) {
		integration1, cleanupIntegration1 := testClientHelper().NotificationIntegration.CreateWithAmazonPush(t)
		t.Cleanup(cleanupIntegration1)
		integration2, cleanupIntegration2 := testClientHelper().NotificationIntegration.CreateWithAmazonPush(t)
		t.Cleanup(cleanupIntegration2)
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithIfExists(true).
			WithSet(*sdk.NewPipeSetRequest().
				WithErrorIntegration(integration1.Name).
				WithComment("changed")))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasErrorIntegration(integration1.Name).
			HasComment("changed"))
		assertThatObject(t, objectassert.PipeDetails(t, pipe.ID()).
			HasErrorIntegration(integration1.Name).
			HasComment("changed"))

		// Snowflake rejects replacing an assigned error integration in one SET.
		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithSet(*sdk.NewPipeSetRequest().
				WithErrorIntegration(integration2.Name)))
		require.ErrorContains(t, err, "Cannot reset ERROR_INTEGRATION property")

		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithUnset(*sdk.NewPipeUnsetRequest().
				WithErrorIntegration(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasErrorIntegration("").
			HasComment("changed"))

		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithUnset(*sdk.NewPipeUnsetRequest().
				WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.PipeDetails(t, pipe.ID()).
			HasErrorIntegration("").
			HasComment(""))

		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithSet(*sdk.NewPipeSetRequest().
				WithErrorIntegration(integration2.Name)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Pipe(t, pipe.ID()).
			HasErrorIntegration(integration2.Name))
	})

	t.Run("alter: set and unset execution paused", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithSet(*sdk.NewPipeSetRequest().
				WithPipeExecutionPaused(true)))
		require.NoError(t, err)

		status, err := client.SystemFunctions.PipeStatus(ctx, sdk.NewPipeStatusRequest(*sdk.NewPipeStatusArgumentsRequest(pipe.ID())))
		require.NoError(t, err)
		require.Equal(t, sdk.PipeExecutionStatePaused, status.ExecutionState)

		err = client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithUnset(*sdk.NewPipeUnsetRequest().
				WithPipeExecutionPaused(true)))
		require.NoError(t, err)

		status, err = client.SystemFunctions.PipeStatus(ctx, sdk.NewPipeStatusRequest(*sdk.NewPipeStatusArgumentsRequest(pipe.ID())))
		require.NoError(t, err)
		require.Equal(t, sdk.PipeExecutionStateRunning, status.ExecutionState)
	})

	t.Run("alter: missing object", func(t *testing.T) {
		req := sdk.NewAlterPipeRequest(NonExistingSchemaObjectIdentifier).
			WithSet(*sdk.NewPipeSetRequest().
				WithComment("changed"))

		err := client.Pipes.Alter(ctx, req)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		err = client.Pipes.Alter(ctx, req.WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("alter: refresh without options", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithRefresh(*sdk.NewPipeRefreshRequest()))
		require.NoError(t, err)
	})

	t.Run("alter: refresh with prefix", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithRefresh(*sdk.NewPipeRefreshRequest().
				WithPrefix("selected/")))
		require.NoError(t, err)
	})

	t.Run("alter: refresh with modified after", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(pipeCleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithRefresh(*sdk.NewPipeRefreshRequest().
				WithModifiedAfter(time.Now().Add(-time.Hour).UTC().Format(time.RFC3339))))
		require.NoError(t, err)
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())
		for _, id := range []sdk.SchemaObjectIdentifier{id1, id2} {
			err := client.Pipes.Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Pipe.DropPipeFunc(t, id))
		}

		e1, err := client.Pipes.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Pipes.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: missing pipe", func(t *testing.T) {
		_, err := client.Pipes.ShowByID(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("refresh: loads selected files", func(t *testing.T) {
		loadTable, loadTableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(loadTableCleanup)
		loadStage, loadStageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(loadStageCleanup)

		// A prefix distinguishes eligible files from unrelated files in the same stage.
		testClientHelper().Stage.PutInLocationWithContent(t, "@"+loadStage.ID().FullyQualifiedName()+"/selected/", "selected.csv", "1\n")
		testClientHelper().Stage.PutInLocationWithContent(t, "@"+loadStage.ID().FullyQualifiedName()+"/other/", "other.csv", "2\n")

		pipe, cleanup := testClientHelper().Pipe.CreatePipe(t, testClientHelper().Pipe.CopyStatement(t, loadTable, loadStage))
		t.Cleanup(cleanup)

		err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(pipe.ID()).
			WithRefresh(*sdk.NewPipeRefreshRequest().
				WithPrefix("selected/")))
		require.NoError(t, err)

		var rows []struct {
			ID int `db:"ID"`
		}
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.NoError(c, client.QueryForTests(ctx, &rows, fmt.Sprintf("SELECT ID FROM %s ORDER BY ID", loadTable.ID().FullyQualifiedName())))
			assert.Equal(c, []struct {
				ID int `db:"ID"`
			}{{ID: 1}}, rows)
		}, 2*time.Minute, 5*time.Second)
	})

	t.Run("safe: show and drop", func(t *testing.T) {
		pipe, cleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
		t.Cleanup(cleanup)

		found, err := client.Pipes.ShowByIDSafely(ctx, pipe.ID())
		require.NoError(t, err)
		require.Equal(t, pipe.ID(), found.ID())

		err = client.Pipes.DropSafely(ctx, pipe.ID())
		require.NoError(t, err)
		err = client.Pipes.DropSafely(ctx, pipe.ID())
		require.NoError(t, err)

		for _, id := range []sdk.SchemaObjectIdentifier{
			pipe.ID(),
			NonExistingSchemaObjectIdentifierWithNonExistingSchema,
			NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema,
		} {
			_, err := client.Pipes.ShowByIDSafely(ctx, id)
			require.ErrorIs(t, err, sdk.ErrObjectNotFound)

			err = client.Pipes.DropSafely(ctx, id)
			require.NoError(t, err)
		}
	})
}
