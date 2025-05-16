# Simple usage
data "snowflake_masking_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_masking_policies.simple.masking_policies
}

# Filtering (like)
data "snowflake_masking_policies" "like" {
  like = "masking-policy-name"
}

output "like_output" {
  value = data.snowflake_masking_policies.like.masking_policies
}

# Filtering by prefix (like)
data "snowflake_masking_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_masking_policies.like_prefix.masking_policies
}

# Filtering (limit)
data "snowflake_masking_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_masking_policies.limit.masking_policies
}

# Filtering (in)
data "snowflake_masking_policies" "in" {
  in {
    database = "database"
  }
}

output "in_output" {
  value = data.snowflake_masking_policies.in.masking_policies
}

# Without additional data (to limit the number of calls make for every found masking policy)
data "snowflake_masking_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE MASKING POLICY for every masking policy found and attaches its output to masking_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_masking_policies.only_show.masking_policies
}

# Ensure the number of masking policies is equal to at least one element (with the use of postcondition)
data "snowflake_masking_policies" "assert_with_postcondition" {
  like = "masking-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.masking_policies) > 0
      error_message = "there should be at least one masking policy"
    }
  }
}

# Ensure the number of masking policies is equal to exactly one element (with the use of check block)
check "masking_policy_check" {
  data "snowflake_masking_policies" "assert_with_check_block" {
    like = "masking-policy-name"
  }

  assert {
    condition     = length(data.snowflake_masking_policies.assert_with_check_block.masking_policies) == 1
    error_message = "masking policies filtered by '${data.snowflake_masking_policies.assert_with_check_block.like}' returned ${length(data.snowflake_masking_policies.assert_with_check_block.masking_policies)} masking policies where one was expected"
  }
}
