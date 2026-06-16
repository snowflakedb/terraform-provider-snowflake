## Minimal
resource "snowflake_storage_lifecycle_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "storage_lifecycle_policy_name"
  argument {
    name = "VAL"
    type = "VARCHAR"
  }
  body = "LENGTH(VAL) > 0"
}

## Complete
resource "snowflake_storage_lifecycle_policy" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "storage_lifecycle_policy_name"
  argument {
    name = "VAL"
    type = "VARCHAR"
  }
  argument {
    name = "CREATED_AT"
    type = "TIMESTAMP_NTZ"
  }
  body             = "LENGTH(VAL) > 0"
  archive_tier     = "COLD"
  archive_for_days = 365
  comment          = "My storage lifecycle policy"
}
