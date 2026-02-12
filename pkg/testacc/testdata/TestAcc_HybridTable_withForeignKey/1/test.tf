resource "snowflake_hybrid_table" "parent" {
  database = var.database
  schema   = var.schema
  name     = var.parent_name

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
    name    = "pk_parent_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  comment = var.comment
}

resource "snowflake_hybrid_table" "child" {
  depends_on = [snowflake_hybrid_table.parent]

  database = var.database
  schema   = var.schema
  name     = var.child_name

  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "parent_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "data"
    type = "VARCHAR(100)"
  }

  constraint {
    name    = "pk_child_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  constraint {
    name    = "fk_parent"
    type    = "FOREIGN KEY"
    columns = ["parent_id"]
    foreign_key {
      table_id = snowflake_hybrid_table.parent.fully_qualified_name
      columns  = ["id"]
    }
  }

  comment = var.comment
}
