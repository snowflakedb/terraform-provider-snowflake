## Minimal
# CREATE requires a schedule, an expiration, or both.
resource "snowflake_backup_policy" "basic" {
  database          = "database_name"
  schema            = "schema_name"
  name              = "backup_policy_name"
  expire_after_days = 7
}

## Complete
# with_retention_lock cannot be changed after creation.
resource "snowflake_backup_policy" "complete" {
  database            = "database_name"
  schema              = "schema_name"
  name                = "backup_policy_name"
  schedule            = "2 HOUR"
  expire_after_days   = 7
  with_retention_lock = "true"
  comment             = "My backup policy"
}
