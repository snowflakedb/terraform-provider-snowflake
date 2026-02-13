resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  column {
    name     = "userId"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "userName"
    type = "VARCHAR(100)"
  }

  column {
    name = "createdAt"
    type = "TIMESTAMP_NTZ"
  }

  constraint {
    name    = "pk_userId"
    type    = "PRIMARY KEY"
    columns = ["userId"]
  }

  comment = var.comment
}
