resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = var.name
  all_privileges    = true
  on_account        = true
  always_apply      = var.always_apply
}
