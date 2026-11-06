//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Connections(t *testing.T) {
	client := testClient(t)
	// Connection failover needs an account in a different region, so the Azure account is used (the default and secondary accounts share a region).
	azureClient := testAzureClient(t)
	ctx := testContext(t)

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)
	accountId := testClientHelper().Account.GetAccountIdentifier(t)

	t.Run("Create minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		require.NoError(t, err)

		err = client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(
			t, objectassert.Connection(t, id).
				HasSnowflakeRegion(strings.TrimPrefix(sessionDetails.Region, "PUBLIC.")).
				HasAccountName(sessionDetails.AccountName).
				HasName(id.Name()).
				HasNoComment().
				HasIsPrimary(true).
				HasPrimary(externalObjectIdentifier).
				HasFailoverAllowedToAccounts(accountId).
				HasOrganizationName(sessionDetails.OrganizationName).
				HasAccountLocator(client.GetAccountLocator()).
				HasConnectionUrl(
					strings.ToLower(
						fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
					),
				),
		)
	})

	t.Run("Create all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).
			WithIfNotExists(true).
			WithComment("test comment for connection"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))

		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(
			t, objectassert.Connection(t, id).
				HasSnowflakeRegion(strings.TrimPrefix(sessionDetails.Region, "PUBLIC.")).
				HasAccountName(sessionDetails.AccountName).
				HasName(id.Name()).
				HasComment("test comment for connection").
				HasIsPrimary(true).
				HasPrimary(externalObjectIdentifier).
				HasFailoverAllowedToAccounts(accountId).
				HasOrganizationName(sessionDetails.OrganizationName).
				HasAccountLocator(client.GetAccountLocator()).
				HasConnectionUrl(
					strings.ToLower(
						fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
					),
				),
		)
	})

	t.Run("Alter enable failover", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		azureAccountId := azureTestClientHelper().Account.GetAccountIdentifier(t)

		_, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		err := client.Connections.Alter(
			ctx, sdk.NewAlterConnectionRequest(id).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{azureAccountId}),
				),
		)
		require.NoError(t, err)

		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(
			t, objectassert.Connection(t, id).
				HasSnowflakeRegion(strings.TrimPrefix(sessionDetails.Region, "PUBLIC.")).
				HasAccountName(sessionDetails.AccountName).
				HasName(id.Name()).
				HasNoComment().
				HasIsPrimary(true).
				HasPrimary(externalObjectIdentifier).
				HasFailoverAllowedToAccountsInAnyOrder(
					accountId,
					azureAccountId,
				).
				HasOrganizationName(sessionDetails.OrganizationName).
				HasAccountLocator(client.GetAccountLocator()).
				HasConnectionUrl(
					strings.ToLower(
						fmt.Sprintf("%s-%s.snowflakecomputing.com", sessionDetails.OrganizationName, id.Name()),
					),
				),
		)
	})

	t.Run("Create as replica of", func(t *testing.T) {
		azureAccountId := azureTestClientHelper().Account.GetAccountIdentifier(t)

		primaryConn, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, testClientHelper().Ids.RandomAccountObjectIdentifier())
		t.Cleanup(connectionCleanup)
		id := primaryConn.ID()

		err := client.Connections.Alter(
			ctx, sdk.NewAlterConnectionRequest(id).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{azureAccountId}),
				),
		)
		require.NoError(t, err)

		// create replica on the Azure account (the failover permission may take a moment to propagate)
		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		require.Eventually(t, func() bool {
			err := azureClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).WithAsReplicaOf(externalObjectIdentifier))
			if err != nil {
				t.Logf("waiting for connection failover to propagate: %s", err)
			}
			return err == nil
		}, 3*time.Minute, 10*time.Second)
		t.Cleanup(azureTestClientHelper().Connection.DropFunc(t, id))

		assertThatObjectAzure(
			t, objectassert.Connection(t, id).
				HasName(id.Name()).
				HasNoComment().
				HasIsPrimary(false).
				HasPrimary(externalObjectIdentifier).
				HasOrganizationName(sessionDetails.OrganizationName).
				HasAccountLocator(azureClient.GetAccountLocator()),
		)
	})

	t.Run("Alter disable failover", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		azureAccountId := azureTestClientHelper().Account.GetAccountIdentifier(t)

		primaryConn, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		// Add the Azure account to failover list
		err := client.Connections.Alter(
			ctx, sdk.NewAlterConnectionRequest(id).
				WithEnableConnectionFailover(
					*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{azureAccountId}),
				),
		)
		require.NoError(t, err)

		// Disable promotion of this connection
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest()))
		require.NoError(t, err)

		// Assert that promotion for other account has been disabled
		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(
			t, objectassert.Connection(t, primaryConn.ID()).
				HasPrimary(externalObjectIdentifier).
				HasFailoverAllowedToAccounts(accountId),
		)

		// Try to create replication on the Azure account (disabling failover may take a moment to propagate,
		// so a replica that slips through is dropped and the attempt is repeated)
		require.Eventually(t, func() bool {
			err := azureClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).WithAsReplicaOf(externalObjectIdentifier))
			if err == nil {
				t.Log("replica was created before disabling failover propagated; dropping it and retrying")
				azureTestClientHelper().Connection.DropFunc(t, id)()
				return false
			}
			return strings.Contains(err.Error(), "This account is not authorized to create a secondary connection of this primary connection")
		}, 3*time.Minute, 10*time.Second)
	})

	t.Run("Alter comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		// Set
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithSet(*sdk.NewConnectionSetRequest().
				WithComment("new integration test comment")))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Connection(t, id).
				HasName(id.Name()).
				HasComment("new integration test comment"),
		)

		// Unset
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithUnset(*sdk.NewConnectionUnsetRequest().
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Connection(t, id).
				HasName(id.Name()).
				HasNoComment(),
		)
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		connection, err := client.Connections.ShowByID(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, connection)

		err = client.Connections.Drop(ctx, sdk.NewDropConnectionRequest(id))
		require.NoError(t, err)

		connection, err = client.Connections.ShowByID(ctx, id)
		require.Nil(t, connection)
		require.Error(t, err)
	})

	t.Run("Drop with if exists", func(t *testing.T) {
		err = client.Connections.Drop(ctx, sdk.NewDropConnectionRequest(NonExistingAccountObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		err = client.Connections.Drop(ctx, sdk.NewDropConnectionRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("Show", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()

		connection1, connectionCleanup1 := testClientHelper().Connection.CreateWithIdentifier(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.CreateWithIdentifier(t, id2)
		t.Cleanup(connectionCleanup2)

		returnedConnections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest())
		require.NoError(t, err)
		require.Contains(t, returnedConnections, *connection1)
		require.Contains(t, returnedConnections, *connection2)
	})

	t.Run("Show with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()

		connection1, connectionCleanup1 := testClientHelper().Connection.CreateWithIdentifier(t, id1)
		t.Cleanup(connectionCleanup1)

		connection2, connectionCleanup2 := testClientHelper().Connection.CreateWithIdentifier(t, id2)
		t.Cleanup(connectionCleanup2)

		returnedConnections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest().
			WithLike(sdk.Like{
				Pattern: sdk.String(id1.Name()),
			}))
		require.NoError(t, err)
		require.Contains(t, returnedConnections, *connection1)
		require.NotContains(t, returnedConnections, *connection2)
	})

	t.Run("ShowByID", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		_, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		returnedConnection, err := client.Connections.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, returnedConnection.ID())
		require.Equal(t, id.Name(), returnedConnection.Name)
	})
}
