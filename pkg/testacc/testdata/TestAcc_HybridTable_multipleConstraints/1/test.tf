resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "email"
    type     = "VARCHAR(255)"
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

  constraint {
    name    = "uq_email"
    type    = "UNIQUE"
    columns = ["email"]
  }

  comment = var.comment
}
