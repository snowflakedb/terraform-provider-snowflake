resource "snowflake_storage_integration" "test_azure" {
  name                      = var.name
  storage_allowed_locations = var.allowed_locations
  storage_provider          = "AZURE"
  azure_tenant_id           = var.azure_tenant_id
  use_private_link_endpoint = false
}
