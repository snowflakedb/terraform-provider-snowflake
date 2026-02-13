resource "snowflake_hybrid_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.name

  # Numeric types
  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "col_int"
    type = "INT"
  }

  column {
    name = "col_float"
    type = "FLOAT"
  }

  column {
    name = "col_decimal"
    type = "DECIMAL(10,2)"
  }

  # String types
  column {
    name = "col_varchar"
    type = "VARCHAR(100)"
  }

  column {
    name = "col_char"
    type = "CHAR(10)"
  }

  column {
    name = "col_text"
    type = "VARCHAR(134217728)"
  }

  # Date/Time types
  column {
    name = "col_date"
    type = "DATE"
  }

  column {
    name = "col_time"
    type = "TIME"
  }

  column {
    name = "col_timestamp_ntz"
    type = "TIMESTAMP_NTZ"
  }

  column {
    name = "col_timestamp_ltz"
    type = "TIMESTAMP_LTZ"
  }

  # Semi-structured types
  column {
    name = "col_variant"
    type = "VARIANT"
  }

  column {
    name = "col_object"
    type = "OBJECT"
  }

  column {
    name = "col_array"
    type = "ARRAY"
  }

  # Boolean type
  column {
    name = "col_boolean"
    type = "BOOLEAN"
  }

  constraint {
    name    = "pk_id"
    type    = "PRIMARY KEY"
    columns = ["id"]
  }

  comment = var.comment
}
