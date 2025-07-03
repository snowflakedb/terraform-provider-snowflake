# basic resource
resource "snowflake_user_programmatic_access_token" "basic" {
  user = "USER"
  name = "TOKEN"
}

# complete resource
resource "snowflake_user_programmatic_access_token" "complete" {
  user                                      = "USER"
  name                                      = "TOKEN"
  role_restriction                          = "ROLE"
  days_to_expiry                            = 30
  mins_to_bypass_network_policy_requirement = 10
  disabled                                  = false
  comment                                   = "COMMENT"
}

# use the token returned from Snowflake
output "token" {
  value     = snowflake_user_programmatic_access_token.complete.token
  sensitive = true
}
