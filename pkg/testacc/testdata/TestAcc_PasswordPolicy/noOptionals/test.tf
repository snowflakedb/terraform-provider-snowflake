resource "snowflake_password_policy" "pa" {
  name     = var.name
  database = var.database
  schema   = var.schema
}