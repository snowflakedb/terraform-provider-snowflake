# Provider configuration for import directory
# This file is kept separate from generated main.tf

terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
      version = ">= 2.0.0"
    }
  }
}

provider "snowflake" {
  # Uses default configuration from ~/.snowflake/config or environment variables
}

