# Simple usage
data "snowflake_iceberg_tables" "simple" {
}

output "simple_output" {
  value = data.snowflake_iceberg_tables.simple.iceberg_tables
}

# Filtering (like)
data "snowflake_iceberg_tables" "like" {
  like = "iceberg-table-name"
}

output "like_output" {
  value = data.snowflake_iceberg_tables.like.iceberg_tables
}

# Filtering by prefix (like)
data "snowflake_iceberg_tables" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_iceberg_tables.like_prefix.iceberg_tables
}

# Filtering (in)
data "snowflake_iceberg_tables" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_iceberg_tables" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "database" : data.snowflake_iceberg_tables.in_database.iceberg_tables,
    "schema" : data.snowflake_iceberg_tables.in_schema.iceberg_tables,
  }
}

# Filtering (starts_with)
data "snowflake_iceberg_tables" "starts_with" {
  starts_with = "prefix"
}

output "starts_with_output" {
  value = data.snowflake_iceberg_tables.starts_with.iceberg_tables
}

# Filtering (limit)
data "snowflake_iceberg_tables" "limit" {
  limit {
    rows = 10
  }
}

output "limit_output" {
  value = data.snowflake_iceberg_tables.limit.iceberg_tables
}

# Without additional data (to limit the number of calls made for every found iceberg table)
data "snowflake_iceberg_tables" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE ICEBERG TABLE for every iceberg table found and attaches its output to iceberg_tables.*.describe_output field
  with_describe = false
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR ICEBERG TABLE for every iceberg table found and attaches its output to iceberg_tables.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_iceberg_tables.only_show.iceberg_tables
}

# Ensure the number of iceberg tables is equal to at least one element (with the use of postcondition)
data "snowflake_iceberg_tables" "assert_with_postcondition" {
  like = "iceberg-table-name%"
  lifecycle {
    postcondition {
      condition     = length(self.iceberg_tables) > 0
      error_message = "there should be at least one iceberg table"
    }
  }
}

# Ensure the number of iceberg tables is equal to exactly one element (with the use of check block)
check "iceberg_table_check" {
  data "snowflake_iceberg_tables" "assert_with_check_block" {
    like = "iceberg-table-name"
  }

  assert {
    condition     = length(data.snowflake_iceberg_tables.assert_with_check_block.iceberg_tables) == 1
    error_message = "iceberg tables filtered by '${data.snowflake_iceberg_tables.assert_with_check_block.like}' returned ${length(data.snowflake_iceberg_tables.assert_with_check_block.iceberg_tables)} iceberg tables where one was expected"
  }
}
