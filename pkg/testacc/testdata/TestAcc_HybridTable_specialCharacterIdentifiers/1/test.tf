resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  column {
    name     = "user-id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "user name"
    type = "VARCHAR(100)"
  }

  column {
    name = "created@time"
    type = "TIMESTAMP_NTZ"
  }

  constraint {
    name    = "pk_user_id"
    type    = "PRIMARY KEY"
    columns = ["user-id"]
  }

  comment = var.comment
}
