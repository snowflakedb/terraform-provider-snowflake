# Basic resource with storage integration
resource "snowflake_external_s3_stage" "basic" {
  name     = "my_s3_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"
}

# Complete resource with all options
resource "snowflake_external_s3_stage" "complete" {
  name                 = "complete_s3_stage"
  database             = "my_database"
  schema               = "my_schema"
  url                  = "s3://mybucket/mypath/"
  storage_integration  = snowflake_storage_integration.s3.name
  aws_access_point_arn = "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"

  encryption {
    aws_cse {
      master_key = var.s3_master_key
    }
  }

  directory {
    enable            = true
    refresh_on_create = true
    auto_refresh      = false
  }

  comment = "Fully configured S3 external stage"
}

# Resource with AWS key credentials instead of storage integration
resource "snowflake_external_s3_stage" "with_key_credentials" {
  name     = "s3_stage_with_keys"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_key_id     = var.aws_access_key_id
    aws_secret_key = var.aws_secret_access_key
    aws_token      = var.aws_token
  }
}

# Resource with AWS IAM role credentials
resource "snowflake_external_s3_stage" "with_role_credentials" {
  name     = "s3_stage_with_role"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_role = var.aws_role_arn
  }
}

# Resource with SSE-S3 encryption
resource "snowflake_external_s3_stage" "sse_s3" {
  name                = "s3_stage_sse_s3"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    aws_sse_s3 {}
  }
}

# Resource with SSE-KMS encryption
resource "snowflake_external_s3_stage" "sse_kms" {
  name                = "s3_stage_sse_kms"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    aws_sse_kms {
      kms_key_id = var.kms_key_id
    }
  }
}

# Resource with encryption set to none
resource "snowflake_external_s3_stage" "no_encryption" {
  name                = "s3_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    none {}
  }
}
