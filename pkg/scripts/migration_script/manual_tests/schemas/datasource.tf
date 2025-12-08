# =============================================================================
# Schemas Data Source - Fetch schemas from Snowflake and generate CSV
# =============================================================================
# This file fetches the test schemas created by objects_def.tf and
# generates a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create schemas: terraform apply (with objects_def.tf)
#   2. Generate CSV: terraform apply (this will create objects.csv)
#   3. Test migration: go run .. -import=block schemas < objects.csv
# =============================================================================

# Fetch all test schemas (filter by MIGRATION_TEST_SCHEMA_ prefix in the test database)
data "snowflake_schemas" "test_schemas" {
  in {
    database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  }
  like = "MIGRATION_TEST_SCHEMA_%"
}

locals {
  # Transform each schema by merging show_output and flattened parameters
  schemas_flattened = [
    for schema in data.snowflake_schemas.test_schemas.schemas : merge(
      schema.show_output[0],
      # Flatten parameters: convert each parameter to {param_name}_value and {param_name}_level
      {
        for param_key, param_values in schema.parameters[0] :
        "${param_key}_value" => param_values[0].value
      },
      {
        for param_key, param_values in schema.parameters[0] :
        "${param_key}_level" => param_values[0].level
      }
    )
  ]

  # Get all unique keys from the first schema to create CSV header
  schemas_csv_header = length(local.schemas_flattened) > 0 ? join(",", [for key in keys(local.schemas_flattened[0]) : "\"${key}\""]) : ""

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = length(local.schemas_flattened) > 0 ? {
    for schema in local.schemas_flattened :
    schema.name => {
      for key in keys(local.schemas_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(schema, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Convert each schema object to CSV row
  schemas_csv_rows = length(local.schemas_flattened) > 0 ? [
    for schema in local.schemas_flattened :
      join(",", [
        for key in keys(local.schemas_flattened[0]) :
        "\"${local.csv_escape[schema.name][key]}\""
      ])
  ] : []

  # Combine header and rows
  schemas_csv_content = length(local.schemas_flattened) > 0 ? join("\n", concat([local.schemas_csv_header], local.schemas_csv_rows)) : "# No schemas found matching MIGRATION_TEST_SCHEMA_%"
}

# Write the CSV file
resource "local_file" "schemas_csv" {
  content  = local.schemas_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "schemas_found" {
  description = "Number of test schemas found"
  value       = length(local.schemas_flattened)
}

output "schema_names" {
  description = "Names of test schemas found"
  value       = [for schema in local.schemas_flattened : schema.name]
}

