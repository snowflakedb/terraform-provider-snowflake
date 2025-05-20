resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = var.database_role_name
  privileges         = var.privileges

  on_schema_object {
    object_type = var.object_type
    object_name = var.hybrid_table_fully_qualified_name
  }
}
