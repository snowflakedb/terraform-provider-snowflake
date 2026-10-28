package schemas

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationPolicyDescribeOutput(t *testing.T) {
	rows := []sdk.AuthenticationPolicyProperty{
		{Property: "NAME", Value: "POLICY"},
		{Property: "OWNER", Value: "ACCOUNTADMIN"},
		{Property: "COMMENT", Value: "null"},
		{Property: "MFA_ENROLLMENT", Value: "OPTIONAL"},
		{Property: "AUTHENTICATION_METHODS", Value: "[PASSWORD, SAML]"},
		{Property: "CLIENT_TYPES", Value: "[ALL]"},
		{Property: "SECURITY_INTEGRATIONS", Value: "[ALL]"},
		{Property: "CLIENT_POLICY", Value: "{}"},
		{Property: "MFA_POLICY", Value: "{ALLOWED_METHODS=[ALL]}"},
		{Property: "PAT_POLICY", Value: "{DEFAULT_EXPIRY_IN_DAYS=15, BLOCKED_ROLES_LIST=[]}"},
		{Property: "WORKLOAD_IDENTITY_POLICY", Value: "{ALLOWED_PROVIDERS=[ALL]}"},
		{Property: "UNEXPECTED", Value: "x"},
	}
	expected := map[string]any{
		"name":                     "POLICY",
		"owner":                    "ACCOUNTADMIN",
		"comment":                  "null",
		"mfa_enrollment":           "OPTIONAL",
		"authentication_methods":   "[PASSWORD, SAML]",
		"client_types":             "[ALL]",
		"security_integrations":    "[ALL]",
		"client_policy":            "{}",
		"mfa_policy":               "{ALLOWED_METHODS=[ALL]}",
		"pat_policy":               "{DEFAULT_EXPIRY_IN_DAYS=15, BLOCKED_ROLES_LIST=[]}",
		"workload_identity_policy": "{ALLOWED_PROVIDERS=[ALL]}",
	}

	t.Run("from rows", func(t *testing.T) {
		require.Equal(t, expected, AuthenticationPolicyDescriptionsToSchema(rows))
		require.Equal(t, map[string]any{}, AuthenticationPolicyDescriptionsToSchema(nil))
	})

	t.Run("from details", func(t *testing.T) {
		details, err := sdk.ParseAuthenticationPolicyDescribed(rows)
		require.NoError(t, err)
		require.Equal(t, expected, AuthenticationPolicyDetailsToSchema(&details))
	})

	t.Run("schema keys", func(t *testing.T) {
		require.Len(t, DescribeAuthenticationPolicyDetailsSchema, len(expected))
		for key := range expected {
			require.Contains(t, DescribeAuthenticationPolicyDetailsSchema, key)
		}
	})
}
