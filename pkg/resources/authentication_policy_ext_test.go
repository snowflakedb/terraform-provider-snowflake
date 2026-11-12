package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_authenticationPolicyApplySetUnset(t *testing.T) {
	t.Run("sets what is configured and nothing else", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, authenticationPolicySchema, map[string]any{
			"authentication_methods": []any{"PASSWORD"},
			"comment":                "c",
		})
		requests, err := authenticationPolicyApplySetUnset(d)
		require.NoError(t, err)
		require.NotNil(t, requests.set.Comment)
		require.Equal(t, "c", *requests.set.Comment)
		require.Len(t, requests.set.AuthenticationMethods, 1)
		require.Nil(t, requests.set.MfaPolicy)
		require.Nil(t, requests.set.PatPolicy)
		require.Nil(t, requests.set.ClientPolicy)
		require.Empty(t, *requests.unset)
	})
}
