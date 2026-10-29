//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Shares(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	currentRole := testClientHelper().Context.CurrentRole(t)

	createShare := func(t *testing.T, req *sdk.CreateShareRequest, id sdk.AccountObjectIdentifier) {
		t.Helper()
		err := client.Shares.Create(ctx, req)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Share.DropShareFunc(t, id))
	}

	grantDatabaseUsage := func(t *testing.T, shareId sdk.AccountObjectIdentifier) {
		t.Helper()
		t.Cleanup(testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), shareId, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))
	}

	t.Run("create: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createShare(t, sdk.NewCreateShareRequest(id), id)

		assertThatObject(t, objectassert.Share(t, id).
			HasCreatedOnNotEmpty().
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasName(id.Name()).
			HasNoTo().
			HasOwner(currentRole.Name()).
			HasComment("").
			HasListingGlobalName("").
			HasSecureObjectsOnly("true").
			HasToTotal(0),
		)
	})

	t.Run("create: all optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		createShare(t, sdk.NewCreateShareRequest(id).
			WithOrReplace(true).
			WithSecureObjectsOnly(false).
			WithComment(comment),
			id,
		)

		assertThatObject(t, objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasComment(comment).
			HasSecureObjectsOnly("false"),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createShare(t, sdk.NewCreateShareRequest(id).WithComment("first"), id)

		err := client.Shares.Create(ctx, sdk.NewCreateShareRequest(id).WithOrReplace(true).WithComment("replaced"))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Share(t, id).HasComment("replaced"))
	})

	t.Run("create: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createShare(t, sdk.NewCreateShareRequest(id).WithComment("first"), id)

		err := client.Shares.Create(ctx, sdk.NewCreateShareRequest(id).WithIfNotExists(true).WithComment("ignored"))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Share(t, id).HasComment("first"))
	})

	t.Run("create: or replace and if not exists are conflicting", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Shares.Create(ctx, sdk.NewCreateShareRequest(id).WithOrReplace(true).WithIfNotExists(true))
		require.ErrorContains(t, err, "are incompatible and cannot be set at the same time")
	})

	t.Run("create: already exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createShare(t, sdk.NewCreateShareRequest(id), id)

		err := client.Shares.Create(ctx, sdk.NewCreateShareRequest(id))
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("alter: add and remove accounts", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		grantDatabaseUsage(t, share.ID())

		accounts := []sdk.AccountIdentifier{secondaryTestClientHelper().Account.GetAccountIdentifier(t)}

		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithIfExists(true).WithAdd(sdk.ShareAddRequest{
			Accounts:          accounts,
			ShareRestrictions: sdk.Bool(false),
		}))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Share(t, share.ID()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId()).
			HasTo(accounts...).
			HasToTotal(0),
		)

		err = client.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithIfExists(true).WithRemove(sdk.ShareRemoveRequest{
			Accounts: accounts,
		}))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Share(t, share.ID()).HasNoTo().HasToTotal(0))
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		secondaryClient := testSecondaryClient(t)
		db, dbCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		share, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		t.Cleanup(secondaryTestClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, db.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

		accounts := []sdk.AccountIdentifier{testClientHelper().Account.GetAccountIdentifier(t)}
		comment := random.Comment()

		// set accounts and comment in one statement
		err := secondaryClient.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithIfExists(true).WithSet(sdk.ShareSetRequest{
			Accounts: accounts,
			Comment:  sdk.String(comment),
		}))
		require.NoError(t, err)

		assertThatObjectSecondary(t, objectassert.Share(t, share.ID()).
			HasDatabaseName(db.ID()).
			HasTo(accounts...).
			HasToTotal(0).
			HasComment(comment),
		)

		// unset the comment; accounts stay
		err = secondaryClient.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithIfExists(true).WithUnset(sdk.ShareUnsetRequest{
			Comment: sdk.Bool(true),
		}))
		require.NoError(t, err)

		assertThatObjectSecondary(t, objectassert.Share(t, share.ID()).
			HasTo(accounts...).
			HasComment(""),
		)
	})

	t.Run("alter: set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithSetTags([]sdk.TagAssociation{{Name: tag.ID(), Value: "v1"}}))
		require.NoError(t, err)

		value, err := testClientHelper().Tag.GetForObject(t, tag.ID(), share.ID(), sdk.ObjectTypeShare)
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, "v1", *value)

		err = client.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)

		value, err = testClientHelper().Tag.GetForObject(t, tag.ID(), share.ID(), sdk.ObjectTypeShare)
		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("alter: non-existing share", func(t *testing.T) {
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(NonExistingAccountObjectIdentifier).WithSet(sdk.ShareSetRequest{Comment: sdk.String("x")}))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: if exists on a non-existing share", func(t *testing.T) {
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(NonExistingAccountObjectIdentifier).WithIfExists(true).WithSet(sdk.ShareSetRequest{Comment: sdk.String("x")}))
		require.NoError(t, err)
	})

	t.Run("describe: provider", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		grantDatabaseUsage(t, share.ID())

		rows, err := client.Shares.DescribeProvider(ctx, share.ID())
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assertThatObject(t, objectassert.ShareDetailsFromObject(t, &rows[0]).
			HasKind(sdk.ObjectTypeDatabase).
			HasName(testClientHelper().Ids.DatabaseId()).
			HasSharedOnNotEmpty(),
		)

		described, err := client.Shares.Describe(ctx, share.ID())
		require.NoError(t, err)
		require.Equal(t, rows, described)
	})

	t.Run("describe: consumer", func(t *testing.T) {
		providerClient := testSecondaryClient(t)
		db, dbCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		share, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		t.Cleanup(secondaryTestClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, db.ID(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

		err := providerClient.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithAdd(sdk.ShareAddRequest{
			Accounts: []sdk.AccountIdentifier{testClientHelper().Account.GetAccountIdentifier(t)},
		}))
		require.NoError(t, err)

		rows, err := client.Shares.DescribeConsumer(ctx, share.ExternalID())
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assertThatObject(t, objectassert.ShareDetailsFromObject(t, &rows[0]).
			HasKind(sdk.ObjectTypeDatabase).
			HasName(ids.DatabasePlaceholder).
			HasSharedOnNotEmpty(),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.Shares.Describe(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: like", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		_, shareCleanup2 := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup2)

		shares, err := client.Shares.Show(ctx, sdk.NewShowShareRequest().WithLike(sdk.Like{Pattern: sdk.String(share.Name)}))
		require.NoError(t, err)
		require.Len(t, shares, 1)
		assertThatObject(t, objectassert.ShareFromObject(t, &shares[0]).
			HasName(share.Name).
			HasKind(sdk.ShareKindOutbound).
			HasListingGlobalName("").
			HasSecureObjectsOnly("true"),
		)
	})

	t.Run("show: without options", func(t *testing.T) {
		_, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		_, shareCleanup2 := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup2)

		shares, err := client.Shares.Show(ctx, sdk.NewShowShareRequest())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(shares), 2)
	})

	t.Run("show: when searching a non-existent share", func(t *testing.T) {
		shares, err := client.Shares.Show(ctx, sdk.NewShowShareRequest().WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		assert.Empty(t, shares)
	})

	t.Run("show: starts with and limit", func(t *testing.T) {
		suffix := testClientHelper().Ids.Alpha()
		for _, name := range []string{"shr_a_" + suffix, "shr_b_" + suffix} {
			id := sdk.NewAccountObjectIdentifier(name)
			createShare(t, sdk.NewCreateShareRequest(id), id)
		}

		shares, err := client.Shares.Show(ctx, sdk.NewShowShareRequest().WithStartsWith("shr_a_"+suffix))
		require.NoError(t, err)
		require.Len(t, shares, 1)
		assert.Equal(t, "shr_a_"+suffix, shares[0].Name)

		shares, err = client.Shares.Show(ctx, sdk.NewShowShareRequest().WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, shares, 1)
	})

	t.Run("show: inbound shares have no secure objects only value", func(t *testing.T) {
		shares, err := client.Shares.Show(ctx, sdk.NewShowShareRequest())
		require.NoError(t, err)
		for _, share := range shares {
			if share.Kind == sdk.ShareKindInbound {
				assert.Empty(t, share.SecureObjectsOnly)
			}
		}
	})

	t.Run("show: target accounts", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)
		grantDatabaseUsage(t, share.ID())

		account := secondaryTestClientHelper().Account.GetAccountIdentifier(t)
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(share.ID()).WithIfExists(true).WithAdd(sdk.ShareAddRequest{
			Accounts:          []sdk.AccountIdentifier{account},
			ShareRestrictions: sdk.Bool(false),
		}))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Share(t, share.ID()).HasTo(account).HasToTotal(0))

		// SNOWFLAKE.ACCOUNT_USAGE.SHARES lags SHOW SHARES by up to 3 hours, so a share created here has no row yet.
		// TODO(SNOW-1348347): add predefined share in our testing env so that we could test ShowTargetAccounts
		accounts, err := client.Shares.ShowTargetAccounts(ctx, share.ID())
		require.NoError(t, err)
		assert.Empty(t, accounts)
	})

	t.Run("show by id: existing", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		found, err := client.Shares.ShowByID(ctx, share.ID())
		require.NoError(t, err)
		assert.Equal(t, share.ID(), found.ID())

		found, err = client.Shares.ShowByIDSafely(ctx, share.ID())
		require.NoError(t, err)
		assert.Equal(t, share.ID(), found.ID())
	})

	t.Run("show by id: non-existing", func(t *testing.T) {
		_, err := client.Shares.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		_, err = client.Shares.ShowByIDSafely(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show by id: like wildcards do not select another share", func(t *testing.T) {
		suffix := testClientHelper().Ids.Alpha()
		for _, name := range []string{"prefix1" + suffix, "prefix_" + suffix} {
			id := sdk.NewAccountObjectIdentifier(name)
			createShare(t, sdk.NewCreateShareRequest(id), id)

			found, err := client.Shares.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, found.ID())
		}
	})

	t.Run("drop: existing", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := client.Shares.Drop(ctx, sdk.NewDropShareRequest(share.ID()))
		require.NoError(t, err)

		_, err = client.Shares.ShowByID(ctx, share.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.Shares.Drop(ctx, sdk.NewDropShareRequest(NonExistingAccountObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on a non-existing share", func(t *testing.T) {
		err := client.Shares.Drop(ctx, sdk.NewDropShareRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing and non-existing", func(t *testing.T) {
		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		require.NoError(t, client.Shares.DropSafely(ctx, share.ID()))
		_, err := client.Shares.ShowByIDSafely(ctx, share.ID())
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		require.NoError(t, client.Shares.DropSafely(ctx, share.ID()))
	})
}
