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
			HasNoComment().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoQueryWarehouse().
			HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(computePoolCleanup)

		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		queryWarehouse, queryWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(queryWarehouseCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		testClientHelper().Stage.PutOnStage(t, stage.ID(), "example.ipynb")
		location := sdk.NewStageLocation(stage.ID(), "")

		// TODO(SNOW-2398051) Some of the fields were omitted due to lack of documentation.
		request := sdk.NewCreateNotebookRequest(id).WithIfNotExists(true).
			WithComment("comment").
			WithTitle("title").
			WithFrom(location).
			WithMainFile("example.ipynb").
			WithComputePool(computePool.ID()).
			WithIdleAutoShutdownTimeSeconds(3600).
			WithDefaultVersion("FIRST").
			WithWarehouse(warehouse.ID()).
			WithQueryWarehouse(queryWarehouse.ID())

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
			HasComment("comment").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasQueryWarehouse(queryWarehouse.ID()).
			HasOwnerRoleType("ROLE").
			HasCodeWarehouse(warehouse.ID()),
		)

		assertThatObject(t, objectassert.NotebookDetails(t, notebook.ID()).
			HasTitle("title").
			HasMainFile("example.ipynb").
			HasQueryWarehouse(queryWarehouse.ID()).
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasComputePool(computePool.ID()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse(warehouse.ID().Name()).
			HasIdleAutoShutdownTimeSeconds(3600).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(id.Name()).
			HasComment("comment").
			HasDefaultVersion("FIRST").
			HasDefaultVersionName("VERSION$1").
			HasNoDefaultVersionAlias().
			HasNonEmptyDefaultVersionLocationUri().
			HasDefaultVersionSourceLocationUri(stage.Location()).
			HasNoDefaultVersionGitCommitHash().
			HasLastVersionName("VERSION$1").
			HasNoLastVersionAlias().
			HasNonEmptyLastVersionLocationUri().
			HasLastVersionSourceLocationUri(stage.Location()).
			HasNoLastVersionGitCommitHash().
			HasNoLiveVersionLocationUri(),
		)
	})

	t.Run("alter: set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateNotebookRequest(id)

		err := client.Notebooks.Create(ctx, request)

		computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(computePoolCleanup)

		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		queryWarehouse, queryWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(queryWarehouseCleanup)

		// secret, secreteCleanup := testClientHelper().Secret.CreateRandomPasswordSecret(t)
		// t.Cleanup(secreteCleanup)

		// secrets := sdk.SecretsListRequest{SecretsList: []sdk.SecretReference{{
		// 	VariableName: "sample_secret",
		// 	Name:         secret,
		// }}}

		// TODO: Investigate the 'Secrets' field (not present in both SHOW and DESC).
		setRequest := sdk.NewNotebookSetRequest().
			WithComment("comment").
			WithQueryWarehouse(queryWarehouse.ID()).
			WithIdleAutoShutdownTimeSeconds(3600).
			WithMainFile("example.ipynb").
			WithWarehouse(warehouse.ID()).
			WithComputePool(computePool.ID())

		alterRequest := sdk.NewAlterNotebookRequest(id).WithSet(*setRequest)

		err = client.Notebooks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedNotebook, err := client.Notebooks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.NotebookFromObject(t, updatedNotebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("comment").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasQueryWarehouse(queryWarehouse.ID()).
			HasOwnerRoleType("ROLE").
			HasCodeWarehouse(warehouse.ID()),
		)

		assertThatObject(t, objectassert.NotebookDetails(t, updatedNotebook.ID()).
			HasNoTitle().
			HasMainFile("example.ipynb").
			HasQueryWarehouse(queryWarehouse.ID()).
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasComputePool(computePool.ID()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse(warehouse.ID().Name()).
			HasIdleAutoShutdownTimeSeconds(3600).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(id.Name()).
			HasComment("comment").
			HasDefaultVersion("LAST").
			HasDefaultVersionName("VERSION$1").
			HasNoDefaultVersionAlias().
			HasNonEmptyDefaultVersionLocationUri().
			HasNoDefaultVersionSourceLocationUri().
			HasNoDefaultVersionGitCommitHash().
			HasLastVersionName("VERSION$1").
			HasNoLastVersionAlias().
			HasNonEmptyLastVersionLocationUri().
			HasNoDefaultVersionSourceLocationUri().
			HasNoLastVersionGitCommitHash().
			HasNoLiveVersionLocationUri(),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(computePoolCleanup)

		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		queryWarehouse, queryWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(queryWarehouseCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		testClientHelper().Stage.PutOnStage(t, stage.ID(), "example.ipynb")
		location := sdk.NewStageLocation(stage.ID(), "")

		createRequest := sdk.NewCreateNotebookRequest(id).WithIfNotExists(true).
			WithComment("comment").
			WithTitle("title").
			WithFrom(location).
			WithMainFile("example.ipynb").
			WithComputePool(computePool.ID()).
			WithIdleAutoShutdownTimeSeconds(3600).
			WithDefaultVersion("FIRST").
			WithWarehouse(warehouse.ID()).
			WithQueryWarehouse(queryWarehouse.ID())

		_, notebookCleanup := testClientHelper().Notebook.CreateWithRequest(t, createRequest)
		t.Cleanup(notebookCleanup)

		// 'QueryWarehouse' and 'Warehouse' are mutually exclusive (SQL execution internal error: Processing aborted due to error 300002:787288943; incident 2110983).
		unsetRequest := sdk.NewNotebookUnsetRequest().
			WithComment(true).
			WithComputePool(true).
			WithExternalAccessIntegrations(true).
			WithQueryWarehouse(true).
			WithRuntimeEnvironmentVersion(true).
			WithRuntimeName(true).
			WithSecrets(true)

		alterRequest := sdk.NewAlterNotebookRequest(id).WithUnset(*unsetRequest)

		err := client.Notebooks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// 'Warehouse' parameter separately.
		alterRequest2 := sdk.NewAlterNotebookRequest(id).WithUnset(*sdk.NewNotebookUnsetRequest().WithQueryWarehouse(true))
		err = client.Notebooks.Alter(ctx, alterRequest2)
		require.NoError(t, err)

		updatedNotebook, err := testClientHelper().Notebook.Show(t, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.NotebookFromObject(t, updatedNotebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasNoComment().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoQueryWarehouse().
			HasOwnerRoleType("ROLE").
			HasCodeWarehouse(warehouse.ID()),
		)

		assertThatObject(t, objectassert.NotebookDetails(t, updatedNotebook.ID()).
			HasTitle("title").
			HasMainFile("example.ipynb").
			HasNoQueryWarehouse().
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasNoComputePool().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse(warehouse.ID().Name()).
			HasIdleAutoShutdownTimeSeconds(3600).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(id.Name()).
			HasNoComment().
			HasDefaultVersion("FIRST").
			HasDefaultVersionName("VERSION$1").
			HasNoDefaultVersionAlias().
			HasNonEmptyDefaultVersionLocationUri().
			HasDefaultVersionSourceLocationUri(stage.Location()).
			HasNoDefaultVersionGitCommitHash().
			HasLastVersionName("VERSION$1").
			HasNoLastVersionAlias().
			HasNonEmptyLastVersionLocationUri().
			HasDefaultVersionSourceLocationUri(stage.Location()).
			HasNoLastVersionGitCommitHash().
			HasNoLiveVersionLocationUri(),
		)
	})

	t.Run("drop", func(t *testing.T) {
		notebook, notebookCleanup := testClientHelper().Notebook.Create(t)
		t.Cleanup(notebookCleanup)

		id := notebook.ID()
		err := client.Notebooks.Drop(ctx, sdk.NewDropNotebookRequest(id).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.Notebooks.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe", func(t *testing.T) {
		notebook, notebookCleanup := testClientHelper().Notebook.Create(t)
		t.Cleanup(notebookCleanup)

		assertThatObject(t, objectassert.NotebookDetails(t, notebook.ID()).
			HasMainFile("notebook_app.ipynb").
			HasNoQueryWarehouse().
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasNoComputePool().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse("SYSTEM$STREAMLIT_NOTEBOOK_WH").
			HasIdleAutoShutdownTimeSeconds(1800).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(notebook.ID().Name()).
			HasNoComment().
			HasDefaultVersion("LAST").
			HasDefaultVersionName("VERSION$1").
			HasNoDefaultVersionAlias().
			HasNonEmptyDefaultVersionLocationUri().
			HasNoDefaultVersionSourceLocationUri().
			HasNoDefaultVersionGitCommitHash().
			HasLastVersionName("VERSION$1").
			HasNoLastVersionAlias().
			HasNonEmptyLastVersionLocationUri().
			HasNoLastVersionSourceLocationUri().
			HasNoLastVersionGitCommitHash().
			HasNoLiveVersionLocationUri(),
		)
	})

	t.Run("show: with like", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createRequest := sdk.NewCreateNotebookRequest(id).WithComment("comment")
		notebook, notebookCleanup := testClientHelper().Notebook.CreateWithRequest(t, createRequest)
		t.Cleanup(notebookCleanup)

		pattern := id.Name()
		notebooks, err := client.Notebooks.Show(ctx, sdk.NewShowNotebookRequest().WithLike(sdk.Like{Pattern: &pattern}))
		require.NoError(t, err)
		require.Len(t, notebooks, 1)
		require.Equal(t, *notebook, notebooks[0])
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())

		_, notebookCleanup1 := testClientHelper().Notebook.CreateWithRequest(t, sdk.NewCreateNotebookRequest(id1))
		t.Cleanup(notebookCleanup1)
		_, notebookCleanup2 := testClientHelper().Notebook.CreateWithRequest(t, sdk.NewCreateNotebookRequest(id2))
		t.Cleanup(notebookCleanup2)

		e1, err := client.Notebooks.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Notebooks.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
