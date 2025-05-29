resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = var.account_role_name
  privileges        = var.privileges

  on_schema_object {
    object_type = var.object_type
    object_name = var.hybrid_table_fully_qualified_name
  }
}
