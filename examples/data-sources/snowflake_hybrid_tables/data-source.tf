# Simple usage
data "snowflake_hybrid_tables" "simple" {
}

output "simple_output" {
  value = data.snowflake_hybrid_tables.simple.hybrid_tables
}

# Filtering (like)
data "snowflake_hybrid_tables" "like" {
  like = "hybrid-table-name"
}

output "like_output" {
  value = data.snowflake_hybrid_tables.like.hybrid_tables
}

# Filtering by prefix (like)
data "snowflake_hybrid_tables" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_hybrid_tables.like_prefix.hybrid_tables
}

# Filtering (in)
data "snowflake_hybrid_tables" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_hybrid_tables" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "database" : data.snowflake_hybrid_tables.in_database.hybrid_tables,
    "schema" : data.snowflake_hybrid_tables.in_schema.hybrid_tables,
  }
}

# Filtering (starts_with)
data "snowflake_hybrid_tables" "starts_with" {
  starts_with = "prefix"
}

output "starts_with_output" {
  value = data.snowflake_hybrid_tables.starts_with.hybrid_tables
}

# Filtering (limit)
data "snowflake_hybrid_tables" "limit" {
  limit {
    rows = 10
    from = "hybrid-table-"
  }
}

output "limit_output" {
  value = data.snowflake_hybrid_tables.limit.hybrid_tables
}

# Without additional data (to limit the number of calls made for every found hybrid table)
data "snowflake_hybrid_tables" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE TABLE for every hybrid table found and attaches its output to hybrid_tables.*.describe_output field
  with_describe = false
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR TABLE for every hybrid table found and attaches its output to hybrid_tables.*.parameters field
  with_parameters = false
  # with_keys is turned off by default; set to true to call SHOW PRIMARY KEYS / UNIQUE KEYS / IMPORTED KEYS for every hybrid table found and attach the merged constraints (grouped by constraint name, ordered by kind then column names) to hybrid_tables.*.show_keys_output field
  with_keys = false
  # with_indexes is turned off by default; set to true to call SHOW INDEXES for every hybrid table found and attach its output to hybrid_tables.*.show_indexes field
  with_indexes = false
}

output "only_show_output" {
  value = data.snowflake_hybrid_tables.only_show.hybrid_tables
}

# Ensure the number of hybrid tables is equal to at least one element (with the use of postcondition)
data "snowflake_hybrid_tables" "assert_with_postcondition" {
  like = "hybrid-table-name%"
  lifecycle {
    postcondition {
      condition     = length(self.hybrid_tables) > 0
      error_message = "there should be at least one hybrid table"
    }
  }
}

# Ensure the number of hybrid tables is equal to exactly one element (with the use of check block)
check "hybrid_table_check" {
  data "snowflake_hybrid_tables" "assert_with_check_block" {
    like = "hybrid-table-name"
  }

  assert {
    condition     = length(data.snowflake_hybrid_tables.assert_with_check_block.hybrid_tables) == 1
    error_message = "hybrid tables filtered by '${data.snowflake_hybrid_tables.assert_with_check_block.like}' returned ${length(data.snowflake_hybrid_tables.assert_with_check_block.hybrid_tables)} hybrid tables where one was expected"
  }
}
