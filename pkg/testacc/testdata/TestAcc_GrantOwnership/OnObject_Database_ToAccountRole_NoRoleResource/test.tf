resource "snowflake_grant_ownership" "test" {
  account_role_name = var.account_role_name
  on {
    object_type = "DATABASE"
    object_name = var.database_name
  }
}
