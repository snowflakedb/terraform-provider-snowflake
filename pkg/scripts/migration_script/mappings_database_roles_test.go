package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleDatabaseRolesMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "minimal database role",
			inputRows: [][]string{
				{"database_name", "name"},
				{"TEST_DB", "TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_database_role" "snowflake_generated_database_role_TEST_DB_TEST_ROLE" {
  database = "TEST_DB"
  name = "TEST_ROLE"
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_TEST_DB_TEST_ROLE
  id = "\"TEST_DB\".\"TEST_ROLE\""
}
`,
		},
		{
			name: "database role with all parameters",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "granted_database_roles", "granted_to_database_roles", "granted_to_roles", "is_current", "is_default", "is_inherited", "name", "owner", "owner_role_type"},
				{"This is a database role", "2024-06-06 00:00:00.000 +0000 UTC", "PROD_DB", "2", "1", "3", "Y", "N", "Y", "ADMIN_ROLE", "ACCOUNTADMIN", "ROLE"},
			},
			expectedOutput: `
resource "snowflake_database_role" "snowflake_generated_database_role_PROD_DB_ADMIN_ROLE" {
  database = "PROD_DB"
  name = "ADMIN_ROLE"
  comment = "This is a database role"
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_PROD_DB_ADMIN_ROLE
  id = "\"PROD_DB\".\"ADMIN_ROLE\""
}
`,
		},
		{
			name: "multiple database roles",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "granted_database_roles", "granted_to_database_roles", "granted_to_roles", "is_current", "is_default", "is_inherited", "name", "owner", "owner_role_type"},
				{"", "2024-06-06 00:00:00.000 +0000 UTC", "DB1", "0", "0", "0", "N", "N", "N", "ROLE_A", "ACCOUNTADMIN", "ROLE"},
				{"Role B comment", "2024-06-06 00:00:00.000 +0000 UTC", "DB1", "0", "0", "0", "N", "N", "N", "ROLE_B", "ACCOUNTADMIN", "ROLE"},
			},
			expectedOutput: `
resource "snowflake_database_role" "snowflake_generated_database_role_DB1_ROLE_A" {
  database = "DB1"
  name = "ROLE_A"
}

resource "snowflake_database_role" "snowflake_generated_database_role_DB1_ROLE_B" {
  database = "DB1"
  name = "ROLE_B"
  comment = "Role B comment"
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_DB1_ROLE_A
  id = "\"DB1\".\"ROLE_A\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_DB1_ROLE_B
  id = "\"DB1\".\"ROLE_B\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleDatabaseRoles(&Config{
				ObjectType: ObjectTypeDatabaseRoles,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
