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

# use the token returned from Snowflake and remember to mark it as sensitive
output "token" {
  value     = snowflake_user_programmatic_access_token.complete.token
  sensitive = true
}

# rotate the token regularly using the keepers field and time_rotating resource
resource "snowflake_user_programmatic_access_token" "rotating" {
  user = "USER"
  name = "TOKEN"
  keepers = {
    rotation_schedule = time_rotating.rotation_schedule.rotation_rfc3339
  }
}

# Note that the fields of this resource are updated only when Terraform is run.
# This means that the schedule may not be respected if Terraform is not run regularly.
resource "time_rotating" "rotation_schedule" {
  rotation_days = 30
}
