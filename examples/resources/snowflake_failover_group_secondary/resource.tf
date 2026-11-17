resource "snowflake_failover_group_secondary" "basic" {
  name    = "failover_group_name"
  primary = "\"<organization_name>\".\"<account_name>\".\"<failover_group_name>\""
}
