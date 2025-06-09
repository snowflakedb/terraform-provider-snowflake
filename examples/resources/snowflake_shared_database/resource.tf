# 1. Preparing database to share
resource "snowflake_share" "test" {
  provider = primary_account # notice the provider fields
  name     = "share_name"
  accounts = ["<secondary_account_organization_name>.<secondary_account_name>"]
}

resource "snowflake_database" "test" {
  provider = primary_account
  name     = "shared_database"
}

resource "snowflake_grant_privileges_to_share" "test" {
  provider    = primary_account
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

# 2. Creating shared database
## 2.1. Minimal version
resource "snowflake_shared_database" "test" {
  provider   = secondary_account
  depends_on = [snowflake_grant_privileges_to_share.test]
  name       = snowflake_database.test.name # shared database should have the same as the "imported" one
  from_share = "<primary_account_organization_name>.<primary_account_name>.${snowflake_share.test.name}"
}

## 2.2. Complete version (with every optional set)
resource "snowflake_shared_database" "test" {
  provider   = secondary_account
  depends_on = [snowflake_grant_privileges_to_share.test]
  name       = snowflake_database.test.name # shared database should have the same as the "imported" one
  from_share = "<primary_account_organization_name>.<primary_account_name>.${snowflake_share.test.name}"
  comment    = "A shared database"

  external_volume                               = "<external_volume_name>"
  catalog                                       = "<catalog_name>"
  replace_invalid_characters                    = false
  default_ddl_collation                         = "en_US"
  storage_serialization_policy                  = "COMPATIBLE"
  log_level                                     = "INFO"
  trace_level                                   = "ALWAYS"
  suspend_task_after_num_failures               = 10
  task_auto_retry_attempts                      = 10
  user_task_managed_initial_warehouse_size      = "LARGE"
  user_task_timeout_ms                          = 3600000
  user_task_minimum_trigger_interval_in_seconds = 120
  quoted_identifiers_ignore_case                = false
  enable_console_output                         = false
}
