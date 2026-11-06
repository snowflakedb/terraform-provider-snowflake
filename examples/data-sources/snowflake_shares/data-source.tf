# Simple usage
data "snowflake_shares" "simple" {
}

output "simple_output" {
  value = data.snowflake_shares.simple.shares
}

# Filtering (like)
data "snowflake_shares" "like" {
  like = "share-name"
}

output "like_output" {
  value = data.snowflake_shares.like.shares
}

# Filtering (starts_with)
data "snowflake_shares" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_shares.starts_with.shares
}

# Filtering (limit)
data "snowflake_shares" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_shares.limit.shares
}

# Without additional data (to limit the number of calls make for every found share)
data "snowflake_shares" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SHARE for every share found and attaches its output to shares.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_shares.only_show.shares
}

# Ensure the number of shares is equal to at least one element (with the use of postcondition)
data "snowflake_shares" "assert_with_postcondition" {
  starts_with = "prefix-"
  lifecycle {
    postcondition {
      condition     = length(self.shares) > 0
      error_message = "there should be at least one share"
    }
  }
}

# Ensure the number of shares is equal to exactly one element (with the use of check block)
check "share_check" {
  data "snowflake_shares" "assert_with_check_block" {
    like = "share-name"
  }

  assert {
    condition     = length(data.snowflake_shares.assert_with_check_block.shares) == 1
    error_message = "shares filtered by '${data.snowflake_shares.assert_with_check_block.like}' returned ${length(data.snowflake_shares.assert_with_check_block.shares)} shares where one was expected"
  }
}
