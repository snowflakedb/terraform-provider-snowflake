# basic resource
resource "snowflake_catalog_integration_object_storage" "basic" {
  name         = "example"
  enabled      = false
  table_format = "DELTA"
}

# complete resource
resource "snowflake_catalog_integration_object_storage" "complete" {
  name                     = "example_complete"
  enabled                  = true
  refresh_interval_seconds = 60
  comment                  = "Lorem ipsum"
  table_format             = "ICEBERG"
}
