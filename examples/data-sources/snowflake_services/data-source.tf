# Simple usage
data "snowflake_services" "simple" {
}

output "simple_output" {
  value = data.snowflake_services.simple.services
}

# Filtering (like)
data "snowflake_services" "like" {
  like = "service-name"
}

output "like_output" {
  value = data.snowflake_services.like.services
}

# Filtering by prefix (like)
data "snowflake_services" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_services.like_prefix.services
}

# Filtering (starts_with)
data "snowflake_services" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_services.starts_with.services
}

# Filtering (in)
data "snowflake_services" "in_account" {
  in {
    account = true
  }
}

data "snowflake_services" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_services" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_services" "in_compute_pool" {
  in {
    compute_pool = "<compute_pool_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_services.in_account.services,
    "database" : data.snowflake_services.in_database.services,
    "schema" : data.snowflake_services.in_schema.services,
    "compute_pool" : data.snowflake_services.in_compute_pool.services,
  }
}

# Filtering (limit)
data "snowflake_services" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_services.limit.services
}

# Filtering (jobs only)
data "snowflake_services" "jobs_only" {
  jobs_only = true
}

output "jobs_only_output" {
  value = data.snowflake_services.jobs_only.services
}

# Filtering (exclude jobs)
data "snowflake_services" "exclude_jobs" {
  exclude_jobs = true
}

output "exclude_jobs_output" {
  value = data.snowflake_services.exclude_jobs.services
}

# Without additional data (to limit the number of calls make for every found service)
data "snowflake_services" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SERVICE for every service found and attaches its output to services.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_services.only_show.services
}

# Ensure the number of services is equal to at least one element (with the use of postcondition)
data "snowflake_services" "assert_with_postcondition" {
  like = "service-name%"
  lifecycle {
    postcondition {
      condition     = length(self.services) > 0
      error_message = "there should be at least one service"
    }
  }
}

# Ensure the number of services is equal to exactly one element (with the use of check block)
check "service_check" {
  data "snowflake_services" "assert_with_check_block" {
    like = "service-name"
  }

  assert {
    condition     = length(data.snowflake_services.assert_with_check_block.services) == 1
    error_message = "services filtered by '${data.snowflake_services.assert_with_check_block.like}' returned ${length(data.snowflake_services.assert_with_check_block.services)} services where one was expected"
  }
}
