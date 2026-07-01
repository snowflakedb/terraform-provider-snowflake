# Basic adaptive warehouse (only required fields)
resource "snowflake_warehouse_adaptive" "basic" {
  name = "my_adaptive_warehouse"
}

# Complete adaptive warehouse (all fields set)
resource "snowflake_warehouse_adaptive" "complete" {
  name                                = "my_adaptive_warehouse_complete"
  comment                             = "My adaptive warehouse with all options set"
  max_query_performance_level         = "MEDIUM"
  query_throughput_multiplier         = 1
  resource_monitor                    = "my_resource_monitor"
  statement_queued_timeout_in_seconds = 30
  statement_timeout_in_seconds        = 3600
}
