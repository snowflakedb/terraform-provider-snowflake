# Simple usage
data "snowflake_code_bundles" "simple" {
}

output "simple_output" {
  value = data.snowflake_code_bundles.simple.code_bundles
}

# Filtering (in)
data "snowflake_code_bundles" "in_account" {
  in {
    account = true
  }
}

data "snowflake_code_bundles" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_code_bundles" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

data "snowflake_code_bundles" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_code_bundles" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_code_bundles.in_account.code_bundles,
    "database" : data.snowflake_code_bundles.in_database.code_bundles,
    "schema" : data.snowflake_code_bundles.in_schema.code_bundles,
    "application" : data.snowflake_code_bundles.in_application.code_bundles,
    "application_package" : data.snowflake_code_bundles.in_application_package.code_bundles,
  }
}

# Without additional data (to limit the number of calls made for every found code bundle)
data "snowflake_code_bundles" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE CODE BUNDLE for every code bundle found and attaches its output to code_bundles.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_code_bundles.only_show.code_bundles
}

# Ensure the number of code bundles is equal to at least one element (with the use of postcondition)
data "snowflake_code_bundles" "assert_with_postcondition" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
  lifecycle {
    postcondition {
      condition     = length(self.code_bundles) > 0
      error_message = "there should be at least one code bundle"
    }
  }
}

# Ensure the number of code bundles is equal to exactly one element (with the use of check block)
check "code_bundle_check" {
  data "snowflake_code_bundles" "assert_with_check_block" {
    in {
      schema = "\"<database_name>\".\"<schema_name>\""
    }
  }

  assert {
    condition     = length(data.snowflake_code_bundles.assert_with_check_block.code_bundles) == 1
    error_message = "code bundles in '${data.snowflake_code_bundles.assert_with_check_block.in.0.schema}' returned ${length(data.snowflake_code_bundles.assert_with_check_block.code_bundles)} code bundles where one was expected"
  }
}
