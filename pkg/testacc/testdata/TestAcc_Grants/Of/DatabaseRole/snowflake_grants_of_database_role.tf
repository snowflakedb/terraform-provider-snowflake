data "snowflake_current_role" "test" {}

resource "snowflake_database_role" "test" {
  name     = var.database_role
  database = var.database
}

resource "snowflake_grant_database_role" "test" {
  database_role_name = snowflake_database_role.test.fully_qualified_name
  parent_role_name   = data.snowflake_current_role.test.name
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_database_role.test]

  grants_of {
    database_role = snowflake_database_role.test.fully_qualified_name
  }
}
