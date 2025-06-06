resource "snowflake_schema" "test" {
  database = var.database
  name     = var.schema_name
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  depends_on        = [snowflake_schema.test]
  account_role_name = var.name
  privileges        = var.privileges
  on_schema {
    schema_name = "${var.database}.${var.schema_name}"
  }
}
