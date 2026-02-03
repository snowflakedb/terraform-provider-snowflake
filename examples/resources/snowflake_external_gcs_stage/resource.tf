# Basic resource with storage integration (required for GCS)
resource "snowflake_external_gcs_stage" "basic" {
  name                = "my_gcs_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name
}

# Complete resource with all options
resource "snowflake_external_gcs_stage" "complete" {
  name                = "complete_gcs_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  encryption {
    gcs_sse_kms {
      kms_key_id = var.gcs_kms_key_id
    }
  }

  directory {
    enable            = true
    refresh_on_create = true
    auto_refresh      = false
  }

  comment = "Fully configured GCS external stage"
}

# Resource with encryption set to none
resource "snowflake_external_gcs_stage" "no_encryption" {
  name                = "gcs_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  encryption {
    none {}
  }
}

# Resource with GCS SSE KMS encryption without specifying key
resource "snowflake_external_gcs_stage" "default_kms" {
  name                = "gcs_stage_default_kms"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  encryption {
    gcs_sse_kms {}
  }
}
