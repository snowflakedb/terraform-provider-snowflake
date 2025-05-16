# Simple usage
data "snowflake_resource_monitors" "simple" {
}

output "simple_output" {
  value = data.snowflake_resource_monitors.simple.resource_monitors
}

# Filtering (like)
data "snowflake_resource_monitors" "like" {
  like = "resource-monitor-name"
}

output "like_output" {
  value = data.snowflake_resource_monitors.like.resource_monitors
}

# Ensure the number of resource monitors is equal to at least one element (with the use of postcondition)
data "snowflake_resource_monitors" "assert_with_postcondition" {
  like = "resource-monitor-name-%"
  lifecycle {
    postcondition {
      condition     = length(self.resource_monitors) > 0
      error_message = "there should be at least one resource monitor"
    }
  }
}

# Ensure the number of resource monitors is equal to exactly one element (with the use of check block)
check "resource_monitor_check" {
  data "snowflake_resource_monitors" "assert_with_check_block" {
    like = "resource-monitor-name"
  }

  assert {
    condition     = length(data.snowflake_resource_monitors.assert_with_check_block.resource_monitors) == 1
    error_message = "Resource monitors filtered by '${data.snowflake_resource_monitors.assert_with_check_block.like}' returned ${length(data.snowflake_resource_monitors.assert_with_check_block.resource_monitors)} resource monitors where one was expected"
  }
}
