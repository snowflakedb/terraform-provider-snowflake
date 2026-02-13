resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  column {
    name     = "tenant_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "user_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "email"
    type = "VARCHAR(255)"
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  constraint {
    name    = "pk_composite"
    type    = "PRIMARY KEY"
    columns = ["tenant_id", "user_id"]
  }

  constraint {
    name    = "uq_email_tenant"
    type    = "UNIQUE"
    columns = ["email", "tenant_id"]
  }

  comment = var.comment
}
