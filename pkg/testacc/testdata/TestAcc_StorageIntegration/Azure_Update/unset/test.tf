resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = false
  storage_provider          = "AZURE"
  storage_allowed_locations = var.allowed_locations
  azure_tenant_id           = var.azure_tenant_id
}
