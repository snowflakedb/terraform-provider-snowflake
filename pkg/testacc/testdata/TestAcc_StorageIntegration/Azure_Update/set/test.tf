resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = true
  storage_provider          = "AZURE"
  comment                   = var.comment
  storage_allowed_locations = var.allowed_locations
  storage_blocked_locations = var.blocked_locations
  azure_tenant_id           = var.azure_tenant_id
}
