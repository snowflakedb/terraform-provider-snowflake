resource "snowflake_storage_integration" "integration" {
  name    = "storage"
  comment = "A storage integration."
  type    = "EXTERNAL_STAGE"

  enabled = true

  #   storage_allowed_locations = [""]
  #   storage_blocked_locations = [""]
  #   storage_aws_object_acl    = "bucket-owner-full-control"

  storage_provider         = "S3"
  storage_aws_external_id  = "ABC12345_DEFRole=2_123ABC459AWQmtAdRqwe/A=="
  storage_aws_iam_user_arn = "..."
  storage_aws_role_arn     = "..."

  # azure_tenant_id
}
