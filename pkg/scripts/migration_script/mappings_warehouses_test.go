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
			name: "warehouse without parameters",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "type", "updated_on"},
				{"true", "600", "71.43", "", "2024-06-06 00:00:00.000 +0000 UTC", "false", "", "false", "false", "3", "1", "WH1", "0", "ADMIN", "ROLE", "0", "0", "1", "0", "", "", "2024-06-06 12:00:00.000 +0000 UTC", "5", "ECONOMY", "XSMALL", "2", "AVAILABLE", "STANDARD", "2024-06-06 00:00:00.000 +0000 UTC"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_WH1" {
  name = "WH1"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_WH1
  id = "\"WH1\""
}
`,
		},
		{
			name: "minimal warehouse",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "type", "updated_on", "max_concurrency_level_level", "max_concurrency_level_value", "statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value"},
				{"true", "600", "71.43", "", "2024-06-06 00:00:00.000 +0000 UTC", "false", "", "false", "false", "3", "1", "COMPUTE_WH", "0", "ADMIN", "ROLE", "0", "0", "1", "0", "", "", "2024-06-06 12:00:00.000 +0000 UTC", "5", "ECONOMY", "XSMALL", "2", "AVAILABLE", "STANDARD", "2024-06-06 00:00:00.000 +0000 UTC", "", "", "", "", "", ""},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_COMPUTE_WH" {
  name = "COMPUTE_WH"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_COMPUTE_WH
  id = "\"COMPUTE_WH\""
}
`,
		},
		{
			name: "warehouse with all parameters",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "type", "updated_on", "max_concurrency_level_level", "max_concurrency_level_value", "statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value"},
				{"true", "300", "85.50", "Production warehouse", "2024-06-06 00:00:00.000 +0000 UTC", "true", "2", "true", "true", "5", "2", "PROD_WH", "0", "ADMIN", "ROLE", "0", "50", "2", "0", "MEMORY_1X", "MONITOR1", "2024-06-06 12:00:00.000 +0000 UTC", "8", "STANDARD", "LARGE", "4", "AVAILABLE", "STANDARD", "2024-06-06 00:00:00.000 +0000 UTC", "WAREHOUSE", "10", "WAREHOUSE", "600", "WAREHOUSE", "300"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_PROD_WH" {
  name = "PROD_WH"
  comment = "Production warehouse"
  max_concurrency_level = 10
  statement_queued_timeout_in_seconds = 600
  statement_timeout_in_seconds = 300
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_PROD_WH
  id = "\"PROD_WH\""
}
`,
		},
		{
			name: "warehouse with all parameters set on higher level",
			inputRows: [][]string{
				{"auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "type", "updated_on", "max_concurrency_level_level", "max_concurrency_level_value", "statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value"},
				{"true", "450", "80.00", "Test warehouse", "2024-06-06 00:00:00.000 +0000 UTC", "true", "1", "false", "false", "4", "1", "TEST_WH", "0", "ADMIN", "ROLE", "0", "30", "1", "0", "MEMORY_2X", "MONITOR2", "2024-06-06 12:00:00.000 +0000 UTC", "3", "CLASSIC", "MEDIUM", "1", "SUSPENDED", "SNOWPARK-OPTIMIZED", "2024-06-06 00:00:00.000 +0000 UTC", "ACCOUNT", "5", "ACCOUNT", "900", "ACCOUNT", "120"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_TEST_WH" {
  name = "TEST_WH"
  comment = "Test warehouse"
}
import {
  to = snowflake_warehouse.snowflake_generated_warehouse_TEST_WH
  id = "\"TEST_WH\""
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
