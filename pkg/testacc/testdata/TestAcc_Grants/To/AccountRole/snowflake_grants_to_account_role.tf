data "snowflake_current_role" "test" {}

data "snowflake_grants" "test" {
  grants_to {
    account_role = data.snowflake_current_role.test.name
  }
}
