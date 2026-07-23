## Minimal
resource "snowflake_file_format_orc" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_orc" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  trim_space                 = "true"
  null_if                    = ["NULL", ""]
  replace_invalid_characters = "false"
  comment                    = "My ORC file format"
}
