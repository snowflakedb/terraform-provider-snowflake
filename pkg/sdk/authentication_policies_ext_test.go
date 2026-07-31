package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseAuthenticationPolicyTargetScopes(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected []AuthenticationPolicyTargetScope
	}{
		{
			Name:     "empty options returns nil",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "missing target_scopes key returns nil",
			Input:    `{"other_field":"value"}`,
			Expected: nil,
		},
		{
			Name:     "empty target_scopes array returns empty slice",
			Input:    `{"target_scopes":[]}`,
			Expected: []AuthenticationPolicyTargetScope{},
		},
		{
			Name:     "single scope",
			Input:    `{"target_scopes":["ACCOUNT"]}`,
			Expected: []AuthenticationPolicyTargetScope{AuthenticationPolicyTargetScopeAccount},
		},
		{
			Name:  "multiple scopes are sorted alphabetically",
			Input: `{"target_scopes":["SERVICE_USERS","ACCOUNT","PERSON_USERS"]}`,
			Expected: []AuthenticationPolicyTargetScope{
				AuthenticationPolicyTargetScopeAccount,
				AuthenticationPolicyTargetScopePersonUsers,
				AuthenticationPolicyTargetScopeServiceUsers,
			},
		},
		{
			Name:  "extra attributes are ignored",
			Input: `{"MFA_ENROLLMENT_REQUIREMENT":"REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY","SYSTEM_MFA_REQUIRED_CLIENT_TYPES":"ALL","target_scopes":["PERSON_USERS"]}`,
			Expected: []AuthenticationPolicyTargetScope{
				AuthenticationPolicyTargetScopePersonUsers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := parseAuthenticationPolicyTargetScopes(tc.Input)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}

	t.Run("invalid json returns error", func(t *testing.T) {
		result, err := parseAuthenticationPolicyTargetScopes(`{"target_scopes":`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
