resource "snowflake_database_role" "test" {
  name     = var.database_role
  database = var.database
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  privileges         = ["CREATE TABLE"]

  on_schema {
    future_schemas_in_database = var.database
  }
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_database_role.test]

  future_grants_to {
    database_role = snowflake_database_role.test.fully_qualified_name
  }
}
