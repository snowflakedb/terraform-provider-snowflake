resource "snowflake_database" "d" {
  name = "some_db"
}

resource "snowflake_schema" "s" {
  name     = "some_schema"
  database = snowflake_database.d.name
}

resource "snowflake_table" "t" {
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  name     = "some_table"

  column {
    name     = "col1"
    type     = "text"
    nullable = false
  }

  column {
    name     = "col2"
    type     = "text"
    nullable = false
  }

  column {
    name     = "col3"
    type     = "text"
    nullable = false
  }
}

resource "snowflake_table" "fk_t" {
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  name     = "fk_table"

  column {
    name     = "fk_col1"
    type     = "text"
    nullable = false
  }

  column {
    name     = "fk_col2"
    type     = "text"
    nullable = false
  }
}

resource "snowflake_table_constraint" "primary_key" {
  name     = "myconstraint"
  type     = "PRIMARY KEY"
  table_id = snowflake_table.t.fully_qualified_name
  columns  = ["col1"]
}

resource "snowflake_table_constraint" "foreign_key" {
  name     = "myconstraintfk"
  type     = "FOREIGN KEY"
  table_id = snowflake_table.t.fully_qualified_name
  columns  = ["col2"]

  foreign_key_properties {
    references {
      table_id = snowflake_table.fk_t.fully_qualified_name
      columns  = ["fk_col1"]
    }
  }

  enforced   = false
  deferrable = false
  initially  = "IMMEDIATE"
}

resource "snowflake_table_constraint" "unique" {
  name     = "unique"
  type     = "UNIQUE"
  table_id = snowflake_table.t.fully_qualified_name
  columns  = ["col3"]
}
