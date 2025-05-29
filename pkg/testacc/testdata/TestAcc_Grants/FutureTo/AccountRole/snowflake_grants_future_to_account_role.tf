data "snowflake_current_role" "test" {}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = data.snowflake_current_role.test.name
  privileges        = ["CREATE TABLE"]

  on_schema {
    future_schemas_in_database = var.database
  }
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_account_role.test]

  future_grants_to {
    account_role = data.snowflake_current_role.test.name
  }
}
