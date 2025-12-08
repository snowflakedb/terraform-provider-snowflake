# =============================================================================
# Warehouses Test Objects
# =============================================================================
# These resources create test warehouses on Snowflake for migration testing.
# Run with: terraform apply
# Cleanup:  terraform destroy
#
# Edge cases covered:
# - Basic warehouse (minimal config with defaults)
# - Different sizes (XSMALL, SMALL, MEDIUM, LARGE)
# - Comment with text and special characters
# - Auto suspend/resume variations
# - Multi-cluster warehouse (min/max cluster count)
# - Query acceleration enabled
# - Different scaling policies (STANDARD, ECONOMY)
# - Different warehouse types (STANDARD, SNOWPARK-OPTIMIZED)
# - Custom statement timeout parameters
# =============================================================================

terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {}

# ------------------------------------------------------------------------------
# Basic Warehouse (minimal configuration - uses all defaults)
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "basic" {
  name = "MIGRATION_TEST_WH_BASIC"
}

# ------------------------------------------------------------------------------
# Warehouse with Comment
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "with_comment" {
  name           = "MIGRATION_TEST_WH_COMMENT"
  warehouse_size = "XSMALL"
  comment        = "This warehouse is used for migration testing purposes"
}

# ------------------------------------------------------------------------------
# Warehouse with Long Comment
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "long_comment" {
  name           = "MIGRATION_TEST_WH_LONG_COMMENT"
  warehouse_size = "XSMALL"
  comment        = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# ------------------------------------------------------------------------------
# Warehouse with Special Characters in Comment
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "special_chars" {
  name           = "MIGRATION_TEST_WH_SPECIAL"
  warehouse_size = "XSMALL"
  comment        = "Comment with special chars: <>&\""
}

# ------------------------------------------------------------------------------
# Small Warehouse
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "small" {
  name           = "MIGRATION_TEST_WH_SMALL"
  warehouse_size = "SMALL"
  comment        = "Small sized warehouse"
}

# ------------------------------------------------------------------------------
# Medium Warehouse
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "medium" {
  name           = "MIGRATION_TEST_WH_MEDIUM"
  warehouse_size = "MEDIUM"
  comment        = "Medium sized warehouse"
}

# ------------------------------------------------------------------------------
# Warehouse with Auto Suspend Disabled (0 = never suspend)
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "no_auto_suspend" {
  name           = "MIGRATION_TEST_WH_NO_SUSPEND"
  warehouse_size = "XSMALL"
  auto_suspend   = 0
  comment        = "Warehouse that never auto-suspends"
}

# ------------------------------------------------------------------------------
# Warehouse with Custom Auto Suspend (10 minutes)
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "custom_suspend" {
  name           = "MIGRATION_TEST_WH_CUSTOM_SUSPEND"
  warehouse_size = "XSMALL"
  auto_suspend   = 600
  comment        = "Warehouse with 10 minute auto-suspend"
}

# ------------------------------------------------------------------------------
# Warehouse with Auto Resume Disabled
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "no_auto_resume" {
  name           = "MIGRATION_TEST_WH_NO_RESUME"
  warehouse_size = "XSMALL"
  auto_resume    = false
  comment        = "Warehouse without auto-resume"
}

# ------------------------------------------------------------------------------
# Multi-Cluster Warehouse with STANDARD scaling policy
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "multi_cluster_standard" {
  name              = "MIGRATION_TEST_WH_MULTI_STD"
  warehouse_size    = "XSMALL"
  warehouse_type    = "STANDARD"
  min_cluster_count = 1
  max_cluster_count = 3
  scaling_policy    = "STANDARD"
  comment           = "Multi-cluster warehouse with STANDARD scaling"
}

# ------------------------------------------------------------------------------
# Multi-Cluster Warehouse with ECONOMY scaling policy
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "multi_cluster_economy" {
  name              = "MIGRATION_TEST_WH_MULTI_ECO"
  warehouse_size    = "XSMALL"
  warehouse_type    = "STANDARD"
  min_cluster_count = 1
  max_cluster_count = 2
  scaling_policy    = "ECONOMY"
  comment           = "Multi-cluster warehouse with ECONOMY scaling"
}

# ------------------------------------------------------------------------------
# Warehouse with Query Acceleration Enabled
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "query_accel" {
  name                               = "MIGRATION_TEST_WH_QUERY_ACCEL"
  warehouse_size                     = "XSMALL"
  enable_query_acceleration          = true
  query_acceleration_max_scale_factor = 4
  comment                            = "Warehouse with query acceleration enabled"
}

# ------------------------------------------------------------------------------
# Warehouse with Custom Statement Timeouts
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "timeouts" {
  name                                = "MIGRATION_TEST_WH_TIMEOUTS"
  warehouse_size                      = "XSMALL"
  statement_timeout_in_seconds        = 3600
  statement_queued_timeout_in_seconds = 300
  max_concurrency_level               = 4
  comment                             = "Warehouse with custom statement timeouts"
}

# ------------------------------------------------------------------------------
# Snowpark-Optimized Warehouse
# Note: Snowpark warehouses have different configurations
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "snowpark" {
  name           = "MIGRATION_TEST_WH_SNOWPARK"
  warehouse_size = "MEDIUM"
  warehouse_type = "SNOWPARK-OPTIMIZED"
  comment        = "Snowpark-optimized warehouse"
}

# ------------------------------------------------------------------------------
# Warehouse with All Common Options
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "complete" {
  name                               = "MIGRATION_TEST_WH_COMPLETE"
  warehouse_size                     = "SMALL"
  warehouse_type                     = "STANDARD"
  auto_suspend                       = 300
  auto_resume                        = true
  min_cluster_count                  = 1
  max_cluster_count                  = 1
  scaling_policy                     = "STANDARD"
  enable_query_acceleration          = false
  query_acceleration_max_scale_factor = 8
  statement_timeout_in_seconds       = 7200
  max_concurrency_level              = 8
  comment                            = "Complete warehouse with all options"
}

# ------------------------------------------------------------------------------
# Warehouse with Underscore in Name
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "underscore_name" {
  name           = "MIGRATION_TEST_WH_WITH_UNDERSCORE"
  warehouse_size = "XSMALL"
  comment        = "Warehouse with underscores in name"
}

# ------------------------------------------------------------------------------
# Warehouse with Numbers in Name
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "with_numbers" {
  name           = "MIGRATION_TEST_WH_123"
  warehouse_size = "XSMALL"
  comment        = "Warehouse with numbers in name"
}

