# minimal
resource "snowflake_storage_integration_aws" "minimal" {
  name                      = "example_aws_storage_integration"
  enabled                   = true
  storage_provider          = "S3"
  storage_allowed_locations = ["s3://mybucket1/path1/"]
  storage_aws_role_arn      = "arn:aws:iam::001234567890:role/myrole"
}

# TODO [next PR]: add all fields example
# all fields
