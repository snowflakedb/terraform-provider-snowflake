//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_McpServers(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	complexSpec := `tools:
  - description: "Executes arbitrary SQL statements against Snowflake, with full DML and DDL support."
    type: "SYSTEM_EXECUTE_SQL"
    name: "sql_exec_tool"
    title: "SQL Execution Tool"
  - name: "sql_query_tool"
    title: "SQL Query Tool"
    description: "Reads data from Snowflake tables, views, and streams."
    type: "SYSTEM_EXECUTE_SQL"
  - title: "Admin SQL Tool"
    description: "Performs administrative operations such as managing warehouses and users."
    name: "admin_sql_tool"
    type: "SYSTEM_EXECUTE_SQL"
`
	normalizedComplexSpec, err := sdk.NormalizeMcpServerSpecification(complexSpec)
	require.NoError(t, err)

	defaultSpec := testClientHelper().McpServer.DefaultSpec()
	normalizedDefaultSpec, err := sdk.NormalizeMcpServerSpecification(defaultSpec)
	require.NoError(t, err)

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.McpServers.Create(ctx, sdk.NewCreateMcpServerRequest(id, defaultSpec))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().McpServer.DropFunc(t, id))

		assertThatObject(
			t, objectassert.McpServer(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment(""),
		)
		assertThatObject(
			t, objectassert.McpServerDetails(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment("").
				HasServerSpec(normalizedDefaultSpec).
				HasCreatedOnNotEmpty(),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		err := client.McpServers.Create(ctx, sdk.NewCreateMcpServerRequest(id, complexSpec).
			WithIfNotExists(true).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().McpServer.DropFunc(t, id))

		assertThatObject(
			t, objectassert.McpServer(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment(comment),
		)
		assertThatObject(
			t, objectassert.McpServerDetails(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment(comment).
				HasServerSpec(normalizedComplexSpec).
				HasCreatedOnNotEmpty(),
		)
	})

	t.Run("show", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		cleanup := testClientHelper().McpServer.Create(t, id)
		t.Cleanup(cleanup)

		findById := func(results []sdk.McpServer) (*sdk.McpServer, error) {
			return collections.FindFirst(results, func(result sdk.McpServer) bool {
				return result.ID().FullyQualifiedName() == id.FullyQualifiedName()
			})
		}

		t.Run("in_account", func(t *testing.T) {
			showResults, err := client.McpServers.Show(ctx, sdk.NewShowMcpServerRequest().WithIn(sdk.In{Account: sdk.Bool(true)}))
			require.NoError(t, err)

			require.GreaterOrEqual(t, len(showResults), 1)
			result, err := findById(showResults)
			require.NoError(t, err)
			require.Equal(t, id.FullyQualifiedName(), result.ID().FullyQualifiedName())
		})

		t.Run("in_database", func(t *testing.T) {
			showResults, err := client.McpServers.Show(ctx, sdk.NewShowMcpServerRequest().WithIn(sdk.In{Database: id.DatabaseId()}))
			require.NoError(t, err)

			require.GreaterOrEqual(t, len(showResults), 1)
			result, err := findById(showResults)
			require.NoError(t, err)
			require.Equal(t, id.FullyQualifiedName(), result.ID().FullyQualifiedName())
		})

		t.Run("in_schema", func(t *testing.T) {
			showResults, err := client.McpServers.Show(ctx, sdk.NewShowMcpServerRequest().WithIn(sdk.In{Schema: id.SchemaId()}))
			require.NoError(t, err)

			require.GreaterOrEqual(t, len(showResults), 1)
			result, err := findById(showResults)
			require.NoError(t, err)
			require.Equal(t, id.FullyQualifiedName(), result.ID().FullyQualifiedName())
		})
	})

	t.Run("show by id", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		cleanup := testClientHelper().McpServer.Create(t, id)
		t.Cleanup(cleanup)

		result, err := client.McpServers.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.FullyQualifiedName(), result.ID().FullyQualifiedName())
	})

	t.Run("drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.McpServers.Create(ctx, sdk.NewCreateMcpServerRequest(id, complexSpec))
		require.NoError(t, err)

		err = client.McpServers.Drop(ctx, sdk.NewDropMcpServerRequest(id).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.McpServers.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		err = client.McpServers.Drop(ctx, sdk.NewDropMcpServerRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})
}
