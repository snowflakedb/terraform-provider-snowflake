resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "some_database.role_name"
  privileges         = ["USAGE"]

  on_schema {
    schema_name             = "some_database.schema_name"
    all_schemas_in_database = "some_database"
  }
}
