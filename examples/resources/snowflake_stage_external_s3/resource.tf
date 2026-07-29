# Basic resource with storage integration
resource "snowflake_stage_external_s3" "basic" {
  name     = "my_s3_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"
}

# Complete resource with all options
resource "snowflake_stage_external_s3" "complete" {
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
resource "snowflake_stage_external_s3" "with_key_credentials" {
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
resource "snowflake_stage_external_s3" "with_role_credentials" {
  name     = "s3_stage_with_role"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_role = var.aws_role_arn
  }
}

# Resource with SSE-S3 encryption
resource "snowflake_stage_external_s3" "sse_s3" {
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
resource "snowflake_stage_external_s3" "sse_kms" {
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
resource "snowflake_stage_external_s3" "no_encryption" {
  name                = "s3_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    none {}
  }
}

# resource with inline CSV file format
resource "snowflake_stage_external_s3" "with_csv_format" {
  name     = "s3_csv_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

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
resource "snowflake_stage_external_s3" "with_json_format" {
  name     = "s3_json_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

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
resource "snowflake_stage_external_s3" "with_avro_format" {
  name     = "s3_avro_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

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
resource "snowflake_stage_external_s3" "with_orc_format" {
  name     = "s3_orc_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    orc {
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline Parquet file format
resource "snowflake_stage_external_s3" "with_parquet_format" {
  name     = "s3_parquet_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

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
resource "snowflake_stage_external_s3" "with_xml_format" {
  name     = "s3_xml_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

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
resource "snowflake_stage_external_s3" "with_named_format" {
  name     = "s3_named_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    format_name = snowflake_file_format.test.fully_qualified_name
  }
}
