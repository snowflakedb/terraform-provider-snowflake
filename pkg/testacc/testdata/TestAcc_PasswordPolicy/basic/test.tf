resource "snowflake_password_policy" "pa" {
  name                 = var.name
  database             = var.database
  schema               = var.schema
  min_length           = var.min_length
  max_length           = var.max_length
  min_upper_case_chars = var.min_upper_case_chars
  min_lower_case_chars = var.min_lower_case_chars
  min_numeric_chars    = var.min_numeric_chars
  min_special_chars    = var.min_special_chars
  min_age_days         = var.min_age_days
  max_age_days         = var.max_age_days
  max_retries          = var.max_retries
  lockout_time_mins    = var.lockout_time_mins
  history              = var.history
  comment              = var.comment
  or_replace           = true
}