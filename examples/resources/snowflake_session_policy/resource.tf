## Minimal
resource "snowflake_session_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "session_policy_name"
}

## Complete (with every optional set)
resource "snowflake_session_policy" "timeouts" {
  database                     = "database_name"
  schema                       = "schema_name"
  name                         = "session_policy_timeouts"
  session_idle_timeout_mins    = 60
  session_ui_idle_timeout_mins = 60
  allowed_secondary_roles      = ["ROLE_A", "ROLE_B"]
  blocked_secondary_roles      = ["ROLE_C", "ROLE_D"]
  comment                      = "My session policy"
}
