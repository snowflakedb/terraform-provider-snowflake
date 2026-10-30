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
  expire_after_days   = 7
  with_retention_lock = "true"
  comment             = "My backup policy"

  # Exactly one of hours, minutes (at least 60), or using_cron.
  schedule {
    hours = 2
  }
}

## Complete with a cron schedule
resource "snowflake_backup_policy" "cron" {
  database          = "database_name"
  schema            = "schema_name"
  name              = "backup_policy_name_cron"
  expire_after_days = 7

  schedule {
    using_cron = "0 0 * * * UTC"
  }
}
