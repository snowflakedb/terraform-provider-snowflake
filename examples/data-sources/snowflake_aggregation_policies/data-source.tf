# Simple usage
data "snowflake_aggregation_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_aggregation_policies.simple.aggregation_policies
}

# Filtering (like)
data "snowflake_aggregation_policies" "like" {
  like = "aggregation-policy-name"
}

output "like_output" {
  value = data.snowflake_aggregation_policies.like.aggregation_policies
}

# Filtering by prefix (like)
data "snowflake_aggregation_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_aggregation_policies.like_prefix.aggregation_policies
}

# Filtering (starts_with)
data "snowflake_aggregation_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_aggregation_policies.starts_with.aggregation_policies
}

# Filtering (limit)
data "snowflake_aggregation_policies" "limit" {
  limit {
    rows = 10
  }
}

output "limit_output" {
  value = data.snowflake_aggregation_policies.limit.aggregation_policies
}

# Filtering (in)
data "snowflake_aggregation_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_aggregation_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_aggregation_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_aggregation_policies.in_account.aggregation_policies,
    "database" : data.snowflake_aggregation_policies.in_database.aggregation_policies,
    "schema" : data.snowflake_aggregation_policies.in_schema.aggregation_policies,
  }
}

# Without additional data (to limit the number of calls made for every found aggregation policy)
data "snowflake_aggregation_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE AGGREGATION POLICY for every aggregation policy found and attaches its output to aggregation_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_aggregation_policies.only_show.aggregation_policies
}

# Ensure the number of aggregation policies is equal to at least one element (with the use of postcondition)
data "snowflake_aggregation_policies" "assert_with_postcondition" {
  like = "aggregation-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.aggregation_policies) > 0
      error_message = "there should be at least one aggregation policy"
    }
  }
}

# Ensure the number of aggregation policies is equal to exactly one element (with the use of check block)
check "aggregation_policy_check" {
  data "snowflake_aggregation_policies" "assert_with_check_block" {
    like = "aggregation-policy-name"
  }

  assert {
    condition     = length(data.snowflake_aggregation_policies.assert_with_check_block.aggregation_policies) == 1
    error_message = "aggregation policies filtered by '${data.snowflake_aggregation_policies.assert_with_check_block.like}' returned ${length(data.snowflake_aggregation_policies.assert_with_check_block.aggregation_policies)} aggregation policies where one was expected"
  }
}
