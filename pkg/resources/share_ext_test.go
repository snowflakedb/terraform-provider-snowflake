package resources

import (
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
