# basic resource
resource "snowflake_notebook" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "NOTEBOOK"
}

# complete resource
resource "snowflake_notebook" "complete" {
  name                            = "NOTEBOOK"
  database                        = "DATABASE"
  schema                          = "SCHEMA"
  from                            = "\"<db_name>\".\"<schema_name>\".\"<stage_name>\""
  main_file                       = "MAIN_FILE.ipynb"
  query_warehouse                 = "\"QUERY_WAREHOUSE\""
  idle_auto_shutdown_time_seconds = number_of_seconds
  warehouse                       = "\"WAREHOUSE\""
  comment                         = "comment"
}
