resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges         = var.privileges
  with_grant_option  = var.with_grant_option

  on_schema_object {
    future {
      object_type_plural = var.object_type_plural
      in_database        = "\"${var.database}\""
    }
  }
}
