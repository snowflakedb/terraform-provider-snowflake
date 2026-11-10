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

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FailoverGroupsCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := testClientHelper().Share.Create(t)
	t.Cleanup(shareCleanup)

	secondaryAccountId := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	t.Run("test complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeShares,
			sdk.FailoverGroupObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			secondaryAccountId,
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()}).
			WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()}).
			WithIgnoreEditionCheck(true).
			WithReplicationSchedule(replicationSchedule).
			WithOptimizedRefresh(true))
		require.NoError(t, err)
		failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(id))
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
		assert.Equal(t, id.Name(), failoverGroup.Name)
		slices.Sort(objectTypes)
		slices.Sort(failoverGroup.ObjectTypes)
		assert.Equal(t, objectTypes, failoverGroup.ObjectTypes)
		assert.Empty(t, failoverGroup.AllowedIntegrationTypes)
		// this is length 2 because it automatically adds the current account to allowed accounts list
		assert.Len(t, failoverGroup.AllowedAccounts, 2)
		for _, allowedAccount := range allowedAccounts {
			assert.Contains(t, failoverGroup.AllowedAccounts, allowedAccount)
		}
		assert.Equal(t, new(replicationSchedule), failoverGroup.ReplicationSchedule)
		assert.True(t, failoverGroup.IsOptimizedRefreshEnabled)
		assert.False(t, failoverGroup.RpoAssurance)

		fgDBS, err := client.FailoverGroups.ShowDatabases(ctx, id)
		require.NoError(t, err)
		assert.Len(t, fgDBS, 1)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), fgDBS[0].Name())

		fgShares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		assert.Len(t, fgShares, 1)
		assert.Equal(t, shareTest.ID().Name(), fgShares[0].Name())
	})

	t.Run("test with identifier containing a dot", func(t *testing.T) {
		shareId := testClientHelper().Ids.RandomAccountObjectIdentifierContaining(".")

		shareWithDot, shareWithDotCleanup := testClientHelper().Share.CreateWithRequest(t, shareId, sdk.NewCreateShareRequest(shareId))
		t.Cleanup(shareWithDotCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeShares,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			secondaryAccountId,
		}
		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithAllowedShares([]sdk.AccountObjectIdentifier{shareWithDot.ID()}).
			WithIgnoreEditionCheck(true))
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(id))
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)

		fgShares, err := client.FailoverGroups.ShowShares(ctx, id)
		require.NoError(t, err)
		assert.Len(t, fgShares, 1)
		assert.Equal(t, shareWithDot.ID().Name(), fgShares[0].Name())
	})

	t.Run("test with allowed integration types", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeIntegrations,
			sdk.FailoverGroupObjectTypeRoles,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			secondaryAccountId,
		}
		allowedIntegrationTypes := []sdk.FailoverGroupIntegrationType{
			sdk.FailoverGroupIntegrationTypeSecurityIntegrations,
			sdk.FailoverGroupIntegrationTypeApiIntegrations,
			sdk.FailoverGroupIntegrationTypeStorageIntegrations,
			sdk.FailoverGroupIntegrationTypeExternalAccessIntegrations,
			sdk.FailoverGroupIntegrationTypeNotificationIntegrations,
		}
		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithAllowedIntegrationTypes(allowedIntegrationTypes))
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(id))
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)
		failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), failoverGroup.Name)
		slices.Sort(failoverGroup.AllowedIntegrationTypes)
		slices.Sort(allowedIntegrationTypes)
		assert.Equal(t, allowedIntegrationTypes, failoverGroup.AllowedIntegrationTypes)
	})
}

func TestInt_Issue2544(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	secondaryAccountId := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	t.Run("alter object types, replication schedule, and allowed integration types at the same time", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeIntegrations,
			sdk.FailoverGroupObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			secondaryAccountId,
		}
		allowedIntegrationTypes := []sdk.FailoverGroupIntegrationType{
			sdk.FailoverGroupIntegrationTypeApiIntegrations,
			sdk.FailoverGroupIntegrationTypeNotificationIntegrations,
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()}).
			WithAllowedIntegrationTypes(allowedIntegrationTypes).
			WithReplicationSchedule(replicationSchedule))
		require.NoError(t, err)
		cleanupFailoverGroup := func() {
			err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(id))
			require.NoError(t, err)
		}
		t.Cleanup(cleanupFailoverGroup)

		newObjectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeIntegrations,
		}
		newAllowedIntegrationTypes := []sdk.FailoverGroupIntegrationType{
			sdk.FailoverGroupIntegrationTypeApiIntegrations,
		}
		newReplicationSchedule := "20 MINUTE"

		// does not work together:
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(newObjectTypes).
				WithAllowedIntegrationTypes(newAllowedIntegrationTypes).
				WithReplicationSchedule(newReplicationSchedule)))
		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected 'REPLICATION_SCHEDULE'")

		// works as two separate alters:
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(newObjectTypes).
				WithAllowedIntegrationTypes(newAllowedIntegrationTypes)))
		require.NoError(t, err)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithReplicationSchedule(newReplicationSchedule)))
		require.NoError(t, err)
	})
}

func TestInt_CreateSecondaryReplicationGroup(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	primaryAccountID := testClientHelper().Account.GetAccountIdentifier(t)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	shareTest, cleanupDatabase := testClientHelper().Share.Create(t)
	t.Cleanup(cleanupDatabase)

	id := testClientHelper().Ids.RandomAccountObjectIdentifier()

	allowedAccounts := []sdk.AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []sdk.FailoverGroupObjectType{
		sdk.FailoverGroupObjectTypeShares,
	}
	err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
		WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()}).
		WithReplicationSchedule("10 MINUTE"))
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, failoverGroup.ID())

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondary(ctx,
		sdk.NewCreateSecondaryFailoverGroupRequest(failoverGroup.ID(), failoverGroup.ExternalID()).
			WithIfNotExists(true))
	require.NoError(t, err)

	// cleanup failover groups with retry (in case of replication delay)
	cleanupFailoverGroups := func() {
		secondaryClientFailoverGroupDropped := func() bool {
			return secondaryClient.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID())) == nil
		}
		assert.Eventually(t, secondaryClientFailoverGroupDropped, 10*time.Second, time.Second)
		failoverGroupDropped := func() bool {
			return client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID())) == nil
		}
		assert.Eventually(t, failoverGroupDropped, 10*time.Second, time.Second)
	}
	t.Cleanup(cleanupFailoverGroups)

	// SHOW is not filtered, so other failover groups may be listed; the primary and the replica of this group are expected.
	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
	require.NoError(t, err)
	assert.Len(t, slices.DeleteFunc(failoverGroups, func(fg sdk.FailoverGroup) bool { return fg.ID() != failoverGroup.ID() }), 2)
}

func TestInt_FailoverGroupsAlterSource(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	t.Run("rename the failover group", func(t *testing.T) {
		failoverGroup, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(failoverGroupCleanup)
		oldID := failoverGroup.ID()
		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(oldID).WithRenameTo(newID))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), failoverGroup.Name)
		t.Cleanup(testClientHelper().FailoverGroup.DropFunc(t, newID))
	})

	t.Run("reset the list of specified object types enabled for replication and failover.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)
		objectTypes := []sdk.FailoverGroupObjectType{
			sdk.FailoverGroupObjectTypeDatabases,
		}
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes(objectTypes)))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, objectTypes, failoverGroup.ObjectTypes)
	})

	t.Run("set and unset", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithReplicationSchedule("10 MINUTE").
				WithOptimizedRefresh(true)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new("10 MINUTE"), failoverGroup.ReplicationSchedule)
		assert.True(t, failoverGroup.IsOptimizedRefreshEnabled)
		assert.False(t, failoverGroup.RpoAssurance)

		// UNSET OPTIMIZED_REFRESH first: unsetting it together with REPLICATION_SCHEDULE
		// fails because OPTIMIZED_REFRESH requires a schedule.
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().
				WithOptimizedRefresh(true)))
		require.NoError(t, err)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().
				WithReplicationSchedule(true)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		require.Empty(t, failoverGroup.ReplicationSchedule)
		assert.False(t, failoverGroup.IsOptimizedRefreshEnabled)

		// RPO_ASSURANCE cannot be set with REPLICATION_SCHEDULE; enable after unset.
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithRpoAssurance(true)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.True(t, failoverGroup.RpoAssurance)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithRpoAssurance(true)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.False(t, failoverGroup.RpoAssurance)
	})

	// This tests checks if unsetting multiple values is possible (docs suggest the comma syntax).
	// The errors are returned so SDK has the exactly one field validation on UNSET.
	t.Run("cannot unset multiple properties in one statement", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)
		id := failoverGroup.ID()

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithReplicationSchedule("10 MINUTE").
				WithOptimizedRefresh(true)))
		require.NoError(t, err)

		_, err = client.ExecForTests(ctx, fmt.Sprintf(`ALTER FAILOVER GROUP %s UNSET REPLICATION_SCHEDULE, OPTIMIZED_REFRESH`, id.FullyQualifiedName()))
		require.ErrorContains(t, err, "Cannot set OPTIMIZED_REFRESH on Failover group without REPLICATION_SCHEDULE")

		_, err = client.ExecForTests(ctx, fmt.Sprintf(`ALTER FAILOVER GROUP %s UNSET REPLICATION_SCHEDULE, OPTIMIZED_REFRESH, RPO_ASSURANCE`, id.FullyQualifiedName()))
		require.ErrorContains(t, err, "Cannot set OPTIMIZED_REFRESH on Failover group without REPLICATION_SCHEDULE")
	})

	t.Run("set or update the replication schedule for automatic refresh of secondary failover groups.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)
		replicationSchedule := "USING CRON 0 0 10-20 * TUE,THU UTC"

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule(replicationSchedule)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new(replicationSchedule), failoverGroup.ReplicationSchedule)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithReplicationSchedule(true)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		require.Empty(t, failoverGroup.ReplicationSchedule)
	})

	t.Run("add and remove database account object", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// first add databases to allowed object types
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, sdk.FailoverGroupObjectTypeDatabases, failoverGroup.ObjectTypes[0])

		// now add database to allowed databases
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
		require.NoError(t, err)
		allowedDBs, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, allowedDBs, 1)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), allowedDBs[0].Name())

		// now remove database from allowed databases
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
		require.NoError(t, err)
		allowedDBs, err = client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, allowedDBs)
	})

	t.Run("add and remove share account object", func(t *testing.T) {
		shareTest, cleanupDatabase := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupDatabase)
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// first add shares to allowed object types
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeShares})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, sdk.FailoverGroupObjectTypeShares, failoverGroup.ObjectTypes[0])

		// now add share to allowed shares
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)
		allowedShares, err := client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, allowedShares, 1)
		assert.Equal(t, shareTest.ID().Name(), allowedShares[0].Name())

		// now remove share from allowed shares
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)
		allowedShares, err = client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, allowedShares)
	})

	t.Run("add and remove security integration account object", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// first add security integrations to allowed object types
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeIntegrations}).
				WithAllowedIntegrationTypes([]sdk.FailoverGroupIntegrationType{
					sdk.FailoverGroupIntegrationTypeApiIntegrations,
					sdk.FailoverGroupIntegrationTypeNotificationIntegrations,
				})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedIntegrationTypes, 2)
		assert.Equal(t, sdk.FailoverGroupIntegrationTypeApiIntegrations, failoverGroup.AllowedIntegrationTypes[0])
		assert.Equal(t, sdk.FailoverGroupIntegrationTypeNotificationIntegrations, failoverGroup.AllowedIntegrationTypes[1])
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, sdk.FailoverGroupObjectTypeIntegrations, failoverGroup.ObjectTypes[0])

		// now remove security integration from allowed security integrations
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeIntegrations}).
				WithAllowedIntegrationTypes([]sdk.FailoverGroupIntegrationType{sdk.FailoverGroupIntegrationTypeApiIntegrations})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedIntegrationTypes, 1)
		assert.Equal(t, sdk.FailoverGroupIntegrationTypeApiIntegrations, failoverGroup.AllowedIntegrationTypes[0])
	})

	t.Run("add or remove target accounts enabled for replication and failover", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		secondaryAccountID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		// first add target account
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedAccounts, 2)
		assert.Contains(t, failoverGroup.AllowedAccounts, secondaryAccountID)

		// now remove target accounts
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedAccounts, 1)
		assert.Contains(t, failoverGroup.AllowedAccounts, testClientHelper().Account.GetAccountIdentifier(t))
	})

	// The page spells the ADD accounts option IGNORE EDITION CHECK; the handwritten SDK used IGNORE_EDITION_CHECK.
	t.Run("add target accounts with ignore edition check", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		secondaryAccountID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID}).WithIgnoreEditionCheck(true)))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Contains(t, failoverGroup.AllowedAccounts, secondaryAccountID)
	})

	t.Run("set and unset comment", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithComment("fg comment")))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new("fg comment"), failoverGroup.Comment)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, failoverGroup.Comment)
	})

	t.Run("move shares to another failover group", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// add "SHARES" to object types of both failover groups
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeShares})))
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup2.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeShares})))
		require.NoError(t, err)

		shareTest, cleanupShare := testClientHelper().Share.Create(t)
		t.Cleanup(cleanupShare)

		// now add share to allowed shares of failover group 1
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)

		// now move share to failover group 2
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithMove(*sdk.NewFailoverGroupMoveRequest(failoverGroup2.ID()).WithShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)

		// verify that share is now in failover group 2
		shares, err := client.FailoverGroups.ShowShares(ctx, failoverGroup2.ID())
		require.NoError(t, err)
		assert.Len(t, shares, 1)

		// verify that share is not in failover group 1
		shares, err = client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, shares)
	})

	t.Run("move database to another failover group", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// add "DATABASES" to object types of both failover groups
		err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases})))
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup2.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases})))
		require.NoError(t, err)

		databaseTest, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)

		// now add database to allowed databases of failover group 1
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseTest.ID()})))
		require.NoError(t, err)

		// now move database to failover group 2
		err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroup.ID()).
			WithMove(*sdk.NewFailoverGroupMoveRequest(failoverGroup2.ID()).WithDatabases([]sdk.AccountObjectIdentifier{databaseTest.ID()})))
		require.NoError(t, err)

		// verify that database is now in failover group 2
		databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup2.ID())
		require.NoError(t, err)
		assert.Len(t, databases, 1)

		// verify that database is not in failover group 1
		databases, err = client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, databases)
	})
}

func TestInt_FailoverGroupsAlterTarget(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	primaryAccountID := testClientHelper().Account.GetAccountIdentifier(t)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	databaseTest, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	id := testClientHelper().Ids.RandomAccountObjectIdentifier()

	allowedAccounts := []sdk.AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []sdk.FailoverGroupObjectType{
		sdk.FailoverGroupObjectTypeDatabases,
	}
	err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
		WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseTest.ID()}).
		WithReplicationSchedule("10 MINUTE"))
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	secondaryTestClientHelper().FailoverGroup.WaitForPrimaryToBeVisible(t, failoverGroup.ID())

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondary(ctx,
		sdk.NewCreateSecondaryFailoverGroupRequest(failoverGroup.ID(), failoverGroup.ExternalID()).
			WithIfNotExists(true))
	require.NoError(t, err)

	// cleanup failover groups with retry (in case of replication delay)
	cleanupFailoverGroups := func() {
		failoverGroupDropped := func() bool {
			return client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID())) == nil
		}
		assert.Eventually(t, failoverGroupDropped, 10*time.Second, time.Second)
		secondaryClientFailoverGroupDropped := func() bool {
			return secondaryClient.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID())) == nil
		}
		assert.Eventually(t, secondaryClientFailoverGroupDropped, 10*time.Second, time.Second)
	}
	t.Cleanup(cleanupFailoverGroups)

	// SHOW is not filtered, so other failover groups may be listed; the primary and the replica of this group are expected.
	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
	require.NoError(t, err)
	assert.Len(t, slices.DeleteFunc(failoverGroups, func(fg sdk.FailoverGroup) bool { return fg.ID() != failoverGroup.ID() }), 2)

	t.Run("perform suspend and resume", func(t *testing.T) {
		// suspend target failover group
		err = secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).WithSuspend(*sdk.NewFailoverGroupSecondarySuspendRequest()))
		require.NoError(t, err)

		// verify that target failover group is suspended
		fg, err := secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new(sdk.FailoverGroupSecondaryStateSuspended), fg.SecondaryState)

		// resume target failover group
		err = secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).WithResume(true))
		require.NoError(t, err)

		// verify that target failover group is resumed
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new(sdk.FailoverGroupSecondaryStateStarted), failoverGroup.SecondaryState)
	})

	t.Run("perform suspend immediate and resume", func(t *testing.T) {
		err = secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).
			WithSuspend(*sdk.NewFailoverGroupSecondarySuspendRequest().WithImmediate(true)))
		require.NoError(t, err)

		fg, err := secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new(sdk.FailoverGroupSecondaryStateSuspended), fg.SecondaryState)

		err = secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).WithResume(true))
		require.NoError(t, err)

		fg, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, new(sdk.FailoverGroupSecondaryStateStarted), fg.SecondaryState)
	})

	// TODO [SNOW-1348343]: Use:
	// SELECT phase_name, start_time, end_time, progress, details
	//  FROM TABLE(
	//  INFORMATION_SCHEMA.REPLICATION_GROUP_REFRESH_PROGRESS(
	//    '...'
	//  )
	// );
	// with the END_TIME = NULL
	// to know when there is no refresh going on
	t.Run("refresh target failover group", func(t *testing.T) {
		// Try a few times to avoid:
		// 003101 (55000): Replication group "..." is already being refreshed. Only one refresh statement can execute at a time.
		refresh := func() bool {
			return secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).WithRefresh(true)) == nil
		}
		assert.Eventually(t, refresh, 2*time.Minute, 1*time.Second)
	})

	t.Run("promote secondary to primary", func(t *testing.T) {
		// promote secondary to primary
		// Try a few times to avoid:
		// 003122 (55000): Replication group "..." cannot currently be set as primary because it is being refreshed. Either wait for the refresh to finish or cancel the refresh and try again.
		promote := func() bool {
			return secondaryClient.FailoverGroups.AlterSecondary(ctx, sdk.NewAlterSecondaryFailoverGroupRequest(failoverGroup.ID()).WithPrimary(true)) == nil
		}
		assert.Eventually(t, promote, 2*time.Minute, 1*time.Second)

		// verify that target failover group is promoted
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.True(t, failoverGroup.IsPrimary)
	})
}

func TestInt_FailoverGroupsDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	t.Run("no options", func(t *testing.T) {
		failoverGroup, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(failoverGroupCleanup)
		err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID()))
		require.NoError(t, err)
	})

	t.Run("with IfExists", func(t *testing.T) {
		failoverGroup, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(failoverGroupCleanup)
		err := client.FailoverGroups.Drop(ctx, sdk.NewDropFailoverGroupRequest(failoverGroup.ID()).WithIfExists(true))
		require.NoError(t, err)
	})
}

func TestInt_FailoverGroupsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	t.Run("without show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, *failoverGroupTest)
	})

	t.Run("with show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest().
			WithInAccount(testClientHelper().Ids.AccountIdentifierWithLocator()))
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, *failoverGroupTest)
	})

	t.Run("when searching a non-existent failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})
}

func TestInt_FailoverGroupsShowDatabases(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroupTest.ID()).
		WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases})))
	require.NoError(t, err)
	err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroupTest.ID()).
		WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
	require.NoError(t, err)
	databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Len(t, databases, 1)
	assert.Equal(t, testClientHelper().Ids.DatabaseId(), databases[0])
}

func TestInt_FailoverGroupsShowShares(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	shareTest, shareCleanup := testClientHelper().Share.Create(t)
	t.Cleanup(shareCleanup)
	err := client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroupTest.ID()).
		WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeShares})))
	require.NoError(t, err)
	err = client.FailoverGroups.Alter(ctx, sdk.NewAlterFailoverGroupRequest(failoverGroupTest.ID()).
		WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
	require.NoError(t, err)
	shares, err := client.FailoverGroups.ShowShares(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Len(t, shares, 1)
	assert.Equal(t, shareTest.ID(), shares[0])
}
