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
			The script always outputs fields that have non-empty default values in Snowflake (they can be removed from the output)
			Caution: Some of the fields are not supported (actives, pendings, failed, suspended, uuid, initially_suspended)
			For more details about using multiple sources, visit https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/scripts/migration_script/README.md#multiple-sources
			Supported resources:
				- snowflake_warehouse
		- "account_roles" which expects input in the form of [SHOW ROLES](https://docs.snowflake.com/en/sql-reference/sql/show-roles) output. Can also be obtained as a converted CSV output from the snowflake_account_roles data source.
			Supported resources:
				- snowflake_account_role
		- "database_roles" which expects input in the form of [SHOW DATABASE ROLES](https://docs.snowflake.com/en/sql-reference/sql/show-database-roles) output. Can also be obtained as a converted CSV output from the snowflake_database_roles data source.
			Supported resources:
				- snowflake_database_role
		- "users" which expects input in the form of [SHOW USERS](https://docs.snowflake.com/en/sql-reference/sql/show-users) output. Can also be obtained as a converted CSV output from the snowflake_users data source.
			To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW USERS output.
			Different user types (PERSON, SERVICE, LEGACY_SERVICE) are mapped to their respective terraform resources:
				- PERSON (or empty) -> snowflake_user
				- SERVICE -> snowflake_service_user
				- LEGACY_SERVICE -> snowflake_legacy_service_user
			Supported resources:
				- snowflake_user
				- snowflake_service_user
				- snowflake_legacy_service_user

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
"true","600","100","","2024-06-06 00:00:00.000 +0000 UTC","false","","false","false","1","1","WH_BASIC","0","ADMIN","ROLE","0","8","0","0","","","2024-06-06 12:00:00.000 +0000 UTC","0","STANDARD","XSMALL","1","AVAILABLE","STANDARD","2024-06-06 00:00:00.000 +0000 UTC","","","","","",""
"false","1200","80.00","Production warehouse with all fields","2024-06-06 00:00:00.000 +0000 UTC","true","","false","false","4","2","WH_COMPLETE","0","ADMIN","ROLE","0","16","1","0","MEMORY_16X","MONITOR1","2024-06-06 12:00:00.000 +0000 UTC","3","ECONOMY","MEDIUM","1","SUSPENDED","SNOWPARK-OPTIMIZED","2024-06-06 00:00:00.000 +0000 UTC","WAREHOUSE","8","WAREHOUSE","300","WAREHOUSE","86400"
"true","600","100","Gen2 warehouse","2024-06-06 00:00:00.000 +0000 UTC","false","2","false","false","1","1","WH_GEN2","0","ADMIN","ROLE","0","8","0","0","","","2024-06-06 12:00:00.000 +0000 UTC","0","STANDARD","LARGE","1","AVAILABLE","STANDARD","2024-06-06 00:00:00.000 +0000 UTC","","","","","",""`,
			expectedOutput: `resource "snowflake_warehouse" "snowflake_generated_warehouse_WH_BASIC" {
  name = "WH_BASIC"
  auto_resume = "true"
  auto_suspend = 600
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_WH_COMPLETE" {
  name = "WH_COMPLETE"
  auto_resume = "false"
  auto_suspend = 1200
  comment = "Production warehouse with all fields"
  enable_query_acceleration = "true"
  max_cluster_count = 4
  max_concurrency_level = 8
  min_cluster_count = 2
  query_acceleration_max_scale_factor = 16
  resource_constraint = "MEMORY_16X"
  resource_monitor = "MONITOR1"
  scaling_policy = "ECONOMY"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 86400
  warehouse_size = "MEDIUM"
  warehouse_type = "SNOWPARK-OPTIMIZED"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_WH_GEN2" {
  name = "WH_GEN2"
  auto_resume = "true"
  auto_suspend = 600
  comment = "Gen2 warehouse"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "LARGE"
  warehouse_type = "STANDARD"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH_BASIC
  id = "\"WH_BASIC\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH_COMPLETE
  id = "\"WH_COMPLETE\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH_GEN2
  id = "\"WH_GEN2\""
}
`,
		},
		{
			name: "basic usage - block import format for account_roles",
			args: []string{"cmd", "-import=block", "account_roles"},
			input: `
"assigned_to_users","comment","created_on","granted_roles","granted_to_roles","is_current","is_default","is_inherited","name","owner"
"0","","2024-06-06 00:00:00.000 +0000 UTC","0","0","N","N","N","MINIMAL_ROLE","ACCOUNTADMIN"
"5","This is a test role","2024-06-06 00:00:00.000 +0000 UTC","2","1","Y","N","Y","ADMIN_ROLE","ACCOUNTADMIN"`,
			expectedOutput: `resource "snowflake_account_role" "snowflake_generated_account_role_MINIMAL_ROLE" {
  name = "MINIMAL_ROLE"
}

resource "snowflake_account_role" "snowflake_generated_account_role_ADMIN_ROLE" {
  name = "ADMIN_ROLE"
  comment = "This is a test role"
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MINIMAL_ROLE
  id = "\"MINIMAL_ROLE\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_ADMIN_ROLE
  id = "\"ADMIN_ROLE\""
}
`,
		},
		{
			name: "basic usage - block import format for database_roles",
			args: []string{"cmd", "-import=block", "database_roles"},
			input: `
"comment","created_on","database_name","granted_database_roles","granted_to_database_roles","granted_to_roles","is_current","is_default","is_inherited","name","owner","owner_role_type"
"","2024-06-06 00:00:00.000 +0000 UTC","TEST_DB","0","0","0","N","N","N","MINIMAL_ROLE","ACCOUNTADMIN","ROLE"
"This is a database role","2024-06-06 00:00:00.000 +0000 UTC","PROD_DB","2","1","3","Y","N","Y","ADMIN_ROLE","ACCOUNTADMIN","ROLE"`,
			expectedOutput: `resource "snowflake_database_role" "snowflake_generated_database_role_TEST_DB_MINIMAL_ROLE" {
  database = "TEST_DB"
  name = "MINIMAL_ROLE"
}

resource "snowflake_database_role" "snowflake_generated_database_role_PROD_DB_ADMIN_ROLE" {
  database = "PROD_DB"
  name = "ADMIN_ROLE"
  comment = "This is a database role"
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_TEST_DB_MINIMAL_ROLE
  id = "\"TEST_DB\".\"MINIMAL_ROLE\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_PROD_DB_ADMIN_ROLE
  id = "\"PROD_DB\".\"ADMIN_ROLE\""
}
`,
		},
		{
			name: "basic usage - block import format for users with all types",
			args: []string{"cmd", "-import=block", "users"},
			input: `
"comment","default_namespace","default_role","default_secondary_roles","default_warehouse","disabled","display_name","email","first_name","last_name","login_name","must_change_password","name","type"
"","","","","","","","","","","","","BASIC_USER","PERSON"
"","","","","","","","","","","","","EMPTY_TYPE_USER",""
"Service user","DB.SCHEMA","SVC_ROLE","ALL","SVC_WH","false","Service User","svc@example.com","","","svc_login","","SERVICE_USER","SERVICE"
"Legacy user","DB.SCHEMA","LEGACY_ROLE","ALL","LEGACY_WH","true","Legacy User","legacy@example.com","","","legacy_login","true","LEGACY_USER","LEGACY_SERVICE"
"Full person","DB.SCHEMA","ANALYST","ALL","COMPUTE_WH","false","John Doe","john@example.com","John","Doe","john_login","false","COMPLETE_USER","PERSON"`,
			expectedOutput: `resource "snowflake_user" "snowflake_generated_user_BASIC_USER" {
  name = "BASIC_USER"
}

resource "snowflake_user" "snowflake_generated_user_EMPTY_TYPE_USER" {
  name = "EMPTY_TYPE_USER"
}

resource "snowflake_service_user" "snowflake_generated_service_user_SERVICE_USER" {
  name = "SERVICE_USER"
  comment = "Service user"
  default_namespace = "DB.SCHEMA"
  default_role = "SVC_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "SVC_WH"
  display_name = "Service User"
  email = "svc@example.com"
  login_name = "svc_login"
}

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_USER" {
  name = "LEGACY_USER"
  comment = "Legacy user"
  default_namespace = "DB.SCHEMA"
  default_role = "LEGACY_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "LEGACY_WH"
  disabled = "true"
  display_name = "Legacy User"
  email = "legacy@example.com"
  login_name = "legacy_login"
  must_change_password = "true"
}

resource "snowflake_user" "snowflake_generated_user_COMPLETE_USER" {
  name = "COMPLETE_USER"
  comment = "Full person"
  default_namespace = "DB.SCHEMA"
  default_role = "ANALYST"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  display_name = "John Doe"
  email = "john@example.com"
  first_name = "John"
  last_name = "Doe"
  login_name = "john_login"
}
import {
  to = snowflake_user.snowflake_generated_user_BASIC_USER
  id = "\"BASIC_USER\""
}
import {
  to = snowflake_user.snowflake_generated_user_EMPTY_TYPE_USER
  id = "\"EMPTY_TYPE_USER\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SERVICE_USER
  id = "\"SERVICE_USER\""
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_USER
  id = "\"LEGACY_USER\""
}
import {
  to = snowflake_user.snowflake_generated_user_COMPLETE_USER
  id = "\"COMPLETE_USER\""
}
`,
		},
		{
			name: "basic usage - statement import format for users",
			args: []string{"cmd", "-import=statement", "users"},
			input: `
"name","type"
"TEST_USER","PERSON"
"SVC_USER","SERVICE"
"LEGACY_USER","LEGACY_SERVICE"`,
			expectedOutput: `resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
}

resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
}

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_USER" {
  name = "LEGACY_USER"
}
# terraform import snowflake_user.snowflake_generated_user_TEST_USER '"TEST_USER"'
# terraform import snowflake_service_user.snowflake_generated_service_user_SVC_USER '"SVC_USER"'
# terraform import snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_USER '"LEGACY_USER"'
`,
		},
		{
			name: "users with lowercase types",
			args: []string{"cmd", "-import=block", "users"},
			input: `
"name","type"
"USER_1","person"
"USER_2","service"
"USER_3","legacy_service"`,
			expectedOutput: `resource "snowflake_user" "snowflake_generated_user_USER_1" {
  name = "USER_1"
}

resource "snowflake_service_user" "snowflake_generated_service_user_USER_2" {
  name = "USER_2"
}

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_USER_3" {
  name = "USER_3"
}
import {
  to = snowflake_user.snowflake_generated_user_USER_1
  id = "\"USER_1\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_USER_2
  id = "\"USER_2\""
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_USER_3
  id = "\"USER_3\""
}
`,
		},
		{
			name: "users with parameters",
			args: []string{"cmd", "-import=block", "users"},
			input: `
"name","type","abort_detached_query_level","abort_detached_query_value","timezone_level","timezone_value","statement_timeout_in_seconds_level","statement_timeout_in_seconds_value"
"PARAM_USER","PERSON","USER","true","USER","America/New_York","",""
"SVC_PARAM_USER","SERVICE","","","","","USER","3600"`,
			expectedOutput: `resource "snowflake_user" "snowflake_generated_user_PARAM_USER" {
  name = "PARAM_USER"
  abort_detached_query = true
  timezone = "America/New_York"
}

resource "snowflake_service_user" "snowflake_generated_service_user_SVC_PARAM_USER" {
  name = "SVC_PARAM_USER"
  statement_timeout_in_seconds = 3600
}
import {
  to = snowflake_user.snowflake_generated_user_PARAM_USER
  id = "\"PARAM_USER\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_PARAM_USER
  id = "\"SVC_PARAM_USER\""
}
`,
		},
		{
			name: "user with all supported parameters",
			args: []string{"cmd", "-import=block", "users"},
			input: `
"comment","default_namespace","default_role","default_secondary_roles","default_warehouse","disabled","display_name","email","first_name","last_name","login_name","must_change_password","name","type","abort_detached_query_level","abort_detached_query_value","autocommit_level","autocommit_value","binary_input_format_level","binary_input_format_value","binary_output_format_level","binary_output_format_value","client_memory_limit_level","client_memory_limit_value","client_metadata_request_use_connection_ctx_level","client_metadata_request_use_connection_ctx_value","client_prefetch_threads_level","client_prefetch_threads_value","client_result_chunk_size_level","client_result_chunk_size_value","client_result_column_case_insensitive_level","client_result_column_case_insensitive_value","client_session_keep_alive_level","client_session_keep_alive_value","client_session_keep_alive_heartbeat_frequency_level","client_session_keep_alive_heartbeat_frequency_value","client_timestamp_type_mapping_level","client_timestamp_type_mapping_value","date_input_format_level","date_input_format_value","date_output_format_level","date_output_format_value","enable_unload_physical_type_optimization_level","enable_unload_physical_type_optimization_value","enable_unredacted_query_syntax_error_level","enable_unredacted_query_syntax_error_value","error_on_nondeterministic_merge_level","error_on_nondeterministic_merge_value","error_on_nondeterministic_update_level","error_on_nondeterministic_update_value","geography_output_format_level","geography_output_format_value","geometry_output_format_level","geometry_output_format_value","jdbc_treat_decimal_as_int_level","jdbc_treat_decimal_as_int_value","jdbc_treat_timestamp_ntz_as_utc_level","jdbc_treat_timestamp_ntz_as_utc_value","jdbc_use_session_timezone_level","jdbc_use_session_timezone_value","json_indent_level","json_indent_value","lock_timeout_level","lock_timeout_value","log_level_level","log_level_value","multi_statement_count_level","multi_statement_count_value","network_policy_level","network_policy_value","noorder_sequence_as_default_level","noorder_sequence_as_default_value","odbc_treat_decimal_as_int_level","odbc_treat_decimal_as_int_value","prevent_unload_to_internal_stages_level","prevent_unload_to_internal_stages_value","query_tag_level","query_tag_value","quoted_identifiers_ignore_case_level","quoted_identifiers_ignore_case_value","rows_per_resultset_level","rows_per_resultset_value","s3_stage_vpce_dns_name_level","s3_stage_vpce_dns_name_value","search_path_level","search_path_value","simulated_data_sharing_consumer_level","simulated_data_sharing_consumer_value","statement_queued_timeout_in_seconds_level","statement_queued_timeout_in_seconds_value","statement_timeout_in_seconds_level","statement_timeout_in_seconds_value","strict_json_output_level","strict_json_output_value","time_input_format_level","time_input_format_value","time_output_format_level","time_output_format_value","timestamp_day_is_always24h_level","timestamp_day_is_always24h_value","timestamp_input_format_level","timestamp_input_format_value","timestamp_l_t_z_output_format_level","timestamp_l_t_z_output_format_value","timestamp_n_t_z_output_format_level","timestamp_n_t_z_output_format_value","timestamp_output_format_level","timestamp_output_format_value","timestamp_t_z_output_format_level","timestamp_t_z_output_format_value","timestamp_type_mapping_level","timestamp_type_mapping_value","timezone_level","timezone_value","trace_level_level","trace_level_value","transaction_abort_on_error_level","transaction_abort_on_error_value","transaction_default_isolation_level_level","transaction_default_isolation_level_value","two_digit_century_start_level","two_digit_century_start_value","unsupported_d_d_l_action_level","unsupported_d_d_l_action_value","use_cached_result_level","use_cached_result_value","week_of_year_policy_level","week_of_year_policy_value","week_start_level","week_start_value"
"Full user with all params","DB.SCHEMA","ANALYST","ALL","COMPUTE_WH","false","John Doe","john@example.com","John","Doe","john_login","false","ALL_PARAMS_USER","PERSON","USER","true","USER","false","USER","HEX","USER","BASE64","USER","2048","USER","true","USER","8","USER","128","USER","true","USER","true","USER","1800","USER","TIMESTAMP_NTZ","USER","YYYY-MM-DD","USER","YYYY-MM-DD","USER","true","USER","true","USER","false","USER","false","USER","GeoJSON","USER","WKT","USER","false","USER","true","USER","true","USER","4","USER","43200","USER","INFO","USER","0","USER","NETWORK_POLICY_1","USER","true","USER","false","USER","true","USER","my_query_tag","USER","true","USER","100000","USER","vpce-dns-name","USER","$current, $public","USER","CONSUMER_ACCOUNT","USER","300","USER","86400","USER","true","USER","HH24:MI:SS","USER","HH24:MI:SS","USER","true","USER","YYYY-MM-DD HH24:MI:SS","USER","YYYY-MM-DD HH24:MI:SS TZHTZM","USER","YYYY-MM-DD HH24:MI:SS","USER","YYYY-MM-DD HH24:MI:SS.FF3","USER","YYYY-MM-DD HH24:MI:SS TZHTZM","USER","TIMESTAMP_NTZ","USER","America/New_York","USER","ON_EVENT","USER","true","USER","READ COMMITTED","USER","1970","USER","IGNORE","USER","false","USER","1","USER","1"`,
			expectedOutput: `resource "snowflake_user" "snowflake_generated_user_ALL_PARAMS_USER" {
  name = "ALL_PARAMS_USER"
  abort_detached_query = true
  autocommit = false
  binary_input_format = "HEX"
  binary_output_format = "BASE64"
  client_memory_limit = 2048
  client_metadata_request_use_connection_ctx = true
  client_prefetch_threads = 8
  client_result_chunk_size = 128
  client_result_column_case_insensitive = true
  client_session_keep_alive = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping = "TIMESTAMP_NTZ"
  comment = "Full user with all params"
  date_input_format = "YYYY-MM-DD"
  date_output_format = "YYYY-MM-DD"
  default_namespace = "DB.SCHEMA"
  default_role = "ANALYST"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  display_name = "John Doe"
  email = "john@example.com"
  enable_unload_physical_type_optimization = true
  enable_unredacted_query_syntax_error = true
  error_on_nondeterministic_merge = false
  error_on_nondeterministic_update = false
  first_name = "John"
  geography_output_format = "GeoJSON"
  geometry_output_format = "WKT"
  jdbc_treat_decimal_as_int = false
  jdbc_treat_timestamp_ntz_as_utc = true
  jdbc_use_session_timezone = true
  json_indent = 4
  last_name = "Doe"
  lock_timeout = 43200
  log_level = "INFO"
  login_name = "john_login"
  multi_statement_count = 0
  network_policy = "NETWORK_POLICY_1"
  noorder_sequence_as_default = true
  odbc_treat_decimal_as_int = false
  prevent_unload_to_internal_stages = true
  query_tag = "my_query_tag"
  quoted_identifiers_ignore_case = true
  rows_per_resultset = 100000
  s3_stage_vpce_dns_name = "vpce-dns-name"
  search_path = "$current, $public"
  simulated_data_sharing_consumer = "CONSUMER_ACCOUNT"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 86400
  strict_json_output = true
  time_input_format = "HH24:MI:SS"
  time_output_format = "HH24:MI:SS"
  timestamp_day_is_always_24h = true
  timestamp_input_format = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ltz_output_format = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_ntz_output_format = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_type_mapping = "TIMESTAMP_NTZ"
  timestamp_tz_output_format = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timezone = "America/New_York"
  trace_level = "ON_EVENT"
  transaction_abort_on_error = true
  transaction_default_isolation_level = "READ COMMITTED"
  two_digit_century_start = 1970
  unsupported_ddl_action = "IGNORE"
  use_cached_result = false
  week_of_year_policy = 1
  week_start = 1
}
import {
  to = snowflake_user.snowflake_generated_user_ALL_PARAMS_USER
  id = "\"ALL_PARAMS_USER\""
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
