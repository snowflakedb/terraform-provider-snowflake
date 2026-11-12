//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such types (the FailoverGroup test helper creates a ROLES group by default).
// A database can also be added to only one failover or replication group, and these tests put the shared test database into groups.
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testint

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The complete create needs Business Critical Edition. A secondary group needs the secondary test account.
func TestInt_FailoverGroups(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	secondaryClient := testSecondaryClient(t)
	secondaryCtx := testSecondaryContext(t)
	sourceAccountID := testClientHelper().Account.GetAccountIdentifier(t)
	secondaryAccountID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	create := func(t *testing.T, request *sdk.CreateFailoverGroupRequest) {
		t.Helper()
		require.NoError(t, client.FailoverGroups.Create(ctx, request))
		t.Cleanup(func() {
			testClientHelper().FailoverGroup.Drop(t, request.GetName())
		})
	}

	createPrimary := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		create(t, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}))
		return id
	}

	createSecondaryOf := func(t *testing.T, primaryID sdk.AccountObjectIdentifier) sdk.AccountObjectIdentifier {
		t.Helper()
		primary := sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID)
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, primary)
		require.NoError(t, secondaryClient.FailoverGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryFailoverGroupRequest(secondaryID, primary)))
		t.Cleanup(func() {
			secondaryTestClientHelper().FailoverGroup.Drop(t, secondaryID)
		})
		return secondaryID
	}

	// OBJECT_TYPES and ALLOWED_ACCOUNTS are required. One database satisfies the allowed-database list.
	t.Run("create: basic", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		create(t, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}))

		group, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Contains(t, group.AllowedAccounts, secondaryAccountID)
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasNoComment().
			HasObjectTypes(sdk.FailoverGroupObjectTypeDatabases).
			HasNoAllowedIntegrationTypes().
			HasNoReplicationSchedule().
			HasNoSecondaryState().
			HasIsListingAutoFulfillmentGroup(false).
			HasIsOptimizedRefreshEnabled(false).
			HasRpoAssurance(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasAccountNameNotEmpty().
			HasAccountLocatorNotEmpty().
			HasOrganizationNameNotEmpty().
			HasSnowflakeRegionNotEmpty().
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccountID, id)),
		)

		databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		require.Equal(t, []sdk.AccountObjectIdentifier{database.ID()}, databases)
	})

	// IF NOT EXISTS is the only existence clause. CREATE has no OR REPLACE.
	// Every optional that can share this statement is set. A second IF NOT EXISTS create leaves the first group in place.
	t.Run("create: complete", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		share, cleanupShare := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupShare)
		externalVolume, cleanupVolume := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanupVolume)
		errorIntegration, cleanupIntegration := testClientHelper().NotificationIntegration.Create(t)
		t.Cleanup(cleanupIntegration)
		tag, cleanupTag := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(cleanupTag)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeDatabases,
			sdk.FailoverGroupObjectTypeExternalVolumes,
			sdk.FailoverGroupObjectTypeShares,
			sdk.FailoverGroupObjectTypeIntegrations,
			sdk.FailoverGroupObjectTypeRoles,
		}
		allowedAccounts := []sdk.AccountIdentifier{sourceAccountID, secondaryAccountID}
		integrationTypes := slices.Clone(sdk.AllFailoverGroupIntegrationTypes)

		create(t, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
			WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}).
			WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}).
			WithAllowedIntegrationTypes(integrationTypes).
			WithIgnoreEditionCheck(true).
			WithReplicationSchedule("10 MINUTE").
			WithOptimizedRefresh(true).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}).
			WithErrorIntegration(errorIntegration.ID()),
		)

		require.NoError(t, client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
			WithReplicationSchedule("60 MINUTE")))

		group, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.ElementsMatch(t, objectTypes, group.ObjectTypes)
		require.ElementsMatch(t, integrationTypes, group.AllowedIntegrationTypes)
		require.Subset(t, group.AllowedAccounts, allowedAccounts)
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasNoComment().
			HasReplicationSchedule("10 MINUTE").
			HasNoSecondaryState().
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false).
			HasOwner(snowflakeroles.Accountadmin.Name()),
		)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeFailoverGroup, "v1")

		databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		require.Equal(t, []sdk.AccountObjectIdentifier{database.ID()}, databases)
		shares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		require.Equal(t, []string{share.ID().Name()}, collections.Map(shares, sdk.AccountObjectIdentifier.Name))
	})

	t.Run("create: share with a dot in the name", func(t *testing.T) {
		shareId := testClientHelper().Ids.RandomAccountObjectIdentifierContaining(".")
		shareWithDot, cleanupShare := testClientHelper().Share.CreateWithRequest(t, shareId, sdk.NewCreateShareRequest(shareId))
		t.Cleanup(cleanupShare)
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		create(t, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeShares},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithAllowedShares([]sdk.AccountObjectIdentifier{shareWithDot.ID()}).WithIgnoreEditionCheck(true))

		shares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		require.Equal(t, []string{shareWithDot.ID().Name()}, collections.Map(shares, sdk.AccountObjectIdentifier.Name))
	})

	// Snowflake validates the statement before it checks existence, so the repeated create must be valid on its own: ROLES needs no allowed list.
	t.Run("create: if not exists", func(t *testing.T) {
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeRoles},
			[]sdk.AccountIdentifier{secondaryAccountID},
		).WithIfNotExists(true).WithReplicationSchedule("20 MINUTE")))

		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasObjectTypes(sdk.FailoverGroupObjectTypeDatabases).
			HasNoReplicationSchedule(),
		)
	})

	t.Run("create: already exists", func(t *testing.T) {
		id := createPrimary(t)

		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(
			id,
			[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeRoles},
			[]sdk.AccountIdentifier{secondaryAccountID},
		))
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("create secondary: basic", func(t *testing.T) {
		primaryID := createPrimary(t)
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		primary := sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID)
		secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, primary)

		require.NoError(t, secondaryClient.FailoverGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryFailoverGroupRequest(secondaryID, primary)))
		t.Cleanup(func() {
			secondaryTestClientHelper().FailoverGroup.Drop(t, secondaryID)
		})

		assertBasicSecondaryFailoverGroup(t, secondaryID, primary)
	})

	// IF NOT EXISTS is the only optional on a secondary create. A second call leaves the first replica in place.
	t.Run("create secondary: all", func(t *testing.T) {
		primaryID := createPrimary(t)
		secondaryID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		primary := sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID)
		secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, primary)
		request := sdk.NewCreateSecondaryFailoverGroupRequest(secondaryID, primary).WithIfNotExists(true)

		require.NoError(t, secondaryClient.FailoverGroups.CreateSecondary(secondaryCtx, request))
		t.Cleanup(func() {
			secondaryTestClientHelper().FailoverGroup.Drop(t, secondaryID)
		})
		require.NoError(t, secondaryClient.FailoverGroups.CreateSecondary(secondaryCtx, request))

		assertBasicSecondaryFailoverGroup(t, secondaryID, primary)
	})

	t.Run("alter: non-existing failover group", func(t *testing.T) {
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(NonExistingAccountObjectIdentifier).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithComment("x")))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: if exists on a non-existing failover group", func(t *testing.T) {
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(NonExistingAccountObjectIdentifier).
			WithIfExists(true).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithComment("x")))
		require.NoError(t, err)
	})

	// OBJECT_TYPES and the allowed lists are one SET.
	t.Run("alter: rename", func(t *testing.T) {
		id := createPrimary(t)
		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).WithRenameTo(newID)))
		t.Cleanup(func() {
			testClientHelper().FailoverGroup.Drop(t, newID)
		})

		_, err := client.FailoverGroups.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assertThatObject(t, objectassert.FailoverGroup(t, newID).
			HasName(newID.Name()).
			HasIsPrimary(true).
			HasObjectTypes(sdk.FailoverGroupObjectTypeDatabases),
		)
	})

	t.Run("alter: set object types", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		share, cleanupShare := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupShare)
		externalVolume, cleanupVolume := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanupVolume)
		id := createPrimary(t)
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeDatabases,
			sdk.FailoverGroupObjectTypeExternalVolumes,
			sdk.FailoverGroupObjectTypeShares,
			sdk.FailoverGroupObjectTypeIntegrations,
			sdk.FailoverGroupObjectTypeRoles,
		}
		integrationTypes := slices.Clone(sdk.AllFailoverGroupIntegrationTypes)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(objectTypes).
				WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
				WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}).
				WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}).
				WithAllowedIntegrationTypes(integrationTypes)),
		))

		group, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.ElementsMatch(t, objectTypes, group.ObjectTypes)
		require.ElementsMatch(t, integrationTypes, group.AllowedIntegrationTypes)

		narrowed := []sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeIntegrations}
		narrowedIntegrationTypes := []sdk.FailoverGroupIntegrationType{sdk.FailoverGroupIntegrationTypeApiIntegrations}
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(narrowed).
				WithAllowedIntegrationTypes(narrowedIntegrationTypes)),
		))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasObjectTypes(narrowed...).
			HasAllowedIntegrationTypes(narrowedIntegrationTypes...),
		)
	})

	// The schedule cannot be set in the same statement as OBJECT_TYPES and ALLOWED_INTEGRATION_TYPES (see GH #2544).
	t.Run("alter: object types and schedule in separate statements", func(t *testing.T) {
		id := createPrimary(t)
		objectTypes := []sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeIntegrations}
		integrationTypes := []sdk.FailoverGroupIntegrationType{sdk.FailoverGroupIntegrationTypeApiIntegrations}

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(objectTypes).
				WithAllowedIntegrationTypes(integrationTypes).
				WithReplicationSchedule("20 MINUTE")))
		require.ErrorContains(t, err, "unexpected 'REPLICATION_SCHEDULE'")

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(objectTypes).
				WithAllowedIntegrationTypes(integrationTypes))))
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule("20 MINUTE"))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasObjectTypes(objectTypes...).
			HasAllowedIntegrationTypes(integrationTypes...).
			HasReplicationSchedule("20 MINUTE"),
		)
	})

	// COMMENT, the schedule, and the error integration are set and unset together.
	t.Run("alter: set and unset properties", func(t *testing.T) {
		errorIntegration, cleanupIntegration := testClientHelper().NotificationIntegration.Create(t)
		t.Cleanup(cleanupIntegration)
		id := createPrimary(t)
		comment := random.Comment()

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithComment(comment).
				WithReplicationSchedule("USING CRON 0 0 10-20 * TUE,THU UTC").
				WithErrorIntegration(errorIntegration.ID())),
		))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasComment(comment).
			HasReplicationSchedule("USING CRON 0 0 10-20 * TUE,THU UTC"),
		)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().
				WithComment(true).
				WithReplicationSchedule(true).
				WithErrorIntegration(true)),
		))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasNoComment().
			HasNoReplicationSchedule(),
		)
	})

	t.Run("alter: set and unset tags", func(t *testing.T) {
		tag, cleanupTag := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(cleanupTag)
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSetTags([]sdk.TagAssociation{{Name: tag.ID(), Value: "v2"}}),
		))
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeFailoverGroup, "v2")

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}),
		))
		assertTagUnset(t, tag.ID(), id, sdk.ObjectTypeFailoverGroup)
	})

	// OPTIMIZED_REFRESH needs a schedule and conflicts with RPO_ASSURANCE, which needs the schedule to be unset.
	t.Run("alter: optimized refresh and rpo assurance", func(t *testing.T) {
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule("10 MINUTE").WithOptimizedRefresh(true))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasReplicationSchedule("10 MINUTE").
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false),
		)

		// Unsetting both in one statement fails because OPTIMIZED_REFRESH requires a schedule.
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER FAILOVER GROUP %s UNSET REPLICATION_SCHEDULE, OPTIMIZED_REFRESH`, id.FullyQualifiedName()))
		require.ErrorContains(t, err, "Cannot set OPTIMIZED_REFRESH on Failover group without REPLICATION_SCHEDULE")

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithOptimizedRefresh(true))))
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithReplicationSchedule(true))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasNoReplicationSchedule().
			HasIsOptimizedRefreshEnabled(false),
		)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).HasRpoAssurance(true))

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithRpoAssurance(true))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).HasRpoAssurance(false))
	})

	t.Run("alter: add and remove databases, shares and external volumes", func(t *testing.T) {
		added, cleanupAdded := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupAdded)
		share, cleanupShare := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupShare)
		externalVolume, cleanupVolume := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanupVolume)
		id := createPrimary(t)
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes([]sdk.FailoverGroupObjectType{
					sdk.FailoverGroupObjectTypeDatabases,
					sdk.FailoverGroupObjectTypeShares,
					sdk.FailoverGroupObjectTypeExternalVolumes,
				}))))

		// Each ADD and REMOVE takes exactly one list.
		for _, add := range []*sdk.FailoverGroupAddRequest{
			sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{added.ID()}),
			sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}),
			sdk.NewFailoverGroupAddRequest().WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}),
		} {
			require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).WithAdd(*add)))
		}
		databases, err := client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		require.Contains(t, databases, added.ID())
		shares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		require.Equal(t, []string{share.ID().Name()}, collections.Map(shares, sdk.AccountObjectIdentifier.Name))

		for _, remove := range []*sdk.FailoverGroupRemoveRequest{
			sdk.NewFailoverGroupRemoveRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{added.ID()}),
			sdk.NewFailoverGroupRemoveRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}),
			sdk.NewFailoverGroupRemoveRequest().WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}),
		} {
			require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).WithRemove(*remove)))
		}
		databases, err = client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		require.NotContains(t, databases, added.ID())
		shares, err = client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		require.Empty(t, shares)
	})

	t.Run("alter: add and remove allowed accounts", func(t *testing.T) {
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID}))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).HasAllowedAccounts(sourceAccountID))

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithAdd(*sdk.NewFailoverGroupAddRequest().
				WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID}).
				WithIgnoreEditionCheck(true))))
		group, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		require.ElementsMatch(t, []sdk.AccountIdentifier{sourceAccountID, secondaryAccountID}, group.AllowedAccounts)
	})

	t.Run("alter: move databases, shares and external volumes", func(t *testing.T) {
		database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)
		share, cleanupShare := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupShare)
		externalVolume, cleanupVolume := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanupVolume)
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeDatabases,
			sdk.FailoverGroupObjectTypeShares,
			sdk.FailoverGroupObjectTypeExternalVolumes,
		}
		source := testClientHelper().Ids.RandomAccountObjectIdentifier()
		create(t, sdk.NewCreateFailoverGroupRequest(source, objectTypes, []sdk.AccountIdentifier{secondaryAccountID}).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
			WithAllowedShares([]sdk.AccountObjectIdentifier{share.ID()}).
			WithAllowedExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}))
		target := createPrimary(t)
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(target).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes(objectTypes))))

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(source).
			WithMove(*sdk.NewFailoverGroupMoveRequest(target).
				WithDatabases([]sdk.AccountObjectIdentifier{database.ID()}).
				WithShares([]sdk.AccountObjectIdentifier{share.ID()}).
				WithExternalVolumes([]sdk.AccountObjectIdentifier{externalVolume}))))

		targetDatabases, err := client.FailoverGroups.ShowDatabases(ctx, target)
		require.NoError(t, err)
		require.Contains(t, targetDatabases, database.ID())
		targetShares, err := client.FailoverGroups.ShowShares(ctx, target)
		require.NoError(t, err)
		require.Equal(t, []string{share.ID().Name()}, collections.Map(targetShares, sdk.AccountObjectIdentifier.Name))
		sourceDatabases, err := client.FailoverGroups.ShowDatabases(ctx, source)
		require.NoError(t, err)
		require.NotContains(t, sourceDatabases, database.ID())
		sourceShares, err := client.FailoverGroups.ShowShares(ctx, source)
		require.NoError(t, err)
		require.Empty(t, sourceShares)
	})

	// Snowflake accepts SET and UNSET READ_ONLY on a failover group, but SHOW keeps read_only null.
	t.Run("alter: read only", func(t *testing.T) {
		id := createPrimary(t)

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReadOnly(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.FailoverGroup(t, id).HasNoReadOnly())

		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithReadOnly(true))))
		assertThatObject(t, objectassert.FailoverGroup(t, id).HasNoReadOnly())
	})

	// The secondary state is set only when the primary has a replication schedule.
	t.Run("alter secondary: suspend, resume, and refresh", func(t *testing.T) {
		primaryID := createPrimary(t)
		require.NoError(t, client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(primaryID).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule("10 MINUTE"))))
		secondaryID := createSecondaryOf(t, primaryID)

		require.NoError(t, secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
			WithSuspend(*sdk.NewFailoverGroupSecondarySuspendRequest())))
		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
			HasSecondaryState(sdk.FailoverGroupSecondaryStateSuspended))

		require.NoError(t, secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
			WithResume(true)))
		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
			HasSecondaryState(sdk.FailoverGroupSecondaryStateStarted))

		require.NoError(t, secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
			WithSuspend(*sdk.NewFailoverGroupSecondarySuspendRequest().WithImmediate(true))))
		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
			HasSecondaryState(sdk.FailoverGroupSecondaryStateSuspended))

		require.NoError(t, secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
			WithResume(true)))
		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
			HasSecondaryState(sdk.FailoverGroupSecondaryStateStarted))

		// TODO [SNOW-1348343]: use INFORMATION_SCHEMA.REPLICATION_GROUP_REFRESH_PROGRESS (END_TIME = NULL) to know when no refresh is going on.
		// A refresh can already be running (003101: Only one refresh statement can execute at a time). Retry until this one is accepted.
		assert.Eventually(t, func() bool {
			err := secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
				WithRefresh(true))
			if err != nil {
				t.Logf("Retrying: %v", err)
			}
			return err == nil
		}, 2*time.Minute, time.Second)
	})

	t.Run("alter secondary: promote to primary", func(t *testing.T) {
		primaryID := createPrimary(t)
		secondaryID := createSecondaryOf(t, primaryID)
		// After the promotion the replica is the primary and cannot be dropped before the former primary (003913).
		t.Cleanup(func() {
			testClientHelper().FailoverGroup.Drop(t, primaryID)
		})

		// A refresh in progress blocks the promotion (003122). Retry until it is accepted.
		assert.Eventually(t, func() bool {
			err := secondaryClient.FailoverGroups.AlterSecondary(secondaryCtx, sdk.NewAlterSecondaryFailoverGroupRequest(secondaryID).
				WithPrimary(true))
			if err != nil {
				t.Logf("Retrying: %v", err)
			}
			return err == nil
		}, 2*time.Minute, time.Second)

		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
			HasIsPrimary(true).
			HasPrimary(sdk.NewExternalObjectIdentifier(secondaryAccountID, secondaryID)))
	})

	t.Run("alter secondary: non-existing failover group", func(t *testing.T) {
		err := client.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(NonExistingAccountObjectIdentifier).WithResume(true))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: lists the failover group", func(t *testing.T) {
		id := createPrimary(t)

		shown, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
		require.NoError(t, err)
		require.True(t, slices.ContainsFunc(shown, func(group sdk.FailoverGroup) bool {
			return group.Name == id.Name() && group.IsPrimary
		}))
	})

	t.Run("show: in account", func(t *testing.T) {
		id := createPrimary(t)

		inSource, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest().WithInAccount(sourceAccountID))
		require.NoError(t, err)
		require.True(t, slices.ContainsFunc(inSource, func(group sdk.FailoverGroup) bool {
			return group.Name == id.Name() && group.IsPrimary
		}))
	})

	t.Run("show by id: existing", func(t *testing.T) {
		id := createPrimary(t)

		assertThatObject(t, objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasIsPrimary(true),
		)

		got, err := client.FailoverGroups.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), got.Name)
	})

	t.Run("show by id: non-existing", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		_, err = client.FailoverGroups.ShowByIDSafely(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show by id: same name in two accounts", func(t *testing.T) {
		primaryID := createPrimary(t)
		primary := sdk.NewExternalObjectIdentifier(sourceAccountID, primaryID)
		secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, primary)
		require.NoError(t, secondaryClient.FailoverGroups.CreateSecondary(secondaryCtx, sdk.NewCreateSecondaryFailoverGroupRequest(primaryID, primary)))
		t.Cleanup(func() {
			secondaryTestClientHelper().FailoverGroup.Drop(t, primaryID)
		})

		assertThatObject(t, objectassert.FailoverGroup(t, primaryID).HasIsPrimary(true))
		assertThatObjectSecondary(t, objectassert.FailoverGroup(t, primaryID).HasIsPrimary(false))
	})

	t.Run("show databases: non-existing failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowDatabases(ctx, NonExistingAccountObjectIdentifier)
		require.Error(t, err)
	})

	t.Run("show shares: non-existing failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowShares(ctx, NonExistingAccountObjectIdentifier)
		require.Error(t, err)
	})

	t.Run("show shares: no allowed shares", func(t *testing.T) {
		id := createPrimary(t)

		shares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		require.Empty(t, shares)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(id)))
		_, err := client.FailoverGroups.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(NonExistingAccountObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := createPrimary(t)

		require.NoError(t, client.FailoverGroups.DropSafely(ctx, id))
		_, err := client.FailoverGroups.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		require.NoError(t, client.FailoverGroups.DropSafely(ctx, NonExistingAccountObjectIdentifier))
	})
}

// SHOW fills object_types, allowed_integration_types, and allowed_accounts on the primary row only.
func assertBasicSecondaryFailoverGroup(t *testing.T, secondaryID sdk.AccountObjectIdentifier, primary sdk.ExternalObjectIdentifier) {
	t.Helper()
	assertThatObjectSecondary(t, objectassert.FailoverGroup(t, secondaryID).
		HasCreatedOnNotEmpty().
		HasName(secondaryID.Name()).
		HasType("FAILOVER").
		HasIsPrimary(false).
		HasPrimary(primary).
		HasNoComment().
		HasNoObjectTypes().
		HasNoAllowedIntegrationTypes().
		HasNoAllowedAccounts().
		HasNoReplicationSchedule().
		HasNoSecondaryState().
		HasIsListingAutoFulfillmentGroup(false).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasAccountNameNotEmpty().
		HasAccountLocatorNotEmpty().
		HasOrganizationNameNotEmpty().
		HasSnowflakeRegionNotEmpty(),
	)
}
