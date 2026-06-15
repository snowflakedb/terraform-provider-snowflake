resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    all {
      object_type_plural = "AGENTS"
      in_schema          = var.schema_fully_qualified_name
    }
  }
}
