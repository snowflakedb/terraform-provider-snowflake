# Simple usage
data "snowflake_quotas" "simple" {
}

output "simple_output" {
  value = data.snowflake_quotas.simple.quotas
}

# Filtering (like)
data "snowflake_quotas" "like" {
  like = "quota-name"
}

output "like_output" {
  value = data.snowflake_quotas.like.quotas
}

# Filtering by prefix (like)
data "snowflake_quotas" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_quotas.like_prefix.quotas
}

# Filtering (in)
data "snowflake_quotas" "in" {
  in {
    schema = "\"database_name\".\"schema_name\""
  }
}

output "in_output" {
  value = data.snowflake_quotas.in.quotas
}

# Filtering (limit)
data "snowflake_quotas" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_quotas.limit.quotas
}

# Without additional data (to limit the number of calls made for every found quota)
data "snowflake_quotas" "only_show" {
  # with_describe is turned on by default and it reads the configuration of every quota found with its GET_* methods and attaches it to quotas.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_quotas.only_show.quotas
}

# Ensure the number of quotas is equal to at least one element (with the use of postcondition)
data "snowflake_quotas" "assert_with_postcondition" {
  like = "quota-name%"
  lifecycle {
    postcondition {
      condition     = length(self.quotas) > 0
      error_message = "there should be at least one quota"
    }
  }
}

# Ensure the number of quotas is equal to exactly one element (with the use of check block)
check "quota_check" {
  data "snowflake_quotas" "assert_with_check_block" {
    like = "quota-name"
  }

  assert {
    condition     = length(data.snowflake_quotas.assert_with_check_block.quotas) == 1
    error_message = "quotas filtered by '${data.snowflake_quotas.assert_with_check_block.like}' returned ${length(data.snowflake_quotas.assert_with_check_block.quotas)} quotas where one was expected"
  }
}
