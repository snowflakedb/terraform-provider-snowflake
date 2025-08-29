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
			expectedErrOutput: `Usage: migration_script [flags] <object_type>

Object types:
	- grants

Flags:
  -import string
    	Determines the output format for import statements.
    	Possible values:
    		- "statement" will print appropriate terraform import statement at the end of generated content
    		- "block" will generate import block next to every generated resource
    	 (default "statement")
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
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.test_resource_name_on_account '"ROLE_NAME"|false|false|CREATE DATABASE|OnAccount'
`,
		},
		{
			name: "basic usage - explicit statement import format",
			args: []string{"cmd", "-import=statement", "grants"},
			input: `privilege,granted_on,name,granted_to,grantee_name,with_grant_option
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
# terraform import snowflake_grant_privileges_to_account_role.test_resource_name_on_account '"ROLE_NAME"|false|false|CREATE DATABASE|OnAccount'
`,
		},
		{
			name: "basic usage - block import format",
			args: []string{"cmd", "-import=block", "grants"},
			input: `privilege,granted_on,name,granted_to,grantee_name,with_grant_option
CREATE DATABASE,ACCOUNT,ACCOUNT_LOCATOR,ROLE,ROLE_NAME,false`,
			expectedOutput: `resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "ROLE_NAME"
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = false
}
import {
  to = snowflake_grant_privileges_to_account_role.test_resource_name_on_account
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
