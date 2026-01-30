# basic resource
resource "snowflake_internal_stage" "basic" {
  name     = "my_internal_stage"
  database = "my_database"
  schema   = "my_schema"
}

# complete resource
resource "snowflake_internal_stage" "complete" {
  name     = "complete_stage"
  database = "my_database"
  schema   = "my_schema"

  encryption {
    snowflake_full {}
  }

  directory {
    enable       = true
    auto_refresh = false
  }

  comment = "Fully configured internal stage"
}
