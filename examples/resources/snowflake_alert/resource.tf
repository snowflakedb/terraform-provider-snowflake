## Basic resource (serverless, on a schedule)
resource "snowflake_alert" "basic" {
  database = "database"
  schema   = "schema"
  name     = "alert"
  alert_schedule {
    interval = 10
  }
  condition = "select 1 as c"
  action    = "select 1 as c"
}

## Complete resource
resource "snowflake_alert" "complete" {
  database  = "database"
  schema    = "schema"
  name      = "alert"
  warehouse = "warehouse"
  alert_schedule {
    cron {
      expression = "0 * * * *"
      time_zone  = "UTC"
    }
  }
  condition = "select 1 as c"
  action    = "select 1 as c"
  enabled   = true
  comment   = "my alert"
}

## Alert on new data (no schedule); the condition table needs change tracking
resource "snowflake_alert" "on_new_data" {
  database  = "database"
  schema    = "schema"
  name      = "alert_on_new_data"
  warehouse = "warehouse"
  condition = "select * from database.schema.table"
  action    = "select 1 as c"
}
