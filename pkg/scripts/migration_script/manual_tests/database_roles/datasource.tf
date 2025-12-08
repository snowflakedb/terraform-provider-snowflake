# =============================================================================
# Database Roles Data Source - Fetch roles from Snowflake and generate CSV
# =============================================================================
# This file fetches the test database roles created by objects_def.tf and
# generates a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create roles: terraform apply (with objects_def.tf)
#   2. Generate CSV: terraform apply (this will create objects.csv)
#   3. Test migration: go run .. -import=block database_roles < objects.csv
# =============================================================================

# Fetch all test database roles (filter by MIGRATION_TEST_DBROLE_ prefix)
data "snowflake_database_roles" "test_roles" {
  in_database = "MIGRATION_TEST_DB_FOR_ROLES"
  like        = "MIGRATION_TEST_DBROLE_%"
}

locals {
  # Transform each role to use show_output fields directly
  roles_flattened = [
    for role in data.snowflake_database_roles.test_roles.database_roles :
    role.show_output[0]
  ]

  # Get all unique keys from the first role to create CSV header
  roles_csv_header = length(local.roles_flattened) > 0 ? join(",", [for key in keys(local.roles_flattened[0]) : "\"${key}\""]) : ""

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = length(local.roles_flattened) > 0 ? {
    for role in local.roles_flattened :
    role.name => {
      for key in keys(local.roles_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(role, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Convert each role object to CSV row
  roles_csv_rows = length(local.roles_flattened) > 0 ? [
    for role in local.roles_flattened :
      join(",", [
        for key in keys(local.roles_flattened[0]) :
        "\"${local.csv_escape[role.name][key]}\""
      ])
  ] : []

  # Combine header and rows
  roles_csv_content = length(local.roles_flattened) > 0 ? join("\n", concat([local.roles_csv_header], local.roles_csv_rows)) : "# No roles found matching MIGRATION_TEST_DBROLE_%"
}

# Write the CSV file
resource "local_file" "roles_csv" {
  content  = local.roles_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "roles_found" {
  description = "Number of test database roles found"
  value       = length(local.roles_flattened)
}

output "role_names" {
  description = "Names of test database roles found"
  value       = [for role in local.roles_flattened : role.name]
}

