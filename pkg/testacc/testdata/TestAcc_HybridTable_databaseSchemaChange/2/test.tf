resource "snowflake_schema" "test_second" {
  database = var.database
  name     = var.second_schema
}

resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = snowflake_schema.test_second.name
  name     = var.name

  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  constraint {
    name    = "pk_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  comment = var.comment
}
