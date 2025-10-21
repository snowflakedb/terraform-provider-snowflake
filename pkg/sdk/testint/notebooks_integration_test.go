//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Notebooks(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateNotebookRequest(id)

		err := client.Notebooks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Notebook.DropFunc(t, id))

		notebook, err := client.Notebooks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.NotebookFromObject(t, notebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
		t.Cleanup(dbCleanup)

		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(schemaCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
		t.Cleanup(stageCleanup)
		location := sdk.NewStageLocation(stage.ID(), "")

		computePool := sdk.NewAccountObjectIdentifier("pool")
		externalAccessIntegrations := []sdk.AccountObjectIdentifier{}

		request := sdk.NewCreateNotebookRequest(id).WithIfNotExists(true).WithFrom(location).
			WithTitle("title").
			WithMainFile("main_file").
			WithComment("comment").
			WithIdleAutoShutdownTimeSeconds(3600).
			WithRuntimeName("rname").
			WithComputePool(computePool).
			WithRuntimeEnvironmentVersion("Last").
			WithDefaultVersion("Last")

		err := client.Notebooks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Notebook.DropFunc(t, id))

		notebook, err := client.Notebooks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.NotebookFromObject(t, notebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasNoQueryWarehouse().
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasOwnerRoleType("ROLE"),
		)
	})
}
