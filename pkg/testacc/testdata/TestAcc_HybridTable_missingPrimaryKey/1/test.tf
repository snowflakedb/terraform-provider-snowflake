resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  constraint {
    name    = "uq_name"
    type    = "UNIQUE"
    columns = ["name"]
  }
}
