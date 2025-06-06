# Simple usage
data "snowflake_compute_pools" "simple" {
}

output "simple_output" {
  value = data.snowflake_compute_pools.simple.compute_pools
}

# Filtering (like)
data "snowflake_compute_pools" "like" {
  like = "compute-pool-name"
}

output "like_output" {
  value = data.snowflake_compute_pools.like.compute_pools
}

# Filtering by prefix (like)
data "snowflake_compute_pools" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_compute_pools.like_prefix.compute_pools
}

# Filtering (starts_with)
data "snowflake_compute_pools" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_compute_pools.starts_with.compute_pools
}

# Filtering (limit)
data "snowflake_compute_pools" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_compute_pools.limit.compute_pools
}

# Without additional data (to limit the number of calls make for every found compute pool)
data "snowflake_compute_pools" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE COMPUTE POOL for every compute pool found and attaches its output to compute_pools.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_compute_pools.only_show.compute_pools
}

# Ensure the number of compute pools is equal to at least one element (with the use of postcondition)
data "snowflake_compute_pools" "assert_with_postcondition" {
  like = "compute-pool-name%"
  lifecycle {
    postcondition {
      condition     = length(self.compute_pools) > 0
      error_message = "there should be at least one compute pool"
    }
  }
}

# Ensure the number of compute pools is equal to exactly one element (with the use of check block)
check "compute_pool_check" {
  data "snowflake_compute_pools" "assert_with_check_block" {
    like = "compute-pool-name"
  }

  assert {
    condition     = length(data.snowflake_compute_pools.assert_with_check_block.compute_pools) == 1
    error_message = "compute pools filtered by '${data.snowflake_compute_pools.assert_with_check_block.like}' returned ${length(data.snowflake_compute_pools.assert_with_check_block.compute_pools)} compute pools where one was expected"
  }
}
