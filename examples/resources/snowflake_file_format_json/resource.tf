## Minimal
resource "snowflake_file_format_json" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_json" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  compression                = "GZIP"
  date_format                = "YYYY-MM-DD"
  time_format                = "HH24:MI:SS"
  timestamp_format           = "YYYY-MM-DD HH24:MI:SS.FF3"
  binary_format              = "BASE64"
  trim_space                 = "true"
  multi_line                 = "true"
  null_if                    = ["NULL", ""]
  file_extension             = ".json"
  enable_octal               = "false"
  allow_duplicate            = "false"
  strip_outer_array          = "true"
  strip_null_values          = "false"
  replace_invalid_characters = "false"
  ignore_utf8_errors         = "false"
  skip_byte_order_mark       = "true"
  comment                    = "My JSON file format"
}
