//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_CodeBundles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// SHOW CODE BUNDLES has only IN. There is no LIKE, STARTS WITH, or LIMIT.
	// ALTER supports ADD VERSION FROM, SET COMMENT, and UNSET COMMENT. RENAME TO is unsupported.
	createBundle := func(t *testing.T, id sdk.SchemaObjectIdentifier, adjust func(*sdk.CreateCodeBundleRequest)) *sdk.Stage {
		t.Helper()
		stage, cleanupStage := testClientHelper().Stage.CreateStageInSchema(t, id.SchemaId())
		t.Cleanup(cleanupStage)
		testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "main.py", "print('code bundle')\n")
		request := sdk.NewCreateCodeBundleRequest(id, stage.Location())
		if adjust != nil {
			adjust(request)
		}
		err := client.CodeBundles.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CodeBundle.DropFunc(t, id))
		return stage
	}

	t.Run("create: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBundle(t, id, nil)

		bundle, err := client.CodeBundles.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.CodeBundleFromObject(t, bundle).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			// SHOW CODE BUNDLES returns an empty owner for a stage-backed bundle. DESCRIBE owner_role is set.
			HasOwner("").
			HasNoComment(),
		)
	})

	t.Run("create: comment and if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		stage := createBundle(t, id, func(request *sdk.CreateCodeBundleRequest) {
			request.WithIfNotExists(true).WithComment("code bundle comment")
		})

		err := client.CodeBundles.Create(ctx, sdk.NewCreateCodeBundleRequest(id, stage.Location()).WithIfNotExists(true))
		require.NoError(t, err)

		bundle, err := client.CodeBundles.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.CodeBundleFromObject(t, bundle).
			HasName(id.Name()).
			HasComment("code bundle comment"),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		stage := createBundle(t, id, func(request *sdk.CreateCodeBundleRequest) {
			request.WithComment("first")
		})

		err := client.CodeBundles.Create(ctx, sdk.NewCreateCodeBundleRequest(id, stage.Location()).
			WithOrReplace(true).
			WithComment("replaced"),
		)
		require.NoError(t, err)

		bundle, err := client.CodeBundles.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.CodeBundleFromObject(t, bundle).
			HasComment("replaced"),
		)
	})

	t.Run("alter: add version", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		stage := createBundle(t, id, nil)

		err := client.CodeBundles.Alter(ctx, sdk.NewAlterCodeBundleRequest(id).WithAddVersionFrom(stage.Location()))
		require.NoError(t, err)

		details, err := client.CodeBundles.DescribeDetails(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, details.ID())
		assertThatObject(t, objectassert.CodeBundleDetailsFromObject(t, details).
			HasLastVersionNameNotEmpty().
			HasLastVersionSourceLocationUriContaining(stage.ID().Name()),
		)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBundle(t, id, nil)

		err := client.CodeBundles.Alter(ctx, sdk.NewAlterCodeBundleRequest(id).
			WithSet(*sdk.NewCodeBundleSetRequest().WithComment("code bundle comment")))
		require.NoError(t, err)

		bundle, err := client.CodeBundles.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.CodeBundleFromObject(t, bundle).
			HasComment("code bundle comment"),
		)

		err = client.CodeBundles.Alter(ctx, sdk.NewAlterCodeBundleRequest(id).
			WithUnset(*sdk.NewCodeBundleUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		bundle, err = client.CodeBundles.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.CodeBundleFromObject(t, bundle).
			HasNoComment(),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBundle(t, id, nil)

		err := client.CodeBundles.Drop(ctx, sdk.NewDropCodeBundleRequest(id))
		require.NoError(t, err)

		_, err = client.CodeBundles.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.CodeBundles.Drop(ctx, sdk.NewDropCodeBundleRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.CodeBundles.Drop(ctx, sdk.NewDropCodeBundleRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBundle(t, id, nil)

		err := client.CodeBundles.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.CodeBundles.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.CodeBundles.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(schemaCleanup)
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(otherSchema.ID())
		createBundle(t, id1, nil)
		createBundle(t, id2, nil)
		createBundle(t, id3, nil)

		t.Run("in account", func(t *testing.T) {
			bundles, err := client.CodeBundles.Show(ctx, sdk.NewShowCodeBundleRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(bundles), 3)
			require.Subset(t, bundleNames(bundles), []string{id1.Name(), id2.Name(), id3.Name()})
		})

		t.Run("in database", func(t *testing.T) {
			bundles, err := client.CodeBundles.Show(ctx, sdk.NewShowCodeBundleRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Database: db.ID()}}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name(), id3.Name()}, bundleNames(bundles))
		})

		t.Run("in schema", func(t *testing.T) {
			bundles, err := client.CodeBundles.Show(ctx, sdk.NewShowCodeBundleRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schema.ID()}}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name()}, bundleNames(bundles))
		})
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())
		createBundle(t, id1, nil)
		createBundle(t, id2, nil)

		got1, err := client.CodeBundles.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, got1.ID())

		got2, err := client.CodeBundles.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, got2.ID())
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBundle(t, id, nil)

		bundle, err := client.CodeBundles.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, bundle.ID())
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		_, err := client.CodeBundles.ShowByIDSafely(ctx, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		stage := createBundle(t, id, nil)

		details, err := client.CodeBundles.DescribeDetails(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, details.ID())
		assertThatObject(t, objectassert.CodeBundleDetailsFromObject(t, details).
			HasName(id.Name()).
			HasNoComment().
			HasDefaultVersion("LAST").
			HasDefaultVersionSourceLocationUriContaining(stage.ID().Name()).
			HasContainerType("Stage").
			HasContainerName(stage.ID().Name()).
			HasCreatedOnNotEmpty().
			HasOwnerRole("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.CodeBundles.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func bundleNames(bundles []sdk.CodeBundle) []string {
	names := make([]string, len(bundles))
	for i, bundle := range bundles {
		names[i] = bundle.Name
	}
	return names
}
