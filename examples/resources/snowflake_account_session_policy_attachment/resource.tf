resource "snowflake_session_policy" "sp" {
  database = "prod"
  schema   = "security"
  name     = "default_session_policy"
}

resource "snowflake_account_session_policy_attachment" "attachment" {
  session_policy_name = snowflake_session_policy.sp.fully_qualified_name
}
