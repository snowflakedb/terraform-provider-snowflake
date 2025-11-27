package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScimIntegrationRunAsRoleToAccountObjectIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test special provisioners in uppercase
		{
			input:    "OKTA_PROVISIONER",
			expected: "OKTA_PROVISIONER",
		},
		{
			input:    "AAD_PROVISIONER",
			expected: "AAD_PROVISIONER",
		},
		{
			input:    "GENERIC_SCIM_PROVISIONER",
			expected: "GENERIC_SCIM_PROVISIONER",
		},
		// Test special provisioners in lowercase (should be converted to uppercase)
		{
			input:    "okta_provisioner",
			expected: "OKTA_PROVISIONER",
		},
		{
			input:    "aad_provisioner",
			expected: "AAD_PROVISIONER",
		},
		{
			input:    "generic_scim_provisioner",
			expected: "GENERIC_SCIM_PROVISIONER",
		},
		// Test special provisioners in mixed case (should be converted to uppercase)
		{
			input:    "Okta_Provisioner",
			expected: "OKTA_PROVISIONER",
		},
		{
			input:    "Aad_Provisioner",
			expected: "AAD_PROVISIONER",
		},
		{
			input:    "Generic_Scim_Provisioner",
			expected: "GENERIC_SCIM_PROVISIONER",
		},
		// Test regular role names (should preserve original casing)
		{
			input:    "MY_ROLE",
			expected: "MY_ROLE",
		},
		{
			input:    "my_role",
			expected: "my_role",
		},
		{
			input:    "My_Role",
			expected: "My_Role",
		},
		{
			input:    "ACCOUNTADMIN",
			expected: "ACCOUNTADMIN",
		},
		{
			input:    "SYSADMIN",
			expected: "SYSADMIN",
		},
		// Test quoted identifiers
		{
			input:    `"my_role"`,
			expected: `my_role`,
		},
		{
			input:    `"okta_provisioner"`,
			expected: `OKTA_PROVISIONER`,
		},
		// Test edge cases
		{
			input:    "",
			expected: "",
		},
		{
			input:    " MY_ROLE ",
			expected: " MY_ROLE ",
		},
		// Test similar but not exact matches (should not be converted)
		{
			input:    "okta_provisioner_custom",
			expected: "okta_provisioner_custom",
		},
		{
			input:    "custom_okta_provisioner",
			expected: "custom_okta_provisioner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := scimIntegrationRunAsRoleToAccountObjectIdentifier(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.Name())
		})
	}
}
