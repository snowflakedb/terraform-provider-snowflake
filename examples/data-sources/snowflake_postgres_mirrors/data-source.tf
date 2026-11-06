# Simple usage
data "snowflake_postgres_mirrors" "simple" {
}

output "simple_output" {
  value = data.snowflake_postgres_mirrors.simple.postgres_mirrors
}

# Filtering by Postgres instance
data "snowflake_postgres_mirrors" "by_instance" {
  postgres_instance = "my_instance"
}

output "by_instance_output" {
  value = data.snowflake_postgres_mirrors.by_instance.postgres_mirrors
}

# Without the additional DESCRIBE_MIRROR call for each mirror
data "snowflake_postgres_mirrors" "no_describe" {
  with_describe = false
}

output "no_describe_output" {
  value = data.snowflake_postgres_mirrors.no_describe.postgres_mirrors
}

# Ensure the number of mirrors is equal to at least one element (with the use of postcondition)
data "snowflake_postgres_mirrors" "assert_with_postcondition" {
  postgres_instance = "my_instance"
  lifecycle {
    postcondition {
      condition     = length(self.postgres_mirrors) > 0
      error_message = "there should be at least one mirror"
    }
  }
}

# Ensure the number of mirrors is equal to exactly one element (with the use of check block)
check "mirror_check" {
  data "snowflake_postgres_mirrors" "assert_with_check_block" {
    postgres_instance = "my_instance"
  }

  assert {
    condition     = length(data.snowflake_postgres_mirrors.assert_with_check_block.postgres_mirrors) == 1
    error_message = "mirrors filtered by '${data.snowflake_postgres_mirrors.assert_with_check_block.postgres_instance}' returned ${length(data.snowflake_postgres_mirrors.assert_with_check_block.postgres_mirrors)} mirrors where one was expected"
  }
}
