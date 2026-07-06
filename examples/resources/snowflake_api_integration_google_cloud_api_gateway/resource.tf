# basic resource
resource "snowflake_api_integration_google_cloud_api_gateway" "basic" {
  name                 = "google_cloud_api_gateway_integration"
  google_audience      = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
  api_allowed_prefixes = ["https://gateway-id-123456.uc.gateway.dev/"]
  enabled              = true
}

# complete resource
resource "snowflake_api_integration_google_cloud_api_gateway" "complete" {
  name                 = "google_cloud_api_gateway_integration_complete"
  google_audience      = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
  api_allowed_prefixes = ["https://gateway-id-123456.uc.gateway.dev/"]
  api_blocked_prefixes = ["https://gateway-id-123456.uc.gateway.dev/blocked/"]
  enabled              = true
  comment              = "Example Google Cloud API Gateway integration"
}
