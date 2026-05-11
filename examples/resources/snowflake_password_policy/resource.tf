## Minimal
resource "snowflake_password_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "password_policy_name"
}

## Complete (with every optional set)
resource "snowflake_password_policy" "complete" {
  database             = "database_name"
  schema               = "schema_name"
  name                 = "password_policy_name"
  min_length           = 10
  max_length           = 30
  min_upper_case_chars = 2
  min_lower_case_chars = 3
  min_numeric_chars    = 4
  min_special_chars    = 5
  min_age_days         = 1
  max_age_days         = 30
  max_retries          = 3
  lockout_time_mins    = 30
  history              = 5
  comment              = "My password policy"
}
