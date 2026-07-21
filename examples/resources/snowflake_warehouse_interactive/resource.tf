# Basic interactive warehouse (only required fields)
resource "snowflake_warehouse_interactive" "basic" {
  name = "my_interactive_warehouse"
}

# Complete interactive warehouse (all fields set)
resource "snowflake_warehouse_interactive" "complete" {
  name                = "my_interactive_warehouse_complete"
  comment             = "My interactive warehouse with all options set"
  warehouse_size      = "XSMALL"
  max_cluster_count   = 2
  min_cluster_count   = 1
  auto_suspend        = 86400
  auto_resume         = "true"
  initially_suspended = true
  resource_monitor    = snowflake_resource_monitor.monitor.fully_qualified_name
  fallback_warehouse  = snowflake_warehouse.fallback.fully_qualified_name
  tables = [
    "\"MY_DB\".\"MY_SCHEMA\".\"MY_TABLE_1\"",
    "\"MY_DB\".\"MY_SCHEMA\".\"MY_TABLE_2\"",
  ]
  max_concurrency_level               = 8
  statement_queued_timeout_in_seconds = 30
  statement_timeout_in_seconds        = 5
}
