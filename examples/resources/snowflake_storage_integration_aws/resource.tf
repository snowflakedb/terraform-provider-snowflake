# minimal
resource "snowflake_storage_integration_aws" "minimal" {
  name                      = "example_aws_storage_integration"
  enabled                   = true
  storage_provider          = "S3"
  storage_allowed_locations = ["s3://mybucket1/path1/"]
  storage_aws_role_arn      = "arn:aws:iam::001234567890:role/myrole"
}

# all fields
resource "snowflake_storage_integration_aws" "all" {
  name                      = "example_aws_storage_integration"
  enabled                   = true
  storage_provider          = "S3"
  storage_allowed_locations = ["s3://mybucket1/allowed-location/", "s3://mybucket1/allowed-location2/"]
  storage_blocked_locations = ["s3://mybucket1/blocked-location/", "s3://mybucket1/blocked-location2/"]
  use_privatelink_endpoint  = "true"
  comment                   = "some comment"

  storage_aws_role_arn    = "arn:aws:iam::001234567890:role/myrole"
  storage_aws_external_id = "some_external_id"
  storage_aws_object_acl  = "bucket-owner-full-control"
}
