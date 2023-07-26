package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createCopyStatement(t *testing.T, table *Table, stage *Stage) string {
	t.Helper()
	require.NotNil(t, table, "table has to be created")
	require.NotNil(t, stage, "stage has to be created")
	return fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
}

func TestInt_IncorrectCreatePipeBehaviour(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := NewSchemaIdentifier("TXR@=9,TBnLj", "tcK1>AJ+")
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, schemaIdentifier.databaseName)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.schemaName)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	t.Run("if we have special characters in db or schema name, create pipe returns error in copy <> from <> section", func(t *testing.T) {
		err := client.Pipes.Create(
			ctx,
			NewSchemaObjectIdentifier(database.Name, schema.Name, randomAlphanumericN(t, 20)),
			createCopyStatement(t, table, stage),
			&PipeCreateOptions{},
		)

		require.ErrorContains(t, err, "(42000): SQL compilation error:\nsyntax error line")
		require.ErrorContains(t, err, "at position")
		require.ErrorContains(t, err, "unexpected ','")
	})

	t.Run("the same works with using db and schema statements", func(t *testing.T) {
		useDatabaseCleanup := useDatabase(t, client, database.ID())
		t.Cleanup(useDatabaseCleanup)
		useSchemaCleanup := useSchema(t, client, schema.ID())
		t.Cleanup(useSchemaCleanup)

		createCopyStatementWithoutQualifiersForStage := func(t *testing.T, table *Table, stage *Stage) string {
			t.Helper()
			require.NotNil(t, table, "table has to be created")
			require.NotNil(t, stage, "stage has to be created")
			return fmt.Sprintf("COPY INTO %s\nFROM @\"%s\"", table.ID().FullyQualifiedName(), stage.Name)
		}

		err := client.Pipes.Create(
			ctx,
			NewSchemaObjectIdentifier(database.Name, schema.Name, randomAlphanumericN(t, 20)),
			createCopyStatementWithoutQualifiersForStage(t, table, stage),
			&PipeCreateOptions{},
		)

		require.NoError(t, err)
	})
}

func TestInt_PipesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericSchemaIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, schemaIdentifier.databaseName)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.schemaName)
	t.Cleanup(schemaCleanup)

	table1, table1Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table2Cleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	pipe1Name := randomAlphanumericN(t, 20)
	pipe1CopyStatement := createCopyStatement(t, table1, stage)
	pipe1, pipe1Cleanup := createPipe(t, client, database, schema, pipe1Name, pipe1CopyStatement)
	t.Cleanup(pipe1Cleanup)

	pipe2Name := randomAlphanumericN(t, 20)
	pipe2CopyStatement := createCopyStatement(t, table2, stage)
	pipe2, pipe2Cleanup := createPipe(t, client, database, schema, pipe2Name, pipe2CopyStatement)
	t.Cleanup(pipe2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		pipes, err := client.Pipes.Show(ctx, &PipeShowOptions{})

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, pipe1)
		assert.Contains(t, pipes, pipe2)
	})

	t.Run("show in schema", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			In: &In{
				Schema: schema.ID(),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, pipe1)
		assert.Contains(t, pipes, pipe2)
	})

	t.Run("show like", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			Like: &Like{
				Pattern: String(pipe1Name),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 1, len(pipes))
		assert.Contains(t, pipes, pipe1)
	})

	t.Run("search for non-existent pipe", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 0, len(pipes))
	})
}

func TestInt_PipeCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericSchemaIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, schemaIdentifier.databaseName)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.schemaName)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	copyStatement := createCopyStatement(t, table, stage)

	assertPipe := func(t *testing.T, pipeDetails *Pipe, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, pipeDetails.CreatedOn)
		assert.Equal(t, expectedName, pipeDetails.Name)
		assert.Equal(t, database.Name, pipeDetails.DatabaseName)
		assert.Equal(t, schema.Name, pipeDetails.SchemaName)
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
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := randomComment(t)

		err := client.Pipes.Create(ctx, id, copyStatement, &PipeCreateOptions{
			OrReplace:   Bool(false),
			IfNotExists: Bool(true),
			AutoIngest:  Bool(false),
			Comment:     String(comment),
		})
		require.NoError(t, err)

		pipe, err := client.Pipes.Describe(ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, comment)
	})

	t.Run("test if not exists and or replace are incompatible", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Pipes.Create(ctx, id, copyStatement, &PipeCreateOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
		})
		require.ErrorContains(t, err, "(0A000): SQL compilation error:\noptions IF NOT EXISTS and OR REPLACE are incompatible")
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Pipes.Create(ctx, id, copyStatement, nil)
		require.NoError(t, err)

		pipe, err := client.Pipes.Describe(ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, "")
	})
}
