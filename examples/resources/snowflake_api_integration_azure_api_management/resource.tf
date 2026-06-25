# basic resource
resource "snowflake_api_integration_azure_api_management" "basic" {
  name                    = "azure_api_management_integration"
  azure_tenant_id         = "00000000-0000-0000-0000-000000000000"
  azure_ad_application_id = "11111111-1111-1111-1111-111111111111"
  api_allowed_prefixes    = ["https://apim-hello-world.azure-api.net/"]
  enabled                 = true
}

# complete resource
resource "snowflake_api_integration_azure_api_management" "complete" {
  name                    = "azure_api_management_integration_complete"
  azure_tenant_id         = "00000000-0000-0000-0000-000000000000"
  azure_ad_application_id = "11111111-1111-1111-1111-111111111111"
  api_allowed_prefixes    = ["https://apim-hello-world.azure-api.net/"]
  api_blocked_prefixes    = ["https://apim-hello-world.azure-api.net/blocked/"]
  api_key                 = "my-api-key"
  enabled                 = true
  comment                 = "Example Azure API Management integration"
}
