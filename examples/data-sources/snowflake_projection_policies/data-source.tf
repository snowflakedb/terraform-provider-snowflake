# Simple usage
data "snowflake_projection_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_projection_policies.simple.projection_policies
}

# Filtering (like)
data "snowflake_projection_policies" "like" {
  like = "projection-policy-name"
}

output "like_output" {
  value = data.snowflake_projection_policies.like.projection_policies
}

# Filtering by prefix (like)
data "snowflake_projection_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_projection_policies.like_prefix.projection_policies
}

# Filtering (starts_with)
data "snowflake_projection_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_projection_policies.starts_with.projection_policies
}

# Filtering (limit)
data "snowflake_projection_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_projection_policies.limit.projection_policies
}

# Filtering (in)
data "snowflake_projection_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_projection_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_projection_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_projection_policies.in_account.projection_policies,
    "database" : data.snowflake_projection_policies.in_database.projection_policies,
    "schema" : data.snowflake_projection_policies.in_schema.projection_policies,
  }
}

# Without additional data (to limit the number of calls made for every found projection policy)
data "snowflake_projection_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE PROJECTION POLICY for every projection policy found and attaches its output to projection_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_projection_policies.only_show.projection_policies
}

# Ensure the number of projection policies is equal to at least one element (with the use of postcondition)
data "snowflake_projection_policies" "assert_with_postcondition" {
  like = "projection-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.projection_policies) > 0
      error_message = "there should be at least one projection policy"
    }
  }
}

# Ensure the number of projection policies is equal to exactly one element (with the use of check block)
check "projection_policy_check" {
  data "snowflake_projection_policies" "assert_with_check_block" {
    like = "projection-policy-name"
  }

  assert {
    condition     = length(data.snowflake_projection_policies.assert_with_check_block.projection_policies) == 1
    error_message = "projection policies filtered by '${data.snowflake_projection_policies.assert_with_check_block.like}' returned ${length(data.snowflake_projection_policies.assert_with_check_block.projection_policies)} projection policies where one was expected"
  }
}
