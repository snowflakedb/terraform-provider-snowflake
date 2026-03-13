# Simple usage
data "snowflake_external_volumes" "simple" {
}

output "simple_output" {
  value = data.snowflake_external_volumes.simple.external_volumes
}

# Filtering (like)
data "snowflake_external_volumes" "like" {
  like = "external-volume-name"
}

output "like_output" {
  value = data.snowflake_external_volumes.like.external_volumes
}

# Filtering by prefix (like)
data "snowflake_external_volumes" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_external_volumes.like_prefix.external_volumes
}

# Without additional data (to limit the number of calls made for every found external volume)
data "snowflake_external_volumes" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE EXTERNAL VOLUME for every external volume found and attaches its output to external_volumes.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_external_volumes.only_show.external_volumes
}

# Ensure the number of external volumes is equal to at least one element (with the use of postcondition)
data "snowflake_external_volumes" "assert_with_postcondition" {
  like = "external-volume-name%"
  lifecycle {
    postcondition {
      condition     = length(self.external_volumes) > 0
      error_message = "there should be at least one external volume"
    }
  }
}

# Ensure the number of external volumes is equal to exactly one element (with the use of check block)
check "external_volume_check" {
  data "snowflake_external_volumes" "assert_with_check_block" {
    like = "external-volume-name"
  }

  assert {
    condition     = length(data.snowflake_external_volumes.assert_with_check_block.external_volumes) == 1
    error_message = "external volumes filtered by '${data.snowflake_external_volumes.assert_with_check_block.like}' returned ${length(data.snowflake_external_volumes.assert_with_check_block.external_volumes)} external volumes where one was expected"
  }
}
