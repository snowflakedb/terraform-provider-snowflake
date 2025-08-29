# Prerequisites:
# - Ensure both accounts, where primary and secondary connections will be created, have:
#   - Correct edition: BUSINESS_CRITICAL
#   - Replication enabled: https://docs.snowflake.com/alias/replication/enable-account-repl
# For more information on how to set up the accounts, refer to the internal documentation
# Commands to run
# - terraform apply

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "2.5.0"
    }
  }
}

provider "snowflake" {}

provider "snowflake" {
  profile = "secondary_test_account"
  alias   = "second_account"
}

resource "snowflake_primary_connection" "primary_connection" {
  name = "TEST_CONNECTION"
  enable_failover_to_accounts = ["<organization_name>.<secondary_connection_account_name>"]
}

resource "snowflake_secondary_connection" "secondary_connection" {
  provider      = snowflake.second_account
  name          = snowflake_primary_connection.primary_connection.name
  as_replica_of = "<organization_name>.<primary_connection_account_name>.${snowflake_primary_connection.primary_connection.name}"
}