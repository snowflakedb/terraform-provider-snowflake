# Basic - only required fields
resource "snowflake_iceberg_table_from_rest" "basic" {
  database           = "DATABASE"
  schema             = "SCHEMA"
  name               = "TABLE"
  catalog_table_name = "my_catalog_table"

}

# Complete - all fields set
resource "snowflake_iceberg_table_from_rest" "complete" {
  database                       = "DATABASE"
  schema                         = "SCHEMA"
  name                           = "TABLE"
  external_volume                = "EXTERNAL_VOLUME"
  catalog                        = "CATALOG"
  catalog_table_name             = "my_catalog_table"
  catalog_namespace              = "my_namespace"
  path_layout                    = "HIERARCHICAL"
  target_file_size               = "128MB"
  replace_invalid_characters     = true
  auto_refresh                   = "true"
  comment                        = "COMMENT"
  storage_serialization_policy   = "OPTIMIZED"
  iceberg_merge_on_read_behavior = "ENABLED"
  enable_iceberg_merge_on_read   = true
}
