resource "snowflake_grant_privileges_to_account_role" "test" {
  privileges        = ["CREATE SCHEMA"]
  account_role_name = "\"${var.name}\""
  on_account_object {
    object_type = "DATABASE"
    object_name = var.database
  }
}