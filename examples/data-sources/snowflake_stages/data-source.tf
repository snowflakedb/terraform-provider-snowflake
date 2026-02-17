# Simple usage
data "snowflake_stages" "simple" {
}

output "simple_output" {
  value = data.snowflake_stages.simple.stages
}

# Filtering (like)
data "snowflake_stages" "like" {
  like = "stage-name"
}

output "like_output" {
  value = data.snowflake_stages.like.stages
}

# Filtering by prefix (like)
data "snowflake_stages" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_stages.like_prefix.stages
}

# Without additional data (to limit the number of calls make for every found stage)
data "snowflake_stages" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE STAGE for every stage found and attaches its output to stages.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_stages.only_show.stages
}

# Ensure the number of stages is equal to at least one element (with the use of postcondition)
data "snowflake_stages" "assert_with_postcondition" {
  like = "stage-name%"
  lifecycle {
    postcondition {
      condition     = length(self.stages) > 0
      error_message = "there should be at least one stage"
    }
  }
}

# Ensure the number of stages is equal to exactly one element (with the use of check block)
check "stage_check" {
  data "snowflake_stages" "assert_with_check_block" {
    like = "stage-name"
  }

  assert {
    condition     = length(data.snowflake_stages.assert_with_check_block.stages) == 1
    error_message = "stages filtered by '${data.snowflake_stages.assert_with_check_block.like}' returned ${length(data.snowflake_stages.assert_with_check_block.stages)} stages where one was expected"
  }
}
