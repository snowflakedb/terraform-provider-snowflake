# Simple usage
data "snowflake_backup_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_backup_policies.simple.backup_policies
}

# Filtering (like)
data "snowflake_backup_policies" "like" {
  like = "backup-policy-name"
}

output "like_output" {
  value = data.snowflake_backup_policies.like.backup_policies
}

# Filtering by prefix (like)
data "snowflake_backup_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_backup_policies.like_prefix.backup_policies
}

# Filtering (starts_with)
data "snowflake_backup_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_backup_policies.starts_with.backup_policies
}

# Filtering (limit)
data "snowflake_backup_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_backup_policies.limit.backup_policies
}

# Filtering (in)
data "snowflake_backup_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_backup_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_backup_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_backup_policies.in_account.backup_policies,
    "database" : data.snowflake_backup_policies.in_database.backup_policies,
    "schema" : data.snowflake_backup_policies.in_schema.backup_policies,
  }
}

# Without additional data (to limit the number of calls made for every found backup policy)
data "snowflake_backup_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE BACKUP POLICY for every backup policy found and attaches its output to backup_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_backup_policies.only_show.backup_policies
}

# Ensure the number of backup policies is equal to at least one element (with the use of postcondition)
data "snowflake_backup_policies" "assert_with_postcondition" {
  like = "backup-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.backup_policies) > 0
      error_message = "there should be at least one backup policy"
    }
  }
}

# Ensure the number of backup policies is equal to exactly one element (with the use of check block)
check "backup_policy_check" {
  data "snowflake_backup_policies" "assert_with_check_block" {
    like = "backup-policy-name"
  }

  assert {
    condition     = length(data.snowflake_backup_policies.assert_with_check_block.backup_policies) == 1
    error_message = "backup policies filtered by '${data.snowflake_backup_policies.assert_with_check_block.like}' returned ${length(data.snowflake_backup_policies.assert_with_check_block.backup_policies)} backup policies where one was expected"
  }
}
