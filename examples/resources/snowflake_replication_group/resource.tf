resource "snowflake_database" "example" {
  name = "database_name"
}

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
  ]
  allowed_accounts = [
    "\"<organization_name>\".\"<account_name>\"",
  ]
  allowed_databases = [
    snowflake_database.example.fully_qualified_name,
  ]
  comment              = "My replication group"
  replication_schedule = "10 MINUTE"
}
