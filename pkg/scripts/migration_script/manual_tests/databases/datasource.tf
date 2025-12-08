# =============================================================================
# Databases Data Source - Fetch databases from Snowflake and generate CSV
# =============================================================================
# This file fetches the test databases created by objects_def.tf and
# generates a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create databases: terraform apply (with objects_def.tf)
#   2. Generate CSV: terraform apply (this will create objects.csv)
#   3. Test migration: go run .. -import=block databases < objects.csv
# =============================================================================

# Fetch all test databases (filter by MIGRATION_TEST_DB_ prefix)
data "snowflake_databases" "test_databases" {
  like = "MIGRATION_TEST_DB_%"
}

locals {
  # Transform each database by merging show_output and flattened parameters
  databases_flattened = [
    for db in data.snowflake_databases.test_databases.databases : merge(
      db.show_output[0],
      # Flatten parameters: convert each parameter to {param_name}_value and {param_name}_level
      {
        for param_key, param_values in db.parameters[0] :
        "${param_key}_value" => param_values[0].value
      },
      {
        for param_key, param_values in db.parameters[0] :
        "${param_key}_level" => param_values[0].level
      }
    )
  ]

  # Get all unique keys from the first database to create CSV header
  databases_csv_header = length(local.databases_flattened) > 0 ? join(",", [for key in keys(local.databases_flattened[0]) : "\"${key}\""]) : ""

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = length(local.databases_flattened) > 0 ? {
    for db in local.databases_flattened :
    db.name => {
      for key in keys(local.databases_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(db, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Convert each database object to CSV row
  databases_csv_rows = length(local.databases_flattened) > 0 ? [
    for db in local.databases_flattened :
      join(",", [
        for key in keys(local.databases_flattened[0]) :
        "\"${local.csv_escape[db.name][key]}\""
      ])
  ] : []

  # Combine header and rows
  databases_csv_content = length(local.databases_flattened) > 0 ? join("\n", concat([local.databases_csv_header], local.databases_csv_rows)) : "# No databases found matching MIGRATION_TEST_DB_%"
}

# Write the CSV file
resource "local_file" "databases_csv" {
  content  = local.databases_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "databases_found" {
  description = "Number of test databases found"
  value       = length(local.databases_flattened)
}

output "database_names" {
  description = "Names of test databases found"
  value       = [for db in local.databases_flattened : db.name]
}

