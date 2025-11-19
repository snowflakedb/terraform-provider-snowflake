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
  from                            = snowflake_stage.test.fully_qualifed_name
  main_file                       = "my_notebook.ipynb"
  query_warehouse                 = snowflake_warehouse.test.name
  idle_auto_shutdown_time_seconds = 2400
  warehouse                       = snowflake_warehouse.test.name
  comment                         = "Lorem ipsum"
}
