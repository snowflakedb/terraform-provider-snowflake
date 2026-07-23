## Minimal
resource "snowflake_file_format_parquet" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_parquet" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  compression                = "LZO"
  binary_as_text             = "false"
  use_logical_type           = "true"
  trim_space                 = "true"
  use_vectorized_scanner     = "true"
  replace_invalid_characters = "false"
  null_if                    = ["NULL", ""]
  comment                    = "My Parquet file format"
}

## Snappy compression (mutually exclusive with `compression`)
resource "snowflake_file_format_parquet" "snappy" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  snappy_compression = "true"
}
