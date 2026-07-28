## Minimal
resource "snowflake_file_format_avro" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_avro" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  compression                = "GZIP"
  trim_space                 = "true"
  replace_invalid_characters = "false"
  null_if                    = ["NULL", ""]
  comment                    = "My AVRO file format"
}
