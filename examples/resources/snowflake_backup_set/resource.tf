## Minimal
# Exactly one of for_table, for_dynamic_table, for_schema, or for_database is required, and it cannot be changed.
resource "snowflake_backup_set" "basic" {
  database  = "database_name"
  schema    = "schema_name"
  name      = "backup_set_name"
  for_table = "\"database_name\".\"schema_name\".\"table_name\""
}

## Complete
# An attached backup policy cannot be removed. It can only be changed to another policy.
resource "snowflake_backup_set" "complete" {
  database      = "database_name"
  schema        = "schema_name"
  name          = "backup_set_name_complete"
  for_schema    = "\"database_name\".\"schema_name\""
  backup_policy = "\"database_name\".\"schema_name\".\"backup_policy_name\""
  comment       = "My backup set"
}
