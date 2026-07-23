## Minimal
resource "snowflake_file_format_xml" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"
}

## Complete (with every optional set)
resource "snowflake_file_format_xml" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "file_format_name"

  compression                = "GZIP"
  preserve_space             = "true"
  strip_outer_element        = "true"
  disable_snowflake_data     = "false"
  disable_auto_convert       = "false"
  replace_invalid_characters = "false"
  ignore_utf8_errors         = "false"
  skip_byte_order_mark       = "true"
  comment                    = "My XML file format"
}
