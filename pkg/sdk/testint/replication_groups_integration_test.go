//go:build non_account_level_tests

package testint

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Database and share replication works on every edition. The complete create includes other object
// types, which need Business Critical Edition or higher. A secondary group needs the secondary test account.
func TestInt_ReplicationGroups(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	secondaryClient := testSecondaryClient(t)
	secondaryCtx := testSecondaryContext(t)
	sourceAccountID := testClientHelper().Account.GetAccountIdentifier(t)
	secondaryAccountID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	create := func(t *testing.T, id sdk.AccountObjectIdentifier, request *sdk.CreateReplicationGroupRequest) {
		t.Helper()
		require.NoError(t, client.ReplicationGroups.Create(ctx, request))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, client, ctx, id)
		})
	}

	// DATABASES and ALLOWED_ACCOUNTS are the required pair. One database satisfies the allowed-database list.
	t.Run("create: basic", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		create(t, id, sdk.NewCreateReplicationGroupRequest(
			id,
			[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}))

		group, err := client.ReplicationGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Contains(t, group.AllowedAccounts, secondaryAccountID)
		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasType(sdk.ReplicationGroupTypeReplication).
			HasIsPrimary(true).
			HasNoComment().
			HasObjectTypes(sdk.ReplicationGroupObjectTypeDatabases).
			HasNoAllowedIntegrationTypes().
			HasNoReplicationSchedule().
			HasNoSecondaryState().
			HasNoReadOnly().
			HasIsListingAutoFulfillmentGroup(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountNameNotEmpty().
			HasAccountLocatorNotEmpty().
			HasOrganizationNameNotEmpty().
			HasSnowflakeRegionNotEmpty().
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccountID, id)),
		)
	})

	// IF NOT EXISTS is the only existence clause. CREATE has no OR REPLACE.
	// Every optional that can share this statement is set. A second IF NOT EXISTS create leaves the first group in place.
	t.Run("create: complete", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		share, cleanupShare := testClientHelper().Share.CreateShare(t)
		t.Cleanup(cleanupShare)
		externalVolume, cleanupVolume := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanupVolume)
		errorIntegration, cleanupIntegration := testClientHelper().NotificationIntegration.Create(t)
		t.Cleanup(cleanupIntegration)
		tag, cleanupTag := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(cleanupTag)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.ReplicationGroupObjectType{
			sdk.ReplicationGroupObjectTypeDatabases,
			sdk.ReplicationGroupObjectTypeExternalVolumes,
			sdk.ReplicationGroupObjectTypeShares,
			sdk.ReplicationGroupObjectTypeIntegrations,
			sdk.ReplicationGroupObjectTypeRoles,
		}
		allowedAccounts := []sdk.AccountIdentifier{sourceAccountID, secondaryAccountID}
		integrationTypes := slices.Clone(sdk.AllReplicationGroupIntegrationTypes)

		create(t, id, sdk.NewCreateReplicationGroupRequest(id, objectTypes, allowedAccounts).
			WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
			WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}).
			WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}).
			WithAllowedIntegrationTypes(integrationTypes).
			WithIgnoreEditionCheck(true).
			WithReplicationSchedule("10 MINUTE").
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}).
			WithErrorIntegration(errorIntegration.ID()),
		)

		err := client.ReplicationGroups.Create(ctx, sdk.NewCreateReplicationGroupRequest(
			id,
			[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
			WithReplicationSchedule("60 MINUTE"))
		require.NoError(t, err)

		group, err := client.ReplicationGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.ElementsMatch(t, objectTypes, group.ObjectTypes)
		require.ElementsMatch(t, integrationTypes, group.AllowedIntegrationTypes)
		require.Subset(t, group.AllowedAccounts, allowedAccounts)
		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasName(id.Name()).
			HasType(sdk.ReplicationGroupTypeReplication).
			HasIsPrimary(true).
			HasNoComment().
			HasReplicationSchedule("10 MINUTE").
			HasNoSecondaryState().
			HasOwner(snowflakeroles.Accountadmin.Name()),
		)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeReplicationGroup, "v1")
	})

	t.Run("create secondary: basic", func(t *testing.T) {
		primaryID, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// A new primary is not always visible to the target account immediately.
		time.Sleep(time.Second)

		require.NoError(t, secondaryClient.ReplicationGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryReplicationGroupRequest(
			secondaryID,
			sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID),
		)))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, secondaryClient, secondaryCtx, secondaryID)
		})

		assertBasicSecondaryReplicationGroup(t, secondaryID, primaryID, sourceAccountID)
	})

	// IF NOT EXISTS is the only optional on a secondary create. A second call leaves the first replica in place.
	t.Run("create secondary: all", func(t *testing.T) {
		primaryID, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		time.Sleep(time.Second)
		request := sdk.NewCreateSecondaryReplicationGroupRequest(
			secondaryID,
			sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID),
		).WithIfNotExists(true)

		require.NoError(t, secondaryClient.ReplicationGroups.CreateSecondary(secondaryCtx, request))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, secondaryClient, secondaryCtx, secondaryID)
		})
		require.NoError(t, secondaryClient.ReplicationGroups.CreateSecondary(secondaryCtx, request))

		assertBasicSecondaryReplicationGroup(t, secondaryID, primaryID, sourceAccountID)
	})

	t.Run("show by id: same name in two accounts", func(t *testing.T) {
		primaryID, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		time.Sleep(time.Second)
		require.NoError(t, secondaryClient.ReplicationGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryReplicationGroupRequest(
			primaryID,
			sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID),
		)))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, secondaryClient, secondaryCtx, primaryID)
		})

		assertThatObject(t, objectassert.ReplicationGroup(t, primaryID).HasIsPrimary(true))
		assertThatObjectSecondary(t, objectassert.ReplicationGroup(t, primaryID).HasIsPrimary(false))
	})

	t.Run("show", func(t *testing.T) {
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)

		shown, err := client.ReplicationGroups.Show(ctx, sdk.NewShowReplicationGroupRequest())
		require.NoError(t, err)
		require.True(t, slices.ContainsFunc(shown, func(group sdk.ReplicationGroup) bool {
			return group.Name == id.Name() && group.IsPrimary
		}))

		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasName(id.Name()).
			HasIsPrimary(true),
		)

		got, err := client.ReplicationGroups.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), got.Name)

		_, err = client.ReplicationGroups.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		_, err = client.ReplicationGroups.ShowByIDSafely(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).WithRenameTo(newID)))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, client, ctx, newID)
		})

		_, err := client.ReplicationGroups.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assertThatObject(t, objectassert.ReplicationGroup(t, newID).
			HasName(newID.Name()).
			HasIsPrimary(true).
			HasObjectTypes(sdk.ReplicationGroupObjectTypeDatabases),
		)
	})

	// OBJECT_TYPES and the allowed lists are one SET. COMMENT, the schedule, and the error integration are another.
	t.Run("alter: set object types", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		share, cleanupShare := testClientHelper().Share.CreateShare(t)
		t.Cleanup(cleanupShare)
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		objectTypes := []sdk.ReplicationGroupObjectType{
			sdk.ReplicationGroupObjectTypeDatabases,
			sdk.ReplicationGroupObjectTypeShares,
			sdk.ReplicationGroupObjectTypeIntegrations,
			sdk.ReplicationGroupObjectTypeRoles,
		}

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithSet(*sdk.NewReplicationGroupSetRequest().
				WithObjectTypes(objectTypes).
				WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
				WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}).
				WithAllowedIntegrationTypes(sdk.AllReplicationGroupIntegrationTypes)),
		))

		group, err := client.ReplicationGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.ElementsMatch(t, objectTypes, group.ObjectTypes)
		require.ElementsMatch(t, sdk.AllReplicationGroupIntegrationTypes, group.AllowedIntegrationTypes)
	})

	// COMMENT, the schedule, and the error integration share one SET. UNSET clears those three together.
	t.Run("alter: set and unset properties", func(t *testing.T) {
		errorIntegration, cleanupIntegration := testClientHelper().NotificationIntegration.Create(t)
		t.Cleanup(cleanupIntegration)
		tag, cleanupTag := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(cleanupTag)
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		comment := random.Comment()

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithSet(*sdk.NewReplicationGroupSetRequest().
				WithComment(comment).
				WithReplicationSchedule("10 MINUTE").
				WithErrorIntegration(errorIntegration.ID())),
		))
		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasComment(comment).
			HasReplicationSchedule("10 MINUTE"),
		)

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithUnset(*sdk.NewReplicationGroupUnsetRequest().
				WithComment(true).
				WithReplicationSchedule(true).
				WithErrorIntegration(true)),
		))
		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasNoComment().
			HasNoReplicationSchedule(),
		)

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithSetTags([]sdk.TagAssociation{{Name: tag.ID(), Value: "v2"}}),
		))
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeReplicationGroup, "v2")

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}),
		))
		assertTagUnset(t, tag.ID(), id, sdk.ObjectTypeReplicationGroup)
	})

	t.Run("alter: add and remove", func(t *testing.T) {
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		added, cleanupAdded := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupAdded)

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithAdd(*sdk.NewReplicationGroupAddRequest().
				WithToAllowedDatabases([]sdk.AccountObjectIdentifier{added.ID()})),
		))
		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).
			WithRemove(*sdk.NewReplicationGroupRemoveRequest().
				WithFromAllowedDatabases([]sdk.AccountObjectIdentifier{added.ID()})),
		))
		assertThatObject(t, objectassert.ReplicationGroup(t, id).
			HasName(id.Name()).
			HasObjectTypes(sdk.ReplicationGroupObjectTypeDatabases),
		)
	})

	t.Run("alter: move", func(t *testing.T) {
		sourceID, databaseID := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		targetDatabase, cleanupTargetDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupTargetDatabase)
		targetID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		create(t, targetID, sdk.NewCreateReplicationGroupRequest(
			targetID,
			[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithAllowedDatabases([]sdk.AccountObjectIdentifier{targetDatabase.ID()}))

		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(sourceID).
			WithMove(*sdk.NewReplicationGroupMoveRequest(targetID).
				WithDatabases([]sdk.AccountObjectIdentifier{databaseID})),
		))
		assertThatObject(t, objectassert.ReplicationGroup(t, sourceID).HasName(sourceID.Name()))
		assertThatObject(t, objectassert.ReplicationGroup(t, targetID).HasName(targetID.Name()))
	})

	t.Run("alter secondary: suspend, resume, and refresh", func(t *testing.T) {
		primaryID, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)
		require.NoError(t, client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(primaryID).
			WithSet(*sdk.NewReplicationGroupSetRequest().WithReplicationSchedule("10 MINUTE")),
		))
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		time.Sleep(time.Second)
		require.NoError(t, secondaryClient.ReplicationGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryReplicationGroupRequest(
			secondaryID,
			sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID),
		)))
		t.Cleanup(func() {
			dropReplicationGroupEventually(t, secondaryClient, secondaryCtx, secondaryID)
		})

		require.NoError(t, secondaryClient.ReplicationGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryReplicationGroupRequest(secondaryID).
			WithSuspend(*sdk.NewReplicationGroupSecondarySuspendRequest()),
		))
		assertThatObjectSecondary(t, objectassert.ReplicationGroup(t, secondaryID).
			HasSecondaryState(sdk.ReplicationGroupSecondaryStateSuspended),
		)

		require.NoError(t, secondaryClient.ReplicationGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryReplicationGroupRequest(secondaryID).
			WithSuspend(*sdk.NewReplicationGroupSecondarySuspendRequest().WithImmediate(true)),
		))
		assertThatObjectSecondary(t, objectassert.ReplicationGroup(t, secondaryID).
			HasSecondaryState(sdk.ReplicationGroupSecondaryStateSuspended),
		)

		require.NoError(t, secondaryClient.ReplicationGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryReplicationGroupRequest(secondaryID).
			WithResume(true),
		))
		assertThatObjectSecondary(t, objectassert.ReplicationGroup(t, secondaryID).
			HasSecondaryState(sdk.ReplicationGroupSecondaryStateStarted),
		)

		// A refresh can already be running. Retry until this one is accepted.
		assert.Eventually(t, func() bool {
			return secondaryClient.ReplicationGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryReplicationGroupRequest(secondaryID).
				WithRefresh(true)) == nil
		}, 20*time.Second, time.Second)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)

		require.NoError(t, client.ReplicationGroups.Drop(ctx, sdk.NewDropReplicationGroupRequest(id)))
		_, err := client.ReplicationGroups.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.ReplicationGroups.Drop(ctx, sdk.NewDropReplicationGroupRequest(NonExistingAccountObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.ReplicationGroups.Drop(ctx, sdk.NewDropReplicationGroupRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id, _ := createPrimaryWithDatabase(t, client, ctx, secondaryAccountID)

		require.NoError(t, client.ReplicationGroups.DropSafely(ctx, id))
		_, err := client.ReplicationGroups.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.ReplicationGroups.DropSafely(ctx, NonExistingAccountObjectIdentifier)
		require.NoError(t, err)
	})
}

// SHOW REPLICATION GROUPS fills object_types, allowed_integration_types, and allowed_accounts on the primary
// row only. The secondary row leaves those lists blank, copies the schedule when one is set, and leaves
// secondary_state null when no schedule is set. read_only is null for every replication group. Owner is set
// here because the replica is read in its own account; the documented null owner is for a group in another region.
func assertBasicSecondaryReplicationGroup(t *testing.T, secondaryID, primaryID sdk.AccountObjectIdentifier, sourceAccountID sdk.AccountIdentifier) {
	t.Helper()
	assertThatObjectSecondary(t, objectassert.ReplicationGroup(t, secondaryID).
		HasCreatedOnNotEmpty().
		HasName(secondaryID.Name()).
		HasType(sdk.ReplicationGroupTypeReplication).
		HasIsPrimary(false).
		HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID)).
		HasNoComment().
		HasNoObjectTypes().
		HasNoAllowedIntegrationTypes().
		HasNoAllowedAccounts().
		HasNoReplicationSchedule().
		HasNoSecondaryState().
		HasNoReadOnly().
		HasIsListingAutoFulfillmentGroup(false).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasAccountNameNotEmpty().
		HasAccountLocatorNotEmpty().
		HasOrganizationNameNotEmpty().
		HasSnowflakeRegionNotEmpty(),
	)
}

func createPrimaryWithDatabase(t *testing.T, client *sdk.Client, ctx context.Context, allowedAccount sdk.AccountIdentifier) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	t.Helper()
	database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	require.NoError(t, client.ReplicationGroups.Create(ctx, sdk.NewCreateReplicationGroupRequest(
		id,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{allowedAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()})))
	t.Cleanup(func() {
		dropReplicationGroupEventually(t, client, ctx, id)
	})
	return id, database.ID()
}

func dropReplicationGroupEventually(t *testing.T, client *sdk.Client, ctx context.Context, id sdk.AccountObjectIdentifier) {
	t.Helper()
	assert.Eventually(t, func() bool {
		return client.ReplicationGroups.DropSafely(ctx, id) == nil
	}, 10*time.Second, time.Second)
}
