resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "AGENT"
    object_name = var.cortex_agent_fully_qualified_name
  }
}
