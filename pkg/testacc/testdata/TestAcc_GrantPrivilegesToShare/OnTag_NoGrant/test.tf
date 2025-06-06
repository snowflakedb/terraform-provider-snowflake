resource "snowflake_tag" "test" {
  name     = var.on_tag
  database = var.database
  schema   = var.schema
}
