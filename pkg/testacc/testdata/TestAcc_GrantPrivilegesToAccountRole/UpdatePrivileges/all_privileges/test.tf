resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = var.name
  all_privileges    = var.all_privileges
  on_account_object {
    object_type = "DATABASE"
    object_name = var.database
  }
}
