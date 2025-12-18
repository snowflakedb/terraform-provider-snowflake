package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleWarehousesMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "minimal warehouse",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "name", "max_cluster_count", "min_cluster_count", "query_acceleration_max_scale_factor", "scaling_policy", "size", "type"},
				{"true", "600", "WAREHOUSE", "1", "1", "8", "STANDARD", "XSMALL", "STANDARD"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_WAREHOUSE" {
  name = "WAREHOUSE"
  auto_resume = "true"
  auto_suspend = 600
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WAREHOUSE
  id = "\"WAREHOUSE\""
}
`,
		},
		{
			name: "minimal warehouse with all parameters set on account level",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "comment", "max_cluster_count", "min_cluster_count", "name", "query_acceleration_max_scale_factor", "scaling_policy", "size", "type", "max_concurrency_level_level", "max_concurrency_level_value", "statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value"},
				{"true", "600", "Warehouse with params set on account level", "1", "1", "WAREHOUSE", "8", "STANDARD", "XSMALL", "STANDARD", "ACCOUNT", "8", "ACCOUNT", "300", "ACCOUNT", "86400"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_WAREHOUSE" {
  name = "WAREHOUSE"
  auto_resume = "true"
  auto_suspend = 600
  comment = "Warehouse with params set on account level"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "XSMALL"
  warehouse_type = "STANDARD"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WAREHOUSE
  id = "\"WAREHOUSE\""
}
`,
		},
		{
			name: "warehouse with all fields",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "comment", "enable_query_acceleration", "generation", "max_cluster_count", "min_cluster_count", "name", "query_acceleration_max_scale_factor", "resource_constraint", "resource_monitor", "scaling_policy", "size", "type", "max_concurrency_level_level", "max_concurrency_level_value", "statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value"},
				{"false", "1200", "An example warehouse.", "true", "", "4", "2", "WAREHOUSE", "16", "MEMORY_16X", "MONITOR", "ECONOMY", "MEDIUM", "SNOWPARK-OPTIMIZED", "WAREHOUSE", "8", "WAREHOUSE", "300", "WAREHOUSE", "86400"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_WAREHOUSE" {
  name = "WAREHOUSE"
  auto_resume = "false"
  auto_suspend = 1200
  comment = "An example warehouse."
  enable_query_acceleration = "true"
  max_cluster_count = 4
  max_concurrency_level = 8
  min_cluster_count = 2
  query_acceleration_max_scale_factor = 16
  resource_constraint = "MEMORY_16X"
  resource_monitor = "MONITOR"
  scaling_policy = "ECONOMY"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 86400
  warehouse_size = "MEDIUM"
  warehouse_type = "SNOWPARK-OPTIMIZED"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WAREHOUSE
  id = "\"WAREHOUSE\""
}
`,
		},
		{
			name: "warehouse with generation",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "comment", "generation", "max_cluster_count", "min_cluster_count", "name", "query_acceleration_max_scale_factor", "scaling_policy", "size", "type"},
				{"true", "600", "Gen2 warehouse", "2", "1", "1", "WAREHOUSE", "8", "STANDARD", "LARGE", "STANDARD"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_WAREHOUSE" {
  name = "WAREHOUSE"
  auto_resume = "true"
  auto_suspend = 600
  comment = "Gen2 warehouse"
  generation = "2"
  max_cluster_count = 1
  min_cluster_count = 1
  query_acceleration_max_scale_factor = 8
  scaling_policy = "STANDARD"
  warehouse_size = "LARGE"
  warehouse_type = "STANDARD"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WAREHOUSE
  id = "\"WAREHOUSE\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleWarehouses(&Config{
				ObjectType: ObjectTypeWarehouses,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
