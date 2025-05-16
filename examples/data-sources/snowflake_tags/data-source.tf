# Simple usage
data "snowflake_tags" "simple" {
}

output "simple_output" {
  value = data.snowflake_tags.simple.tags
}

# Filtering (like)
data "snowflake_tags" "like" {
  like = "tag-name"
}

output "like_output" {
  value = data.snowflake_tags.like.tags
}

# Filtering by prefix (like)
data "snowflake_tags" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_tags.like_prefix.tags
}

# Filtering (in)
data "snowflake_tags" "in_account" {
  in {
    account = true
  }
}

data "snowflake_tags" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_tags" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_tags" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_tags" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_tags.in_account.tags,
    "database" : data.snowflake_tags.in_database.tags,
    "schema" : data.snowflake_tags.in_schema.tags,
    "application" : data.snowflake_tags.in_application.tags,
    "application_package" : data.snowflake_tags.in_application_package.tags,
  }
}

output "in_output" {
  value = data.snowflake_tags.in.tags
}

# Ensure the number of tags is equal to at least one element (with the use of postcondition)
data "snowflake_tags" "assert_with_postcondition" {
  like = "tag-name%"
  lifecycle {
    postcondition {
      condition     = length(self.tags) > 0
      error_message = "there should be at least one tag"
    }
  }
}

# Ensure the number of tags is equal to exactly one element (with the use of check block)
check "tag_check" {
  data "snowflake_tags" "assert_with_check_block" {
    like = "tag-name"
  }

  assert {
    condition     = length(data.snowflake_tags.assert_with_check_block.tags) == 1
    error_message = "tags filtered by '${data.snowflake_tags.assert_with_check_block.like}' returned ${length(data.snowflake_tags.assert_with_check_block.tags)} tags where one was expected"
  }
}
