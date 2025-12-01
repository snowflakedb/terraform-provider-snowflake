package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
		- "schemas" which expects a converted CSV output from the snowflake_schemas data source
			To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW SCHEMAS output, so the CSV header looks like "comment","created_on",...,"catalog_value","catalog_level","data_retention_time_in_days_value","data_retention_time_in_days_level",...
			When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "SCHEMA".
			For more details about using multiple sources, visit https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/scripts/migration_script/README.md#multiple-sources
			Supported resources:
				- snowflake_schema
		- "databases" which expects a converted CSV output from the snowflake_databases data source
			To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW DATABASES output, so the CSV header looks like "comment","created_on",...,"catalog_value","catalog_level","data_retention_time_in_days_value","data_retention_time_in_days_level",...
			When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "DATABASE".
			For more details about using multiple sources, visit https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/scripts/migration_script/README.md#multiple-sources
			Warning: currently secondary databases and shared databases are treated as plain databases.
			Supported resources:
				- snowflake_database
		- "warehouses" which expects a converted CSV output from the snowflake_warehouses data source
			To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW WAREHOUSES output, so the CSV header looks like "comment","created_on",...,"max_cluster_count","min_cluster_count","name","other",...
			When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "WAREHOUSE".
			For more details about using multiple sources, visit https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/scripts/migration_script/README.md#multiple-sources
			Supported resources:
				- snowflake_warehouse

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
		{
			name: "basic usage - block import format for schemas",
			args: []string{"cmd", "-import=block", "schemas"},
			input: `
"comment","created_on","database_name","dropped_on","is_current","is_default","name","options","owner","owner_role_type","retention_time","catalog_value","catalog_level","data_retention_time_in_days_value","data_retention_time_in_days_level","default_ddl_collation_value","default_ddl_collation_level","enable_console_output_value","enable_console_output_level","external_volume_value","external_volume_level","log_level_value","log_level_level","max_data_extension_time_in_days_value","max_data_extension_time_in_days_level","pipe_execution_paused_value","pipe_execution_paused_level","quoted_identifiers_ignore_case_value","quoted_identifiers_ignore_case_level","replace_invalid_characters_value","replace_invalid_characters_level","storage_serialization_policy_value","storage_serialization_policy_level","suspend_task_after_num_failures_value","suspend_task_after_num_failures_level","task_auto_retry_attempts_value","task_auto_retry_attempts_level","trace_level_value","trace_level_level","user_task_managed_initial_warehouse_size_value","user_task_managed_initial_warehouse_size_level","user_task_minimum_trigger_interval_in_seconds_value","user_task_minimum_trigger_interval_in_seconds_level","user_task_timeout_ms_value","user_task_timeout_ms_level"
"","2025-11-20 04:53:26.906 -0800 PST","DATABASE","0001-01-01 00:00:00 +0000 UTC",false,false,"BASIC","","ACCOUNTADMIN","ROLE","1","","","1","","","","false","","","","OFF","","14","","false","","false","","false","","OPTIMIZED","","10","","0","","OFF","","Medium","","30","","3600000",""
"comment","2025-11-20 04:53:25.625 -0800 PST","DATABASE","0001-01-01 00:00:00 +0000 UTC",true,true,"COMPLETE","TRANSIENT, MANAGED ACCESS","ACCOUNTADMIN","ROLE","1","CATALOG","SCHEMA","1","SCHEMA","en_US-trim","SCHEMA","true","SCHEMA","EXTERNAL_VOLUME","SCHEMA","INFO","SCHEMA","10","SCHEMA","true","SCHEMA","true","SCHEMA","true","SCHEMA","COMPATIBLE","SCHEMA","10","SCHEMA","10","SCHEMA","PROPAGATE","SCHEMA","MEDIUM","SCHEMA","30","SCHEMA","3600000","SCHEMA"`,
			expectedOutput: `resource "snowflake_schema" "snowflake_generated_schema_DATABASE_BASIC" {
  database = "DATABASE"
  name = "BASIC"
}

resource "snowflake_schema" "snowflake_generated_schema_DATABASE_COMPLETE" {
  database = "DATABASE"
  name = "COMPLETE"
  catalog = "CATALOG"
  comment = "comment"
  data_retention_time_in_days = 1
  default_ddl_collation = "en_US-trim"
  enable_console_output = true
  external_volume = "EXTERNAL_VOLUME"
  is_transient = "true"
  log_level = "INFO"
  max_data_extension_time_in_days = 10
  pipe_execution_paused = true
  quoted_identifiers_ignore_case = true
  replace_invalid_characters = true
  storage_serialization_policy = "COMPATIBLE"
  suspend_task_after_num_failures = 10
  task_auto_retry_attempts = 10
  trace_level = "PROPAGATE"
  user_task_managed_initial_warehouse_size = "MEDIUM"
  user_task_minimum_trigger_interval_in_seconds = 30
  user_task_timeout_ms = 3600000
  with_managed_access = "true"
}
import {
  to = snowflake_schema.snowflake_generated_schema_DATABASE_BASIC
  id = "\"DATABASE\".\"BASIC\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_DATABASE_COMPLETE
  id = "\"DATABASE\".\"COMPLETE\""
}
`,
		},
		{
			name: "basic usage - block import format for databases",
			args: []string{"cmd", "-import=block", "databases"},
			input: `
"comment","created_on","dropped_on","is_current","is_default","kind","name","options","origin","owner","owner_role_type","resource_group","retention_time","catalog_level","catalog_value","data_retention_time_in_days_level","data_retention_time_in_days_value","default_ddl_collation_level","default_ddl_collation_value","enable_console_output_level","enable_console_output_value","external_volume_level","external_volume_value","log_level_level","log_level_value","max_data_extension_time_in_days_level","max_data_extension_time_in_days_value","quoted_identifiers_ignore_case_level","quoted_identifiers_ignore_case_value","replace_invalid_characters_level","replace_invalid_characters_value","storage_serialization_policy_level","storage_serialization_policy_value","suspend_task_after_num_failures_level","suspend_task_after_num_failures_value","task_auto_retry_attempts_level","task_auto_retry_attempts_value","trace_level_level","trace_level_value","user_task_managed_initial_warehouse_size_level","user_task_managed_initial_warehouse_size_value","user_task_minimum_trigger_interval_in_seconds_level","user_task_minimum_trigger_interval_in_seconds_value","user_task_timeout_ms_level","user_task_timeout_ms_value"
"","2025-11-20 04:53:26.906 -0800 PST","0001-01-01 00:00:00 +0000 UTC","false","false","STANDARD","BASIC","","","ACCOUNTADMIN","ROLE","","1","","","","1","","","","false","","","","OFF","","14","","false","","false","","OPTIMIZED","","10","","0","","OFF","","Medium","","30","","3600000"
"comment","2025-11-20 04:53:25.625 -0800 PST","0001-01-01 00:00:00 +0000 UTC","true","true","STANDARD","COMPLETE","TRANSIENT","","ACCOUNTADMIN","ROLE","","1","DATABASE","CATALOG","DATABASE","1","DATABASE","en_US-trim","DATABASE","true","DATABASE","EXTERNAL_VOLUME","DATABASE","INFO","DATABASE","10","DATABASE","true","DATABASE","true","DATABASE","COMPATIBLE","DATABASE","10","DATABASE","10","DATABASE","PROPAGATE","DATABASE","MEDIUM","DATABASE","30","DATABASE","3600000"`,
			expectedOutput: `resource "snowflake_database" "snowflake_generated_database_BASIC" {
  name = "BASIC"
}

resource "snowflake_database" "snowflake_generated_database_COMPLETE" {
  name = "COMPLETE"
  catalog = "CATALOG"
  comment = "comment"
  data_retention_time_in_days = 1
  default_ddl_collation = "en_US-trim"
  enable_console_output = true
  external_volume = "EXTERNAL_VOLUME"
  is_transient = true
  log_level = "INFO"
  max_data_extension_time_in_days = 10
  quoted_identifiers_ignore_case = true
  replace_invalid_characters = true
  storage_serialization_policy = "COMPATIBLE"
  suspend_task_after_num_failures = 10
  task_auto_retry_attempts = 10
  trace_level = "PROPAGATE"
  user_task_managed_initial_warehouse_size = "MEDIUM"
  user_task_minimum_trigger_interval_in_seconds = 30
  user_task_timeout_ms = 3600000
}
import {
  to = snowflake_database.snowflake_generated_database_BASIC
  id = "\"BASIC\""
}
import {
  to = snowflake_database.snowflake_generated_database_COMPLETE
  id = "\"COMPLETE\""
}
`,
		},
		{
			name: "basic usage - block import format for warehouses",
			args: []string{"cmd", "-import=block", "warehouses"},
			input: `
"auto_resume","auto_suspend","available","comment","created_on","enable_query_acceleration","generation","is_current","is_default","max_cluster_count","min_cluster_count","name","other","owner","owner_role_type","provisioning","query_acceleration_max_scale_factor","queued","quiescing","resource_constraint","resource_monitor","resumed_on","running","scaling_policy","size","started_clusters","state","type","updated_on","max_concurrency_level_level","max_concurrency_level_value","statement_queued_timeout_in_seconds_level","statement_queued_timeout_in_seconds_value","statement_timeout_in_seconds_level","statement_timeout_in_seconds_value"
"true","600","71.43","","2024-06-06 00:00:00.000 +0000 UTC","false","","false","false","3","1","WH_BASIC","0","ADMIN","ROLE","0","0","1","0","","","2024-06-06 12:00:00.000 +0000 UTC","5","ECONOMY","XSMALL","2","AVAILABLE","STANDARD","2024-06-06 00:00:00.000 +0000 UTC","","","","","",""
"true","450","80.00","Production warehouse","2024-06-06 00:00:00.000 +0000 UTC","true","1","false","false","4","1","WH_PROD","0","ADMIN","ROLE","0","30","1","0","MEMORY_2X","MONITOR1","2024-06-06 12:00:00.000 +0000 UTC","3","CLASSIC","MEDIUM","1","SUSPENDED","SNOWPARK-OPTIMIZED","2024-06-06 00:00:00.000 +0000 UTC","WAREHOUSE","10","WAREHOUSE","600","WAREHOUSE","300"`,
			expectedOutput: `resource "snowflake_warehouse" "snowflake_generated_warehouse_WH_BASIC" {
  name = "WH_BASIC"
  max_cluster_count = 3
  query_acceleration_max_scale_factor = 0
  scaling_policy = "ECONOMY"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_WH_PROD" {
  name = "WH_PROD"
  auto_suspend = 450
  comment = "Production warehouse"
  enable_query_acceleration = "true"
  max_cluster_count = 4
  max_concurrency_level = 10
  query_acceleration_max_scale_factor = 30
  resource_monitor = "MONITOR1"
  scaling_policy = "CLASSIC"
  statement_queued_timeout_in_seconds = 600
  statement_timeout_in_seconds = 300
  warehouse_size = "MEDIUM"
  warehouse_type = "SNOWPARK-OPTIMIZED"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH_BASIC
  id = "\"WH_BASIC\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH_PROD
  id = "\"WH_PROD\""
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
