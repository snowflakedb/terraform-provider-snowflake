# Resource with required fields
resource "snowflake_warehouse" "warehouse" {
  name = "WAREHOUSE"
}

# Resource with all fields
resource "snowflake_warehouse" "warehouse" {
  name                                = "WAREHOUSE"
  warehouse_type                      = "SNOWPARK-OPTIMIZED"
  warehouse_size                      = "MEDIUM"
  max_cluster_count                   = 4
  min_cluster_count                   = 2
  scaling_policy                      = "ECONOMY"
  auto_suspend                        = 1200
  auto_resume                         = false
  initially_suspended                 = false
  resource_monitor                    = snowflake_resource_monitor.monitor.fully_qualified_name
  comment                             = "An example warehouse."
  enable_query_acceleration           = true
  query_acceleration_max_scale_factor = 4
  resource_constraint                 = "MEMORY_16X"

  max_concurrency_level               = 4
  statement_queued_timeout_in_seconds = 5
  statement_timeout_in_seconds        = 86400
}
