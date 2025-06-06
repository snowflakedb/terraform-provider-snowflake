resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_resource_monitor" "test" {
  name = var.resource_monitor_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "RESOURCE MONITOR"
    object_name = snowflake_resource_monitor.test.name
  }
}
