resource "snowflake_authentication_policy" "default" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}
resource "snowflake_account_authentication_policy_attachment" "attachment" {
  authentication_policy = snowflake_authentication_policy.default.fully_qualified_name
}
