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
