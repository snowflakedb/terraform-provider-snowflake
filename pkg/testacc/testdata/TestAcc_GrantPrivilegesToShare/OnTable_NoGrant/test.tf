resource "snowflake_table" "test" {
  name     = var.on_table
  database = var.database
  schema   = var.schema
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}
