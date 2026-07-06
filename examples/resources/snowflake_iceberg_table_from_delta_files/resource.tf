# Basic - only required fields
resource "snowflake_iceberg_table_from_delta_files" "basic" {
  database      = "DATABASE"
  schema        = "SCHEMA"
  name          = "TABLE"
  base_location = "path/to/delta/table"
}

# Complete - all fields set
resource "snowflake_iceberg_table_from_delta_files" "complete" {
  database                   = "DATABASE"
  schema                     = "SCHEMA"
  name                       = "TABLE"
  base_location              = "path/to/delta/table"
  external_volume            = "EXTERNAL_VOLUME"
  catalog                    = "CATALOG"
  auto_refresh               = "true"
  comment                    = "COMMENT"
  replace_invalid_characters = true
}
