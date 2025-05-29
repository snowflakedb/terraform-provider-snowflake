# Simple usage
data "snowflake_accounts" "simple" {
}

output "simple_output" {
  value = data.snowflake_accounts.simple.accounts
}

# Filtering (like)
data "snowflake_accounts" "like" {
  like = "account-name"
}

output "like_output" {
  value = data.snowflake_accounts.like.accounts
}

# With history
data "snowflake_accounts" "with_history" {
  with_history = true
}

output "with_history_output" {
  value = data.snowflake_accounts.like.accounts
}

# Ensure the number of accounts is equal to at least one element (with the use of postcondition)
data "snowflake_accounts" "assert_with_postcondition" {
  like = "account-name"
  lifecycle {
    postcondition {
      condition     = length(self.accounts) > 0
      error_message = "there should be at least one account"
    }
  }
}

# Ensure the number of accounts is equal to exactly one element (with the use of check block)
check "account_check" {
  data "snowflake_accounts" "assert_with_check_block" {
    like = "account-name"
  }

  assert {
    condition     = length(data.snowflake_accounts.assert_with_check_block.accounts) == 1
    error_message = "accounts filtered by '${data.snowflake_accounts.assert_with_check_block.like}' returned ${length(data.snowflake_accounts.assert_with_check_block.accounts)} accounts where one was expected"
  }
}
