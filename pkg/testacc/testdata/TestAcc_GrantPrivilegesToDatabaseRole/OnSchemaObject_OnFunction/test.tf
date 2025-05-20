resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = var.name
  privileges         = var.privileges
  with_grant_option  = var.with_grant_option

  on_schema_object {
    object_type = "FUNCTION"
    object_name = "\"${var.database}\".\"${var.schema}\".\"${var.function_name}\"(${var.argument_type})"
  }
}
