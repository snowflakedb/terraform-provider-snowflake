# Simple usage
data "snowflake_views" "simple" {
}

output "simple_output" {
  value = data.snowflake_views.simple.views
}

# Filtering (like)
data "snowflake_views" "like" {
  like = "view-name"
}

output "like_output" {
  value = data.snowflake_views.like.views
}

# Filtering by prefix (like)
data "snowflake_views" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_views.like_prefix.views
}

# Filtering (limit)
data "snowflake_views" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_views.limit.views
}

# Filtering (in)
data "snowflake_views" "in_account" {
  in {
    account = true
  }
}

data "snowflake_views" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_views" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

# Without additional data (to limit the number of calls make for every found view)
data "snowflake_views" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE VIEW for every view found and attaches its output to views.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_views.only_show.views
}

# Ensure the number of views is equal to at least one element (with the use of postcondition)
data "snowflake_views" "assert_with_postcondition" {
  like = "view-name%"
  lifecycle {
    postcondition {
      condition     = length(self.views) > 0
      error_message = "there should be at least one view"
    }
  }
}

# Ensure the number of views is equal to exactly one element (with the use of check block)
check "view_check" {
  data "snowflake_views" "assert_with_check_block" {
    like = "view-name"
  }

  assert {
    condition     = length(data.snowflake_views.assert_with_check_block.views) == 1
    error_message = "views filtered by '${data.snowflake_views.assert_with_check_block.like}' returned ${length(data.snowflake_views.assert_with_check_block.views)} views where one was expected"
  }
}
