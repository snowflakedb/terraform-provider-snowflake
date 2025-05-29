resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"some_database\".\"some_name\""
  privileges         = ["OWNERSHIP"]
  with_grant_option  = false

  on_schema_object {
    object_type = "TABLE"
    object_name = "some_database.some_schema.some_table"
  }
}
