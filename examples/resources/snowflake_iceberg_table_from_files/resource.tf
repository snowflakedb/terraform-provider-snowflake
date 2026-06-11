# Basic - only required fields
resource "snowflake_iceberg_table_from_files" "basic" {
  database           = "DATABASE"
  schema             = "SCHEMA"
  name               = "TABLE"
  metadata_file_path = "path/to/metadata/v1.metadata.json"
  external_volume    = "my_external_volume"
}

# Complete - all fields set
resource "snowflake_iceberg_table_from_files" "complete" {
  database                   = "DATABASE"
  schema                     = "SCHEMA"
  name                       = "TABLE"
  metadata_file_path         = "path/to/metadata/v1.metadata.json"
  external_volume            = "EXTERNAL_VOLUME"
  catalog                    = "CATALOG"
  comment                    = "COMMENT"
  replace_invalid_characters = true
}
