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

func TestInt_BackupSets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	create := func(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateBackupSetRequest) {
		t.Helper()
		require.NoError(t, client.BackupSets.Create(ctx, request))
		t.Cleanup(func() {
			require.NoError(t, client.BackupSets.DropSafely(ctx, id))
		})
	}

	tableTarget := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		table, cleanup := testClientHelper().Table.CreateWithPredefinedColumns(t)
		t.Cleanup(cleanup)
		return table.ID()
	}

	backupSetNames := func(sets []sdk.BackupSet) []string {
		names := make([]string, len(sets))
		for i, set := range sets {
			names[i] = set.Name
		}
		return names
	}

	assertTableSet := func(t *testing.T, id sdk.SchemaObjectIdentifier, tableId sdk.SchemaObjectIdentifier) {
		t.Helper()
		assertThatObject(t, objectassert.BackupSet(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasObjectKind("TABLE").
			HasObjectName(tableId.Name()).
			HasObjectDatabaseName(tableId.DatabaseName()).
			HasObjectSchemaName(tableId.SchemaName()).
			HasNoBackupPolicyName().
			HasNoBackupPolicyDatabaseName().
			HasNoBackupPolicySchemaName().
			HasNoBackupPolicyState().
			HasOwnerRole(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasNoComment(),
		)
		assertThatObject(t, objectassert.BackupSetDetails(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasObjectKind("TABLE").
			HasObjectName(tableId.Name()).
			HasObjectDatabaseName(tableId.DatabaseName()).
			HasObjectSchemaName(tableId.SchemaName()).
			HasNoBackupPolicyName().
			HasNoComment(),
		)
	}

	t.Run("create: for table", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		tableId := tableTarget(t)
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableId))
		assertTableSet(t, id, tableId)
	})

	t.Run("create: for dynamic table", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		base, baseCleanup := testClientHelper().Table.CreateWithPredefinedColumns(t)
		t.Cleanup(baseCleanup)
		dynamicTable, dynamicCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, base.ID())
		t.Cleanup(dynamicCleanup)

		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForDynamicTable(dynamicTable.ID()))

		assertThatObject(t, objectassert.BackupSet(t, id).
			HasObjectKind("TABLE").
			HasObjectName(dynamicTable.ID().Name()).
			HasObjectDatabaseName(dynamicTable.ID().DatabaseName()).
			HasObjectSchemaName(dynamicTable.ID().SchemaName()),
		)
		assertThatObject(t, objectassert.BackupSetDetails(t, id).
			HasObjectKind("TABLE").
			HasObjectName(dynamicTable.ID().Name()),
		)
	})

	t.Run("create: for schema", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForSchema(schema.ID()))

		assertThatObject(t, objectassert.BackupSet(t, id).
			HasObjectKind("SCHEMA").
			HasObjectName(schema.ID().Name()).
			HasObjectDatabaseName(schema.ID().DatabaseName()).
			HasNoObjectSchemaName(),
		)
		assertThatObject(t, objectassert.BackupSetDetails(t, id).
			HasObjectKind("SCHEMA").
			HasObjectName(schema.ID().Name()).
			HasNoObjectSchemaName(),
		)
	})

	t.Run("create: for database", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForDatabase(database.ID()))

		assertThatObject(t, objectassert.BackupSet(t, id).
			HasObjectKind("DATABASE").
			HasObjectName(database.ID().Name()).
			HasNoObjectDatabaseName().
			HasNoObjectSchemaName(),
		)
	})

	// OR REPLACE conflicts with IF NOT EXISTS, so the complete create uses IF NOT EXISTS.
	t.Run("create: all optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		tableId := tableTarget(t)
		comment := random.Comment()
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		policyId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		require.NoError(t, client.BackupPolicies.Create(ctx, sdk.NewCreateBackupPolicyRequest(policyId).WithExpireAfterDays(7)))
		t.Cleanup(func() {
			require.NoError(t, client.BackupPolicies.DropSafely(ctx, policyId))
		})

		create(t, id, sdk.NewCreateBackupSetRequest(id).
			WithIfNotExists(true).
			WithForTable(tableId).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}).
			WithWithBackupPolicy(policyId).
			WithComment(comment),
		)

		err := client.BackupSets.Create(ctx, sdk.NewCreateBackupSetRequest(id).
			WithIfNotExists(true).
			WithForTable(tableId).
			WithComment("ignored"),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupSet(t, id).
			HasComment(comment).
			HasBackupPolicyName(policyId.Name()).
			HasBackupPolicyDatabaseName(policyId.DatabaseName()).
			HasBackupPolicySchemaName(policyId.SchemaName()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive),
		)
		assertThatObject(t, objectassert.BackupSetDetails(t, id).
			HasComment(comment).
			HasBackupPolicyName(policyId.Name()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive),
		)
		value, err := testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupSet)
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, "v1", *value)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		tableId := tableTarget(t)
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableId).WithComment("first"))

		err := client.BackupSets.Create(ctx, sdk.NewCreateBackupSetRequest(id).
			WithOrReplace(true).
			WithForTable(tableId).
			WithComment("replaced"),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.BackupSet(t, id).HasComment("replaced"))
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, client.BackupSets.DropSafely(ctx, newId))
		})

		_, err = client.BackupSets.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		got, err := client.BackupSets.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.Equal(t, newId, got.ID())
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))
		comment := random.Comment()

		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithSet(*sdk.NewBackupSetSetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasComment(comment))

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithUnset(*sdk.NewBackupSetUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasNoComment())
	})

	t.Run("alter: set and unset tags", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithSetTags([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}),
		)
		require.NoError(t, err)
		value, err := testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupSet)
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, "v1", *value)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)
		value, err = testClientHelper().Tag.GetForObject(t, tag.ID(), id, sdk.ObjectTypeBackupSet)
		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("alter: add backup and modify it", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))

		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithAddBackup(true))
		require.NoError(t, err)
		backupID := testClientHelper().BackupSet.BackupID(t, id)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithModifyBackup(*sdk.NewBackupSetModifyBackupRequest(backupID).
				WithSetCommentModifyBackupParams(*sdk.NewSetCommentModifyBackupParamsRequest("backup comment"))),
		)
		require.NoError(t, err)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithModifyBackup(*sdk.NewBackupSetModifyBackupRequest(backupID).
				WithUnsetCommentModifyBackupParams(*sdk.NewUnsetCommentModifyBackupParamsRequest())),
		)
		require.NoError(t, err)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithModifyBackup(*sdk.NewBackupSetModifyBackupRequest(backupID).
				WithAddLegalHoldModifyBackupParams(*sdk.NewAddLegalHoldModifyBackupParamsRequest())),
		)
		require.NoError(t, err)
		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithModifyBackup(*sdk.NewBackupSetModifyBackupRequest(backupID).
				WithRemoveLegalHoldModifyBackupParams(*sdk.NewRemoveLegalHoldModifyBackupParamsRequest())),
		)
		require.NoError(t, err)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithDeleteBackupIdentifier(backupID),
		)
		require.NoError(t, err)
	})

	t.Run("alter: apply policy, suspend, and resume", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))
		policyId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		replacementId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		require.NoError(t, client.BackupPolicies.Create(ctx, sdk.NewCreateBackupPolicyRequest(policyId).WithExpireAfterDays(7)))
		t.Cleanup(func() {
			require.NoError(t, client.BackupPolicies.DropSafely(ctx, policyId))
		})
		require.NoError(t, client.BackupPolicies.Create(ctx, sdk.NewCreateBackupPolicyRequest(replacementId).WithExpireAfterDays(8)))
		t.Cleanup(func() {
			require.NoError(t, client.BackupPolicies.DropSafely(ctx, replacementId))
		})

		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithApplyBackupPolicy(*sdk.NewBackupSetApplyBackupPolicyRequest(policyId)),
		)
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).
			HasBackupPolicyName(policyId.Name()).
			HasBackupPolicyState(sdk.BackupPolicyStateActive),
		)

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).
			WithApplyBackupPolicy(*sdk.NewBackupSetApplyBackupPolicyRequest(replacementId).WithForce(true)),
		)
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasBackupPolicyName(replacementId.Name()))

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithSuspendBackupCreationPolicy(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasBackupPolicyState(sdk.BackupPolicyStateCreationSuspended))

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithResumeBackupCreationPolicy(true))
		require.NoError(t, err)
		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithSuspendBackupExpirationPolicy(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasBackupPolicyState(sdk.BackupPolicyStateExpirationSuspended))

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithResumeBackupExpirationPolicy(true))
		require.NoError(t, err)
		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithSuspendBackupPolicy(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasBackupPolicyState(sdk.BackupPolicyStateSuspended))

		err = client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(id).WithResumeBackupPolicy(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.BackupSet(t, id).HasBackupPolicyState(sdk.BackupPolicyStateActive))
	})

	t.Run("alter: if exists on non-existing", func(t *testing.T) {
		err := client.BackupSets.Alter(ctx, sdk.NewAlterBackupSetRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true).
			WithAddBackup(true),
		)
		require.NoError(t, err)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))

		err := client.BackupSets.Drop(ctx, sdk.NewDropBackupSetRequest(id))
		require.NoError(t, err)

		_, err = client.BackupSets.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.BackupSets.Drop(ctx, sdk.NewDropBackupSetRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.BackupSets.Drop(ctx, sdk.NewDropBackupSetRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))

		err := client.BackupSets.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.BackupSets.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.BackupSets.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(schemaCleanup)
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsInSchema(t, schema.ID())
		t.Cleanup(tableCleanup)
		otherTable, otherTableCleanup := testClientHelper().Table.CreateWithPredefinedColumnsInSchema(t, otherSchema.ID())
		t.Cleanup(otherTableCleanup)

		prefix := "BSHOW"
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(otherSchema.ID())
		create(t, id1, sdk.NewCreateBackupSetRequest(id1).WithForTable(table.ID()))
		create(t, id2, sdk.NewCreateBackupSetRequest(id2).WithForTable(table.ID()))
		create(t, id3, sdk.NewCreateBackupSetRequest(id3).WithForTable(otherTable.ID()))

		t.Run("like", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().
				WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.Equal(t, []string{id1.Name()}, backupSetNames(sets))
		})

		t.Run("in account", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().WithIn(sdk.In{Account: sdk.Bool(true)}))
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(sets), 3)
			require.Subset(t, backupSetNames(sets), []string{id1.Name(), id2.Name(), id3.Name()})
		})

		t.Run("in database", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().WithIn(sdk.In{Database: db.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name(), id3.Name()}, backupSetNames(sets))
		})

		t.Run("in schema", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name()}, backupSetNames(sets))
		})

		t.Run("starts with", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().
				WithStartsWith(prefix).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.ElementsMatch(t, []string{id1.Name(), id2.Name()}, backupSetNames(sets))
		})

		t.Run("limit", func(t *testing.T) {
			sets, err := client.BackupSets.Show(ctx, sdk.NewShowBackupSetRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.In{Schema: schema.ID()}))
			require.NoError(t, err)
			require.Len(t, sets, 1)
		})
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)
		table, tableCleanup := testClientHelper().Table.CreateWithPredefinedColumns(t)
		t.Cleanup(tableCleanup)
		otherTable, otherTableCleanup := testClientHelper().Table.CreateInSchema(t, otherSchema.ID())
		t.Cleanup(otherTableCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())
		create(t, id1, sdk.NewCreateBackupSetRequest(id1).WithForTable(table.ID()))
		create(t, id2, sdk.NewCreateBackupSetRequest(id2).WithForTable(otherTable.ID()))

		got1, err := client.BackupSets.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, got1.ID())

		got2, err := client.BackupSets.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, got2.ID())
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		create(t, id, sdk.NewCreateBackupSetRequest(id).WithForTable(tableTarget(t)))

		got, err := client.BackupSets.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, got.ID())
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		_, err := client.BackupSets.ShowByIDSafely(ctx, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.BackupSets.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
