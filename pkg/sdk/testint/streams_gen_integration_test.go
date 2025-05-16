//go:build !account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Streams(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseId := testClientHelper().Ids.DatabaseId()
	schemaId := testClientHelper().Ids.SchemaId()

	// There is no way to check at/before fields in show and describe. That's why in Create tests we try creating with these values, but do not assert them.
	t.Run("CreateOnTable", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateWithChangeTracking(t)
		t.Cleanup(cleanupTable)
		tableId := table.ID()

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOn(*sdk.NewOnStreamRequest().WithAt(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithOffset("0"))).
			WithAppendOnly(true).
			WithShowInitialRows(true).
			WithComment("some comment").
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
			})
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStream)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v1"), tag1Value)

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeTable).
			HasMode(sdk.StreamModeAppendOnly).
			HasTableId(tableId),
		)

		// at stream
		req = sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOrReplace(true).
			WithOn(*sdk.NewOnStreamRequest().WithAt(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithStream(id.FullyQualifiedName())))
		err = client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		// at statement
		_, err = testClient(t).ExecForTests(ctx, fmt.Sprintf("INSERT INTO %s VALUES(1);", table.ID().FullyQualifiedName()))
		require.NoError(t, err)

		lastQueryId, err := testClient(t).ContextFunctions.LastQueryId(ctx)
		require.NoError(t, err)

		req = sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOrReplace(true).
			WithOn(*sdk.NewOnStreamRequest().WithAt(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithStatement(lastQueryId)))
		err = client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		// before offset
		req = sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOrReplace(true).
			WithOn(*sdk.NewOnStreamRequest().WithBefore(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithOffset("0")))
		err = client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		// before stream
		req = sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOrReplace(true).
			WithOn(*sdk.NewOnStreamRequest().WithBefore(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithStream(id.FullyQualifiedName())))
		err = client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		// before statement
		_, err = testClient(t).ExecForTests(ctx, fmt.Sprintf("INSERT INTO %s VALUES(1);", table.ID().FullyQualifiedName()))
		require.NoError(t, err)

		lastQueryId, err = testClient(t).ContextFunctions.LastQueryId(ctx)
		require.NoError(t, err)

		req = sdk.NewCreateOnTableStreamRequest(id, tableId).
			WithOrReplace(true).
			WithOn(*sdk.NewOnStreamRequest().WithBefore(true).WithStatement(*sdk.NewOnStreamStatementRequest().WithStatement(lastQueryId)))
		err = client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		// TODO(SNOW-1689111): test timestamps
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		stage, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t)
		t.Cleanup(stageCleanup)

		externalTableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		externalTableReq := sdk.NewCreateExternalTableRequest(externalTableId, stage.Location()).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON))
		_, externalTableCleanup := testClientHelper().ExternalTable.CreateWithRequest(t, externalTableReq)
		t.Cleanup(externalTableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnExternalTableStreamRequest(id, externalTableId).WithInsertOnly(true).WithComment("some comment")
		err := client.Streams.CreateOnExternalTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeExternalTable).
			HasMode(sdk.StreamModeInsertOnly).
			HasTableId(externalTableId),
		)
	})

	t.Run("CreateOnDirectoryTable", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStageWithDirectory(t)
		t.Cleanup(cleanupStage)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnDirectoryTableStreamRequest(id, stage.ID()).WithComment("some comment")
		err := client.Streams.CreateOnDirectoryTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeStage).
			HasMode(sdk.StreamModeDefault).
			HasStageName(stage.ID().Name()).
			HasBaseTablesPartiallyQualified(stage.ID().Name()),
		)
	})

	t.Run("CreateOnView", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		view, cleanupView := testClientHelper().View.CreateView(t, fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName()))
		t.Cleanup(cleanupView)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnViewStreamRequest(id, view.ID()).
			WithAppendOnly(true).
			WithShowInitialRows(true).
			WithComment("some comment")
		err := client.Streams.CreateOnView(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeView).
			HasMode(sdk.StreamModeAppendOnly).
			HasTableId(view.ID()),
		)
	})

	t.Run("Clone", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		cloneId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.Streams.Clone(ctx, sdk.NewCloneStreamRequest(cloneId, id).WithCopyGrants(true))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, cloneId))

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeTable).
			HasMode(sdk.StreamModeDefault).
			HasTableId(table.ID()),
		)
	})

	t.Run("Alter comment", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID())
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment("some_comment"))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "some_comment", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(true))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		_, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("Show terse", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().WithTerse(true))
		require.NoError(t, err)

		stream, err := collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assert.NotNil(t, stream)
		assert.Equal(t, id.Name(), stream.Name)
		assert.Equal(t, databaseId.Name(), stream.DatabaseName)
		assert.Equal(t, schemaId.Name(), stream.SchemaName)
		assert.Nil(t, stream.Comment)
		assert.Nil(t, stream.SourceType)
		assert.Nil(t, stream.Mode)
	})

	t.Run("Show single with options", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().
			WithTerse(false).
			WithIn(sdk.ExtendedIn{
				In: sdk.In{
					Schema: schemaId,
				},
			}).
			WithLike(sdk.Like{
				Pattern: sdk.String(id.Name()),
			}).
			WithStartsWith(id.Name()).
			WithLimit(sdk.LimitFrom{
				Rows: sdk.Int(1),
			}))
		require.NoError(t, err)
		assert.Equal(t, 1, len(s))

		_, err = collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeTable).
			HasMode(sdk.StreamModeDefault).
			HasTableId(table.ID()),
		)
	})

	t.Run("Show multiple", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req2 := sdk.NewCreateOnTableStreamRequest(id2, table.ID()).WithComment("some comment")
		err = client.Streams.CreateOnTable(ctx, req2)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id2))

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest())
		require.NoError(t, err)
		assert.Equal(t, 2, len(s))

		_, err = collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		_, err = collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id2.Name() == stream.Name })
		require.NoError(t, err)
	})

	t.Run("Show multiple with options", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		idPrefix := "stream_show_"

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(idPrefix)
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(idPrefix)
		req2 := sdk.NewCreateOnTableStreamRequest(id2, table.ID()).WithComment("some comment")
		err = client.Streams.CreateOnTable(ctx, req2)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id2))

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().
			WithTerse(false).
			WithIn(sdk.ExtendedIn{
				In: sdk.In{
					Schema: schemaId,
				},
			}).
			WithStartsWith(idPrefix).
			WithLimit(sdk.LimitFrom{
				Rows: sdk.Int(2),
			}))
		require.NoError(t, err)
		assert.Equal(t, 2, len(s))

		_, err = collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		_, err = collections.FindFirst[sdk.Stream](s, func(stream sdk.Stream) bool { return id2.Name() == stream.Name })
		require.NoError(t, err)
	})

	t.Run("Describe", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schemaId)
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment("some comment")
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Stream.DropFunc(t, id))

		s, err := client.Streams.Describe(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, s)

		assertThatObject(t, objectassert.Stream(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("some comment").
			HasSourceType(sdk.StreamSourceTypeTable).
			HasMode(sdk.StreamModeDefault).
			HasTableId(table.ID()),
		)
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		table, cleanupTable := testClientHelper().Table.Create(t)
		t.Cleanup(cleanupTable)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, stream1Cleanup := testClientHelper().Stream.CreateOnTableWithRequest(t, sdk.NewCreateOnTableStreamRequest(id1, table.ID()))
		t.Cleanup(stream1Cleanup)
		_, stream2Cleanup := testClientHelper().Stream.CreateOnTableWithRequest(t, sdk.NewCreateOnTableStreamRequest(id2, table.ID()))
		t.Cleanup(stream2Cleanup)

		e1, err := client.Streams.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Streams.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
