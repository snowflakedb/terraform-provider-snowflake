resource "snowflake_api_integration" "test_aws_int" {
  name                 = var.name
  api_provider         = var.api_provider
  api_aws_role_arn     = var.api_aws_role_arn
  api_allowed_prefixes = var.api_allowed_prefixes
  api_blocked_prefixes = var.api_blocked_prefixes
  api_key              = var.api_key
  comment              = var.comment
  enabled              = var.enabled
}
