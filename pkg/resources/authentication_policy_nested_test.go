package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseAuthenticationPolicyNestedBlocks(t *testing.T) {
	d := schema.TestResourceDataRaw(t, authenticationPolicySchema, map[string]any{
		"mfa_policy": []any{
			map[string]any{
				"enforce_mfa_on_external_authentication": "all",
				"allowed_methods":                        []any{"passkey", "duo"},
			},
		},
		"pat_policy": []any{
			map[string]any{
				"default_expiry_in_days":                     20,
				"network_policy_evaluation":                  "enforced_not_required",
				"require_role_restriction_for_service_users": "false",
			},
		},
		"workload_identity_policy": []any{
			map[string]any{
				"allowed_providers":    []any{"aws", "oidc"},
				"allowed_aws_accounts": []any{"123"},
			},
		},
	})

	mfa, err := parseAuthenticationPolicyMfaPolicy(d, "mfa_policy.0.")
	require.NoError(t, err)
	require.NotNil(t, mfa.EnforceMfaOnExternalAuthentication)
	assert.Equal(t, sdk.EnforceMfaOnExternalAuthenticationOptionAll, *mfa.EnforceMfaOnExternalAuthentication)
	assert.ElementsMatch(t, []sdk.MfaPolicyAllowedMethodsOptionListItem{
		{Value: sdk.MfaPolicyAllowedMethodsOptionPasskey},
		{Value: sdk.MfaPolicyAllowedMethodsOptionDuo},
	}, mfa.AllowedMethods)

	pat, err := parseAuthenticationPolicyPatPolicy(d, "pat_policy.0.")
	require.NoError(t, err)
	require.NotNil(t, pat.DefaultExpiryInDays)
	assert.Equal(t, 20, *pat.DefaultExpiryInDays)
	assert.Nil(t, pat.MaxExpiryInDays)
	require.NotNil(t, pat.NetworkPolicyEvaluation)
	assert.Equal(t, sdk.NetworkPolicyEvaluationOptionEnforcedNotRequired, *pat.NetworkPolicyEvaluation)
	require.NotNil(t, pat.RequireRoleRestrictionForServiceUsers)
	assert.False(t, *pat.RequireRoleRestrictionForServiceUsers)

	workload, err := parseAuthenticationPolicyWorkloadIdentityPolicy(d, "workload_identity_policy.0.")
	require.NoError(t, err)
	assert.ElementsMatch(t, []sdk.AllowedProviderOptionListItem{
		{Value: sdk.AllowedProviderOptionAws},
		{Value: sdk.AllowedProviderOptionOidc},
	}, workload.AllowedProviders)
	assert.Equal(t, []sdk.StringListItemWrapper{{Value: "123"}}, workload.AllowedAwsAccounts)
	assert.Nil(t, workload.AllowedAzureIssuers)
	assert.Nil(t, workload.AllowedOidcIssuers)
}

func Test_parseAuthenticationPolicyPatPolicy_leavesRequireRoleRestrictionToSnowflake(t *testing.T) {
	d := schema.TestResourceDataRaw(t, authenticationPolicySchema, map[string]any{
		"pat_policy": []any{
			map[string]any{
				"default_expiry_in_days":                     10,
				"require_role_restriction_for_service_users": BooleanDefault,
			},
		},
	})

	pat, err := parseAuthenticationPolicyPatPolicy(d, "pat_policy.0.")
	require.NoError(t, err)
	assert.Nil(t, pat.RequireRoleRestrictionForServiceUsers)
	require.NotNil(t, pat.DefaultExpiryInDays)
	assert.Equal(t, 10, *pat.DefaultExpiryInDays)
}

func Test_authenticationPolicyNestedBlockUpdate(t *testing.T) {
	t.Run("changes only the nested field that differs", func(t *testing.T) {
		d := resourceDataWithPlan(t, authenticationPolicySchema, map[string]string{
			"name":                                   "n",
			"database":                               "db",
			"schema":                                 "sc",
			"pat_policy.#":                           "1",
			"pat_policy.0.default_expiry_in_days":    "15",
			"pat_policy.0.max_expiry_in_days":        "0",
			"pat_policy.0.network_policy_evaluation": "",
			"pat_policy.0.require_role_restriction_for_service_users": BooleanDefault,
		}, map[string]any{
			"name":     "n",
			"database": "db",
			"schema":   "sc",
			"pat_policy": []any{
				map[string]any{
					"default_expiry_in_days":                     30,
					"require_role_restriction_for_service_users": BooleanDefault,
				},
			},
		})

		var got *sdk.AuthenticationPolicyPatPolicyRequest
		var unset *bool
		require.NoError(t, singleBlockAttributeUpdate(d, "pat_policy", &got, &unset, applyPatPolicyUpdate))
		require.NotNil(t, got)
		assert.Nil(t, unset)
		require.NotNil(t, got.DefaultExpiryInDays)
		assert.Equal(t, 30, *got.DefaultExpiryInDays)
		assert.Nil(t, got.MaxExpiryInDays)
		assert.Nil(t, got.NetworkPolicyEvaluation)
		assert.Nil(t, got.RequireRoleRestrictionForServiceUsers)
	})

	t.Run("removes the block", func(t *testing.T) {
		d := resourceDataWithPlan(t, authenticationPolicySchema, map[string]string{
			"name":         "n",
			"database":     "db",
			"schema":       "sc",
			"mfa_policy.#": "1",
			"mfa_policy.0.enforce_mfa_on_external_authentication": "ALL",
			"mfa_policy.0.allowed_methods.#":                      "1",
			"mfa_policy.0.allowed_methods.0":                      "PASSKEY",
		}, map[string]any{
			"name":     "n",
			"database": "db",
			"schema":   "sc",
		})

		var got *sdk.AuthenticationPolicyMfaPolicyRequest
		var unset *bool
		require.NoError(t, singleBlockAttributeUpdate(d, "mfa_policy", &got, &unset, applyMfaPolicyUpdate))
		assert.Nil(t, got)
		require.NotNil(t, unset)
		assert.True(t, *unset)
	})
}
