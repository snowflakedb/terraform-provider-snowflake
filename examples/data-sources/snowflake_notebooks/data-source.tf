# Simple usage
data "snowflake_notebooks" "simple" {
}

output "simple_output" {
  value = data.snowflake_notebooks.simple.notebooks
}

# Filtering (like)
data "snowflake_notebooks" "like" {
  like = "notebook-name"
}

output "like_output" {
  value = data.snowflake_notebooks.like.notebooks
}

# Filtering by prefix (like)
data "snowflake_notebooks" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_notebooks.like_prefix.notebooks
}

# Filtering (starts_with)
data "snowflake_notebooks" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_notebooks.starts_with.notebooks
}

# Filtering (limit)
data "snowflake_notebooks" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_notebooks.limit.notebooks
}

# Without additional data (to limit the number of calls make for every found notebook)
data "snowflake_notebooks" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE NOTEBOOK for every notebook found and attaches its output to notebooks.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_notebooks.only_show.notebooks
}

# Ensure the number of notebooks is equal to at least one element (with the use of postcondition)
data "snowflake_notebooks" "assert_with_postcondition" {
  like = "notebook-name%"
  lifecycle {
    postcondition {
      condition     = length(self.notebooks) > 0
      error_message = "there should be at least one notebook"
    }
  }
}

# Ensure the number of notebooks is equal to exactly one element (with the use of check block)
check "notebook_check" {
  data "snowflake_notebooks" "assert_with_check_block" {
    like = "notebook-name"
  }

  assert {
    condition     = length(data.snowflake_notebooks.assert_with_check_block.notebooks) == 1
    error_message = "notebooks filtered by '${data.snowflake_notebooks.assert_with_check_block.like}' returned ${length(data.snowflake_notebooks.assert_with_check_block.notebooks)} notebooks where one was expected"
  }
}
