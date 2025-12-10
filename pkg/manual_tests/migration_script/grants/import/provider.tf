# Provider configuration for import directory
# This file is kept separate from generated main.tf

terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {
  # Uses default configuration from ~/.snowflake/config or environment variables
}

