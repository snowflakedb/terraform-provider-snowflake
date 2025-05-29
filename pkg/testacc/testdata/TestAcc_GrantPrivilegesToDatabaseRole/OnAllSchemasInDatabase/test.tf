resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges         = var.privileges
  with_grant_option  = var.with_grant_option

  on_schema {
    all_schemas_in_database = "\"${var.database}\""
  }
}
