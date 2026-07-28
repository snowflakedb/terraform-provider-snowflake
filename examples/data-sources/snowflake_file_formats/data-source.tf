# Simple usage
data "snowflake_file_formats" "simple" {
}

output "simple_output" {
  value = data.snowflake_file_formats.simple.file_formats
}

# Filtering (like)
data "snowflake_file_formats" "like" {
  like = "file-format-name"
}

output "like_output" {
  value = data.snowflake_file_formats.like.file_formats
}

# Filtering (in)
data "snowflake_file_formats" "in" {
  in {
    database = "database"
  }
}

output "in_output" {
  value = data.snowflake_file_formats.in.file_formats
}

# Without additional data (to limit the number of calls make for every found file format)
data "snowflake_file_formats" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE FILE FORMAT for every file format found and attaches its output to file_formats.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_file_formats.only_show.file_formats
}

# Ensure the number of file formats is equal to at least one element (with the use of postcondition)
data "snowflake_file_formats" "assert_with_postcondition" {
  like = "file-format-name%"
  lifecycle {
    postcondition {
      condition     = length(self.file_formats) > 0
      error_message = "there should be at least one file format"
    }
  }
}

# Ensure the number of file formats is equal to exactly one element (with the use of check block)
check "file_format_check" {
  data "snowflake_file_formats" "assert_with_check_block" {
    like = "file-format-name"
  }

  assert {
    condition     = length(data.snowflake_file_formats.assert_with_check_block.file_formats) == 1
    error_message = "file formats filtered by '${data.snowflake_file_formats.assert_with_check_block.like}' returned ${length(data.snowflake_file_formats.assert_with_check_block.file_formats)} file formats where one was expected"
  }
}
