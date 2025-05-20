resource "snowflake_task" "root" {
  name          = var.tasks[0].name
  database      = var.tasks[0].database
  schema        = var.tasks[0].schema
  started       = var.tasks[0].started
  sql_statement = var.tasks[0].sql_statement

  # Optionals
  dynamic "schedule" {
    for_each = [for element in [var.tasks[0].schedule] : element if element != null]
    content {
      minutes    = lookup(var.tasks[0].schedule, "minutes", null)
      using_cron = lookup(var.tasks[0].schedule, "cron", null)
    }
  }

  comment  = var.tasks[0].comment
  after    = var.tasks[0].after
  finalize = var.tasks[0].finalize

  # Parameters
  suspend_task_after_num_failures = var.tasks[0].suspend_task_after_num_failures
}

resource "snowflake_task" "child" {
  depends_on = [snowflake_task.root]

  name          = var.tasks[1].name
  database      = var.tasks[1].database
  schema        = var.tasks[1].schema
  started       = var.tasks[1].started
  sql_statement = var.tasks[1].sql_statement

  # Optionals
  dynamic "schedule" {
    for_each = [for element in [var.tasks[1].schedule] : element if element != null]
    content {
      minutes    = lookup(var.tasks[1].schedule, "minutes", null)
      using_cron = lookup(var.tasks[1].schedule, "cron", null)
    }
  }

  comment  = var.tasks[1].comment
  after    = var.tasks[1].after
  finalize = var.tasks[1].finalize

  # Parameters
  suspend_task_after_num_failures = var.tasks[1].suspend_task_after_num_failures
}
