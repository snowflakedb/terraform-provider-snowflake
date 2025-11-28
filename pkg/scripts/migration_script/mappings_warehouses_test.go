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
				{"actives", "auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "failed", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "pendings", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "suspended", "type", "u_u_i_d", "updated_on", "auto_resume_level", "auto_resume_value", "auto_suspend_level", "auto_suspend_value", "comment_level", "comment_value", "enable_query_acceleration_level", "enable_query_acceleration_value", "generation_level", "generation_value", "initially_suspended_level", "initially_suspended_value", "max_cluster_count_level", "max_cluster_count_value", "min_cluster_count_level", "min_cluster_count_value", "query_acceleration_max_scale_factor_level", "query_acceleration_max_scale_factor_value", "resource_constraint_level", "resource_constraint_value", "resource_monitor_level", "resource_monitor_value", "scaling_policy_level", "scaling_policy_value", "warehouse_size_level", "warehouse_size_value", "warehouse_type_level", "warehouse_type_value"},
				{"0", "true", "600", "71.43", "", "2024-06-06 00:00:00.000 +0000 UTC", "false", "0", "", "false", "false", "3", "1", "WH1", "0", "ADMIN", "ROLE", "0", "0", "0", "1", "0", "", "", "2024-06-06 12:00:00.000 +0000 UTC", "5", "ECONOMY", "XSMALL", "2", "AVAILABLE", "0", "STANDARD", "abc-123", "2024-06-06 00:00:00.000 +0000 UTC", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
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
				{"actives", "auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "failed", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "pendings", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "suspended", "type", "u_u_i_d", "updated_on", "auto_resume_level", "auto_resume_value", "auto_suspend_level", "auto_suspend_value", "comment_level", "comment_value", "enable_query_acceleration_level", "enable_query_acceleration_value", "generation_level", "generation_value", "initially_suspended_level", "initially_suspended_value", "max_cluster_count_level", "max_cluster_count_value", "min_cluster_count_level", "min_cluster_count_value", "query_acceleration_max_scale_factor_level", "query_acceleration_max_scale_factor_value", "resource_constraint_level", "resource_constraint_value", "resource_monitor_level", "resource_monitor_value", "scaling_policy_level", "scaling_policy_value", "warehouse_size_level", "warehouse_size_value", "warehouse_type_level", "warehouse_type_value"},
				{"0", "true", "600", "71.43", "", "2024-06-06 00:00:00.000 +0000 UTC", "false", "0", "", "false", "false", "3", "1", "COMPUTE_WH", "0", "ADMIN", "ROLE", "0", "0", "0", "1", "0", "", "", "2024-06-06 12:00:00.000 +0000 UTC", "5", "ECONOMY", "XSMALL", "2", "AVAILABLE", "0", "STANDARD", "def-456", "2024-06-06 00:00:00.000 +0000 UTC", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
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
			name: "warehouse with all fields",
			inputRows: [][]string{
				{"actives", "auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "failed", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "pendings", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "suspended", "type", "u_u_i_d", "updated_on", "auto_resume_level", "auto_resume_value", "auto_suspend_level", "auto_suspend_value", "comment_level", "comment_value", "enable_query_acceleration_level", "enable_query_acceleration_value", "generation_level", "generation_value", "initially_suspended_level", "initially_suspended_value", "max_cluster_count_level", "max_cluster_count_value", "min_cluster_count_level", "min_cluster_count_value", "query_acceleration_max_scale_factor_level", "query_acceleration_max_scale_factor_value", "resource_constraint_level", "resource_constraint_value", "resource_monitor_level", "resource_monitor_value", "scaling_policy_level", "scaling_policy_value", "warehouse_size_level", "warehouse_size_value", "warehouse_type_level", "warehouse_type_value"},
				{"2", "true", "300", "85.50", "Production warehouse", "2024-06-06 00:00:00.000 +0000 UTC", "true", "0", "2", "true", "true", "5", "2", "PROD_WH", "0", "ADMIN", "ROLE", "0", "0", "50", "2", "0", "MEMORY_1X", "MONITOR1", "2024-06-06 12:00:00.000 +0000 UTC", "8", "STANDARD", "LARGE", "4", "AVAILABLE", "0", "STANDARD", "ghi-789", "2024-06-06 00:00:00.000 +0000 UTC", "WAREHOUSE", "true", "WAREHOUSE", "300", "WAREHOUSE", "Production warehouse", "WAREHOUSE", "true", "WAREHOUSE", "2", "WAREHOUSE", "false", "WAREHOUSE", "5", "WAREHOUSE", "2", "WAREHOUSE", "50", "WAREHOUSE", "MEMORY_1X", "WAREHOUSE", "MONITOR1", "WAREHOUSE", "STANDARD", "WAREHOUSE", "LARGE", "WAREHOUSE", "STANDARD"},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_PROD_WH" {
  name = "PROD_WH"
  auto_resume = "true"
  auto_suspend = 300
  comment = "Production warehouse"
  enable_query_acceleration = "true"
  generation = "2"
  initially_suspended = false
  max_cluster_count = 5
  min_cluster_count = 2
  query_acceleration_max_scale_factor = 50
  resource_constraint = "MEMORY_1X"
  resource_monitor = "MONITOR1"
  scaling_policy = "STANDARD"
  warehouse_size = "LARGE"
  warehouse_type = "STANDARD"
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
				{"actives", "auto_resume", "auto_suspend", "available", "comment", "created_on", "enable_query_acceleration", "failed", "generation", "is_current", "is_default", "max_cluster_count", "min_cluster_count", "name", "other", "owner", "owner_role_type", "pendings", "provisioning", "query_acceleration_max_scale_factor", "queued", "quiescing", "resource_constraint", "resource_monitor", "resumed_on", "running", "scaling_policy", "size", "started_clusters", "state", "suspended", "type", "u_u_i_d", "updated_on", "auto_resume_level", "auto_resume_value", "auto_suspend_level", "auto_suspend_value", "comment_level", "comment_value", "enable_query_acceleration_level", "enable_query_acceleration_value", "generation_level", "generation_value", "initially_suspended_level", "initially_suspended_value", "max_cluster_count_level", "max_cluster_count_value", "min_cluster_count_level", "min_cluster_count_value", "query_acceleration_max_scale_factor_level", "query_acceleration_max_scale_factor_value", "resource_constraint_level", "resource_constraint_value", "resource_monitor_level", "resource_monitor_value", "scaling_policy_level", "scaling_policy_value", "warehouse_size_level", "warehouse_size_value", "warehouse_type_level", "warehouse_type_value"},
				{"1", "true", "450", "80.00", "Test warehouse", "2024-06-06 00:00:00.000 +0000 UTC", "true", "0", "1", "false", "false", "4", "1", "TEST_WH", "0", "ADMIN", "ROLE", "0", "0", "30", "1", "0", "MEMORY_2X", "MONITOR2", "2024-06-06 12:00:00.000 +0000 UTC", "3", "CLASSIC", "MEDIUM", "1", "SUSPENDED", "0", "SNOWPARK-OPTIMIZED", "xyz-999", "2024-06-06 00:00:00.000 +0000 UTC", "ACCOUNT", "true", "ACCOUNT", "450", "ACCOUNT", "Test warehouse", "ACCOUNT", "true", "ACCOUNT", "1", "ACCOUNT", "true", "ACCOUNT", "4", "ACCOUNT", "1", "ACCOUNT", "30", "ACCOUNT", "MEMORY_2X", "ACCOUNT", "MONITOR2", "ACCOUNT", "CLASSIC", "ACCOUNT", "MEDIUM", "ACCOUNT", ""},
			},
			expectedOutput: `
resource "snowflake_warehouse" "snowflake_generated_warehouse_TEST_WH" {
  name = "TEST_WH"
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
