# Simple usage
data "snowflake_backup_sets" "simple" {
}

output "simple_output" {
  value = data.snowflake_backup_sets.simple.backup_sets
}

# Filtering (like)
data "snowflake_backup_sets" "like" {
  like = "backup-set-name"
}

output "like_output" {
  value = data.snowflake_backup_sets.like.backup_sets
}

# Filtering by prefix (like)
data "snowflake_backup_sets" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_backup_sets.like_prefix.backup_sets
}

# Filtering (starts_with)
data "snowflake_backup_sets" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_backup_sets.starts_with.backup_sets
}

# Filtering (limit)
data "snowflake_backup_sets" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_backup_sets.limit.backup_sets
}

# Filtering (in)
data "snowflake_backup_sets" "in_account" {
  in {
    account = true
  }
}

data "snowflake_backup_sets" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_backup_sets" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_backup_sets.in_account.backup_sets,
    "database" : data.snowflake_backup_sets.in_database.backup_sets,
    "schema" : data.snowflake_backup_sets.in_schema.backup_sets,
  }
}

# Without additional data (to limit the number of calls made for every found backup set)
data "snowflake_backup_sets" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE BACKUP SET for every backup set found and attaches its output to backup_sets.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_backup_sets.only_show.backup_sets
}

# Ensure the number of backup sets is equal to at least one element (with the use of postcondition)
data "snowflake_backup_sets" "assert_with_postcondition" {
  like = "backup-set-name%"
  lifecycle {
    postcondition {
      condition     = length(self.backup_sets) > 0
      error_message = "there should be at least one backup set"
    }
  }
}

# Ensure the number of backup sets is equal to exactly one element (with the use of check block)
check "backup_set_check" {
  data "snowflake_backup_sets" "assert_with_check_block" {
    like = "backup-set-name"
  }

  assert {
    condition     = length(data.snowflake_backup_sets.assert_with_check_block.backup_sets) == 1
    error_message = "backup sets filtered by '${data.snowflake_backup_sets.assert_with_check_block.like}' returned ${length(data.snowflake_backup_sets.assert_with_check_block.backup_sets)} backup sets where one was expected"
  }
}
