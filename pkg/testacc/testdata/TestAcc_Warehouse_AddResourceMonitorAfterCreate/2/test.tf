resource "snowflake_resource_monitor" "monitor" {
  name         = var.resource_monitor_name
  credit_quota = 100
}

resource "snowflake_warehouse" "test" {
  name             = var.warehouse_name
  resource_monitor = snowflake_resource_monitor.monitor.fully_qualified_name
}
