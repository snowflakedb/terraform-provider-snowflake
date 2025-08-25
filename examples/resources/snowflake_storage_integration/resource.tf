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

# This example demonstrates how to create a Snowflake storage integration for AWS S3 with private link enabled.
resource "snowflake_storage_integration" "integration" {
  name    = "s3_integration"
  comment = "AWS S3 integration with private link"
  type    = "EXTERNAL_STAGE"

  enabled = true

  storage_provider          = "S3"
  storage_aws_role_arn      = "arn:aws:iam::001234567890:role/myrole"
  storage_allowed_locations = ["s3://mybucket1/path1/", "s3://mybucket2/path2/"]
  storage_blocked_locations = ["s3://mybucket1/path1/blocked/", "s3://mybucket2/path2/blocked/"]
  use_private_link_endpoint = true
}

# This example demonstrates how to create a Snowflake storage integration for Azure with private link enabled.
resource "snowflake_storage_integration" "azure_integration" {
  name    = "azure_integration"
  comment = "Azure integration with private link"
  type    = "EXTERNAL_STAGE"

  enabled = true

  storage_provider          = "AZURE"
  azure_tenant_id           = "a123b4c5-1234-123a-a12b-1a23b45678c9"
  storage_allowed_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/"]
  use_private_link_endpoint = true
}
