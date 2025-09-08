package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProgram(t *testing.T) {
	testCases := []struct {
		name   string
		args   []string
		input  string
		config *Config

		expectedExitCode  ExitCode
		expectedErrOutput string
		expectedOutput    string
	}{
		{
			name: "help flag",
			args: []string{"cmd", "-h"},
			expectedErrOutput: `Migration script's purpose is to generate terraform resources from existing Snowflake objects.
It operates on STDIN input and expects output from Snowflake commands in the CSV format.
The script writes generated terraform resources to STDOUT in case you want to redirect it to a file.
Any logs or errors are written to STDERR. You should separate outputs from STDOUT and STDERR when running the script (e.g. by redirecting STDOUT to a file)
to clearly see in case of any errors or skipped objects (due to, for example, incorrect or unexpected format).

usage: migration_script [-import=<statement|block>] <object_type>

import optional flag determines the output format for import statements. The possible values are:
	- "statement" will print appropriate terraform import command at the end of generated content (default) (see https://developer.hashicorp.com/terraform/cli/commands/import)
	- "block" will generate import block at the end of generated content (see https://developer.hashicorp.com/terraform/language/import)
	
object_type represents the type of Snowflake object you want to generate terraform resources for.
	It is a required positional argument and possible values are listed below.
	A given object_type corresponds to a specific Snowflake output expected as input to the script.
	Currently supported object types are:
		- "grants" which expects output from SHOW GRANTS command (see https://docs.snowflake.com/en/sql-reference/sql/show-grants) to generate new grant resources (see https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#mapping-from-old-grant-resources-to-the-new-ones).
			The allowed SHOW GRANTS commands are:
				- 'SHOW GRANTS ON ACCOUNT'
				- 'SHOW GRANTS ON <object_type>'
				- 'SHOW GRANTS TO ROLE <role_name>'
				- 'SHOW GRANTS TO DATABASE ROLE <database_role_name>'
			Supported resources:
				- snowflake_grant_privileges_to_account_role
				- snowflake_grant_privileges_to_database_role
				- snowflake_grant_account_role
				- snowflake_grant_database_role
			Limitations:
				- grants on 'future' or on 'all' objects are not supported
				- all_privileges and always_apply fields are not supported
		
example usage:
	migration_script -import=block grants < show_grants_output.csv > generated_output.tf
`,
		},
		{
			name:              "validation: missing object type",
			args:              []string{"cmd"},
			expectedErrOutput: `Error parsing input arguments: no object type specified, use -h for help, run -h to get more information on running the script`,
			expectedExitCode:  ExitCodeFailedInputArgumentParsing,
		},
		{
			name:              "validation: invalid object type",
			args:              []string{"cmd", "invalid-object-type"},
			expectedErrOutput: `Error parsing input arguments: error parsing object type: unsupported object type: invalid-object-type, run -h to get more information on running the script`,
			expectedExitCode:  ExitCodeFailedInputArgumentParsing,
		},
		{
			name:              "validation: invalid import format",
			args:              []string{"cmd", "-import=invalid_import_format", "grants"},
			expectedErrOutput: `Error parsing input arguments: error parsing import flag: invalid import statement type: invalid_import_format, run -h to get more information on running the script`,
			expectedExitCode:  ExitCodeFailedInputArgumentParsing,
		},
		{
			name:              "validation: invalid arg order",
			args:              []string{"cmd", "grants", "-import=invalid_import_format"},
			expectedErrOutput: `Error parsing input arguments: no object type specified, use -h for help, run -h to get more information on running the script`,
			expectedExitCode:  ExitCodeFailedInputArgumentParsing,
		},
		{
			name: "validation: invalid csv input",
			args: []string{"cmd", "grants"},
			input: `privilege,granted_on
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedErrOutput: `Error reading CSV input: record on line 2: wrong number of fields`,
			expectedExitCode:  ExitCodeFailedCsvInputParsing,
		},
		{
			name: "basic usage",
			args: []string{"cmd", "grants"},
			input: `privilege,granted_on,name,granted_to,grantee_name,with_grant_option
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option '"ROLE_NAME"|false|false|CREATE DATABASE|OnAccount'
`,
		},
		{
			name: "basic usage - explicit statement import format",
			args: []string{"cmd", "-import=statement", "grants"},
			input: `privilege,granted_on,name,granted_to,grantee_name,with_grant_option
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option '"ROLE_NAME"|false|false|CREATE DATABASE|OnAccount'
`,
		},
		{
			name: "basic usage - block import format",
			args: []string{"cmd", "-import=block", "grants"},
			input: `privilege,granted_on,name,granted_to,grantee_name,with_grant_option
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_ROLE_NAME_without_grant_option
  id = "\"ROLE_NAME\"|false|false|CREATE DATABASE|OnAccount"
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			program := Program{
				Args:   tc.args,
				StdOut: bytes.NewBuffer(nil),
				StdErr: bytes.NewBuffer(nil),
				StdIn:  bytes.NewBufferString(tc.input),
				Config: tc.config,
			}

			assert.Equal(t, tc.expectedExitCode, program.Run())
			assert.Equal(t, tc.expectedOutput, program.StdOut.(*bytes.Buffer).String())
			assert.Equal(t, tc.expectedErrOutput, program.StdErr.(*bytes.Buffer).String())
		})
	}
}
