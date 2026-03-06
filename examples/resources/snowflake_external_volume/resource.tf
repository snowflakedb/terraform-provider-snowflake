# Basic - S3 storage location with required fields only
resource "snowflake_external_volume" "s3_basic" {
  name = "my_external_volume"

  storage_location {
    storage_location_name = "my-s3-location"
    storage_provider      = "S3"
    storage_base_url      = "s3://mybucket/"
    storage_aws_role_arn  = "arn:aws:iam::123456789012:role/myrole"
  }
}

# Complete - S3 with all optional fields
resource "snowflake_external_volume" "s3_complete" {
  name         = "my_external_volume_complete"
  comment      = "my external volume"
  allow_writes = "true"

  storage_location {
    storage_location_name        = "my-s3-location"
    storage_provider             = "S3"
    storage_base_url             = "s3://mybucket/"
    storage_aws_role_arn         = "arn:aws:iam::123456789012:role/myrole"
    storage_aws_access_point_arn = "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"
    use_privatelink_endpoint     = "true"
    encryption_type              = "AWS_SSE_KMS"
    encryption_kms_key_id        = "1234abcd-12ab-34cd-56ef-1234567890ab"
  }
}

# GCS storage location
resource "snowflake_external_volume" "gcs" {
  name = "my_gcs_external_volume"

  storage_location {
    storage_location_name = "my-gcs-location"
    storage_provider      = "GCS"
    storage_base_url      = "gcs://mybucket/"
    encryption_type       = "GCS_SSE_KMS"
    encryption_kms_key_id = "1234abcd-12ab-34cd-56ef-1234567890ab"
  }
}

# Azure storage location
resource "snowflake_external_volume" "azure" {
  name = "my_azure_external_volume"

  storage_location {
    storage_location_name = "my-azure-location"
    storage_provider      = "AZURE"
    storage_base_url      = "azure://myaccount.blob.core.windows.net/mycontainer/"
    azure_tenant_id       = "123e4567-e89b-12d3-a456-426614174000"
  }
}

# S3-compatible storage location
resource "snowflake_external_volume" "s3compat" {
  name = "my_s3compat_external_volume"

  storage_location {
    storage_location_name  = "my-s3compat-location"
    storage_provider       = "S3COMPAT"
    storage_base_url       = "s3compat://mybucket/"
    storage_endpoint       = "https://s3-compatible.example.com"
    storage_aws_key_id     = var.aws_key_id
    storage_aws_secret_key = var.aws_secret_key
  }
}
