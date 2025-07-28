resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = false
  storage_provider          = "GCS"
  storage_allowed_locations = var.allowed_locations
}
