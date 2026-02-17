# Basic hybrid table example with inline primary key
#
# This example demonstrates the minimal configuration needed to create
# a hybrid table with a simple structure.
#
# Note: Hybrid tables require preview feature enablement in the provider:
# provider "snowflake" {
#   preview_features_enabled = ["snowflake_hybrid_table_resource"]
# }

resource "snowflake_hybrid_table" "orders" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "ORDERS"
  comment  = "Orders table with basic structure"

  # Simple column with inline primary key constraint
  column {
    name        = "order_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
    comment     = "Unique order identifier"
  }

  column {
    name     = "customer_name"
    type     = "VARCHAR(100)"
    nullable = true
  }

  column {
    name    = "order_date"
    type    = "DATE"
    comment = "Date when order was placed"
  }

  column {
    name = "order_total"
    type = "NUMBER(10,2)"
  }

  # Optional: Set data retention period (0-90 days)
  data_retention_time_in_days = 7
}

# Output the fully qualified name
output "hybrid_table_fqn" {
  value = snowflake_hybrid_table.orders.fully_qualified_name
}
