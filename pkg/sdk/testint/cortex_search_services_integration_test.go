//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
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

	createCortexSearchService := func(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.CortexSearchService {
		t.Helper()

		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumns(t)
		t.Cleanup(tableCleanup)

		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseId, targetLag, buildQuery(table.ID())))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CortexSearchService.DropCortexSearchServiceFunc(t, id))

		cortexSearchService, err := client.CortexSearchServices.ShowByID(ctx, id)
		require.NoError(t, err)

		return cortexSearchService
	}

	t.Run("create: test complete", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumns(t)
		t.Cleanup(tableCleanup)

		name := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(name, on, testClientHelper().Ids.WarehouseId(), targetLag, buildQuery(table.ID())).WithOrReplace(true).WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.CortexSearchServices.Drop(ctx, sdk.NewDropCortexSearchServiceRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.CortexSearchServices.Show(ctx, sdk.NewShowCortexSearchServiceRequest().WithLike(sdk.Like{Pattern: sdk.String(name.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)

		cortexSearchServiceById, err := client.CortexSearchServices.ShowByID(ctx, name)
		require.NoError(t, err)
		require.NotNil(t, cortexSearchServiceById)
		require.NotEmpty(t, cortexSearchServiceById.CreatedOn)
		require.Equal(t, name.Name(), cortexSearchServiceById.Name)
		require.Equal(t, name.DatabaseName(), cortexSearchServiceById.DatabaseName)
		require.Equal(t, name.SchemaName(), cortexSearchServiceById.SchemaName)
		require.Equal(t, comment, cortexSearchServiceById.Comment)
	})

	t.Run("describe: when cortex search service exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		cortexSearchService := createCortexSearchService(t, id)

		cortexSearchServiceDetails, err := client.CortexSearchServices.Describe(ctx, cortexSearchService.ID())
		require.NoError(t, err)
		assert.NotEmpty(t, cortexSearchServiceDetails.CreatedOn)
		assert.Equal(t, cortexSearchService.Name, cortexSearchServiceDetails.Name)
		// Yes, the names are exchanged on purpose, because now it works like this
		assert.Equal(t, cortexSearchService.DatabaseName, cortexSearchServiceDetails.DatabaseName)
		assert.Equal(t, cortexSearchService.SchemaName, cortexSearchServiceDetails.SchemaName)
		assert.Equal(t, targetLag, cortexSearchServiceDetails.TargetLag)
		assert.NotEmpty(t, cortexSearchServiceDetails.Warehouse)
		assert.Equal(t, strings.ToUpper(on), *cortexSearchServiceDetails.SearchColumn)
		assert.NotEmpty(t, cortexSearchServiceDetails.AttributeColumns)
		assert.NotEmpty(t, cortexSearchServiceDetails.Columns)
		assert.NotEmpty(t, cortexSearchServiceDetails.Definition)
		assert.Nil(t, cortexSearchServiceDetails.Comment)
		assert.NotEmpty(t, cortexSearchServiceDetails.ServiceQueryUrl)
		assert.NotEmpty(t, cortexSearchServiceDetails.DataTimestamp)
		assert.GreaterOrEqual(t, cortexSearchServiceDetails.SourceDataNumRows, 0)
		assert.NotEmpty(t, cortexSearchServiceDetails.IndexingState)
		assert.Empty(t, cortexSearchServiceDetails.IndexingError)
	})

	t.Run("describe: when cortex search service does not exist", func(t *testing.T) {
		_, err := client.CortexSearchServices.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: with set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createCortexSearchService(t, id)

		newComment := "new comment"
		newTargetLag := "10 minutes"

		err := client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(id).WithSet(*sdk.NewCortexSearchServiceSetRequest().
			WithTargetLag(newTargetLag).
			WithComment(newComment),
		))
		require.NoError(t, err)

		alteredService, err := client.CortexSearchServices.ShowByID(ctx, id)
		require.NoError(t, err)

		require.Equal(t, newComment, alteredService.Comment)

		cortexSearchServiceDetails, err := client.CortexSearchServices.Describe(ctx, id)
		require.NoError(t, err)

		require.Equal(t, newComment, *cortexSearchServiceDetails.Comment)
		require.Equal(t, newTargetLag, cortexSearchServiceDetails.TargetLag)
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		// order matters in this test, creating the schema first and then trying to create cortex search service in the default test schema fails with a strange error
		// (probably caused by the implicit use schema after schema creation)
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createCortexSearchService(t, id1)

		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())
		createCortexSearchService(t, id2)

		e1, err := client.CortexSearchServices.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.CortexSearchServices.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
