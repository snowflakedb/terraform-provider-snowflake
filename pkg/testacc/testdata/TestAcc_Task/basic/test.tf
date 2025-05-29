resource "snowflake_task" "test" {
  name          = var.name
  database      = var.database
  schema        = var.schema
  started       = var.started
  sql_statement = var.sql_statement

  # Optionals
  warehouse                   = var.warehouse
  config                      = var.config
  allow_overlapping_execution = var.allow_overlapping_execution
  error_integration           = var.error_integration
  when                        = var.when
  comment                     = var.comment

  dynamic "schedule" {
    for_each = [for element in [var.schedule] : element if element != null]
    content {
      minutes    = lookup(var.schedule, "minutes", null)
      using_cron = lookup(var.schedule, "cron", null)
    }
  }

  # Parameters
  suspend_task_after_num_failures               = var.suspend_task_after_num_failures
  task_auto_retry_attempts                      = var.task_auto_retry_attempts
  user_task_managed_initial_warehouse_size      = var.user_task_managed_initial_warehouse_size
  user_task_minimum_trigger_interval_in_seconds = var.user_task_minimum_trigger_interval_in_seconds
  user_task_timeout_ms                          = var.user_task_timeout_ms
  abort_detached_query                          = var.abort_detached_query
  autocommit                                    = var.autocommit
  binary_input_format                           = var.binary_input_format
  binary_output_format                          = var.binary_output_format
  client_memory_limit                           = var.client_memory_limit
  client_metadata_request_use_connection_ctx    = var.client_metadata_request_use_connection_ctx
  client_prefetch_threads                       = var.client_prefetch_threads
  client_result_chunk_size                      = var.client_result_chunk_size
  client_result_column_case_insensitive         = var.client_result_column_case_insensitive
  client_session_keep_alive                     = var.client_session_keep_alive
  client_session_keep_alive_heartbeat_frequency = var.client_session_keep_alive_heartbeat_frequency
  client_timestamp_type_mapping                 = var.client_timestamp_type_mapping
  date_input_format                             = var.date_input_format
  date_output_format                            = var.date_output_format
  enable_unload_physical_type_optimization      = var.enable_unload_physical_type_optimization
  error_on_nondeterministic_merge               = var.error_on_nondeterministic_merge
  error_on_nondeterministic_update              = var.error_on_nondeterministic_update
  geography_output_format                       = var.geography_output_format
  geometry_output_format                        = var.geometry_output_format
  jdbc_use_session_timezone                     = var.jdbc_use_session_timezone
  json_indent                                   = var.json_indent
  lock_timeout                                  = var.lock_timeout
  log_level                                     = var.log_level
  multi_statement_count                         = var.multi_statement_count
  noorder_sequence_as_default                   = var.noorder_sequence_as_default
  odbc_treat_decimal_as_int                     = var.odbc_treat_decimal_as_int
  query_tag                                     = var.query_tag
  quoted_identifiers_ignore_case                = var.quoted_identifiers_ignore_case
  rows_per_resultset                            = var.rows_per_resultset
  s3_stage_vpce_dns_name                        = var.s3_stage_vpce_dns_name
  search_path                                   = var.search_path
  statement_queued_timeout_in_seconds           = var.statement_queued_timeout_in_seconds
  statement_timeout_in_seconds                  = var.statement_timeout_in_seconds
  strict_json_output                            = var.strict_json_output
  timestamp_day_is_always_24h                   = var.timestamp_day_is_always_24h
  timestamp_input_format                        = var.timestamp_input_format
  timestamp_ltz_output_format                   = var.timestamp_ltz_output_format
  timestamp_ntz_output_format                   = var.timestamp_ntz_output_format
  timestamp_output_format                       = var.timestamp_output_format
  timestamp_type_mapping                        = var.timestamp_type_mapping
  timestamp_tz_output_format                    = var.timestamp_tz_output_format
  timezone                                      = var.timezone
  time_input_format                             = var.time_input_format
  time_output_format                            = var.time_output_format
  trace_level                                   = var.trace_level
  transaction_abort_on_error                    = var.transaction_abort_on_error
  transaction_default_isolation_level           = var.transaction_default_isolation_level
  two_digit_century_start                       = var.two_digit_century_start
  unsupported_ddl_action                        = var.unsupported_ddl_action
  use_cached_result                             = var.use_cached_result
  week_of_year_policy                           = var.week_of_year_policy
  week_start                                    = var.week_start
}
