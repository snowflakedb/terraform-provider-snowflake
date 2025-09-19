resource "snowflake_storage_integration" "test" {
  name                      = var.name
  enabled                   = true
  storage_provider          = "S3"
  comment                   = var.comment
  storage_allowed_locations = var.allowed_locations
  storage_blocked_locations = var.blocked_locations
  storage_aws_role_arn      = var.aws_role_arn
  storage_aws_object_acl    = var.aws_object_acl
  storage_aws_external_id   = var.external_id
  use_privatelink_endpoint  = var.use_privatelink_endpoint
}
