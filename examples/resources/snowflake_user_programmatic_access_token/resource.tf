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

# Set up dependencies and reference them from the token resource.
resource "snowflake_account_role" "role" {
  name = "ROLE"
}

resource "snowflake_user" "user" {
  name = "USER"
}

# Grant the role to the user. This is required to authenticate with PAT with role restriction.
resource "snowflake_grant_account_role" "grant_role_to_user" {
  role_name = snowflake_account_role.role.name
  user_name = snowflake_user.user.name
}

# complete resource with external references
resource "snowflake_user_programmatic_access_token" "complete_with_external_references" {
  user                                      = snowflake_user.user.name
  name                                      = "TOKEN"
  role_restriction                          = snowflake_account_role.role.name
  days_to_expiry                            = 30
  mins_to_bypass_network_policy_requirement = 10
  disabled                                  = false
  comment                                   = "COMMENT"

  # Use the keepers map to force token rotation. If any key or value in the map changes, the token will be rotated.
  keepers = {
    # here we use the time_rotating's rotation_rfc3339 field which provides a new timestamp every 30 days.
    rotation_time = time_rotating.my_token_rotation.rotation_rfc3339
  }
}

# note this requires the terraform to be run regularly
resource "time_rotating" "my_token_rotation" {
  rotation_days = 30
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
