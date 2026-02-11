# Basic resource with storage integration (required for GCS)
resource "snowflake_stage_external_gcs" "basic" {
  name                = "my_gcs_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name
}

# Complete resource with all options
resource "snowflake_stage_external_gcs" "complete" {
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
resource "snowflake_stage_external_gcs" "no_encryption" {
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
resource "snowflake_stage_external_gcs" "default_kms" {
  name                = "gcs_stage_default_kms"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  encryption {
    gcs_sse_kms {}
  }
}

# resource with inline CSV file format
resource "snowflake_stage_external_gcs" "with_csv_format" {
  name                = "gcs_csv_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

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
resource "snowflake_stage_external_gcs" "with_json_format" {
  name                = "gcs_json_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

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
resource "snowflake_stage_external_gcs" "with_avro_format" {
  name                = "gcs_avro_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

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
resource "snowflake_stage_external_gcs" "with_orc_format" {
  name                = "gcs_orc_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  file_format {
    orc {
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline Parquet file format
resource "snowflake_stage_external_gcs" "with_parquet_format" {
  name                = "gcs_parquet_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

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
resource "snowflake_stage_external_gcs" "with_xml_format" {
  name                = "gcs_xml_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

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
resource "snowflake_stage_external_gcs" "with_named_format" {
  name                = "gcs_named_format_stage"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "gcs://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.gcs.name

  file_format {
    format_name = snowflake_file_format.test.fully_qualified_name
  }
}
