resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "role_name"
  privileges        = ["USAGE"]

  on_schema {
    schema_name             = "some_database.schema_name"
    all_schemas_in_database = "some_database"
  }
}
