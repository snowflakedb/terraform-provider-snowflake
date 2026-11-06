# Simple usage
data "snowflake_replication_groups" "simple" {
}

output "simple_output" {
  value = data.snowflake_replication_groups.simple.replication_groups
}

# Filtering (in_account)
data "snowflake_replication_groups" "in_account" {
  in_account = "\"<organization_name>\".\"<account_name>\""
}

output "in_account_output" {
  value = data.snowflake_replication_groups.in_account.replication_groups
}
