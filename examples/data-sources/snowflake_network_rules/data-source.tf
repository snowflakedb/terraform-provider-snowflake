# Simple usage
data "snowflake_network_rules" "simple" {
}

output "simple_output" {
  value = data.snowflake_network_rules.simple.network_rules
}

# Filtering (like)
data "snowflake_network_rules" "like" {
  like = "network-rule-name"
}

output "like_output" {
  value = data.snowflake_network_rules.like.network_rules
}

# Filtering by prefix (like)
data "snowflake_network_rules" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_network_rules.like_prefix.network_rules
}

# Without additional data (to limit the number of calls make for every found network rule)
data "snowflake_network_rules" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE NETWORK RULE for every network rule found and attaches its output to network_rules.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_network_rules.only_show.network_rules
}

# Ensure the number of network rules is equal to at least one element (with the use of postcondition)
data "snowflake_network_rules" "assert_with_postcondition" {
  like = "network-rule-name%"
  lifecycle {
    postcondition {
      condition     = length(self.network_rules) > 0
      error_message = "there should be at least one network rule"
    }
  }
}

# Ensure the number of network rules is equal to exactly one element (with the use of check block)
check "network_rule_check" {
  data "snowflake_network_rules" "assert_with_check_block" {
    like = "network-rule-name"
  }

  assert {
    condition     = length(data.snowflake_network_rules.assert_with_check_block.network_rules) == 1
    error_message = "network rules filtered by '${data.snowflake_network_rules.assert_with_check_block.like}' returned ${length(data.snowflake_network_rules.assert_with_check_block.network_rules)} network rules where one was expected"
  }
}
