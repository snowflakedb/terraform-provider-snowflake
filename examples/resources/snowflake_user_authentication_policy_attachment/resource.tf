resource "snowflake_user" "user" {
  name = "USER_NAME"
}
resource "snowflake_authentication_policy" "ap" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}
resource "snowflake_user_authentication_policy_attachment" "apa" {
  authentication_policy_name = snowflake_authentication_policy.ap.fully_qualified_name
  user_name                  = snowflake_user.user.name
}
