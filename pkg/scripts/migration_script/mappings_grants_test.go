package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGrantValidations(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		outputContains []string
	}{
		{
			name: "same rows",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE ROLE", "ACCOUNT", "ACC_NAME", "ROLE", "TEST_ROLE", "true"},
				{"CREATE ROLE", "ACCOUNT", "ACC_NAME", "ROLE", "TEST_ROLE", "true"},
			},
			outputContains: []string{`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_with_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account = true
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`},
		},
		{
			name: "valid grant grouping",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE ROLE", "DATABASE", "TEST_DATABASE", "ROLE", "TEST_ROLE", "true"},
				{"CREATE ROLE", "USER", "TEST_DATABASE", "ROLE", "TEST_ROLE", "true"},         // different granted on
				{"CREATE ROLE", "DATABASE", "TEST_DATABASE_2", "ROLE", "TEST_ROLE", "true"},   // different name
				{"CREATE ROLE", "DATABASE", "TEST_DATABASE", "DATABASE", "TEST_ROLE", "true"}, // different granted to (won't be generated as we are not able to map it)
				{"CREATE ROLE", "DATABASE", "TEST_DATABASE", "ROLE", "TEST_ROLE_2", "true"},   // different grantee
				{"CREATE ROLE", "DATABASE", "TEST_DATABASE", "ROLE", "TEST_ROLE", "false"},    // different grant option
			},
			outputContains: []string{
				`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_TEST_DATABASE_to_TEST_ROLE_with_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_USER_TEST_DATABASE_to_TEST_ROLE_with_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "USER"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_TEST_DATABASE_2_to_TEST_ROLE_with_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE_2"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_TEST_DATABASE_to_TEST_ROLE_2_with_grant_option" {
  account_role_name = "TEST_ROLE_2"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_TEST_DATABASE_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = false
}`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(&Config{
				ObjectType: ObjectTypeGrants,
				ImportFlag: ImportStatementTypeStatement,
			}, tc.inputRows)

			assert.NoError(t, err)
			for _, expectedOutput := range tc.outputContains {
				assert.Contains(t, output, expectedOutput)
			}
		})
	}
}

func TestHandleGrantPrivilegeToAccountRoleMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "grants on account",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE DATABASE", "ACCOUNT", "ACC_NAME", "ROLE", "TEST_ROLE", "true"},
				{"CREATE ROLE", "ACCOUNT", "ACC_NAME", "ROLE", "TEST_ROLE", "true"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_with_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account = true
  privileges = ["CREATE DATABASE", "CREATE ROLE"]
  with_grant_option = true
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_ROLE_with_grant_option '"TEST_ROLE"|true|false|CREATE DATABASE,CREATE ROLE|OnAccount'
`,
		},
		{
			name: "grants on account object",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE DATABASE ROLE", "DATABASE", "TEST_DATABASE", "ROLE", "TEST_ROLE", "false"},
				{"CREATE SCHEMA", "DATABASE", "TEST_DATABASE", "ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_TEST_DATABASE_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE DATABASE ROLE", "CREATE SCHEMA"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_DATABASE_TEST_DATABASE_to_TEST_ROLE_without_grant_option '"TEST_ROLE"|false|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnAccountObject|DATABASE|"TEST_DATABASE"'
`,
		},
		{
			name: "grants on schema",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE TABLE", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "ROLE", "TEST_ROLE", "false"},
				{"CREATE VIEW", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_schema_TEST_DATABASE_TEST_SCHEMA_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_schema {
    schema_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\""
  }
  privileges = ["CREATE TABLE", "CREATE VIEW"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_schema_TEST_DATABASE_TEST_SCHEMA_to_TEST_ROLE_without_grant_option '"TEST_ROLE"|false|false|CREATE TABLE,CREATE VIEW|OnSchema|OnSchema|"TEST_DATABASE"."TEST_SCHEMA"'
`,
		},
		{
			name: "grants on schema object",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"INSERT", "TABLE", "TEST_DATABASE.TEST_SCHEMA.TEST_TABLE", "ROLE", "TEST_ROLE", "false"},
				{"SELECT", "TABLE", "TEST_DATABASE.TEST_SCHEMA.TEST_TABLE", "ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_TABLE_TEST_DATABASE_TEST_SCHEMA_TEST_TABLE_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_schema_object {
    object_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\".\"TEST_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["INSERT", "SELECT"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_TABLE_TEST_DATABASE_TEST_SCHEMA_TEST_TABLE_to_TEST_ROLE_without_grant_option '"TEST_ROLE"|false|false|INSERT,SELECT|OnSchemaObject|OnObject|TABLE|"TEST_DATABASE"."TEST_SCHEMA"."TEST_TABLE"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(&Config{
				ObjectType: ObjectTypeGrants,
				ImportFlag: ImportStatementTypeStatement,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), output)
		})
	}
}

func TestHandleGrantPrivilegeToDatabaseRoleMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "grants on database",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE SCHEMA", "DATABASE", "TEST_DATABASE", "DATABASE_ROLE", "TEST_ROLE", "true"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_database_TEST_DATABASE_to_TEST_DATABASE_TEST_ROLE_with_grant_option" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  on_database = "TEST_DATABASE"
  privileges = ["CREATE SCHEMA"]
  with_grant_option = true
}
# terraform import snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_database_TEST_DATABASE_to_TEST_DATABASE_TEST_ROLE_with_grant_option '"TEST_DATABASE"."TEST_ROLE"|true|false|CREATE SCHEMA|OnDatabase|"TEST_DATABASE"'
`,
		},
		{
			name: "grants on schema",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE TABLE", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "DATABASE_ROLE", "TEST_ROLE", "false"},
				{"CREATE VIEW", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "DATABASE_ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_schema_TEST_DATABASE_TEST_SCHEMA_to_TEST_DATABASE_TEST_ROLE_without_grant_option" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  on_schema {
    schema_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\""
  }
  privileges = ["CREATE TABLE", "CREATE VIEW"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_schema_TEST_DATABASE_TEST_SCHEMA_to_TEST_DATABASE_TEST_ROLE_without_grant_option '"TEST_DATABASE"."TEST_ROLE"|false|false|CREATE TABLE,CREATE VIEW|OnSchema|OnSchema|"TEST_DATABASE"."TEST_SCHEMA"'
`,
		},
		{
			name: "grants on schema",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"INSERT", "TABLE", "TEST_DATABASE.TEST_SCHEMA.TEST_TABLE", "DATABASE_ROLE", "TEST_ROLE", "false"},
				{"SELECT", "TABLE", "TEST_DATABASE.TEST_SCHEMA.TEST_TABLE", "DATABASE_ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_TABLE_TEST_DATABASE_TEST_SCHEMA_TEST_TABLE_to_TEST_DATABASE_TEST_ROLE_without_grant_option" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  on_schema_object {
    object_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\".\"TEST_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["INSERT", "SELECT"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_TABLE_TEST_DATABASE_TEST_SCHEMA_TEST_TABLE_to_TEST_DATABASE_TEST_ROLE_without_grant_option '"TEST_DATABASE"."TEST_ROLE"|false|false|INSERT,SELECT|OnSchemaObject|OnObject|TABLE|"TEST_DATABASE"."TEST_SCHEMA"."TEST_TABLE"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(&Config{
				ObjectType: ObjectTypeGrants,
				ImportFlag: ImportStatementTypeStatement,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), output)
		})
	}
}

func TestHandleGrantAccountRoleMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "grant role to role (SHOW GRANTS TO ROLE output)",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"USAGE", "ROLE", "TEST_ROLE", "ROLE", "PARENT_TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_account_role" "snowflake_generated_grant_TEST_ROLE_to_role_PARENT_TEST_ROLE" {
  parent_role_name = "PARENT_TEST_ROLE"
  role_name = "TEST_ROLE"
}
# terraform import snowflake_grant_account_role.snowflake_generated_grant_TEST_ROLE_to_role_PARENT_TEST_ROLE '"TEST_ROLE"|ROLE|"PARENT_TEST_ROLE"'
`,
		},
		{
			name: "grant role to user (SHOW GRANTS TO ROLE output)",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "role", "granted_to", "grantee_name", "grant_option"},
				{"USAGE", "ROLE", "TEST_ROLE", "TEST_ROLE", "USER", "TEST_USER", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_account_role" "snowflake_generated_grant_TEST_ROLE_to_user_TEST_USER" {
  role_name = "TEST_ROLE"
  user_name = "TEST_USER"
}
# terraform import snowflake_grant_account_role.snowflake_generated_grant_TEST_ROLE_to_user_TEST_USER '"TEST_ROLE"|USER|"TEST_USER"'
`,
		},
		{
			name: "grant role to role (SHOW GRANTS OF ROLE output)",
			inputRows: [][]string{
				{"role", "granted_to", "grantee_name"},
				{"TEST_ROLE", "ROLE", "PARENT_TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_grant_account_role" "snowflake_generated_grant_TEST_ROLE_to_role_PARENT_TEST_ROLE" {
  parent_role_name = "PARENT_TEST_ROLE"
  role_name = "TEST_ROLE"
}
# terraform import snowflake_grant_account_role.snowflake_generated_grant_TEST_ROLE_to_role_PARENT_TEST_ROLE '"TEST_ROLE"|ROLE|"PARENT_TEST_ROLE"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(&Config{
				ObjectType: ObjectTypeGrants,
				ImportFlag: ImportStatementTypeStatement,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), output)
		})
	}
}

func TestHandleGrantDatabaseRoleMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "grant database role to role (SHOW GRANTS TO DATABASE ROLE output)",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"USAGE", "DATABASE_ROLE", "TEST_DATABASE.TEST_ROLE", "ROLE", "PARENT_TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_database_role" "snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_role_PARENT_TEST_ROLE" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  parent_role_name = "PARENT_TEST_ROLE"
}
# terraform import snowflake_grant_database_role.snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_role_PARENT_TEST_ROLE '"TEST_DATABASE"."TEST_ROLE"|ROLE|"PARENT_TEST_ROLE"'
`,
		},
		{
			name: "grant database role to database role (SHOW GRANTS TO DATABASE ROLE output)",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"USAGE", "DATABASE_ROLE", "TEST_DATABASE.TEST_ROLE", "DATABASE_ROLE", "PARENT_TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_database_role" "snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  parent_database_role_name = "\"TEST_DATABASE\".\"PARENT_TEST_ROLE\""
}
# terraform import snowflake_grant_database_role.snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE '"TEST_DATABASE"."TEST_ROLE"|DATABASE ROLE|"TEST_DATABASE"."PARENT_TEST_ROLE"'
`,
		},
		{
			name: "grant database role to database role (SHOW GRANTS OF DATABASE ROLE output)",
			inputRows: [][]string{
				{"role", "granted_to", "grantee_name"},
				{"TEST_DATABASE.TEST_ROLE", "DATABASE_ROLE", "PARENT_TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_grant_database_role" "snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  parent_database_role_name = "\"TEST_DATABASE\".\"PARENT_TEST_ROLE\""
}
# terraform import snowflake_grant_database_role.snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE '"TEST_DATABASE"."TEST_ROLE"|DATABASE ROLE|"TEST_DATABASE"."PARENT_TEST_ROLE"'
`,
		},
		{
			name: "grant database role to database role (SHOW GRANTS OF ROLE output)",
			inputRows: [][]string{
				{"role", "granted_to", "grantee_name"},
				{"TEST_DATABASE.TEST_ROLE", "DATABASE_ROLE", "TEST_DATABASE.PARENT_TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_grant_database_role" "snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  parent_database_role_name = "\"TEST_DATABASE\".\"PARENT_TEST_ROLE\""
}
# terraform import snowflake_grant_database_role.snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_database_role_TEST_DATABASE_PARENT_TEST_ROLE '"TEST_DATABASE"."TEST_ROLE"|DATABASE ROLE|"TEST_DATABASE"."PARENT_TEST_ROLE"'
`,
		},
		{
			name: "grant database role to role (SHOW GRANTS OF ROLE output)",
			inputRows: [][]string{
				{"role", "granted_to", "grantee_name"},
				{"TEST_DATABASE.TEST_ROLE", "ROLE", "PARENT_TEST_ROLE"},
			},
			expectedOutput: `
resource "snowflake_grant_database_role" "snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_role_PARENT_TEST_ROLE" {
  database_role_name = "\"TEST_DATABASE\".\"TEST_ROLE\""
  parent_role_name = "PARENT_TEST_ROLE"
}
# terraform import snowflake_grant_database_role.snowflake_generated_grant_TEST_DATABASE_TEST_ROLE_to_role_PARENT_TEST_ROLE '"TEST_DATABASE"."TEST_ROLE"|ROLE|"PARENT_TEST_ROLE"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(&Config{
				ObjectType: ObjectTypeGrants,
				ImportFlag: ImportStatementTypeStatement,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), output)
		})
	}
}
