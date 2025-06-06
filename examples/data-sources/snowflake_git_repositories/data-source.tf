# Simple usage
data "snowflake_git_repositories" "simple" {
}

output "simple_output" {
  value = data.snowflake_git_repositories.simple.git_repositories
}

# Filtering (like)
data "snowflake_git_repositories" "like" {
  like = "git-repository-name"
}

output "like_output" {
  value = data.snowflake_git_repositories.like.git_repositories
}

# Filtering by prefix (like)
data "snowflake_git_repositories" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_git_repositories.like_prefix.git_repositories
}

# Filtering (in)
data "snowflake_git_repositories" "in_account" {
  in {
    account = true
  }
}

data "snowflake_git_repositories" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_git_repositories" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_git_repositories.in_account.git_repositories,
    "database" : data.snowflake_git_repositories.in_database.git_repositories,
    "schema" : data.snowflake_git_repositories.in_schema.git_repositories,
  }
}

# Filtering (limit)
data "snowflake_git_repositories" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_git_repositories.limit.git_repositories
}

# Without additional data (to limit the number of calls make for every found git repository)
data "snowflake_git_repositories" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE GIT REPOSITORY for every git repository found and attaches its output to git_repositories.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_git_repositories.only_show.git_repositories
}

# Ensure the number of git repositories is equal to at least one element (with the use of postcondition)
data "snowflake_git_repositories" "assert_with_postcondition" {
  like = "git-repository-name%"
  lifecycle {
    postcondition {
      condition     = length(self.git_repositories) > 0
      error_message = "there should be at least one git repository"
    }
  }
}

# Ensure the number of git repositories is equal to exactly one element (with the use of check block)
check "git_repository_check" {
  data "snowflake_git_repositories" "assert_with_check_block" {
    like = "git-repository-name"
  }

  assert {
    condition     = length(data.snowflake_git_repositories.assert_with_check_block.git_repositories) == 1
    error_message = "git repositories filtered by '${data.snowflake_git_repositories.assert_with_check_block.like}' returned ${length(data.snowflake_git_repositories.assert_with_check_block.git_repositories)} git repositories where one was expected"
  }
}
