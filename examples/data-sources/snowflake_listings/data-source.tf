# Simple usage
data "snowflake_listings" "simple" {
}

output "simple_output" {
  value = data.snowflake_listings.simple.listings
}

# Filtering (like)
data "snowflake_listings" "like" {
  like = "listing-name"
}

output "like_output" {
  value = data.snowflake_listings.like.listings
}

# Filtering by prefix (like)
data "snowflake_listings" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_listings.like_prefix.listings
}

# Filtering (starts_with)
data "snowflake_listings" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_listings.starts_with.listings
}

# Filtering (limit)
data "snowflake_listings" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_listings.limit.listings
}

# Without additional data (to limit the number of calls make for every found listing)
data "snowflake_listings" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE LISTING for every listing found and attaches its output to listings.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_listings.only_show.listings
}

# Ensure the number of listings is equal to at least one element (with the use of postcondition)
data "snowflake_listings" "assert_with_postcondition" {
  like = "listing-name%"
  lifecycle {
    postcondition {
      condition     = length(self.listings) > 0
      error_message = "there should be at least one listing"
    }
  }
}

# Ensure the number of listings is equal to exactly one element (with the use of check block)
check "listing_check" {
  data "snowflake_listings" "assert_with_check_block" {
    like = "listing-name"
  }

  assert {
    condition     = length(data.snowflake_listings.assert_with_check_block.listings) == 1
    error_message = "listings filtered by '${data.snowflake_listings.assert_with_check_block.like}' returned ${length(data.snowflake_listings.assert_with_check_block.listings)} listings where one was expected"
  }
}

