resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = false
  storage_provider          = "S3"
  storage_allowed_locations = var.allowed_locations
  storage_aws_role_arn      = var.aws_role_arn
}
