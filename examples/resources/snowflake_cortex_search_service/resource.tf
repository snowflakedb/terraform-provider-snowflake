## Basic
resource "snowflake_database" "test" {
  name = "some_database"
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.name
  name     = "some_schema"
}

resource "snowflake_table" "test" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "some_table"
  change_tracking = true
  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }

  column {
    name = "SOME_TEXT"
    type = "VARCHAR"
  }
}

resource "snowflake_cortex_search_service" "test" {
  depends_on = [snowflake_table.test]

  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "some_name"
  on              = "SOME_TEXT"
  target_lag      = "2 minutes"
  warehouse       = "some_warehouse"
  query           = "SELECT SOME_TEXT FROM \"some_database\".\"some_schema\".\"some_table\""
  comment         = "some comment"
  embedding_model = "snowflake-arctic-embed-m-v1.5"
}
