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

# Resource with inline CSV file format
resource "snowflake_external_s3_compatible_stage" "with_csv_format" {
  name     = "s3_compat_csv_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

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

# Resource with inline JSON file format
resource "snowflake_external_s3_compatible_stage" "with_json_format" {
  name     = "s3_compat_json_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

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

# Resource with inline AVRO file format
resource "snowflake_external_s3_compatible_stage" "with_avro_format" {
  name     = "s3_compat_avro_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

  file_format {
    avro {
      compression                = "GZIP"
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# Resource with inline ORC file format
resource "snowflake_external_s3_compatible_stage" "with_orc_format" {
  name     = "s3_compat_orc_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

  file_format {
    orc {
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# Resource with inline Parquet file format
resource "snowflake_external_s3_compatible_stage" "with_parquet_format" {
  name     = "s3_compat_parquet_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

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

# Resource with inline XML file format
resource "snowflake_external_s3_compatible_stage" "with_xml_format" {
  name     = "s3_compat_xml_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

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

# Resource with named file format
resource "snowflake_external_s3_compatible_stage" "with_named_format" {
  name     = "s3_compat_named_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3compat://bucket/path/"
  endpoint = "s3.my-provider.com"

  file_format {
    format_name = snowflake_file_format.test.fully_qualified_name
  }
}
