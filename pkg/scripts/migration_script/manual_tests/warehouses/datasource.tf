# =============================================================================
# Warehouses Data Source - Fetch warehouses from Snowflake and generate CSV
# =============================================================================
# This file fetches the test warehouses created by objects_def.tf and
# generates a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create warehouses: terraform apply (with objects_def.tf)
#   2. Generate CSV: terraform apply (this will create objects.csv)
#   3. Test migration: go run .. -import=block warehouses < objects.csv
# =============================================================================

# Fetch all test warehouses (filter by MIGRATION_TEST_WH_ prefix)
data "snowflake_warehouses" "test_warehouses" {
  like = "MIGRATION_TEST_WH_%"
}

locals {
  # Transform each warehouse by merging show_output and flattened parameters
  warehouses_flattened = [
    for wh in data.snowflake_warehouses.test_warehouses.warehouses : merge(
      wh.show_output[0],
      # Flatten parameters: convert each parameter to {param_name}_value and {param_name}_level
      {
        for param_key, param_values in wh.parameters[0] :
        "${param_key}_value" => param_values[0].value
      },
      {
        for param_key, param_values in wh.parameters[0] :
        "${param_key}_level" => param_values[0].level
      }
    )
  ]

  # Get all unique keys from the first warehouse to create CSV header
  warehouses_csv_header = length(local.warehouses_flattened) > 0 ? join(",", [for key in keys(local.warehouses_flattened[0]) : "\"${key}\""]) : ""

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = length(local.warehouses_flattened) > 0 ? {
    for wh in local.warehouses_flattened :
    wh.name => {
      for key in keys(local.warehouses_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(wh, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Convert each warehouse object to CSV row
  warehouses_csv_rows = length(local.warehouses_flattened) > 0 ? [
    for wh in local.warehouses_flattened :
      join(",", [
        for key in keys(local.warehouses_flattened[0]) :
        "\"${local.csv_escape[wh.name][key]}\""
      ])
  ] : []

  # Combine header and rows
  warehouses_csv_content = length(local.warehouses_flattened) > 0 ? join("\n", concat([local.warehouses_csv_header], local.warehouses_csv_rows)) : "# No warehouses found matching MIGRATION_TEST_WH_%"
}

# Write the CSV file
resource "local_file" "warehouses_csv" {
  content  = local.warehouses_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "warehouses_found" {
  description = "Number of test warehouses found"
  value       = length(local.warehouses_flattened)
}

output "warehouse_names" {
  description = "Names of test warehouses found"
  value       = [for wh in local.warehouses_flattened : wh.name]
}

