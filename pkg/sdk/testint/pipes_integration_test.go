//go:build non_account_level_tests

package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPipeCopyStatement(t *testing.T, table *sdk.Table, stage *sdk.Stage) string {
	t.Helper()
	require.NotNil(t, table, "table has to be created")
	require.NotNil(t, stage, "stage has to be created")
	return fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
}

// TestInt_CreatePipeWithStrangeSchemaName documented previous bad behavior. It changed with Snowflake 8.3.1 release.
// We leave the test for future reference.
func TestInt_CreatePipeWithStrangeSchemaName(t *testing.T) {
	schemaIdentifier := testClientHelper().Ids.NewDatabaseObjectIdentifier("tcK1>AJ+")

	// creating a new schema on purpose
	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithName(t, schemaIdentifier.Name())
	t.Cleanup(schemaCleanup)

	table, tableCleanup := testClientHelper().Table.CreateInSchema(t, schema.ID())
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	t.Run("if we have special characters in db or schema name, create pipe succeeds", func(t *testing.T) {
		err := itc.client.Pipes.Create(
			itc.ctx,
			testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()),
			createPipeCopyStatement(t, table, stage),
			&sdk.CreatePipeOptions{},
		)

		require.NoError(t, err)
	})

	t.Run("the same does not work when using non fully qualified name for table", func(t *testing.T) {
		createCopyStatementWithoutQualifiersForStage := func(t *testing.T, table *sdk.Table, stage *sdk.Stage) string {
			t.Helper()
			require.NotNil(t, table, "table has to be created")
			require.NotNil(t, stage, "stage has to be created")
			return fmt.Sprintf("COPY INTO %s\nFROM @\"%s\"", table.ID().FullyQualifiedName(), stage.Name)
		}

		err := itc.client.Pipes.Create(
			itc.ctx,
			testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()),
			createCopyStatementWithoutQualifiersForStage(t, table, stage),
			&sdk.CreatePipeOptions{},
		)

		require.Error(t, err)
		require.ErrorContains(t, err, "object does not exist or not authorized")
	})
}

func TestInt_PipesShowAndDescribe(t *testing.T) {
	table1, table1Cleanup := testClientHelper().Table.Create(t)
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := testClientHelper().Table.Create(t)
	t.Cleanup(table2Cleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	pipe1CopyStatement := createPipeCopyStatement(t, table1, stage)
	pipe1, pipe1Cleanup := testClientHelper().Pipe.CreatePipe(t, pipe1CopyStatement)
	t.Cleanup(pipe1Cleanup)

	pipe2CopyStatement := createPipeCopyStatement(t, table2, stage)
	pipe2, pipe2Cleanup := testClientHelper().Pipe.CreatePipe(t, pipe2CopyStatement)
	t.Cleanup(pipe2Cleanup)

	t.Run("show: without options", func(t *testing.T) {
		pipes, err := itc.client.Pipes.Show(itc.ctx, &sdk.ShowPipeOptions{})

		require.NoError(t, err)
		assert.Len(t, pipes, 2)
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: in schema", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Len(t, pipes, 2)
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: like", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(pipe1.Name),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Len(t, pipes, 1)
		assert.Contains(t, pipes, *pipe1)
	})

	t.Run("show: non-existent pipe", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Empty(t, pipes)
	})

	t.Run("describe: existing pipe", func(t *testing.T) {
		pipe, err := itc.client.Pipes.Describe(itc.ctx, pipe1.ID())

		require.NoError(t, err)
		assert.Equal(t, pipe1.Name, pipe.Name)
	})

	t.Run("describe: non-existing pipe", func(t *testing.T) {
		_, err := itc.client.Pipes.Describe(itc.ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeCreate(t *testing.T) {
	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	assertPipe := func(t *testing.T, pipeDetails *sdk.Pipe, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, pipeDetails.CreatedOn)
		assert.Equal(t, expectedName, pipeDetails.Name)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), pipeDetails.DatabaseName)
		assert.Equal(t, testClientHelper().Ids.SchemaId().Name(), pipeDetails.SchemaName)
		assert.Equal(t, copyStatement, pipeDetails.Definition)
		assert.Equal(t, "ACCOUNTADMIN", pipeDetails.Owner)
		assert.Empty(t, pipeDetails.NotificationChannel)
		assert.Equal(t, expectedComment, pipeDetails.Comment)
		assert.Empty(t, pipeDetails.Integration)
		assert.Empty(t, pipeDetails.Pattern)
		assert.Empty(t, pipeDetails.ErrorIntegration)
		assert.Equal(t, "ROLE", pipeDetails.OwnerRoleType)
		assert.Empty(t, pipeDetails.InvalidReason)
	}

	// TODO: test error integration, aws sns topic and integration when we have them in project
	t.Run("test complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, &sdk.CreatePipeOptions{
			OrReplace:   sdk.Bool(false),
			IfNotExists: sdk.Bool(true),
			AutoIngest:  sdk.Bool(false),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)

		pipe, err := itc.client.Pipes.Describe(itc.ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, id.Name(), comment)
	})

	t.Run("test if not exists and or replace are incompatible", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, &sdk.CreatePipeOptions{
			OrReplace:   sdk.Bool(true),
			IfNotExists: sdk.Bool(true),
		})
		require.ErrorContains(t, err, "(0A000): SQL compilation error:\noptions IF NOT EXISTS and OR REPLACE are incompatible")
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, nil)
		require.NoError(t, err)

		pipe, err := itc.client.Pipes.Describe(itc.ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, id.Name(), "")
	})
}

func TestInt_PipeDrop(t *testing.T) {
	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	t.Run("pipe exists", func(t *testing.T) {
		pipeCopyStatement := createPipeCopyStatement(t, table, stage)
		pipe, _ := testClientHelper().Pipe.CreatePipe(t, pipeCopyStatement)

		err := itc.client.Pipes.Drop(itc.ctx, pipe.ID(), nil)

		require.NoError(t, err)
		_, err = itc.client.Pipes.Describe(itc.ctx, pipe.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("pipe does not exist", func(t *testing.T) {
		err := itc.client.Pipes.Drop(itc.ctx, NonExistingSchemaObjectIdentifier, &sdk.DropPipeOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeAlter(t *testing.T) {
	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	pipeCopyStatement := createPipeCopyStatement(t, table, stage)

	// TODO: test error integration when we have them in project
	t.Run("set value and unset value", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				Comment:             sdk.String("new comment"),
				PipeExecutionPaused: sdk.Bool(true),
			},
		}

		err := itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err := itc.client.Pipes.ShowByID(itc.ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredPipe.Comment)

		alterOptions = &sdk.AlterPipeOptions{
			Unset: &sdk.PipeUnset{
				Comment:             sdk.Bool(true),
				PipeExecutionPaused: sdk.Bool(true),
			},
		}

		err = itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err = itc.client.Pipes.ShowByID(itc.ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "", alteredPipe.Comment)
	})

	t.Run("refresh with all", func(t *testing.T) {
		pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &sdk.AlterPipeOptions{
			Refresh: &sdk.PipeRefresh{
				Prefix:        sdk.String("/d1"),
				ModifiedAfter: sdk.String("2018-07-30T13:56:46-07:00"),
			},
		}

		err := itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)
	})
}

func TestInt_PipesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	cleanupPipeHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Pipes.Drop(ctx, id, nil)
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createPipeHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		statement := createPipeCopyStatement(t, table, stage)
		err := client.Pipes.Create(ctx, id, statement, &sdk.CreatePipeOptions{})
		require.NoError(t, err)
		t.Cleanup(cleanupPipeHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createPipeHandle(t, id1)
		createPipeHandle(t, id2)

		e1, err := client.Pipes.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Pipes.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
