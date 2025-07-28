resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "role_name"
  privileges        = ["OWNERSHIP"]
  with_grant_option = false

  on_schema_object {
    object_type = "TABLE"
    object_name = "some_database.schema_name.some_table"
  }
}
