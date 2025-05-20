resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_database        = var.database_name
    }
  }
}
