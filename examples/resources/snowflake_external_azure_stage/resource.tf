# Basic resource with storage integration
resource "snowflake_external_azure_stage" "basic" {
  name                = "my_azure_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name
}

# Complete resource with all options
resource "snowflake_external_azure_stage" "complete" {
  name                = "complete_azure_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer/path/"
  storage_integration = snowflake_storage_integration.azure.name

  encryption {
    azure_cse {
      master_key = var.azure_master_key
    }
  }

  directory {
    enable            = true
    refresh_on_create = true
    auto_refresh      = false
  }

  comment = "Fully configured Azure external stage"
}

# Resource with SAS token credentials instead of storage integration
resource "snowflake_external_azure_stage" "with_credentials" {
  name     = "azure_stage_with_sas"
  database = "my_database"
  schema   = "my_schema"
  url      = "azure://myaccount.blob.core.windows.net/mycontainer"

  credentials {
    azure_sas_token = var.azure_sas_token
  }
}

# Resource with encryption set to none
resource "snowflake_external_azure_stage" "no_encryption" {
  name                = "azure_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  encryption {
    none {}
  }
}
