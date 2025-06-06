resource "snowflake_database_role" "test" {
  name     = var.database_role_name
  database = var.database_name
}

resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = var.database_name
}

resource "snowflake_table" "test" {
  name     = var.table_name
  database = var.database_name
  schema   = snowflake_schema.test.name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_ownership" "test" {
  database_role_name = "\"${var.database_name}\".\"${snowflake_database_role.test.name}\""
  on {
    object_type = "TABLE"
    object_name = "\"${var.database_name}\".\"${snowflake_schema.test.name}\".\"${snowflake_table.test.name}\""
  }
}
