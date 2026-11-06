# Simple usage
data "snowflake_alerts" "simple" {
}

output "simple_output" {
  value = data.snowflake_alerts.simple.alerts
}

# Filtering (like)
data "snowflake_alerts" "like" {
  like = "alert-name"
}

output "like_output" {
  value = data.snowflake_alerts.like.alerts
}

# Filtering by prefix (like)
data "snowflake_alerts" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_alerts.like_prefix.alerts
}

# Filtering (starts_with)
data "snowflake_alerts" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_alerts.starts_with.alerts
}

# Filtering (limit)
data "snowflake_alerts" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_alerts.limit.alerts
}

# Filtering (in)
data "snowflake_alerts" "in_account" {
  in {
    account = true
  }
}

data "snowflake_alerts" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_alerts" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

data "snowflake_alerts" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_alerts" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_alerts.in_account.alerts,
    "database" : data.snowflake_alerts.in_database.alerts,
    "schema" : data.snowflake_alerts.in_schema.alerts,
    "application" : data.snowflake_alerts.in_application.alerts,
    "application_package" : data.snowflake_alerts.in_application_package.alerts,
  }
}

# Without additional data (to limit the number of calls made for every found alert)
data "snowflake_alerts" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE ALERT for every alert found and attaches its output to alerts.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_alerts.only_show.alerts
}

# Ensure the number of alerts is equal to at least one element (with the use of postcondition)
data "snowflake_alerts" "assert_with_postcondition" {
  like = "alert-name%"
  lifecycle {
    postcondition {
      condition     = length(self.alerts) > 0
      error_message = "there should be at least one alert"
    }
  }
}

# Ensure the number of alerts is equal to exactly one element (with the use of check block)
check "alert_check" {
  data "snowflake_alerts" "assert_with_check_block" {
    like = "alert-name"
  }

  assert {
    condition     = length(data.snowflake_alerts.assert_with_check_block.alerts) == 1
    error_message = "alerts filtered by '${data.snowflake_alerts.assert_with_check_block.like}' returned ${length(data.snowflake_alerts.assert_with_check_block.alerts)} alerts where one was expected"
  }
}
