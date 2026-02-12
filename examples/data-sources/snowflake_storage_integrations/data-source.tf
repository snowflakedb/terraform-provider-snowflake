# Simple usage
data "snowflake_storage_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_storage_integrations.simple.storage_integrations
}

# Filtering (like)
data "snowflake_storage_integrations" "like" {
  like = "storage-integration-name"
}

output "like_output" {
  value = data.snowflake_storage_integrations.like.storage_integrations
}

# Filtering by prefix (like)
data "snowflake_storage_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_storage_integrations.like_prefix.storage_integrations
}

# Without additional data (to limit the number of calls make for every found storage integration)
data "snowflake_storage_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE STORAGE INTEGRATION for every storage integration found and attaches its output to storage_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_storage_integrations.only_show.storage_integrations
}

# Ensure the number of storage_integrations is equal to at least one element (with the use of postcondition)
data "snowflake_storage_integrations" "assert_with_postcondition" {
  like = "storage-integration-name%"
  lifecycle {
    postcondition {
      condition     = length(self.storage_integrations) > 0
      error_message = "there should be at least one storage integration"
    }
  }
}

# Ensure the number of storage_integrations is equal to exactly one element (with the use of check block)
check "storage_integration_check" {
  data "snowflake_storage_integrations" "assert_with_check_block" {
    like = "storage-integration-name"
  }

  assert {
    condition     = length(data.snowflake_storage_integrations.assert_with_check_block.storage_integrations) == 1
    error_message = "storage integrations filtered by '${data.snowflake_storage_integrations.assert_with_check_block.like}' returned ${length(data.snowflake_storage_integrations.assert_with_check_block.storage_integrations)} storage integrations where one was expected"
  }
}
