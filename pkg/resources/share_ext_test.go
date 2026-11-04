package resources

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_shareSetFieldsNotSetByReadExt(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("SHARE")

	t.Run("outbound share, enforcing mode", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, shareSchema, map[string]any{})
		err := shareSetFieldsNotSetByReadExt(d, id, &sdk.Share{Kind: sdk.ShareKindOutbound, SecureObjectsOnly: sdk.Bool(false)}, false)
		require.NoError(t, err)
		assert.Equal(t, "SHARE", d.Get("name"))
		assert.Equal(t, "false", d.Get("secure_objects_only"))
	})

	t.Run("outbound share, permissive mode", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, shareSchema, map[string]any{})
		err := shareSetFieldsNotSetByReadExt(d, id, &sdk.Share{Kind: sdk.ShareKindOutbound, SecureObjectsOnly: sdk.Bool(false)}, true)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, d.Get("secure_objects_only"))
	})

	t.Run("outbound share without a value", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, shareSchema, map[string]any{})
		err := shareSetFieldsNotSetByReadExt(d, id, &sdk.Share{Kind: sdk.ShareKindOutbound}, false)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, d.Get("secure_objects_only"))
	})

	t.Run("inbound share", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, shareSchema, map[string]any{})
		err := shareSetFieldsNotSetByReadExt(d, id, &sdk.Share{Kind: sdk.ShareKindInbound}, false)
		require.ErrorContains(t, err, "invalid share kind, expected OUTBOUND, got INBOUND")
	})
}

func Test_accountIdentifiersFromSlice(t *testing.T) {
	t.Run("valid accounts", func(t *testing.T) {
		result, err := accountIdentifiersFromSlice([]string{"org1.account1", `"org2"."account2"`})
		require.NoError(t, err)
		assert.Equal(t, []sdk.AccountIdentifier{sdk.NewAccountIdentifier("org1", "account1"), sdk.NewAccountIdentifier("org2", "account2")}, result)
	})

	t.Run("empty", func(t *testing.T) {
		result, err := accountIdentifiersFromSlice(nil)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("invalid account does not panic", func(t *testing.T) {
		_, err := accountIdentifiersFromSlice([]string{"org1.account1", "incorrect"})
		require.Error(t, err)
	})
}

func Test_shareAccountsDifference(t *testing.T) {
	account := func(name string) sdk.AccountIdentifier { return sdk.NewAccountIdentifier("org", name) }

	t.Run("accounts of the first list that are not in the second one", func(t *testing.T) {
		result := shareAccountsDifference(
			[]sdk.AccountIdentifier{account("a"), account("b"), account("c")},
			[]sdk.AccountIdentifier{account("b")},
		)
		assert.Equal(t, []sdk.AccountIdentifier{account("a"), account("c")}, result)
	})

	t.Run("the accounts are compared ignoring case", func(t *testing.T) {
		result := shareAccountsDifference(
			[]sdk.AccountIdentifier{account("ACCOUNT_A"), account("account_b")},
			[]sdk.AccountIdentifier{account("account_a"), account("ACCOUNT_B")},
		)
		assert.Empty(t, result)
	})

	t.Run("empty lists", func(t *testing.T) {
		assert.Empty(t, shareAccountsDifference(nil, []sdk.AccountIdentifier{account("a")}))
		assert.Equal(t, []sdk.AccountIdentifier{account("a")}, shareAccountsDifference([]sdk.AccountIdentifier{account("a")}, nil))
		assert.Empty(t, shareAccountsDifference(nil, nil))
	})
}

func Test_shareRecipientsExt(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("SHARE")
	account := func(name string) sdk.AccountIdentifier { return sdk.NewAccountIdentifier("org", name) }
	notCalled := func(context.Context, sdk.AccountObjectIdentifier) ([]sdk.AccountIdentifier, error) {
		t.Fatal("the target accounts should not be queried")
		return nil, nil
	}

	t.Run("all the recipients are listed, so the view is not queried", func(t *testing.T) {
		share := &sdk.Share{To: []sdk.AccountIdentifier{account("a"), account("b")}, ToTotal: 0}
		result, err := shareRecipientsExt(context.Background(), notCalled, id, share)
		require.NoError(t, err)
		assert.Equal(t, share.To, result)

		share = &sdk.Share{To: []sdk.AccountIdentifier{account("a"), account("b"), account("c")}, ToTotal: 3}
		result, err = shareRecipientsExt(context.Background(), notCalled, id, share)
		require.NoError(t, err)
		assert.Equal(t, share.To, result)
	})

	t.Run("the listed recipients are merged with the view, sorted and deduplicated", func(t *testing.T) {
		share := &sdk.Share{To: []sdk.AccountIdentifier{account("c"), account("a"), account("z_new")}, ToTotal: 5}
		lookup := func(_ context.Context, shareId sdk.AccountObjectIdentifier) ([]sdk.AccountIdentifier, error) {
			assert.Equal(t, id, shareId)
			return []sdk.AccountIdentifier{account("a"), account("b"), account("c"), account("d")}, nil
		}
		result, err := shareRecipientsExt(context.Background(), lookup, id, share)
		require.NoError(t, err)
		assert.Equal(t, []sdk.AccountIdentifier{account("a"), account("b"), account("c"), account("d"), account("z_new")}, result)
		// the listed recipients of the SHOW row are not modified
		assert.Equal(t, []sdk.AccountIdentifier{account("c"), account("a"), account("z_new")}, share.To)
	})

	t.Run("a failed lookup returns the listed recipients with the error", func(t *testing.T) {
		share := &sdk.Share{To: []sdk.AccountIdentifier{account("a"), account("b"), account("c")}, ToTotal: 4}
		lookupErr := errors.New("insufficient privileges")
		lookup := func(context.Context, sdk.AccountObjectIdentifier) ([]sdk.AccountIdentifier, error) {
			return nil, lookupErr
		}
		result, err := shareRecipientsExt(context.Background(), lookup, id, share)
		require.ErrorIs(t, err, lookupErr)
		assert.Equal(t, share.To, result)
	})
}
