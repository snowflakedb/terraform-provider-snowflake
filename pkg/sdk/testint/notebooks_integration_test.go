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

		computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(computePoolCleanup)

		// TODO(SNOW-2398051) Some of the fields were omitted due to lack of documentation.
		request := sdk.NewCreateNotebookRequest(id).WithIfNotExists(true).
			WithComment("comment").
			WithTitle("title").
			WithMainFile("main_file").
			WithComputePool(computePool.ID()).
			WithIdleAutoShutdownTimeSeconds(3600).
			WithDefaultVersion("FIRST")

		err := client.Notebooks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Notebook.DropFunc(t, id))

		notebook, err := client.Notebooks.ShowByID(ctx, id)
		require.NoError(t, err)

		comment := "comment"

		assertThatObject(t, objectassert.NotebookFromObject(t, notebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment(&comment).
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasNoQueryWarehouse().
			HasOwnerRoleType("ROLE"),
		)

		assertThatObject(t, objectassert.NotebookDetails(t, notebook.ID()).
			HasTitle("title").
			HasMainFile("main_file").
			HasNoQueryWarehouse().
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasComputePool(computePool.ID()).
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse("SYSTEM$STREAMLIT_NOTEBOOK_WH").
			HasIdleAutoShutdownTimeSeconds(3600).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(id.Name()).
			HasComment(&comment).
			HasDefaultVersion("FIRST").
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

	t.Run("alter: set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		_, notebookCleanup := testClientHelper().Notebook.Create(t, id)
		t.Cleanup(notebookCleanup)

		setRequest := sdk.NewNotebookSetRequest().WithComment("comment")
		alterRequest := sdk.NewAlterNotebookRequest(id).WithSet(*setRequest)

		err := client.Notebooks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedNotebook, err := client.Notebooks.ShowByID(ctx, id)
		require.NoError(t, err)

		comment := "comment"

		assertThatObject(t, objectassert.NotebookFromObject(t, updatedNotebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment(&comment).
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasNoQueryWarehouse().
			HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		createRequest := sdk.NewCreateNotebookRequest(id).WithComment("comment")
		_, notebookCleanup := testClientHelper().Notebook.CreateWithRequest(t, createRequest)
		t.Cleanup(notebookCleanup)

		unsetRequest := sdk.NewNotebookUnsetRequest().WithComment(true)
		alterRequest := sdk.NewAlterNotebookRequest(id).WithUnset(*unsetRequest)

		err := client.Notebooks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedNotebook, err := testClientHelper().Notebook.Show(t, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.NotebookFromObject(t, updatedNotebook).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment(nil).
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasNoQueryWarehouse().
			HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, notebookCleanup := testClientHelper().Notebook.Create(t, id)
		t.Cleanup(notebookCleanup)

		err := client.Notebooks.Drop(ctx, sdk.NewDropNotebookRequest(id).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.Notebooks.ShowByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		notebook, notebookCleanup := testClientHelper().Notebook.Create(t, id)
		t.Cleanup(notebookCleanup)

		assertThatObject(t, objectassert.NotebookDetails(t, notebook.ID()).
			HasMainFile("notebook_app.ipynb").
			HasNoQueryWarehouse().
			HasUrlId().
			HasNonEmptyDefaultPackages().
			HasUserPackages("").
			HasNoComputePool().
			HasOwner(snowflakeroles.PentestingRole.Name()).
			HasImportUrls("[]").
			HasExternalAccessIntegrations("[]").
			HasExternalAccessSecrets("{}").
			HasCodeWarehouse("SYSTEM$STREAMLIT_NOTEBOOK_WH").
			HasIdleAutoShutdownTimeSeconds(1800).
			HasRuntimeEnvironmentVersion("WH-RUNTIME-2.0").
			HasName(id.Name()).
			HasComment(nil).
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
		_, notebookCleanup := testClientHelper().Notebook.CreateWithRequest(t, createRequest)
		t.Cleanup(notebookCleanup)

		notebook, err := testClientHelper().Notebook.Show(t, id)
		require.NoError(t, err)

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

		_, notebookCleanup1 := testClientHelper().Notebook.Create(t, id1)
		t.Cleanup(notebookCleanup1)
		_, notebookCleanup2 := testClientHelper().Notebook.Create(t, id2)
		t.Cleanup(notebookCleanup2)

		e1, err := client.Notebooks.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Notebooks.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
