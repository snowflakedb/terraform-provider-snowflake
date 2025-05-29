# Simple usage
data "snowflake_warehouses" "simple" {
}

output "simple_output" {
  value = data.snowflake_warehouses.simple.warehouses
}

# Filtering (like)
data "snowflake_warehouses" "like" {
  like = "warehouse-name"
}

output "like_output" {
  value = data.snowflake_warehouses.like.warehouses
}

# Filtering by prefix (like)
data "snowflake_warehouses" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_warehouses.like_prefix.warehouses
}

# Without additional data (to limit the number of calls make for every found warehouse)
data "snowflake_warehouses" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE WAREHOUSE for every warehouse found and attaches its output to warehouses.*.describe_output field
  with_describe = false

  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR WAREHOUSE for every warehouse found and attaches its output to warehouses.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_warehouses.only_show.warehouses
}

# Ensure the number of warehouses is equal to at least one element (with the use of postcondition)
data "snowflake_warehouses" "assert_with_postcondition" {
  like = "warehouse-name%"
  lifecycle {
    postcondition {
      condition     = length(self.warehouses) > 0
      error_message = "there should be at least one warehouse"
    }
  }
}

# Ensure the number of warehouses is equal to exactly one element (with the use of check block)
check "warehouse_check" {
  data "snowflake_warehouses" "assert_with_check_block" {
    like = "warehouse-name"
  }

  assert {
    condition     = length(data.snowflake_warehouses.assert_with_check_block.warehouses) == 1
    error_message = "warehouses filtered by '${data.snowflake_warehouses.assert_with_check_block.like}' returned ${length(data.snowflake_warehouses.assert_with_check_block.warehouses)} warehouses where one was expected"
  }
}
