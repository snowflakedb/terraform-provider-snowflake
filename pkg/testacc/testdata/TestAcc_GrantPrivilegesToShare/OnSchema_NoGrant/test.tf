resource "snowflake_schema" "test" {
  name     = var.schema
  database = var.database
}

