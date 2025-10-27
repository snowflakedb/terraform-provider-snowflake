# Basic standalone task
resource "snowflake_task" "task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  schedule {
    minutes = 5
  }
  sql_statement = "select 1"
}

# Basic serverless task
resource "snowflake_task" "serverless_task" {
  database                                 = "database"
  schema                                   = "schema"
  name                                     = "task"
  user_task_managed_initial_warehouse_size = "XSMALL"
  started                                  = true
  schedule {
    minutes = 5
  }
  sql_statement = "select 1"
}

# Basic child task
resource "snowflake_task" "child_task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  # You can do it by referring to task by computed fully_qualified_name field or write the task name in manually if it's not managed by Terraform
  after         = [snowflake_task.root_task.fully_qualified_name, "<database_name>.<schema_name>.<root_task_name>"]
  sql_statement = "select 1"
}

# Basic finalizer task
resource "snowflake_task" "child_task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  # You can do it by referring to task by computed fully_qualified_name field or write the task name in manually if it's not managed by Terraform
  finalize      = snowflake_task.root_task.fully_qualified_name
  sql_statement = "select 1"
}

# Complete child task
resource "snowflake_task" "complete_child_task" {
  database      = "database"
  schema        = "schema"
  name          = "complete_child_task"
  started       = true
  sql_statement = "CALL my_stored_procedure()"

  after                       = [snowflake_task.root_task.fully_qualified_name]
  config                      = "{\"output_dir\": \"/temp/\", \"learning_rate\": 0.1}"
  allow_overlapping_execution = true
  error_integration           = snowflake_notification_integration.example.fully_qualified_name
  when                        = "SYSTEM$STREAM_HAS_DATA('my_stream')"
  comment                     = "complete child task with all parameters"
  target_completion_interval  = "10 MINUTES"

  # Task Parameters
  suspend_task_after_num_failures               = 5
  task_auto_retry_attempts                      = 3
  user_task_managed_initial_warehouse_size      = "MEDIUM"
  user_task_minimum_trigger_interval_in_seconds = 60
  user_task_timeout_ms                          = 7200000
  serverless_task_min_statement_size            = "SMALL"
  serverless_task_max_statement_size            = "LARGE"

  # Session Parameters
  abort_detached_query         = false
  autocommit                   = true
  binary_input_format          = "HEX"
  binary_output_format         = "HEX"
  json_indent                  = 4
  lock_timeout                 = 21600
  log_level                    = "INFO"
  query_tag                    = "child_task_query"
  statement_timeout_in_seconds = 86400
  timezone                     = "America/New_York"
}

# Complete standalone task
resource "snowflake_task" "test" {
  database      = "database"
  schema        = "schema"
  name          = "task"
  warehouse     = snowflake_warehouse.example.fully_qualified_name
  started       = true
  sql_statement = "select 1"

  config                      = "{\"key\":\"value\"}"
  allow_overlapping_execution = true
  error_integration           = snowflake_notification_integration.example.fully_qualified_name
  when                        = "SYSTEM$STREAM_HAS_DATA('<stream_name>')"
  comment                     = "complete task"

  schedule {
    minutes = 10
  }

  # Session Parameters
  suspend_task_after_num_failures               = 10
  task_auto_retry_attempts                      = 0
  user_task_managed_initial_warehouse_size      = "Medium"
  user_task_minimum_trigger_interval_in_seconds = 30
  user_task_timeout_ms                          = 3600000
  abort_detached_query                          = false
  autocommit                                    = true
  binary_input_format                           = "HEX"
  binary_output_format                          = "HEX"
  client_memory_limit                           = 1536
  client_metadata_request_use_connection_ctx    = false
  client_prefetch_threads                       = 4
  client_result_chunk_size                      = 160
  client_result_column_case_insensitive         = false
  client_session_keep_alive                     = false
  client_session_keep_alive_heartbeat_frequency = 3600
  client_timestamp_type_mapping                 = "TIMESTAMP_LTZ"
  date_input_format                             = "AUTO"
  date_output_format                            = "YYYY-MM-DD"
  enable_unload_physical_type_optimization      = true
  error_on_nondeterministic_merge               = true
  error_on_nondeterministic_update              = false
  geography_output_format                       = "GeoJSON"
  geometry_output_format                        = "GeoJSON"
  jdbc_use_session_timezone                     = true
  json_indent                                   = 2
  lock_timeout                                  = 43200
  log_level                                     = "OFF"
  multi_statement_count                         = 1
  noorder_sequence_as_default                   = true
  odbc_treat_decimal_as_int                     = false
  query_tag                                     = ""
  quoted_identifiers_ignore_case                = false
  rows_per_resultset                            = 0
  s3_stage_vpce_dns_name                        = ""
  search_path                                   = "$current, $public"
  statement_queued_timeout_in_seconds           = 0
  statement_timeout_in_seconds                  = 172800
  strict_json_output                            = false
  timestamp_day_is_always_24h                   = false
  timestamp_input_format                        = "AUTO"
  timestamp_ltz_output_format                   = ""
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM"
  timestamp_type_mapping                        = "TIMESTAMP_NTZ"
  timestamp_tz_output_format                    = ""
  timezone                                      = "America/Los_Angeles"
  time_input_format                             = "AUTO"
  time_output_format                            = "HH24:MI:SS"
  trace_level                                   = "OFF"
  transaction_abort_on_error                    = false
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1970
  unsupported_ddl_action                        = "ignore"
  use_cached_result                             = true
  week_of_year_policy                           = 0
  week_start                                    = 0
}
