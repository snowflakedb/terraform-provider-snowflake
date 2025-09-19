resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = true
  storage_provider          = "S3"
  storage_allowed_locations = var.allowed_locations
  storage_aws_role_arn      = var.aws_role_arn
  storage_aws_external_id   = var.external_id
}
