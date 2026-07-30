# Basic - only required fields
resource "snowflake_iceberg_table_from_aws_glue" "basic" {
  database           = "DATABASE"
  schema             = "SCHEMA"
  name               = "TABLE"
  catalog_table_name = "my_catalog_table"
}

# Complete - all fields set
resource "snowflake_iceberg_table_from_aws_glue" "complete" {
  database                   = "DATABASE"
  schema                     = "SCHEMA"
  name                       = "TABLE"
  external_volume            = "EXTERNAL_VOLUME"
  catalog                    = "CATALOG"
  catalog_table_name         = "my_catalog_table"
  catalog_namespace          = "my_namespace"
  replace_invalid_characters = true
  auto_refresh               = "true"
  comment                    = "COMMENT"
}
