# Simple usage
data "snowflake_network_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_network_policies.simple.network_policies
}

# Filtering (like)
data "snowflake_network_policies" "like" {
  like = "network-policy-name"
}

output "like_output" {
  value = data.snowflake_network_policies.like.network_policies
}

# Without additional data (to limit the number of calls make for every found network policy)
data "snowflake_network_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE NETWORK POLICY for every network policy found and attaches its output to network_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_network_policies.only_show.network_policies
}

# Ensure the number of network policies is equal to at least one element (with the use of postcondition)
data "snowflake_network_policies" "assert_with_postcondition" {
  like = "network-policy-name"
  lifecycle {
    postcondition {
      condition     = length(self.network_policies) > 0
      error_message = "there should be at least one network policy"
    }
  }
}

# Ensure the number of network policies is equal to exactly one element (with the use of check block)
check "network_policy_check" {
  data "snowflake_network_policies" "assert_with_check_block" {
    like = "network-policy-name"
  }

  assert {
    condition     = length(data.snowflake_network_policies.assert_with_check_block.network_policies) == 1
    error_message = "Network policies filtered by '${data.snowflake_network_policies.assert_with_check_block.like}' returned ${length(data.snowflake_network_policies.assert_with_check_block.network_policies)} network policies where one was expected"
  }
}
