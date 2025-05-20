resource "snowflake_api_integration" "test_gcp_int" {
  name                 = var.name
  api_provider         = "google_api_gateway"
  google_audience      = var.google_audience
  api_allowed_prefixes = var.api_allowed_prefixes
  api_blocked_prefixes = var.api_blocked_prefixes
  comment              = var.comment
  enabled              = var.enabled
}
