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
