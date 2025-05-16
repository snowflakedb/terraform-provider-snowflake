# Simple usage
data "snowflake_row_access_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_row_access_policies.simple.row_access_policies
}

# Filtering (like)
data "snowflake_row_access_policies" "like" {
  like = "row-access-policy-name"
}

output "like_output" {
  value = data.snowflake_row_access_policies.like.row_access_policies
}

# Filtering by prefix (like)
data "snowflake_row_access_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_row_access_policies.like_prefix.row_access_policies
}

# Filtering (limit)
data "snowflake_row_access_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_row_access_policies.limit.row_access_policies
}

# Filtering (in)
data "snowflake_row_access_policies" "in" {
  in {
    database = "database"
  }
}

output "in_output" {
  value = data.snowflake_row_access_policies.in.row_access_policies
}

# Without additional data (to limit the number of calls make for every found row access policy)
data "snowflake_row_access_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE ROW ACCESS POLICY for every row access policy found and attaches its output to row_access_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_row_access_policies.only_show.row_access_policies
}

# Ensure the number of row access policies is equal to at least one element (with the use of postcondition)
data "snowflake_row_access_policies" "assert_with_postcondition" {
  like = "row-access-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.row_access_policies) > 0
      error_message = "there should be at least one row access policy"
    }
  }
}

# Ensure the number of row access policies is equal to exactly one element (with the use of check block)
check "row_access_policy_check" {
  data "snowflake_row_access_policies" "assert_with_check_block" {
    like = "row-access-policy-name"
  }

  assert {
    condition     = length(data.snowflake_row_access_policies.assert_with_check_block.row_access_policies) == 1
    error_message = "row access policies filtered by '${data.snowflake_row_access_policies.assert_with_check_block.like}' returned ${length(data.snowflake_row_access_policies.assert_with_check_block.row_access_policies)} row access policies where one was expected"
  }
}
