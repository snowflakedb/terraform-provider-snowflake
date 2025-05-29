# Simple usage
data "snowflake_streamlits" "simple" {
}

output "simple_output" {
  value = data.snowflake_streamlits.simple.streamlits
}

# Filtering (like)
data "snowflake_streamlits" "like" {
  like = "streamlit-name"
}

output "like_output" {
  value = data.snowflake_streamlits.like.streamlits
}

# Filtering by prefix (like)
data "snowflake_streamlits" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_streamlits.like_prefix.streamlits
}

# Filtering (limit)
data "snowflake_streamlits" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_streamlits.limit.streamlits
}

# Filtering (in)
data "snowflake_streamlits" "in" {
  in {
    database = "database"
  }
}

output "in_output" {
  value = data.snowflake_streamlits.in.streamlits
}

# Without additional data (to limit the number of calls make for every found streamlit)
data "snowflake_streamlits" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE STREAMLIT for every streamlit found and attaches its output to streamlits.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_streamlits.only_show.streamlits
}

# Ensure the number of streamlits is equal to at least one element (with the use of postcondition)
data "snowflake_streamlits" "assert_with_postcondition" {
  like = "streamlit-name%"
  lifecycle {
    postcondition {
      condition     = length(self.streamlits) > 0
      error_message = "there should be at least one streamlit"
    }
  }
}

# Ensure the number of streamlits is equal to exactly one element (with the use of check block)
check "streamlit_check" {
  data "snowflake_streamlits" "assert_with_check_block" {
    like = "streamlit-name"
  }

  assert {
    condition     = length(data.snowflake_streamlits.assert_with_check_block.streamlits) == 1
    error_message = "streamlits filtered by '${data.snowflake_streamlits.assert_with_check_block.like}' returned ${length(data.snowflake_streamlits.assert_with_check_block.streamlits)} streamlits where one was expected"
  }
}
