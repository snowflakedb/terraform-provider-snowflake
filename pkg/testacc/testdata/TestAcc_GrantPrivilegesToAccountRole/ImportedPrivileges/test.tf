resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "\"${var.role_name}\""
  privileges        = var.privileges
  on_account_object {
    object_type = "DATABASE"
    object_name = "\"${var.shared_database_name}\""
  }
}
