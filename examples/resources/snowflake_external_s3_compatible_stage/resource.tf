# Basic resource with credentials
resource "snowflake_external_s3_compatible_stage" "basic" {
  name     = "my_s3_compatible_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"
}

# Complete resource with all options
resource "snowflake_external_s3_compatible_stage" "complete" {
  name     = "complete_s3_compatible_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

  credentials {
    aws_key_id     = var.aws_key_id
    aws_secret_key = var.aws_secret_key
  }

  directory {
    enable            = true
    refresh_on_create = true
    auto_refresh      = false
  }

  comment = "Fully configured S3-compatible external stage"
}
