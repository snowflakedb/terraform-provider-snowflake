resource "snowflake_database" "test" {
  name    = var.name
  catalog = var.catalog
}
