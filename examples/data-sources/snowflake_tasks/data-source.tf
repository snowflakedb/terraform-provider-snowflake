# Simple usage
data "snowflake_tasks" "simple" {
}

output "simple_output" {
  value = data.snowflake_tasks.simple.tasks
}

# Filtering (like)
data "snowflake_tasks" "like" {
  like = "task-name"
}

output "like_output" {
  value = data.snowflake_tasks.like.tasks
}

# Filtering (in - account - database - schema - application - application package)
data "snowflake_tasks" "in_account" {
  in {
    account = true
  }
}

data "snowflake_tasks" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_tasks" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_tasks" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_tasks" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_tasks.in_account.tasks,
    "database" : data.snowflake_tasks.in_database.tasks,
    "schema" : data.snowflake_tasks.in_schema.tasks,
    "application" : data.snowflake_tasks.in_application.tasks,
    "application_package" : data.snowflake_tasks.in_application_package.tasks,
  }
}

# Filtering (root only tasks)
data "snowflake_tasks" "root_only" {
  root_only = true
}

output "root_only_output" {
  value = data.snowflake_tasks.root_only.tasks
}

# Filtering (starts_with)
data "snowflake_tasks" "starts_with" {
  starts_with = "task-"
}

output "starts_with_output" {
  value = data.snowflake_tasks.starts_with.tasks
}

# Filtering (limit)
data "snowflake_tasks" "limit" {
  limit {
    rows = 10
    from = "task-"
  }
}

output "limit_output" {
  value = data.snowflake_tasks.limit.tasks
}

# Without additional data (to limit the number of calls make for every found task)
data "snowflake_tasks" "only_show" {
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR task for every task found and attaches its output to tasks.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_tasks.only_show.tasks
}

# Ensure the number of tasks is equal to at least one element (with the use of postcondition)
data "snowflake_tasks" "assert_with_postcondition" {
  starts_with = "task-name"
  lifecycle {
    postcondition {
      condition     = length(self.tasks) > 0
      error_message = "there should be at least one task"
    }
  }
}

# Ensure the number of tasks is equal to exactly one element (with the use of check block)
check "task_check" {
  data "snowflake_tasks" "assert_with_check_block" {
    like = "task-name"
  }

  assert {
    condition     = length(data.snowflake_tasks.assert_with_check_block.tasks) == 1
    error_message = "tasks filtered by '${data.snowflake_tasks.assert_with_check_block.like}' returned ${length(data.snowflake_tasks.assert_with_check_block.tasks)} tasks where one was expected"
  }
}
