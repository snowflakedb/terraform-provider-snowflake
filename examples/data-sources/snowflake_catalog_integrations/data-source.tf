# Simple usage
data "snowflake_catalog_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_catalog_integrations.simple.catalog_integrations
}

# Filtering (like)
data "snowflake_catalog_integrations" "like" {
  like = "catalog-integration-name"
}

output "like_output" {
  value = data.snowflake_catalog_integrations.like.catalog_integrations
}

# Filtering by prefix (like)
data "snowflake_catalog_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_catalog_integrations.like_prefix.catalog_integrations
}

# Without additional data (to limit the number of calls made for every found catalog integration)
data "snowflake_catalog_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE CATALOG INTEGRATION for every catalog integration found and attaches its output to catalog_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_catalog_integrations.only_show.catalog_integrations
}

# Ensure the number of catalog_integrations is equal to at least one element (with the use of postcondition)
data "snowflake_catalog_integrations" "assert_with_postcondition" {
  like = "catalog-integration-name%"
  lifecycle {
    postcondition {
      condition     = length(self.catalog_integrations) > 0
      error_message = "there should be at least one catalog integration"
    }
  }
}

# Ensure the number of catalog_integrations is equal to exactly one element (with the use of check block)
check "catalog_integration_check" {
  data "snowflake_catalog_integrations" "assert_with_check_block" {
    like = "catalog-integration-name"
  }

  assert {
    condition     = length(data.snowflake_catalog_integrations.assert_with_check_block.catalog_integrations) == 1
    error_message = "catalog integrations filtered by '${data.snowflake_catalog_integrations.assert_with_check_block.like}' returned ${length(data.snowflake_catalog_integrations.assert_with_check_block.catalog_integrations)} catalog integrations where one was expected"
  }
}
