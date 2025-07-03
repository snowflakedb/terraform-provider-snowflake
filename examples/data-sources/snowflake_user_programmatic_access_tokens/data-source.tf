# Simple usage
data "snowflake_user_programmatic_access_tokens" "simple" {
  for_user = "<user_name>"
}

output "simple_output" {
  value = data.snowflake_user_programmatic_access_tokens.simple.user_programmatic_access_tokens
}

# Ensure the number of user programmatic access tokens is equal to at least one element (with the use of postcondition)
data "snowflake_user_programmatic_access_tokens" "assert_with_postcondition" {
  for_user = "<user_name>"
  lifecycle {
    postcondition {
      condition     = length(self.user_programmatic_access_tokens) > 0
      error_message = "there should be at least one user programmatic access token"
    }
  }
}

# Ensure the number of user programmatic access tokens is equal to exactly one element (with the use of check block)
check "user_programmatic_access_token_check" {
  data "snowflake_user_programmatic_access_tokens" "assert_with_check_block" {
    for_user = "<user_name>"
  }

  assert {
    condition     = length(data.snowflake_user_programmatic_access_tokens.assert_with_check_block.user_programmatic_access_tokens) == 1
    error_message = "user programmatic access tokens filtered by '${data.snowflake_user_programmatic_access_tokens.assert_with_check_block.for_user}' returned ${length(data.snowflake_user_programmatic_access_tokens.assert_with_check_block.user_programmatic_access_tokens)} user programmatic access tokens where one was expected"
  }
}
