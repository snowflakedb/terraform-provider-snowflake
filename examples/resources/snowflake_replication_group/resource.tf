## Minimal
resource "snowflake_replication_group" "basic" {
  name = "replication_group_name"
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
resource "snowflake_replication_group" "complete" {
  name = "replication_group_name"
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
  ignore_edition_check = "true"
  replication_schedule = "10 MINUTE"
  error_integration    = snowflake_notification_integration.example.fully_qualified_name
  comment              = "My replication group"
}
