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
    name = "email"
    type = "VARCHAR(255)"
  }

  column {
    name = "username"
    type = "VARCHAR(100)"
  }

  column {
    name = "phone"
    type = "VARCHAR(50)"
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

  constraint {
    name    = "uq_username"
    type    = "UNIQUE"
    columns = ["username"]
  }

  comment = var.comment
}
