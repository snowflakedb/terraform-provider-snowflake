# Simple usage
data "snowflake_security_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_security_integrations.simple.security_integrations
}

# Filtering (like)
data "snowflake_security_integrations" "like" {
  like = "security-integration-name"
}

output "like_output" {
  value = data.snowflake_security_integrations.like.security_integrations
}

# Filtering by prefix (like)
data "snowflake_security_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_security_integrations.like_prefix.security_integrations
}

# Without additional data (to limit the number of calls make for every found security integration)
data "snowflake_security_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SECURITY INTEGRATION for every security integration found and attaches its output to security_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_security_integrations.only_show.security_integrations
}

# Ensure the number of security_integrations is equal to at least one element (with the use of postcondition)
data "snowflake_security_integrations" "assert_with_postcondition" {
  like = "security-integration-name%"
  lifecycle {
    postcondition {
      condition     = length(self.security_integrations) > 0
      error_message = "there should be at least one security integration"
    }
  }
}

# Ensure the number of security_integrations is equal to exactly one element (with the use of check block)
check "security_integration_check" {
  data "snowflake_security_integrations" "assert_with_check_block" {
    like = "security-integration-name"
  }

  assert {
    condition     = length(data.snowflake_security_integrations.assert_with_check_block.security_integrations) == 1
    error_message = "security integrations filtered by '${data.snowflake_security_integrations.assert_with_check_block.like}' returned ${length(data.snowflake_security_integrations.assert_with_check_block.security_integrations)} security integrations where one was expected"
  }
}
