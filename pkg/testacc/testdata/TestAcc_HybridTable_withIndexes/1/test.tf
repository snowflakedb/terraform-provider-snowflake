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
    name = "name"
    type = "VARCHAR(100)"
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
  }

  constraint {
    name    = "pk_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  index {
    name    = "idx_name"
    columns = ["name"]
  }

  index {
    name    = "idx_created"
    columns = ["created_at"]
  }

  comment = var.comment
}
