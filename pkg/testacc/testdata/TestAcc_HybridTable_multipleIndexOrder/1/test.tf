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
    name = "email"
    type = "VARCHAR(255)"
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

  # Indexes in specific order: idx_name, idx_email, idx_created
  index {
    name    = "idx_name"
    columns = ["name"]
  }

  index {
    name    = "idx_email"
    columns = ["email"]
  }

  index {
    name    = "idx_created"
    columns = ["created_at"]
  }

  comment = var.comment
}
