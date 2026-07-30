## Minimal
resource "snowflake_file_format_csv" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_csv" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  compression                    = "GZIP"
  record_delimiter               = ";"
  field_delimiter                = "|"
  multi_line                     = "true"
  file_extension                 = ".csv"
  skip_header                    = 1
  skip_blank_lines               = "true"
  date_format                    = "YYYY-MM-DD"
  time_format                    = "HH24:MI:SS"
  timestamp_format               = "YYYY-MM-DD HH24:MI:SS.FF3"
  binary_format                  = "BASE64"
  escape                         = "NONE"
  escape_unenclosed_field        = "NONE"
  trim_space                     = "true"
  field_optionally_enclosed_by   = "\""
  null_if                        = ["NULL", ""]
  error_on_column_count_mismatch = "false"
  replace_invalid_characters     = "false"
  empty_field_as_null            = "true"
  skip_byte_order_mark           = "true"
  encoding                       = "UTF8"
  comment                        = "My CSV file format"
}
