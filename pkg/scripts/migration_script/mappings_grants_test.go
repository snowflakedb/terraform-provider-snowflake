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
			outputContains: []string{`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
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
				`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "USER"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE_2"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE_2"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE ROLE"]
  with_grant_option = true
}`,
				`resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
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
			output, err := HandleGrants(tc.inputRows)

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
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "TEST_ROLE"
  on_account = true
  privileges = ["CREATE DATABASE", "CREATE ROLE"]
  with_grant_option = true
}
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
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["CREATE DATABASE ROLE", "CREATE SCHEMA"]
  with_grant_option = false
}
`,
		},
		{
			name: "grants on schema",
			inputRows: [][]string{
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option"},
				{"CREATE TABLE", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "ROLE", "TEST_ROLE", "false"},
				{"CREATE VIEW", "SCHEMA", "TEST_DATABASE.TEST_SCHEMA", "ROLE", "TEST_ROLE", "false"},
			},
			expectedOutput: `
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_schema" {
  account_role_name = "TEST_ROLE"
  on_schema {
    schema_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\""
  }
  privileges = ["CREATE TABLE", "CREATE VIEW"]
  with_grant_option = false
}
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
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_schema_object" {
  account_role_name = "TEST_ROLE"
  on_schema_object {
    object_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\".\"TEST_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["INSERT", "SELECT"]
  with_grant_option = false
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleGrants(tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), output)
		})
	}
}
