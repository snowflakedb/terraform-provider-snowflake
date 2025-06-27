//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SharesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	_, shareCleanup2 := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup2)

	t.Run("without show options", func(t *testing.T) {
		shares, err := client.Shares.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(shares))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, shares, 1)
		assert.Contains(t, shares, *shareTest)
	})

	t.Run("when searching a non-existent share", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Empty(t, shares)
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowShareOptions{
			Limit: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		}
		shares, err := client.Shares.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, shares, 1)
	})
}

func TestInt_SharesCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Shares.Create(ctx, id, &sdk.CreateShareOptions{
			OrReplace: sdk.Bool(true),
			Comment:   sdk.String("test comment"),
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
			Limit: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)
		assert.Len(t, shares, 1)
		assert.Equal(t, id.Name(), shares[0].Name.Name())
		assert.Equal(t, "test comment", shares[0].Comment)

		t.Cleanup(testClientHelper().Share.DropShareFunc(t, id))
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Shares.Create(ctx, id, &sdk.CreateShareOptions{
			OrReplace: sdk.Bool(true),
			Comment:   sdk.String("test comment"),
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(shares), 1)

		t.Cleanup(testClientHelper().Share.DropShareFunc(t, id))
	})
}

func TestInt_SharesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when share exists", func(t *testing.T) {
		shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		err := client.Shares.Drop(ctx, shareTest.ID(), &sdk.DropShareOptions{})
		require.NoError(t, err)
	})

	t.Run("when share does not exist", func(t *testing.T) {
		err := client.Shares.Drop(ctx, NonExistingAccountObjectIdentifier, &sdk.DropShareOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_SharesAlter(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	t.Run("add and remove accounts", func(t *testing.T) {
		shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testClientHelper().Ids.DatabaseId(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: testClientHelper().Ids.DatabaseId(),
			}, shareTest.ID())
		})
		require.NoError(t, err)
		accountsToAdd := []sdk.AccountIdentifier{
			secondaryTestClientHelper().Account.GetAccountIdentifier(t),
		}
		// first add the account.
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Add: &sdk.ShareAdd{
				Accounts:          accountsToAdd,
				ShareRestrictions: sdk.Bool(false),
			},
		})
		require.NoError(t, err)
		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Len(t, shares, 1)
		share := shares[0]
		assert.Equal(t, accountsToAdd, share.To)

		// now remove the account that was added.
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Remove: &sdk.ShareRemove{
				Accounts: accountsToAdd,
			},
		})
		require.NoError(t, err)
		shares, err = client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)
		assert.Len(t, shares, 1)
		share = shares[0]
		assert.Empty(t, share.To)
	})

	t.Run("set accounts", func(t *testing.T) {
		db, dbCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := secondaryClient.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: db.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err := secondaryClient.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: db.ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		accountsToSet := []sdk.AccountIdentifier{
			testClientHelper().Account.GetAccountIdentifier(t),
		}

		// first add the account.
		err = secondaryClient.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Set: &sdk.ShareSet{
				Accounts: accountsToSet,
			},
		})
		require.NoError(t, err)

		shares, err := secondaryClient.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)

		assert.Len(t, shares, 1)
		share := shares[0]
		assert.Equal(t, accountsToSet, share.To)
	})

	t.Run("set and unset comment", func(t *testing.T) {
		shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testClientHelper().Ids.DatabaseId(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: testClientHelper().Ids.DatabaseId(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		comment := random.Comment()
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Set: &sdk.ShareSet{
				Comment: sdk.String(comment),
			},
		})
		require.NoError(t, err)

		shares, err := client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)

		assert.Len(t, shares, 1)
		share := shares[0]
		assert.Equal(t, comment, share.Comment)

		// reset comment
		err = client.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			IfExists: sdk.Bool(true),
			Unset: &sdk.ShareUnset{
				Comment: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		shares, err = client.Shares.Show(ctx, &sdk.ShowShareOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(shareTest.Name.Name()),
			},
		})
		require.NoError(t, err)

		assert.Len(t, shares, 1)
		share = shares[0]
		assert.Equal(t, "", share.Comment)
	})
}

func TestInt_ShareDescribeProvider(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("describe share", func(t *testing.T) {
		shareTest, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testClientHelper().Ids.DatabaseId(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: testClientHelper().Ids.DatabaseId(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		shareDetails, err := client.Shares.DescribeProvider(ctx, shareTest.ID())
		require.NoError(t, err)

		assert.Len(t, shareDetails.SharedObjects, 1)
		sharedObject := shareDetails.SharedObjects[0]
		assert.Equal(t, sdk.ObjectTypeDatabase, sharedObject.Kind)
		assert.Equal(t, testClientHelper().Ids.DatabaseId(), sharedObject.Name)
	})
}

func TestInt_ShareDescribeConsumer(t *testing.T) {
	ctx := testContext(t)
	providerClient := testSecondaryClient(t)
	consumerClient := testClient(t)

	t.Run("describe share", func(t *testing.T) {
		db, dbCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := providerClient.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: db.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
		t.Cleanup(func() {
			err = providerClient.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
				Database: db.ID(),
			}, shareTest.ID())
			require.NoError(t, err)
		})

		// add a consumer account to share.
		err = providerClient.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
			Add: &sdk.ShareAdd{
				Accounts: []sdk.AccountIdentifier{
					testClientHelper().Account.GetAccountIdentifier(t),
				},
			},
		})
		require.NoError(t, err)

		shareDetails, err := consumerClient.Shares.DescribeConsumer(ctx, shareTest.ExternalID())
		require.NoError(t, err)

		assert.Len(t, shareDetails.SharedObjects, 1)
		sharedObject := shareDetails.SharedObjects[0]
		assert.Equal(t, sdk.ObjectTypeDatabase, sharedObject.Kind)
		assert.Equal(t, ids.DatabasePlaceholder, sharedObject.Name)
	})
}
