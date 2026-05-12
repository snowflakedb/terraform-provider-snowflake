package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleAccountRolesMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "minimal account role",
			inputRows: [][]string{
				{"name"},
				{"TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_account_role" "snowflake_generated_account_role_TEST_ROLE" {
  name = "TEST_ROLE"
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_TEST_ROLE
  id = "\"TEST_ROLE\""
}
`,
		},
		{
			name: "account role with all parameters",
			inputRows: [][]string{
				{"assigned_to_users", "comment", "created_on", "granted_roles", "granted_to_roles", "is_current", "is_default", "is_inherited", "name", "owner"},
				{"5", "This is a test role", "2024-06-06 00:00:00.000 +0000 UTC", "2", "1", "Y", "N", "Y", "ADMIN_ROLE", "ACCOUNTADMIN"},
			},
			expectedOutput: `
resource "snowflake_account_role" "snowflake_generated_account_role_ADMIN_ROLE" {
  name = "ADMIN_ROLE"
  comment = "This is a test role"
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_ADMIN_ROLE
  id = "\"ADMIN_ROLE\""
}
`,
		},
		{
			name: "multiple account roles",
			inputRows: [][]string{
				{"assigned_to_users", "comment", "created_on", "granted_roles", "granted_to_roles", "is_current", "is_default", "is_inherited", "name", "owner"},
				{"0", "", "2024-06-06 00:00:00.000 +0000 UTC", "0", "0", "N", "N", "N", "ROLE_A", "ACCOUNTADMIN"},
				{"0", "Role B comment", "2024-06-06 00:00:00.000 +0000 UTC", "0", "0", "N", "N", "N", "ROLE_B", "ACCOUNTADMIN"},
			},
			expectedOutput: `
resource "snowflake_account_role" "snowflake_generated_account_role_ROLE_A" {
  name = "ROLE_A"
}

resource "snowflake_account_role" "snowflake_generated_account_role_ROLE_B" {
  name = "ROLE_B"
  comment = "Role B comment"
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_ROLE_A
  id = "\"ROLE_A\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_ROLE_B
  id = "\"ROLE_B\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleAccountRoles(&Config{
				ObjectType: ObjectTypeAccountRoles,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
