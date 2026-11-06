resource "snowflake_replication_group_secondary" "basic" {
  name    = "replication_group_name"
  primary = "\"<organization_name>\".\"<account_name>\".\"<replication_group_name>\""
}
