# Basic example with OBJECT_STORE catalog source
resource "snowflake_catalog_integration" "object_store" {
  name           = "my_object_store_catalog_integration"
  catalog_source = "OBJECT_STORE"
  table_format   = "ICEBERG"
  enabled        = true
}

# Complete example with OBJECT_STORE and optional fields
resource "snowflake_catalog_integration" "object_store_complete" {
  name              = "my_object_store_catalog_integration"
  catalog_source    = "OBJECT_STORE"
  table_format      = "ICEBERG"
  enabled           = true
  comment           = "Catalog integration for Iceberg tables with object store"
  catalog_namespace = "my_namespace"
}

# Example with AWS Glue catalog source
resource "snowflake_catalog_integration" "glue" {
  name              = "my_glue_catalog_integration"
  catalog_source    = "GLUE"
  table_format      = "ICEBERG"
  enabled           = true
  glue_aws_role_arn = "arn:aws:iam::123456789012:role/SnowflakeGlueRole"
  glue_catalog_id   = "123456789012"
  glue_region       = "us-west-2"
  comment           = "Catalog integration for AWS Glue Data Catalog"
}

# Example with DELTA table format
resource "snowflake_catalog_integration" "delta" {
  name           = "my_delta_catalog_integration"
  catalog_source = "OBJECT_STORE"
  table_format   = "DELTA"
  enabled        = true
  comment        = "Catalog integration for Delta Lake tables"
}

# Example with disabled integration
resource "snowflake_catalog_integration" "disabled" {
  name           = "my_disabled_catalog_integration"
  catalog_source = "OBJECT_STORE"
  table_format   = "ICEBERG"
  enabled        = false
  comment        = "Disabled catalog integration"
}
