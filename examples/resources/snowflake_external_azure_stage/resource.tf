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

# resource with inline CSV file format
resource "snowflake_external_azure_stage" "with_csv_format" {
  name                = "azure_csv_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    csv {
      compression                    = "GZIP"
      record_delimiter               = "\n"
      field_delimiter                = "|"
      multi_line                     = "false"
      file_extension                 = ".csv"
      skip_header                    = 1 # or parse_header = true
      skip_blank_lines               = "true"
      date_format                    = "AUTO"
      time_format                    = "AUTO"
      timestamp_format               = "AUTO"
      binary_format                  = "HEX"
      escape                         = "\\"
      escape_unenclosed_field        = "\\"
      trim_space                     = "false"
      field_optionally_enclosed_by   = "\""
      null_if                        = ["NULL", ""]
      error_on_column_count_mismatch = "true"
      replace_invalid_characters     = "false"
      empty_field_as_null            = "true"
      skip_byte_order_mark           = "true"
      encoding                       = "UTF8"
    }
  }
}

# resource with inline JSON file format
resource "snowflake_external_azure_stage" "with_json_format" {
  name                = "azure_json_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    json {
      compression                = "AUTO"
      date_format                = "AUTO"
      time_format                = "AUTO"
      timestamp_format           = "AUTO"
      binary_format              = "HEX"
      trim_space                 = "false"
      multi_line                 = "false"
      null_if                    = ["NULL", ""]
      file_extension             = ".json"
      enable_octal               = "false"
      allow_duplicate            = "false"
      strip_outer_array          = "false"
      strip_null_values          = "false"
      replace_invalid_characters = "false" # or ignore_utf8_errors = true
      skip_byte_order_mark       = "false"
    }
  }
}

# resource with inline AVRO file format
resource "snowflake_external_azure_stage" "with_avro_format" {
  name                = "azure_avro_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    avro {
      compression                = "GZIP"
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline ORC file format
resource "snowflake_external_azure_stage" "with_orc_format" {
  name                = "azure_orc_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    orc {
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline Parquet file format
resource "snowflake_external_azure_stage" "with_parquet_format" {
  name                = "azure_parquet_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    parquet {
      compression                = "SNAPPY"
      binary_as_text             = "true"
      use_logical_type           = "true"
      trim_space                 = "false"
      use_vectorized_scanner     = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline XML file format
resource "snowflake_external_azure_stage" "with_xml_format" {
  name                = "azure_xml_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    xml {
      compression                = "AUTO"
      preserve_space             = "false"
      strip_outer_element        = "false"
      disable_auto_convert       = "false"
      replace_invalid_characters = "false" # or ignore_utf8_errors = true
      skip_byte_order_mark       = "false"
    }
  }
}

# resource with named file format
resource "snowflake_external_azure_stage" "with_named_format" {
  name                = "azure_named_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "azure://myaccount.blob.core.windows.net/mycontainer"
  storage_integration = snowflake_storage_integration.azure.name

  file_format {
    format_name = snowflake_file_format.test.fully_qualified_name
  }
}
