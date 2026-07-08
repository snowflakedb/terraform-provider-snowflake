//go:build non_account_level_tests

package testint

import (
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_FailoverGroupsCreate(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	t.Run("test complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeShares,
			sdk.PluralObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		replicationSchedule := "10 MINUTE"
		err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
			WithIfNotExists(true).
			WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()}).
			WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()}).
			WithIgnoreEditionCheck(true).
			WithReplicationSchedule(replicationSchedule))
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
		assert.Equal(t, replicationSchedule, failoverGroup.ReplicationSchedule)

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

		shareWithDot, shareWithDotCleanup := testClientHelper().Share.CreateShareWithRequest(t, shareId, sdk.NewCreateShareRequest(shareId))
		t.Cleanup(shareWithDotCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeShares,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
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
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
			sdk.PluralObjectTypeRoles,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		allowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeSecurityIntegrations,
			sdk.IntegrationTypeAPIIntegrations,
			sdk.IntegrationTypeStorageIntegrations,
			sdk.IntegrationTypeExternalAccessIntegrations,
			sdk.IntegrationTypeNotificationIntegrations,
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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)

	accountName := testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
	businessCriticalAccountId := sdk.NewAccountIdentifierFromFullyQualifiedName(accountName)

	t.Run("alter object types, replication schedule, and allowed integration types at the same time", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
			sdk.PluralObjectTypeDatabases,
		}
		allowedAccounts := []sdk.AccountIdentifier{
			businessCriticalAccountId,
		}
		allowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeAPIIntegrations,
			sdk.IntegrationTypeNotificationIntegrations,
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

		newObjectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeIntegrations,
		}
		newAllowedIntegrationTypes := []sdk.IntegrationType{
			sdk.IntegrationTypeAPIIntegrations,
		}
		newReplicationSchedule := "20 MINUTE"

		// does not work together:
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(newObjectTypes).
				WithAllowedIntegrationTypes(newAllowedIntegrationTypes).
				WithReplicationSchedule(newReplicationSchedule)))
		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected 'REPLICATION_SCHEDULE'")

		// works as two separate alters:
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes(newObjectTypes).
				WithAllowedIntegrationTypes(newAllowedIntegrationTypes)))
		require.NoError(t, err)

		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(id).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithReplicationSchedule(newReplicationSchedule)))
		require.NoError(t, err)
	})
}

func TestInt_CreateSecondaryReplicationGroup(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	primaryAccountID := testClientHelper().Account.GetAccountIdentifier(t)
	secondaryClient := testSecondaryClient(t)
	secondaryClientID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

	shareTest, cleanupDatabase := testClientHelper().Share.CreateShare(t)
	t.Cleanup(cleanupDatabase)

	id := testClientHelper().Ids.RandomAccountObjectIdentifier()

	allowedAccounts := []sdk.AccountIdentifier{
		primaryAccountID,
		secondaryClientID,
	}
	objectTypes := []sdk.PluralObjectType{
		sdk.PluralObjectTypeShares,
	}
	err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
		WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()}).
		WithReplicationSchedule("10 MINUTE"))
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx,
		sdk.NewCreateSecondaryReplicationGroupFailoverGroupRequest(failoverGroup.ID(), failoverGroup.ExternalID()).
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

	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
	require.NoError(t, err)
	assert.Len(t, failoverGroups, 2)
}

func TestInt_FailoverGroupsAlterSource(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	t.Run("rename the failover group", func(t *testing.T) {
		failoverGroup, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(failoverGroupCleanup)
		oldID := failoverGroup.ID()
		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(oldID).WithNewName(newID))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), failoverGroup.Name)
		t.Cleanup(testClientHelper().FailoverGroup.DropFunc(t, newID))
	})

	t.Run("reset the list of specified object types enabled for replication and failover.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)
		objectTypes := []sdk.PluralObjectType{
			sdk.PluralObjectTypeDatabases,
		}
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes(objectTypes)))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, objectTypes, failoverGroup.ObjectTypes)
	})

	t.Run("set or update the replication schedule for automatic refresh of secondary failover groups.", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)
		replicationSchedule := "USING CRON 0 0 10-20 * TUE,THU UTC"

		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule(replicationSchedule)))
		require.NoError(t, err)

		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, replicationSchedule, failoverGroup.ReplicationSchedule)

		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
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
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeDatabases})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, sdk.PluralObjectTypeDatabases, failoverGroup.ObjectTypes[0])

		// now add database to allowed databases
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
		require.NoError(t, err)
		allowedDBs, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, allowedDBs, 1)
		assert.Equal(t, testClientHelper().Ids.DatabaseId().Name(), allowedDBs[0].Name())

		// now remove database from allowed databases
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
		require.NoError(t, err)
		allowedDBs, err = client.FailoverGroups.ShowDatabases(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Empty(t, allowedDBs)
	})

	t.Run("add and remove share account object", func(t *testing.T) {
		shareTest, cleanupDatabase := testClientHelper().Share.CreateShare(t)
		t.Cleanup(cleanupDatabase)
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// first add shares to allowed object types
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeShares})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, shareTest.ObjectType().Plural(), failoverGroup.ObjectTypes[0])

		// now add share to allowed shares
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)
		allowedShares, err := client.FailoverGroups.ShowShares(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, allowedShares, 1)
		assert.Equal(t, shareTest.ID().Name(), allowedShares[0].Name())

		// now remove share from allowed shares
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
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
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeIntegrations}).
				WithAllowedIntegrationTypes([]sdk.IntegrationType{
					sdk.IntegrationTypeAPIIntegrations,
					sdk.IntegrationTypeNotificationIntegrations,
				})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedIntegrationTypes, 2)
		assert.Equal(t, sdk.IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
		assert.Equal(t, sdk.IntegrationTypeNotificationIntegrations, failoverGroup.AllowedIntegrationTypes[1])
		assert.Len(t, failoverGroup.ObjectTypes, 1)
		assert.Equal(t, sdk.PluralObjectTypeIntegrations, failoverGroup.ObjectTypes[0])

		// now remove security integration from allowed security integrations
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().
				WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeIntegrations}).
				WithAllowedIntegrationTypes([]sdk.IntegrationType{sdk.IntegrationTypeAPIIntegrations})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedIntegrationTypes, 1)
		assert.Equal(t, sdk.IntegrationTypeAPIIntegrations, failoverGroup.AllowedIntegrationTypes[0])
	})

	t.Run("add or remove target accounts enabled for replication and failover", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		secondaryAccountID := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		// first add target account
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedAccounts, 2)
		assert.Contains(t, failoverGroup.AllowedAccounts, secondaryAccountID)

		// now remove target accounts
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts([]sdk.AccountIdentifier{secondaryAccountID})))
		require.NoError(t, err)
		failoverGroup, err = client.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Len(t, failoverGroup.AllowedAccounts, 1)
		assert.Contains(t, failoverGroup.AllowedAccounts, testClientHelper().Account.GetAccountIdentifier(t))
	})

	t.Run("move shares to another failover group", func(t *testing.T) {
		failoverGroup, cleanupFailoverGroup := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup)

		// add "SHARES" to object types of both failover groups
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeShares})))
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup2.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeShares})))
		require.NoError(t, err)

		shareTest, cleanupShare := testClientHelper().Share.CreateShare(t)
		t.Cleanup(cleanupShare)

		// now add share to allowed shares of failover group 1
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
		require.NoError(t, err)

		// now move share to failover group 2
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
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
		err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeDatabases})))
		require.NoError(t, err)

		failoverGroup2, cleanupFailoverGroup2 := testClientHelper().FailoverGroup.Create(t)
		t.Cleanup(cleanupFailoverGroup2)

		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup2.ID()).
			WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeDatabases})))
		require.NoError(t, err)

		databaseTest, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cleanupDatabase)

		// now add database to allowed databases of failover group 1
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
			WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseTest.ID()})))
		require.NoError(t, err)

		// now move database to failover group 2
		err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroup.ID()).
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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

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
	objectTypes := []sdk.PluralObjectType{
		sdk.PluralObjectTypeDatabases,
	}
	err := client.FailoverGroups.Create(ctx, sdk.NewCreateFailoverGroupRequest(id, objectTypes, allowedAccounts).
		WithAllowedDatabases([]sdk.AccountObjectIdentifier{databaseTest.ID()}).
		WithReplicationSchedule("10 MINUTE"))
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)

	// there is a delay between creating a failover group and it being available for replication
	time.Sleep(1 * time.Second)

	// create a replica of failover group in target account
	err = secondaryClient.FailoverGroups.CreateSecondaryReplicationGroup(ctx,
		sdk.NewCreateSecondaryReplicationGroupFailoverGroupRequest(failoverGroup.ID(), failoverGroup.ExternalID()).
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

	failoverGroups, err := secondaryClient.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
	require.NoError(t, err)
	assert.Len(t, failoverGroups, 2)

	t.Run("perform suspend and resume", func(t *testing.T) {
		// suspend target failover group
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, sdk.NewAlterTargetFailoverGroupRequest(failoverGroup.ID()).WithSuspend(true))
		require.NoError(t, err)

		// verify that target failover group is suspended
		fg, err := secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.FailoverGroupSecondaryStateSuspended, fg.SecondaryState)

		// resume target failover group
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, sdk.NewAlterTargetFailoverGroupRequest(failoverGroup.ID()).WithResume(true))
		require.NoError(t, err)

		// verify that target failover group is resumed
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.FailoverGroupSecondaryStateStarted, failoverGroup.SecondaryState)
	})

	t.Run("refresh target failover group", func(t *testing.T) {
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, sdk.NewAlterTargetFailoverGroupRequest(failoverGroup.ID()).WithRefresh(true))
		require.NoError(t, err)
	})

	t.Run("promote secondary to primary", func(t *testing.T) {
		// promote secondary to primary
		err = secondaryClient.FailoverGroups.AlterTarget(ctx, sdk.NewAlterTargetFailoverGroupRequest(failoverGroup.ID()).WithPrimary(true))
		require.NoError(t, err)

		// verify that target failover group is promoted
		failoverGroup, err = secondaryClient.FailoverGroups.ShowByID(ctx, failoverGroup.ID())
		require.NoError(t, err)
		assert.True(t, failoverGroup.IsPrimary)
	})
}

func TestInt_FailoverGroupsDrop(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

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
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	t.Run("without show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest())
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("with show options", func(t *testing.T) {
		failoverGroups, err := client.FailoverGroups.Show(ctx, sdk.NewShowFailoverGroupRequest().
			WithInAccount(testClientHelper().Ids.AccountIdentifierWithLocator()))
		require.NoError(t, err)
		assert.LessOrEqual(t, 1, len(failoverGroups))
		assert.Contains(t, failoverGroups, failoverGroupTest)
	})

	t.Run("when searching a non-existent failover group", func(t *testing.T) {
		_, err := client.FailoverGroups.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_FailoverGroupsShowDatabases(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroupTest.ID()).
		WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeDatabases})))
	require.NoError(t, err)
	err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroupTest.ID()).
		WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedDatabases([]sdk.AccountObjectIdentifier{testClientHelper().Ids.DatabaseId()})))
	require.NoError(t, err)
	databases, err := client.FailoverGroups.ShowDatabases(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Len(t, databases, 1)
	assert.Equal(t, testClientHelper().Ids.DatabaseId(), databases[0])
}

func TestInt_FailoverGroupsShowShares(t *testing.T) {
	// TODO [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	ctx := testContext(t)
	failoverGroupTest, failoverGroupCleanup := testClientHelper().FailoverGroup.Create(t)
	t.Cleanup(failoverGroupCleanup)

	shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	err := client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroupTest.ID()).
		WithSet(*sdk.NewFailoverGroupSetRequest().WithObjectTypes([]sdk.PluralObjectType{sdk.PluralObjectTypeShares})))
	require.NoError(t, err)
	err = client.FailoverGroups.AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(failoverGroupTest.ID()).
		WithAdd(*sdk.NewFailoverGroupAddRequest().WithAllowedShares([]sdk.AccountObjectIdentifier{shareTest.ID()})))
	require.NoError(t, err)
	shares, err := client.FailoverGroups.ShowShares(ctx, failoverGroupTest.ID())
	require.NoError(t, err)
	assert.Len(t, shares, 1)
	assert.Equal(t, shareTest.ID(), shares[0])
}
