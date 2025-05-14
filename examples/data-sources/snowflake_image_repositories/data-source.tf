# Simple usage
data "snowflake_image_repositories" "simple" {
}

output "simple_output" {
  value = data.snowflake_image_repositories.simple.image_repositories
}

# Filtering (like)
data "snowflake_image_repositories" "like" {
  like = "image-repository-name"
}

output "like_output" {
  value = data.snowflake_image_repositories.like.image_repositories
}

# Filtering by prefix (like)
data "snowflake_image_repositories" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_image_repositories.like_prefix.image_repositories
}

# Filtering (in)
data "snowflake_image_repositories" "in_account" {
  in {
    account = true
  }
}

data "snowflake_image_repositories" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_image_repositories" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_image_repositories.in_account.image_repositories,
    "database" : data.snowflake_image_repositories.in_database.image_repositories,
    "schema" : data.snowflake_image_repositories.in_schema.image_repositories,
  }
}

output "in_output_simple" {
  value = data.snowflake_image_repositories.in.image_repositories
}

# Ensure the number of image repositories is equal to at least one element (with the use of postcondition)
data "snowflake_image_repositories" "assert_with_postcondition" {
  like = "image-repository-name%"
  lifecycle {
    postcondition {
      condition     = length(self.image_repositories) > 0
      error_message = "there should be at least one image repository"
    }
  }
}

# Ensure the number of image repositories is equal to at exactly one element (with the use of check block)
check "image_repository_check" {
  data "snowflake_image_repositories" "assert_with_check_block" {
    like = "image-repository-name"
  }

  assert {
    condition     = length(data.snowflake_image_repositories.assert_with_check_block.image_repositories) == 1
    error_message = "image repositories filtered by '${data.snowflake_image_repositories.assert_with_check_block.like}' returned ${length(data.snowflake_image_repositories.assert_with_check_block.image_repositories)} image repositories where one was expected"
  }
}
