# =============================================================================
# Users Data Source - Fetch users from Snowflake and generate CSV
# =============================================================================
# This file fetches the test users created by user_objects_def.tf and generates
# a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create users: terraform apply (with user_objects_def.tf)
#   2. Generate CSV: terraform apply (this will create users.csv)
#   3. Test migration: go run .. -import=block users < users.csv
# =============================================================================

# Fetch all test users (filter by MIGRATION_TEST_ prefix)
data "snowflake_users" "test_users" {
  like = "MIGRATION_TEST_%"
}

locals {
  # Transform each user by merging show_output, describe_output, and flattened parameters
  users_flattened = [
    for user in data.snowflake_users.test_users.users : merge(
      user.show_output[0],
      # Include describe output fields (if describe_output is present)
      length(user.describe_output) > 0 ? user.describe_output[0] : {},
      # Flatten parameters: convert each parameter to {param_name}_value and {param_name}_level
      {
        for param_key, param_values in user.parameters[0] :
        "${param_key}_value" => param_values[0].value
      },
      {
        for param_key, param_values in user.parameters[0] :
        "${param_key}_level" => param_values[0].level
      }
    )
  ]

  # Get all unique keys from the first user to create CSV header
  users_csv_header = length(local.users_flattened) > 0 ? join(",", [for key in keys(local.users_flattened[0]) : "\"${key}\""]) : ""

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = length(local.users_flattened) > 0 ? {
    for user in local.users_flattened :
    user.name => {
      for key in keys(local.users_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(user, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Convert each user object to CSV row
  users_csv_rows = length(local.users_flattened) > 0 ? [
    for user in local.users_flattened :
      join(",", [
        for key in keys(local.users_flattened[0]) :
        "\"${local.csv_escape[user.name][key]}\""
      ])
  ] : []

  # Combine header and rows
  users_csv_content = length(local.users_flattened) > 0 ? join("\n", concat([local.users_csv_header], local.users_csv_rows)) : "# No users found matching MIGRATION_TEST_%"
}

# Write the CSV file
resource "local_file" "users_csv" {
  content  = local.users_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "users_found" {
  description = "Number of test users found"
  value       = length(local.users_flattened)
}

output "user_names" {
  description = "Names of test users found"
  value       = [for user in local.users_flattened : user.name]
}

