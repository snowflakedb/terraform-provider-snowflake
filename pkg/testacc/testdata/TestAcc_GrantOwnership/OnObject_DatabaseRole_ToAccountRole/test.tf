resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_database_role" "test" {
  name     = var.database_role_name
  database = var.database_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE ROLE"
    object_name = snowflake_database_role.test.fully_qualified_name
  }
}
