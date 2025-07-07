resource "snowflake_user_programmatic_access_token" "test" {
  name                             = var.name
  user                             = var.user
  keepers                          = try(var.keepers, null)
  expire_rotated_token_after_hours = try(var.expire_rotated_token_after_hours, null)
}
