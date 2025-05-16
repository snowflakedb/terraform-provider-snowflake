# Simple usage
data "snowflake_users" "simple" {
}

output "simple_output" {
  value = data.snowflake_users.simple.users
}

# Filtering (like)
data "snowflake_users" "like" {
  like = "user-name"
}

output "like_output" {
  value = data.snowflake_users.like.users
}

# Filtering (starts_with)
data "snowflake_users" "starts_with" {
  starts_with = "user-"
}

output "starts_with_output" {
  value = data.snowflake_users.starts_with.users
}

# Filtering (limit)
data "snowflake_users" "limit" {
  limit {
    rows = 10
    from = "user-"
  }
}

output "limit_output" {
  value = data.snowflake_users.limit.users
}

# Without additional data (to limit the number of calls make for every found user)
data "snowflake_users" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE USER for every user found and attaches its output to users.*.describe_output field
  with_describe = false

  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR USER for every user found and attaches its output to users.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_users.only_show.users
}

# Ensure the number of users is equal to at least one element (with the use of postcondition)
data "snowflake_users" "assert_with_postcondition" {
  starts_with = "user-name"
  lifecycle {
    postcondition {
      condition     = length(self.users) > 0
      error_message = "there should be at least one user"
    }
  }
}

# Ensure the number of users is equal to exactly one element (with the use of check block)
check "user_check" {
  data "snowflake_users" "assert_with_check_block" {
    like = "user-name"
  }

  assert {
    condition     = length(data.snowflake_users.assert_with_check_block.users) == 1
    error_message = "users filtered by '${data.snowflake_users.assert_with_check_block.like}' returned ${length(data.snowflake_users.assert_with_check_block.users)} users where one was expected"
  }
}
