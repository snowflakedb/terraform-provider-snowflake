# =============================================================================
# Expected Migration Script Output - Warehouses
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block warehouses < objects.csv
#
# Use this to compare with actual output:
#   go run .. -import=block warehouses < objects.csv > actual_output.tf
#   diff expected_output.tf actual_output.tf
#
# NOTE:
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# - Special characters are Unicode-escaped in HCL output
# - Warehouses always include: warehouse_type, warehouse_size, auto_suspend,
#   auto_resume, min/max_cluster_count, scaling_policy, query_acceleration_max_scale_factor
# - Standard warehouses include generation = "2"
# - Snowpark warehouses include resource_constraint instead of generation
# - auto_suspend default is account-specific (may vary)
# =============================================================================

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_123" {
  name = "MIGRATION_TEST_WH_123"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Warehouse with numbers in name"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_BASIC" {
  name = "MIGRATION_TEST_WH_BASIC"
  auto_resume = "true"
  auto_suspend = 34
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_COMMENT" {
  name = "MIGRATION_TEST_WH_COMMENT"
  auto_resume = "true"
  auto_suspend = 34
  comment = "This warehouse is used for migration testing purposes"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_COMPLETE" {
  name = "MIGRATION_TEST_WH_COMPLETE"
  auto_resume = "true"
  auto_suspend = 300
  comment = "Complete warehouse with all options"
  generation = "2"
  max_cluster_count = 1
  max_concurrency_level = 8
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  statement_timeout_in_seconds = 7200
  warehouse_size = "SMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_CUSTOM_SUSPEND" {
  name = "MIGRATION_TEST_WH_CUSTOM_SUSPEND"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Warehouse with 10 minute auto-suspend"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_LONG_COMMENT" {
  name = "MIGRATION_TEST_WH_LONG_COMMENT"
  auto_resume = "true"
  auto_suspend = 34
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_MEDIUM" {
  name = "MIGRATION_TEST_WH_MEDIUM"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Medium sized warehouse"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "MEDIUM"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_MULTI_ECO" {
  name = "MIGRATION_TEST_WH_MULTI_ECO"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Multi-cluster warehouse with ECONOMY scaling"
  generation = "2"
  max_cluster_count = 2
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "ECONOMY"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_MULTI_STD" {
  name = "MIGRATION_TEST_WH_MULTI_STD"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Multi-cluster warehouse with STANDARD scaling"
  generation = "2"
  max_cluster_count = 3
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_NO_RESUME" {
  name = "MIGRATION_TEST_WH_NO_RESUME"
  auto_resume = "false"
  auto_suspend = 34
  comment = "Warehouse without auto-resume"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_NO_SUSPEND" {
  name = "MIGRATION_TEST_WH_NO_SUSPEND"
  auto_resume = "true"
  auto_suspend = 0
  comment = "Warehouse that never auto-suspends"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_QUERY_ACCEL" {
  name = "MIGRATION_TEST_WH_QUERY_ACCEL"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Warehouse with query acceleration enabled"
  enable_query_acceleration = "true"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 4
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_SMALL" {
  name = "MIGRATION_TEST_WH_SMALL"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Small sized warehouse"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "SMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_SNOWPARK" {
  name = "MIGRATION_TEST_WH_SNOWPARK"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Snowpark-optimized warehouse"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  resource_constraint = "MEMORY_16X"
  scaling_policy = "STANDARD"
  warehouse_size = "MEDIUM"
  warehouse_type = "SNOWPARK-OPTIMIZED"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_SPECIAL" {
  name = "MIGRATION_TEST_WH_SPECIAL"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Comment with special chars: \u003c\u003e\u0026\""
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_TIMEOUTS" {
  name = "MIGRATION_TEST_WH_TIMEOUTS"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Warehouse with custom statement timeouts"
  generation = "2"
  max_cluster_count = 1
  max_concurrency_level = 4
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 3600
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_WITH_UNDERSCORE" {
  name = "MIGRATION_TEST_WH_WITH_UNDERSCORE"
  auto_resume = "true"
  auto_suspend = 34
  comment = "Warehouse with underscores in name"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}

import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_123
  id = "\"MIGRATION_TEST_WH_123\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_BASIC
  id = "\"MIGRATION_TEST_WH_BASIC\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_COMMENT
  id = "\"MIGRATION_TEST_WH_COMMENT\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_COMPLETE
  id = "\"MIGRATION_TEST_WH_COMPLETE\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_CUSTOM_SUSPEND
  id = "\"MIGRATION_TEST_WH_CUSTOM_SUSPEND\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_LONG_COMMENT
  id = "\"MIGRATION_TEST_WH_LONG_COMMENT\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_MEDIUM
  id = "\"MIGRATION_TEST_WH_MEDIUM\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_MULTI_ECO
  id = "\"MIGRATION_TEST_WH_MULTI_ECO\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_MULTI_STD
  id = "\"MIGRATION_TEST_WH_MULTI_STD\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_NO_RESUME
  id = "\"MIGRATION_TEST_WH_NO_RESUME\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_NO_SUSPEND
  id = "\"MIGRATION_TEST_WH_NO_SUSPEND\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_QUERY_ACCEL
  id = "\"MIGRATION_TEST_WH_QUERY_ACCEL\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_SMALL
  id = "\"MIGRATION_TEST_WH_SMALL\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_SNOWPARK
  id = "\"MIGRATION_TEST_WH_SNOWPARK\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_SPECIAL
  id = "\"MIGRATION_TEST_WH_SPECIAL\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_TIMEOUTS
  id = "\"MIGRATION_TEST_WH_TIMEOUTS\""
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_WITH_UNDERSCORE
  id = "\"MIGRATION_TEST_WH_WITH_UNDERSCORE\""
}
