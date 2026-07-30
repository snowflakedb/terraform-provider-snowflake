//go:build non_account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_CortexSearchServices(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouseId := testClientHelper().Ids.WarehouseId()

	on := "some_text_column"
	targetLag := "2 minutes"

	buildQuery := func(tableId sdk.SchemaObjectIdentifier) string {
		return fmt.Sprintf(`select %s from %s`, on, tableId.FullyQualifiedName())
	}

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, buildQuery(table.ID())))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		return id
	}

	t.Run("create: minimal", func(t *testing.T) {
		id := createBasic(t)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasWarehouseNotEmpty().
			HasTargetLag(targetLag).
			HasSearchColumn(strings.ToUpper(on)).
			HasNoAttributeColumns().
			HasColumns(strings.ToUpper(on)).
			HasNoPrimaryKeyColumns().
			HasScoringProfileCount(0).
			HasComment(""))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTargetLag(targetLag).
			HasWarehouseNotEmpty().
			HasSearchColumn(strings.ToUpper(on)).
			HasNoAttributeColumns().
			HasColumns(strings.ToUpper(on)).
			HasNoComment().
			HasServiceQueryUrlNotEmpty().
			HasDataTimestampNotEmpty().
			HasIndexingStateNotEmpty().
			HasIndexingError("").
			HasEmbeddingModel("snowflake-arctic-embed-m-v1.5").
			HasServingStateNotEmpty().
			HasScoringProfileCount(0).
			HasNoPrimaryKeyColumns().
			HasNoFullIndexBuildIntervalDays())
	})

	t.Run("create: with primary key", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`select %s, id from %s`, on, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, query).
			WithPrimaryKey([]string{"id"}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasPrimaryKeyColumns("ID").
			HasColumns("SOME_TEXT_COLUMN", "ID"))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasPrimaryKeyColumns("ID").
			HasColumns("SOME_TEXT_COLUMN", "ID").
			HasFullIndexBuildIntervalDays(1))
	})

	t.Run("create: complete", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		embeddingModel := "snowflake-arctic-embed-m-v1.5"
		query := fmt.Sprintf(`select %s, id, some_other_text_column, another_text_column from %s`, on, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, query).
			WithOrReplace(true).
			WithPrimaryKey([]string{"id", "another_text_column"}).
			WithAttributes(*sdk.NewAttributesRequest().WithColumns([]string{"some_other_text_column", "another_text_column"})).
			WithEmbeddingModel(embeddingModel).
			WithRefreshMode(sdk.CortexSearchServiceRefreshModeIncremental).
			WithInitialize(sdk.CortexSearchServiceInitializeOnCreate).
			WithFullIndexBuildIntervalDays(30).
			WithRequestLogging(true).
			WithAutoSuspend(3600).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasWarehouse(warehouseId.Name()).
			HasTargetLag(targetLag).
			HasSearchColumn(strings.ToUpper(on)).
			HasAttributeColumns("SOME_OTHER_TEXT_COLUMN", "ANOTHER_TEXT_COLUMN").
			HasColumns("SOME_TEXT_COLUMN", "ID", "SOME_OTHER_TEXT_COLUMN", "ANOTHER_TEXT_COLUMN").
			HasPrimaryKeyColumns("ID", "ANOTHER_TEXT_COLUMN").
			HasScoringProfileCount(0).
			HasComment(comment))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTargetLag(targetLag).
			HasWarehouse(warehouseId.Name()).
			HasSearchColumn(strings.ToUpper(on)).
			HasAttributeColumns("SOME_OTHER_TEXT_COLUMN", "ANOTHER_TEXT_COLUMN").
			HasColumns("SOME_TEXT_COLUMN", "ID", "SOME_OTHER_TEXT_COLUMN", "ANOTHER_TEXT_COLUMN").
			HasPrimaryKeyColumns("ID", "ANOTHER_TEXT_COLUMN").
			HasComment(comment).
			HasServiceQueryUrlNotEmpty().
			HasDataTimestampNotEmpty().
			HasIndexingStateNotEmpty().
			HasIndexingError("").
			HasEmbeddingModel(embeddingModel).
			HasServingStateNotEmpty().
			HasScoringProfileCount(0).
			HasFullIndexBuildIntervalDays(30))
	})

	// There are no UNSETs for the attributes used in set, so we don't run "normal" UNSETs in this test but SET to defaults instead
	t.Run("alter: set and unset", func(t *testing.T) {
		newWarehouse, newWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(newWarehouseCleanup)

		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`select %s, id from %s`, on, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, query).
			WithPrimaryKey([]string{"id"}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		newComment := random.Comment()
		newTargetLag := "10 minutes"

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSet(
			*sdk.NewCortexSearchServiceSetRequest().
				WithTargetLag(newTargetLag).
				WithWarehouse(newWarehouse.ID()).
				WithFullIndexBuildIntervalDays(30).
				WithRequestLogging(true).
				WithAutoSuspend(3600).
				WithComment(newComment),
		))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasWarehouse(newWarehouse.ID().Name()).
			HasTargetLag(newTargetLag).
			HasComment(newComment).
			HasPrimaryKeyColumns("ID"))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasWarehouse(newWarehouse.ID().Name()).
			HasTargetLag(newTargetLag).
			HasComment(newComment).
			HasPrimaryKeyColumns("ID").
			HasFullIndexBuildIntervalDays(30))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSet(
			*sdk.NewCortexSearchServiceSetRequest().
				WithFullIndexBuildIntervalDays(0).
				WithRequestLogging(false).
				WithComment(""),
		))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasComment(""))

		// REQUEST_LOGGING and AUTO_SUSPEND are not exposed in SHOW or DESCRIBE output
		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasComment("").
			HasFullIndexBuildIntervalDays(0))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSetDefaults(
			*sdk.NewCortexSearchServiceSetDefaultsRequest().WithAutoSuspend(true),
		))
		require.NoError(t, err)
	})

	t.Run("alter: set and unset primary key", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`select %s, id, another_text_column from %s`, on, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, query))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSetPrimaryKey(
			*sdk.NewCortexSearchServiceSetPrimaryKeyRequest().WithPrimaryKey([]string{"id", "another_text_column"}),
		))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasPrimaryKeyColumns("ID", "ANOTHER_TEXT_COLUMN"))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasPrimaryKeyColumns("ID", "ANOTHER_TEXT_COLUMN"))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithUnsetPrimaryKey(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasNoPrimaryKeyColumns())

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasNoPrimaryKeyColumns())
	})

	t.Run("alter: set and unset attributes", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsForCortexSearchService(t)
		t.Cleanup(tableCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`select %s, some_other_text_column from %s`, on, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, query))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSetAttributes(
			*sdk.NewCortexSearchServiceSetAttributesRequest().WithColumns([]string{"some_other_text_column"}),
		))
		require.NoError(t, err)

		// columns lists all columns of the source query and is not affected by SET/UNSET ATTRIBUTES
		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasAttributeColumns("SOME_OTHER_TEXT_COLUMN").
			HasColumns("SOME_TEXT_COLUMN", "SOME_OTHER_TEXT_COLUMN"))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasAttributeColumns("SOME_OTHER_TEXT_COLUMN").
			HasColumns("SOME_TEXT_COLUMN", "SOME_OTHER_TEXT_COLUMN"))

		err = client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithUnsetAttributes(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CortexSearchService(t, id).
			HasNoAttributeColumns().
			HasColumns("SOME_TEXT_COLUMN", "SOME_OTHER_TEXT_COLUMN"))

		assertThatObject(t, objectassert.CortexSearchServiceDetails(t, id).
			HasNoAttributeColumns().
			HasColumns("SOME_TEXT_COLUMN", "SOME_OTHER_TEXT_COLUMN"))
	})

	t.Run("describe: when cortex search service does not exist", func(t *testing.T) {
		_, err := client.CortexSearchServices.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	// Documents a Snowflake limitation: cortex search services reject double-quoted identifiers
	// in the source query. This means lowercase (case-sensitive) column names
	// cannot be used, and all column references must use the uppercase form that Snowflake stores for unquoted identifiers.
	t.Run("create: does not support double-quoted column names in source query", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsLowercased(t)
		t.Cleanup(tableCleanup)

		name := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		query := fmt.Sprintf(`select "some_text_column" from %s`, table.ID().FullyQualifiedName())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(name, `"some_text_column"`, warehouseId, targetLag, query))
		require.ErrorContains(t, err, "399115 (42601): Invalid source query: quoted identifier or reserved word \"some_text_column\" not allowed.")
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsInSchema(t, schema.ID())
		t.Cleanup(tableCleanup)

		id1 := createBasic(t)
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id2, on, warehouseId, targetLag, buildQuery(table.ID())))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id2))

		e1, err := client.CortexSearchServices.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.CortexSearchServices.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
