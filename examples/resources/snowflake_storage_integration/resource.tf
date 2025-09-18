# AWS S3 integration with minimal fields set.
resource "snowflake_storage_integration" "integration" {
  name                      = "storage"
  type                      = "EXTERNAL_STAGE"
  enabled                   = true
  storage_allowed_locations = ["s3://mybucket1/path1/", "s3://mybucket2/path2/"]

  storage_provider     = "S3"
  storage_aws_role_arn = "arn:aws:iam::001234567890:role/myrole"
}

# AWS S3 integration with all fields set.
resource "snowflake_storage_integration" "integration" {
  name                      = "S3_INTEGRATION"
  comment                   = "AWS S3 integration with private link"
  type                      = "EXTERNAL_STAGE"
  enabled                   = true
  storage_allowed_locations = ["s3://mybucket1/path1/", "s3://mybucket2/path2/"]
  storage_blocked_locations = ["s3://mybucket1/path1/blocked/", "s3://mybucket2/path2/blocked/"]

  # AWS S3-specific fields.
  storage_provider          = "S3"
  storage_aws_role_arn      = "arn:aws:iam::001234567890:role/myrole"
  storage_aws_external_id   = "ABC12345_DEFRole=2_123ABC459AWQmtAdRqwe/A=="
  storage_aws_object_acl    = "bucket-owner-full-control"
  use_private_link_endpoint = "true"
}

# Azure integration with all fields set.
resource "snowflake_storage_integration" "azure_integration" {
  name                      = "AZURE_INTEGRATION"
  comment                   = "Azure integration with private link"
  type                      = "EXTERNAL_STAGE"
  enabled                   = true
  storage_allowed_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/"]
  storage_blocked_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/blocked/"]

  # Azure-specific fields.
  storage_provider          = "AZURE"
  azure_tenant_id           = "a123b4c5-1234-123a-a12b-1a23b45678c9"
  use_private_link_endpoint = "true"
}
