package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseSessionPolicyTargetScopes(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected []SessionPolicyTargetScope
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
			Expected: []SessionPolicyTargetScope{},
		},
		{
			Name:     "single scope",
			Input:    `{"target_scopes":["ACCOUNT"]}`,
			Expected: []SessionPolicyTargetScope{SessionPolicyTargetScopeAccount},
		},
		{
			Name:  "multiple scopes are sorted alphabetically",
			Input: `{"target_scopes":["SERVICE_USERS","ACCOUNT","PERSON_USERS"]}`,
			Expected: []SessionPolicyTargetScope{
				SessionPolicyTargetScopeAccount,
				SessionPolicyTargetScopePersonUsers,
				SessionPolicyTargetScopeServiceUsers,
			},
		},
		{
			Name:  "extra attributes are ignored",
			Input: `{"target_scopes":["PERSON_USERS","ACCOUNT"],"comment":"ignored","nested":{"a":1},"list":[1,2,3]}`,
			Expected: []SessionPolicyTargetScope{
				SessionPolicyTargetScopeAccount,
				SessionPolicyTargetScopePersonUsers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := parseSessionPolicyTargetScopes(tc.Input)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}

	t.Run("invalid json returns error", func(t *testing.T) {
		result, err := parseSessionPolicyTargetScopes(`{"target_scopes":`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
