//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_StorageLifecyclePolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	defaultArgs := testClientHelper().StorageLifecyclePolicy.DefaultArgs()
	defaultSignature := []sdk.TableColumnSignature{
		{Name: "VAL", Type: testdatatypes.DataTypeVarchar},
	}
	defaultBody := testClientHelper().StorageLifecyclePolicy.DefaultBody()

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		cleanup := testClientHelper().StorageLifecyclePolicy.CreateWithId(t, id)
		t.Cleanup(cleanup)

		return id
	}

	t.Run("create: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.StorageLifecyclePolicies.Create(ctx, sdk.NewCreateStorageLifecyclePolicyRequest(id, defaultArgs, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().StorageLifecyclePolicy.DropFunc(t, id))

		assertThatObject(
			t, objectassert.StorageLifecyclePolicy(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasKind("STORAGE_LIFECYCLE_POLICY").
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment("").
				HasOwnerRoleType("ROLE"),
		)
		assertThatObject(
			t, objectassert.StorageLifecyclePolicyDetails(t, id).
				HasName(id.Name()).
				HasSignature(defaultSignature...).
				HasReturnType(testdatatypes.DataTypeBoolean).
				HasBody(defaultBody).
				HasNoArchiveForDays().
				HasArchiveTier(""),
		)
	})

	t.Run("create: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		err := client.StorageLifecyclePolicies.Create(ctx, sdk.NewCreateStorageLifecyclePolicyRequest(id, defaultArgs, defaultBody).
			WithIfNotExists(true).
			WithArchiveTier(sdk.StorageLifecyclePolicyArchiveTierCold).
			WithArchiveForDays(365).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().StorageLifecyclePolicy.DropFunc(t, id))

		assertThatObject(
			t, objectassert.StorageLifecyclePolicy(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasKind("STORAGE_LIFECYCLE_POLICY").
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasComment(comment).
				HasOwnerRoleType("ROLE"),
		)
		assertThatObject(
			t, objectassert.StorageLifecyclePolicyDetails(t, id).
				HasName(id.Name()).
				HasSignature(defaultSignature...).
				HasReturnType(testdatatypes.DataTypeBoolean).
				HasBody(defaultBody).
				HasArchiveForDays(365).
				HasArchiveTier(string(sdk.StorageLifecyclePolicyArchiveTierCold)),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).
			WithRenameTo(newId))
		require.NoError(t, err)

		t.Cleanup(testClientHelper().StorageLifecyclePolicy.DropFunc(t, newId))

		_, err = client.StorageLifecyclePolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		renamed, err := client.StorageLifecyclePolicies.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), renamed.Name)
	})

	t.Run("alter: set body", func(t *testing.T) {
		id := createBasic(t)
		newBody := "LENGTH(VAL) > 5"

		err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).
			WithSetBody(newBody))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.StorageLifecyclePolicyDetails(t, id).
				HasBody(newBody),
		)
	})

	t.Run("alter: set and unset archive_tier, archive_for_days, comment", func(t *testing.T) {
		id := createBasic(t)
		comment := random.Comment()

		err := client.StorageLifecyclePolicies.Alter(
			ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).
				WithSet(*sdk.NewStorageLifecyclePolicySetRequest().
					WithArchiveTier(sdk.StorageLifecyclePolicyArchiveTierCool).
					WithArchiveForDays(120).
					WithComment(comment)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.StorageLifecyclePolicy(t, id).
			HasComment(comment))
		assertThatObject(
			t, objectassert.StorageLifecyclePolicyDetails(t, id).
				HasArchiveForDays(120).
				HasArchiveTier(string(sdk.StorageLifecyclePolicyArchiveTierCool)),
		)

		err = client.StorageLifecyclePolicies.Alter(
			ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).
				WithUnset(*sdk.NewStorageLifecyclePolicyUnsetRequest().
					WithArchiveForDays(true).
					WithComment(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.StorageLifecyclePolicy(t, id).
			HasComment(""))
		assertThatObject(
			t, objectassert.StorageLifecyclePolicyDetails(t, id).
				HasNoArchiveForDays().
				HasArchiveTier(string(sdk.StorageLifecyclePolicyArchiveTierCool)),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.StorageLifecyclePolicies.Drop(ctx, sdk.NewDropStorageLifecyclePolicyRequest(id))
		require.NoError(t, err)

		_, err = client.StorageLifecyclePolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.StorageLifecyclePolicies.Drop(ctx, sdk.NewDropStorageLifecyclePolicyRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.StorageLifecyclePolicies.Drop(ctx, sdk.NewDropStorageLifecyclePolicyRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("show", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(sdk.NewDatabaseObjectIdentifier(db.Name, "PUBLIC"))
		ids := []sdk.SchemaObjectIdentifier{id1, id2, id3}
		for _, id := range ids {
			cleanup := testClientHelper().StorageLifecyclePolicy.CreateWithId(t, id)
			t.Cleanup(cleanup)
		}

		t.Run("like", func(t *testing.T) {
			policies, err := client.StorageLifecyclePolicies.Show(ctx, sdk.NewShowStorageLifecyclePolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
				WithIn(sdk.In{Schema: id1.SchemaId()}))
			require.NoError(t, err)
			assert.Len(t, policies, 1)
			assert.Equal(t, id1.Name(), policies[0].Name)
		})

		t.Run("in_account", func(t *testing.T) {
			policies, err := client.StorageLifecyclePolicies.Show(ctx, sdk.NewShowStorageLifecyclePolicyRequest().
				WithIn(sdk.In{Account: sdk.Bool(true)}))
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(policies), 3)
		})

		t.Run("in_database", func(t *testing.T) {
			policies, err := client.StorageLifecyclePolicies.Show(ctx, sdk.NewShowStorageLifecyclePolicyRequest().
				WithIn(sdk.In{Database: id1.DatabaseId()}))
			require.NoError(t, err)
			assert.Len(t, policies, 2)
		})

		t.Run("in_schema", func(t *testing.T) {
			policies, err := client.StorageLifecyclePolicies.Show(ctx, sdk.NewShowStorageLifecyclePolicyRequest().
				WithIn(sdk.In{Schema: id1.SchemaId()}))
			require.NoError(t, err)
			assert.Len(t, policies, 2)
		})
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.StorageLifecyclePolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
