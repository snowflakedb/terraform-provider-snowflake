resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  # Columns in specific order: id, created_at, updated_at, name, email
  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
  }

  column {
    name = "updated_at"
    type = "TIMESTAMP_NTZ"
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  column {
    name = "email"
    type = "VARCHAR(255)"
  }

  constraint {
    name    = "pk_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  comment = var.comment
}
