resource "snowflake_api_integration" "test_change" {
  name                    = var.name
  api_provider            = "azure_api_management"
  azure_tenant_id         = var.azure_tenant_id
  azure_ad_application_id = var.azure_ad_application_id
  api_allowed_prefixes    = var.api_allowed_prefixes
  api_blocked_prefixes    = var.api_blocked_prefixes
  api_key                 = var.api_key
  comment                 = var.comment
  enabled                 = var.enabled
}
