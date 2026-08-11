# Simple usage
data "snowflake_external_access_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_external_access_integrations.simple.external_access_integrations
}

# Filtering (like)
data "snowflake_external_access_integrations" "like" {
  like = "external-access-integration-name"
}

output "like_output" {
  value = data.snowflake_external_access_integrations.like.external_access_integrations
}

# Filtering by prefix (like)
data "snowflake_external_access_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_external_access_integrations.like_prefix.external_access_integrations
}

# Without additional data (to limit the number of calls made for every found external access integration)
data "snowflake_external_access_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE EXTERNAL ACCESS INTEGRATION for every external access integration found and attaches its output to external_access_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_external_access_integrations.only_show.external_access_integrations
}

# Ensure the number of external access integrations is equal to at least one element (with the use of postcondition)
data "snowflake_external_access_integrations" "assert_with_postcondition" {
  like = "external-access-integration-name%"
  lifecycle {
    postcondition {
      condition     = length(self.external_access_integrations) > 0
      error_message = "there should be at least one external access integration"
    }
  }
}

# Ensure the number of external access integrations is equal to exactly one element (with the use of check block)
check "external_access_integration_check" {
  data "snowflake_external_access_integrations" "assert_with_check_block" {
    like = "external-access-integration-name"
  }

  assert {
    condition     = length(data.snowflake_external_access_integrations.assert_with_check_block.external_access_integrations) == 1
    error_message = "external access integrations filtered by '${data.snowflake_external_access_integrations.assert_with_check_block.like}' returned ${length(data.snowflake_external_access_integrations.assert_with_check_block.external_access_integrations)} external access integrations where one was expected"
  }
}
