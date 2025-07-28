resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = var.name
  privileges        = var.privileges
  with_grant_option = var.with_grant_option

  on_schema {
    schema_name = "\"${var.database}\".\"${var.schema}\""
  }
}
