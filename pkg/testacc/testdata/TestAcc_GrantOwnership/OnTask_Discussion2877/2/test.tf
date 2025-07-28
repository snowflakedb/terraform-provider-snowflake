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

resource "snowflake_task" "child" {
  database      = var.database
  schema        = var.schema
  name          = var.child
  warehouse     = var.warehouse
  after         = [snowflake_task.test.fully_qualified_name]
  started       = false
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_grant_ownership" "test" {
  depends_on = [snowflake_task.child]

  account_role_name = snowflake_account_role.test.name

  on {
    object_type = "TASK"
    object_name = snowflake_task.test.fully_qualified_name
  }
}

resource "snowflake_grant_ownership" "child" {
  account_role_name = snowflake_account_role.test.name

  on {
    object_type = "TASK"
    object_name = snowflake_task.child.fully_qualified_name
  }
}
