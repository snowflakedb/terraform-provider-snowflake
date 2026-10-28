//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_BackupPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	create := func(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateBackupPolicyRequest) {
		t.Helper()
		require.NoError(t, client.BackupPolicies.Create(ctx, request))
		t.Cleanup(func() {
			require.NoError(t, client.BackupPolicies.DropSafely(ctx, id))
		})
	}

	// Either SCHEDULE or EXPIRE_AFTER_DAYS satisfies create. The minimal case uses expiration.
	t.Run("create: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasNoSchedule().
			HasExpireAfterDays(7).
			HasHasRetentionLock(false),
		)
		assertThatObject(t, objectassert.BackupPolicyDetails(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasComment("").
			HasNoSchedule().
			HasExpireAfterDays(7).
			HasHasRetentionLock(false),
		)
	})

	// OR REPLACE conflicts with IF NOT EXISTS, so the complete create uses IF NOT EXISTS.
	t.Run("create: all optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		create(t, id, sdk.NewCreateBackupPolicyRequest(id).
			WithIfNotExists(true).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}).
			WithWithRetentionLock(true).
			WithSchedule("2 HOUR").
			WithExpireAfterDays(30).
			WithComment(comment),
		)

		err := client.BackupPolicies.Create(ctx, sdk.NewCreateBackupPolicyRequest(id).
			WithIfNotExists(true).
			WithExpireAfterDays(1).
			WithComment("ignored"),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasComment(comment).
			HasSchedule("2 HOUR").
			HasExpireAfterDays(30).
			HasHasRetentionLock(true),
		)
		assertThatObject(t, objectassert.BackupPolicyDetails(t, id).
			HasComment(comment).
			HasSchedule("2 HOUR").
			HasExpireAfterDays(30).
			HasHasRetentionLock(true),
		)
		value, err := testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupPolicy)
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, "v1", *value)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).
			WithExpireAfterDays(7).
			WithComment("first"),
		)

		err := client.BackupPolicies.Create(ctx, sdk.NewCreateBackupPolicyRequest(id).
			WithOrReplace(true).
			WithExpireAfterDays(8).
			WithComment("replaced"),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasExpireAfterDays(8).
			HasComment("replaced"),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, client.BackupPolicies.DropSafely(ctx, newId))
		})

		_, err = client.BackupPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		assertThatObject(t, objectassert.BackupPolicy(t, newId).
			HasName(newId.Name()).
			HasExpireAfterDays(7),
		)
	})

	// SET can change comment, schedule, and expiration together.
	// UNSET can clear comment with one of schedule or expiration, not both of those.
	t.Run("alter: set and unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))
		comment := random.Comment()
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err := client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithSet(*sdk.NewBackupPolicySetRequest().
				WithComment(comment).
				WithSchedule("2 HOUR").
				WithExpireAfterDays(30)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasComment(comment).
			HasSchedule("2 HOUR").
			HasExpireAfterDays(30),
		)
		assertThatObject(t, objectassert.BackupPolicyDetails(t, id).
			HasComment(comment).
			HasSchedule("2 HOUR").
			HasExpireAfterDays(30),
		)

		err = client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithUnset(*sdk.NewBackupPolicyUnsetRequest().
				WithComment(true).
				WithSchedule(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasComment("").
			HasNoSchedule().
			HasExpireAfterDays(30),
		)

		err = client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithSet(*sdk.NewBackupPolicySetRequest().WithSchedule("60 MINUTE")),
		)
		require.NoError(t, err)
		err = client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithUnset(*sdk.NewBackupPolicyUnsetRequest().WithExpireAfterDays(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupPolicy(t, id).
			HasSchedule("60 MINUTE").
			HasNoExpireAfterDays(),
		)

		err = client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithSetTags([]sdk.TagAssociation{{Name: tag.ID(), Value: "v2"}}),
		)
		require.NoError(t, err)
		value, err := testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupPolicy)
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, "v2", *value)

		err = client.BackupPolicies.Alter(ctx, sdk.NewAlterBackupPolicyRequest(id).
			WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}),
		)
		require.NoError(t, err)
		value, err = testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupPolicy)
		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))

		err := client.BackupPolicies.Drop(ctx, sdk.NewDropBackupPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.BackupPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.BackupPolicies.Drop(ctx, sdk.NewDropBackupPolicyRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.BackupPolicies.Drop(ctx, sdk.NewDropBackupPolicyRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))

		err := client.BackupPolicies.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.BackupPolicies.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.BackupPolicies.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(schemaCleanup)
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		prefix := "BPSHOW"
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(otherSchema.ID())
		create(t, id1, sdk.NewCreateBackupPolicyRequest(id1).WithExpireAfterDays(7))
		create(t, id2, sdk.NewCreateBackupPolicyRequest(id2).WithExpireAfterDays(7))
		create(t, id3, sdk.NewCreateBackupPolicyRequest(id3).WithExpireAfterDays(7))

		t.Run("like", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.Equal(t, []string{id1.Name()}, backupPolicyNames(policies))
		})

		t.Run("in account", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithIn(sdk.In{Account: sdk.Bool(true)}))
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(policies), 3)
			require.Subset(t, backupPolicyNames(policies), []string{id1.Name(), id2.Name(), id3.Name()})
		})

		t.Run("in database", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithIn(sdk.In{Database: db.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name(), id3.Name()}, backupPolicyNames(policies))
		})

		t.Run("in schema", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name()}, backupPolicyNames(policies))
		})

		t.Run("starts with", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithStartsWith(prefix).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name()}, backupPolicyNames(policies))
		})

		t.Run("limit", func(t *testing.T) {
			policies, err := client.BackupPolicies.Show(ctx, sdk.NewShowBackupPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.Len(t, policies, 1)
		})
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())
		create(t, id1, sdk.NewCreateBackupPolicyRequest(id1).WithExpireAfterDays(7))
		create(t, id2, sdk.NewCreateBackupPolicyRequest(id2).WithExpireAfterDays(7))

		got1, err := client.BackupPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, got1.ID())

		got2, err := client.BackupPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, got2.ID())
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupPolicyRequest(id).WithExpireAfterDays(7))

		got, err := client.BackupPolicies.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, got.ID())
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		_, err := client.BackupPolicies.ShowByIDSafely(ctx, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.BackupPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func backupPolicyNames(policies []sdk.BackupPolicy) []string {
	names := make([]string, len(policies))
	for i, policy := range policies {
		names[i] = policy.Name
	}
	return names
}
