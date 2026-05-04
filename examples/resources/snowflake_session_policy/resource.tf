## Minimal
resource "snowflake_session_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "session_policy_name"
}

## Complete (with every optional set)
resource "snowflake_session_policy" "complete" {
  database                     = "database_name"
  schema                       = "schema_name"
  name                         = "session_policy_name"
  session_idle_timeout_mins    = 60
  session_ui_idle_timeout_mins = 60
  allowed_secondary_roles {
    roles = ["ROLE_A", "ROLE_B"]
  }
  blocked_secondary_roles {
    roles = ["ROLE_C", "ROLE_D"]
  }
  comment = "My session policy"
}

## Secondary roles using all / none
resource "snowflake_session_policy" "secondary_roles_all_none" {
  database = "database_name"
  schema   = "schema_name"
  name     = "session_policy_name"

  allowed_secondary_roles {
    all = true
  }
  blocked_secondary_roles {
    none = true
  }
}
