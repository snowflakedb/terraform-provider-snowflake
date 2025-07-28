resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = true
  storage_provider          = "GCS"
  comment                   = var.comment
  storage_allowed_locations = var.allowed_locations
  storage_blocked_locations = var.blocked_locations
}
