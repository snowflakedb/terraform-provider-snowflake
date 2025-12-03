package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleUserValidations(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		outputContains []string
	}{
		{
			name: "multiple users of different types",
			inputRows: [][]string{
				{"comment", "name", "type", "login_name"},
				{"person comment", "TEST_USER_PERSON", "PERSON", "test_person_login"},
				{"service comment", "TEST_USER_SERVICE", "SERVICE", "test_service_login"},
				{"legacy comment", "TEST_USER_LEGACY", "LEGACY_SERVICE", "test_legacy_login"},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER_PERSON"`,
				`resource "snowflake_service_user" "snowflake_generated_service_user_TEST_USER_SERVICE"`,
				`resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_TEST_USER_LEGACY"`,
			},
		},
		{
			name: "empty type defaults to PERSON",
			inputRows: [][]string{
				{"comment", "name", "type"},
				{"", "TEST_USER", ""},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER"`,
			},
		},
		{
			name: "lowercase type is handled correctly",
			inputRows: [][]string{
				{"comment", "name", "type"},
				{"", "TEST_USER_1", "person"},
				{"", "TEST_USER_2", "service"},
				{"", "TEST_USER_3", "legacy_service"},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER_1"`,
				`resource "snowflake_service_user" "snowflake_generated_service_user_TEST_USER_2"`,
				`resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_TEST_USER_3"`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			for _, expectedOutput := range tc.outputContains {
				assert.Contains(t, output, expectedOutput)
			}
		})
	}
}

func TestHandleUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic user (PERSON)",
			inputRows: [][]string{
				{"name", "type"},
				{"TEST_USER", "PERSON"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
		{
			name: "complete user (PERSON) with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "first_name", "last_name", "login_name", "must_change_password", "name", "type"},
				{"User comment", "DB.SCHEMA", "ANALYST", "ALL", "COMPUTE_WH", "true", "Test User Display", "test@example.com", "John", "Doe", "test_login", "true", "TEST_USER", "PERSON"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
  comment = "User comment"
  default_namespace = "DB.SCHEMA"
  default_role = "ANALYST"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  disabled = "true"
  display_name = "Test User Display"
  email = "test@example.com"
  first_name = "John"
  last_name = "Doe"
  login_name = "test_login"
  must_change_password = "true"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
		{
			name: "user with parameters",
			inputRows: [][]string{
				{"name", "type", "abort_detached_query_level", "abort_detached_query_value", "timezone_level", "timezone_value"},
				{"TEST_USER", "PERSON", "USER", "true", "USER", "America/New_York"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
  abort_detached_query = true
  timezone = "America/New_York"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleServiceUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic service user",
			inputRows: [][]string{
				{"name", "type"},
				{"SVC_USER", "SERVICE"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
		{
			name: "complete service user with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "login_name", "name", "type"},
				{"Service user comment", "DB.SCHEMA", "SVC_ROLE", "ALL", "SVC_WH", "false", "Service User", "svc@example.com", "svc_login", "SVC_USER", "SERVICE"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
  comment = "Service user comment"
  default_namespace = "DB.SCHEMA"
  default_role = "SVC_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "SVC_WH"
  display_name = "Service User"
  email = "svc@example.com"
  login_name = "svc_login"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
		{
			name: "service user with parameters",
			inputRows: [][]string{
				{"name", "type", "autocommit_level", "autocommit_value", "log_level_level", "log_level_value"},
				{"SVC_USER", "SERVICE", "USER", "false", "USER", "INFO"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
  autocommit = false
  log_level = "INFO"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleLegacyServiceUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic legacy service user",
			inputRows: [][]string{
				{"name", "type"},
				{"LEGACY_SVC_USER", "LEGACY_SERVICE"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
		{
			name: "complete legacy service user with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "login_name", "must_change_password", "name", "type"},
				{"Legacy service comment", "DB.SCHEMA", "LEGACY_ROLE", "ALL", "LEGACY_WH", "true", "Legacy Service User", "legacy@example.com", "legacy_login", "true", "LEGACY_SVC_USER", "LEGACY_SERVICE"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
  comment = "Legacy service comment"
  default_namespace = "DB.SCHEMA"
  default_role = "LEGACY_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "LEGACY_WH"
  disabled = "true"
  display_name = "Legacy Service User"
  email = "legacy@example.com"
  login_name = "legacy_login"
  must_change_password = "true"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
		{
			name: "legacy service user with parameters",
			inputRows: [][]string{
				{"name", "type", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value", "trace_level_level", "trace_level_value"},
				{"LEGACY_SVC_USER", "LEGACY_SERVICE", "USER", "3600", "USER", "ON_EVENT"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
  statement_timeout_in_seconds = 3600
  trace_level = "ON_EVENT"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleUserImportStatementFormat(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		importFlag     ImportStatementType
		expectedOutput string
	}{
		{
			name: "user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"TEST_USER", "PERSON"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
}
# terraform import snowflake_user.snowflake_generated_user_TEST_USER '"TEST_USER"'
`,
		},
		{
			name: "service user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"SVC_USER", "SERVICE"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
}
# terraform import snowflake_service_user.snowflake_generated_service_user_SVC_USER '"SVC_USER"'
`,
		},
		{
			name: "legacy service user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"LEGACY_USER", "LEGACY_SERVICE"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_USER" {
  name = "LEGACY_USER"
}
# terraform import snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_USER '"LEGACY_USER"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: tc.importFlag,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
