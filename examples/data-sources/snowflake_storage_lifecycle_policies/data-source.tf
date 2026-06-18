# Simple usage
data "snowflake_storage_lifecycle_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_storage_lifecycle_policies.simple.storage_lifecycle_policies
}

# Filtering (like)
data "snowflake_storage_lifecycle_policies" "like" {
  like = "storage-lifecycle-policy-name"
}

output "like_output" {
  value = data.snowflake_storage_lifecycle_policies.like.storage_lifecycle_policies
}

# Filtering by prefix (like)
data "snowflake_storage_lifecycle_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_storage_lifecycle_policies.like_prefix.storage_lifecycle_policies
}

# Filtering (in)
data "snowflake_storage_lifecycle_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_storage_lifecycle_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_storage_lifecycle_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_storage_lifecycle_policies.in_account.storage_lifecycle_policies,
    "database" : data.snowflake_storage_lifecycle_policies.in_database.storage_lifecycle_policies,
    "schema" : data.snowflake_storage_lifecycle_policies.in_schema.storage_lifecycle_policies,
  }
}

# Without additional data (to limit the number of calls made for every found storage lifecycle policy)
data "snowflake_storage_lifecycle_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE STORAGE LIFECYCLE POLICY for every storage lifecycle policy found and attaches its output to storage_lifecycle_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_storage_lifecycle_policies.only_show.storage_lifecycle_policies
}

# Ensure the number of storage lifecycle policies is equal to at least one element (with the use of postcondition)
data "snowflake_storage_lifecycle_policies" "assert_with_postcondition" {
  like = "storage-lifecycle-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.storage_lifecycle_policies) > 0
      error_message = "there should be at least one storage lifecycle policy"
    }
  }
}

# Ensure the number of storage lifecycle policies is equal to exactly one element (with the use of check block)
check "storage_lifecycle_policy_check" {
  data "snowflake_storage_lifecycle_policies" "assert_with_check_block" {
    like = "storage-lifecycle-policy-name"
  }

  assert {
    condition     = length(data.snowflake_storage_lifecycle_policies.assert_with_check_block.storage_lifecycle_policies) == 1
    error_message = "storage lifecycle policies filtered by '${data.snowflake_storage_lifecycle_policies.assert_with_check_block.like}' returned ${length(data.snowflake_storage_lifecycle_policies.assert_with_check_block.storage_lifecycle_policies)} storage lifecycle policies where one was expected"
  }
}
