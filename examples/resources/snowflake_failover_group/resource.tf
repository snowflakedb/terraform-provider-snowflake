## Minimal
resource "snowflake_failover_group" "basic" {
  name = "failover_group_name"
  object_types = [
    "DATABASES",
  ]
  allowed_accounts = [
    "\"<organization_name>\".\"<account_name>\"",
  ]
  allowed_databases = [
    snowflake_database.example.fully_qualified_name,
  ]
}

## Complete
resource "snowflake_failover_group" "complete" {
  name = "failover_group_name"
  object_types = [
    "DATABASES",
    "ROLES",
    "EXTERNAL VOLUMES",
    "SHARES",
    "INTEGRATIONS",
  ]
  allowed_accounts = [
    "\"<organization_name>\".\"<account_name>\"",
  ]
  allowed_databases = [
    snowflake_database.example.fully_qualified_name,
  ]
  allowed_external_volumes = [
    snowflake_external_volume.example.fully_qualified_name,
  ]
  allowed_shares = [
    snowflake_share.example.fully_qualified_name,
  ]
  allowed_integration_types = [
    "API INTEGRATIONS",
    "EXTERNAL ACCESS INTEGRATIONS",
    "NOTIFICATION INTEGRATIONS",
    "SECURITY INTEGRATIONS",
    "STORAGE INTEGRATIONS",
  ]
  ignore_edition_check = true
  # an interval in minutes ("10 MINUTE") or a cron expression with the time zone ("USING CRON 0 0 10-20 * TUE,THU UTC")
  replication_schedule = "10 MINUTE"
  # needs a replication schedule with an interval of 6 hours or less
  optimized_refresh = true
  error_integration = snowflake_notification_integration.example.fully_qualified_name
  comment           = "My failover group"
}
