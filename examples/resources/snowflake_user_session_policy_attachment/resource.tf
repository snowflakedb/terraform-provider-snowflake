resource "snowflake_session_policy" "sp" {
  database = "prod"
  schema   = "security"
  name     = "default_session_policy"
}

resource "snowflake_user" "user" {
  name = "USER_NAME"
}

resource "snowflake_user_session_policy_attachment" "spa" {
  session_policy_name = snowflake_session_policy.sp.fully_qualified_name
  user_name           = snowflake_user.user.name
}
