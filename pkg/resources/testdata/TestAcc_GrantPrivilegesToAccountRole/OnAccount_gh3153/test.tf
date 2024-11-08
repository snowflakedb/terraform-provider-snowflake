resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = var.name
  privileges        = var.privileges
  on_account        = true
}
