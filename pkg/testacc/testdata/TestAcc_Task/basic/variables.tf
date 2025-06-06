variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "name" {
  type = string
}

variable "started" {
  type = bool
}

variable "sql_statement" {
  type = string
}

# Optionals
variable "comment" {
  type    = string
  default = null
}

variable "warehouse" {
  type    = string
  default = null
}

variable "config" {
  type    = string
  default = null
}

variable "allow_overlapping_execution" {
  type    = string
  default = null
}

variable "error_integration" {
  type    = string
  default = null
}

variable "when" {
  type    = string
  default = null
}

variable "schedule" {
  default = null
  type    = map(string)
}

# Parameters
variable "suspend_task_after_num_failures" {
  default = null
  type    = number
}

variable "task_auto_retry_attempts" {
  default = null
  type    = number
}

variable "user_task_managed_initial_warehouse_size" {
  default = null
  type    = string
}

variable "user_task_minimum_trigger_interval_in_seconds" {
  default = null
  type    = number
}

variable "user_task_timeout_ms" {
  default = null
  type    = number
}

variable "abort_detached_query" {
  default = null
  type    = bool
}

variable "autocommit" {
  default = null
  type    = bool
}

variable "binary_input_format" {
  default = null
  type    = string
}

variable "binary_output_format" {
  default = null
  type    = string
}

variable "client_memory_limit" {
  default = null
  type    = number
}

variable "client_metadata_request_use_connection_ctx" {
  default = null
  type    = bool
}

variable "client_prefetch_threads" {
  default = null
  type    = number
}

variable "client_result_chunk_size" {
  default = null
  type    = number
}

variable "client_result_column_case_insensitive" {
  default = null
  type    = bool
}

variable "client_session_keep_alive" {
  default = null
  type    = bool
}

variable "client_session_keep_alive_heartbeat_frequency" {
  default = null
  type    = number
}

variable "client_timestamp_type_mapping" {
  default = null
  type    = string
}

variable "date_input_format" {
  default = null
  type    = string
}

variable "date_output_format" {
  default = null
  type    = string
}

variable "enable_unload_physical_type_optimization" {
  default = null
  type    = bool
}

variable "error_on_nondeterministic_merge" {
  default = null
  type    = bool
}

variable "error_on_nondeterministic_update" {
  default = null
  type    = bool
}

variable "geography_output_format" {
  default = null
  type    = string
}

variable "geometry_output_format" {
  default = null
  type    = string
}

variable "jdbc_use_session_timezone" {
  default = null
  type    = bool
}

variable "json_indent" {
  default = null
  type    = number
}

variable "lock_timeout" {
  default = null
  type    = number
}

variable "log_level" {
  default = null
  type    = string
}

variable "multi_statement_count" {
  default = null
  type    = number
}

variable "noorder_sequence_as_default" {
  default = null
  type    = bool
}

variable "odbc_treat_decimal_as_int" {
  default = null
  type    = bool
}

variable "query_tag" {
  default = null
  type    = string
}

variable "quoted_identifiers_ignore_case" {
  default = null
  type    = bool
}

variable "rows_per_resultset" {
  default = null
  type    = number
}

variable "s3_stage_vpce_dns_name" {
  default = null
  type    = string
}

variable "search_path" {
  default = null
  type    = string
}

variable "statement_queued_timeout_in_seconds" {
  default = null
  type    = number
}

variable "statement_timeout_in_seconds" {
  default = null
  type    = number
}

variable "strict_json_output" {
  default = null
  type    = bool
}

variable "timestamp_day_is_always_24h" {
  default = null
  type    = bool
}

variable "timestamp_input_format" {
  default = null
  type    = string
}

variable "timestamp_ltz_output_format" {
  default = null
  type    = string
}

variable "timestamp_ntz_output_format" {
  default = null
  type    = string
}

variable "timestamp_output_format" {
  default = null
  type    = string
}

variable "timestamp_type_mapping" {
  default = null
  type    = string
}

variable "timestamp_tz_output_format" {
  default = null
  type    = string
}

variable "timezone" {
  default = null
  type    = string
}

variable "time_input_format" {
  default = null
  type    = string
}

variable "time_output_format" {
  default = null
  type    = string
}

variable "trace_level" {
  default = null
  type    = string
}

variable "transaction_abort_on_error" {
  default = null
  type    = bool
}

variable "transaction_default_isolation_level" {
  default = null
  type    = string
}

variable "two_digit_century_start" {
  default = null
  type    = number
}

variable "unsupported_ddl_action" {
  default = null
  type    = string
}

variable "use_cached_result" {
  default = null
  type    = bool
}

variable "week_of_year_policy" {
  default = null
  type    = number
}

variable "week_start" {
  default = null
  type    = number
}
