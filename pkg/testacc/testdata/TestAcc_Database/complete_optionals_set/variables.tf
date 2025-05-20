variable "name" {
  type = string
}

variable "transient" {
  type = bool
}

variable "comment" {
  type = string
}

variable "account_identifier" {
  type = string
}

variable "with_failover" {
  type = bool
}

variable "ignore_edition_check" {
  type = bool
}

variable "data_retention_time_in_days" {
  type = string
}

variable "max_data_extension_time_in_days" {
  type = string
}

variable "external_volume" {
  type = string
}

variable "catalog" {
  type = string
}

variable "replace_invalid_characters" {
  type = string
}

variable "default_ddl_collation" {
  type = string
}

variable "storage_serialization_policy" {
  type = string
}

variable "log_level" {
  type = string
}

variable "trace_level" {
  type = string
}

variable "suspend_task_after_num_failures" {
  type = number
}

variable "task_auto_retry_attempts" {
  type = number
}

variable "user_task_managed_initial_warehouse_size" {
  type = string
}

variable "user_task_timeout_ms" {
  type = number
}

variable "user_task_minimum_trigger_interval_in_seconds" {
  type = number
}

variable "quoted_identifiers_ignore_case" {
  type = bool
}

variable "enable_console_output" {
  type = bool
}
