//go:build !account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const ConnectionFailoverToAccountInSameRegionErrorMessage = "The connection cannot be failed over to an account in the same region"

func TestInt_Connections(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
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
		assertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasNoComment().
			HasIsPrimary(true).
			HasPrimaryIdentifier(externalObjectIdentifier).
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
		assertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasComment("test comment for connection").
			HasIsPrimary(true).
			HasPrimaryIdentifier(externalObjectIdentifier).
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

		secondaryAccountId := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		_, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{secondaryAccountId})),
		)
		// TODO: [SNOW-1763442]
		// require.NoError(t, err)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)

		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasNoComment().
			HasIsPrimary(true).
			HasPrimaryIdentifier(externalObjectIdentifier).
			HasFailoverAllowedToAccounts(
				accountId,
				secondaryAccountId,
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
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		secondaryAccountId := secondaryTestClientHelper().Ids.AccountIdentifierWithLocator()

		primaryConn, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, testClientHelper().Ids.RandomAccountObjectIdentifier())
		t.Cleanup(connectionCleanup)

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(primaryConn.ID()).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{secondaryAccountId})),
		)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)
		// TODO: [SNOW-1763442]
		//
		// require.NoError(t, err)

		// create replica on secondary account
		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		err = secondaryClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).WithAsReplicaOf(externalObjectIdentifier))
		t.Cleanup(testClientHelper().Connection.DropFunc(t, id))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Connection(t, id).
			HasSnowflakeRegion(sessionDetails.Region).
			HasAccountName(sessionDetails.AccountName).
			HasName(id.Name()).
			HasNoComment().
			HasIsPrimary(false).
			HasPrimaryIdentifier(externalObjectIdentifier).
			HasFailoverAllowedToAccounts(
				accountId,
				secondaryAccountId,
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

	t.Run("Alter disable failover", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		secondaryAccountId := secondaryTestClientHelper().Account.GetAccountIdentifier(t)

		primaryConn, connectionCleanup := testClientHelper().Connection.CreateWithIdentifier(t, id)
		t.Cleanup(connectionCleanup)

		// Add secondary account to failover list
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(
				*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{secondaryAccountId})),
		)
		require.ErrorContains(t, err, ConnectionFailoverToAccountInSameRegionErrorMessage)
		// TODO: [SNOW-1763442]
		//
		// require.NoError(t, err)

		// Disable promotion of this connection
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest()))
		require.NoError(t, err)

		// Assert that promotion for other account has been disabled
		externalObjectIdentifier := sdk.NewExternalObjectIdentifier(accountId, id)
		assertThatObject(t, objectassert.Connection(t, primaryConn.ID()).
			HasPrimaryIdentifier(externalObjectIdentifier).
			HasFailoverAllowedToAccounts(accountId),
		)

		// Try to create repllication on secondary account
		err = secondaryClient.Connections.Create(ctx, sdk.NewCreateConnectionRequest(id).WithAsReplicaOf(externalObjectIdentifier))
		require.ErrorContains(t, err, "This account is not authorized to create a secondary connection of this primary connection")
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

		assertThatObject(t, objectassert.Connection(t, id).
			HasName(id.Name()).
			HasComment("new integration test comment"),
		)

		// Unset
		err = client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithUnset(*sdk.NewConnectionUnsetRequest().
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Connection(t, id).
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
