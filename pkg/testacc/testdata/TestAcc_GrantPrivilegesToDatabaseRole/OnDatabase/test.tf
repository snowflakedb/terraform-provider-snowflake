resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges         = var.privileges
  on_database        = "\"${var.database}\""
  with_grant_option  = var.with_grant_option
}
