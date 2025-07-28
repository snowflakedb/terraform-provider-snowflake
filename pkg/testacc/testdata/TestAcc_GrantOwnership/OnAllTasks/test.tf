resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_task" "test" {
  database      = var.database
  schema        = var.schema
  name          = var.task
  started       = false
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_task" "second_test" {
  database      = var.database
  schema        = var.schema
  name          = var.second_task
  started       = false
  sql_statement = "SELECT CURRENT_TIMESTAMP"
}

resource "snowflake_grant_ownership" "test" {
  depends_on        = [snowflake_task.test, snowflake_task.second_test]
  account_role_name = snowflake_account_role.test.name

  on {
    all {
      object_type_plural = "TASKS"
      in_schema          = "\"${var.database}\".\"${var.schema}\""
    }
  }

  outbound_privileges = "REVOKE"
}
