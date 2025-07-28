resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_task" "test" {
  database      = var.database
  schema        = var.schema
  name          = var.task
  warehouse     = var.warehouse
  started       = false
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}
