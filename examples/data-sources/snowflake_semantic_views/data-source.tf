# Simple usage
data "snowflake_semantic_views" "simple" {
}

output "simple_output" {
  value = data.snowflake_semantic_views.simple.semantic_views
}

# Filtering (like)
data "snowflake_semantic_views" "like" {
  like = "semantic-view-name"
}

output "like_output" {
  value = data.snowflake_semantic_views.like.semantic_views
}

# Filtering by prefix (like)
data "snowflake_semantic_views" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_semantic_views.like_prefix.semantic_views
}

# Filtering (starts_with)
data "snowflake_semantic_views" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_semantic_views.starts_with.semantic_views
}

# Filtering (in)
data "snowflake_semantic_views" "in_account" {
  in {
    account = true
  }
}

data "snowflake_semantic_views" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_semantic_views" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_semantic_views.in_account.semantic_views,
    "database" : data.snowflake_semantic_views.in_database.semantic_views,
    "schema" : data.snowflake_semantic_views.in_schema.semantic_views,
  }
}

# Filtering (limit)
data "snowflake_semantic_views" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_semantic_views.limit.semantic_views
}

# Without additional data (to limit the number of calls made for every found semantic view)
data "snowflake_semantic_views" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SEMANTIC VIEW for every semantic view found and attaches its output to semantic_views.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_semantic_views.only_show.semantic_views
}

# Ensure the number of semantic views is equal to at least one element (with the use of postcondition)
data "snowflake_semantic_views" "assert_with_postcondition" {
  like = "semantic-view-name%"
  lifecycle {
    postcondition {
      condition     = length(self.semantic_views) > 0
      error_message = "there should be at least one semantic view"
    }
  }
}

# Ensure the number of semantic views is equal to exactly one element (with the use of check block)
check "semantic_view_check" {
  data "snowflake_semantic_views" "assert_with_check_block" {
    like = "semantic-view-name"
  }

  assert {
    condition     = length(data.snowflake_semantic_views.assert_with_check_block.semantic_views) == 1
    error_message = "semantic views filtered by '${data.snowflake_semantic_views.assert_with_check_block.like}' returned ${length(data.snowflake_semantic_views.assert_with_check_block.semantic_views)} semantic views where one was expected"
  }
}