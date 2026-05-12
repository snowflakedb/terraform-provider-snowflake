# Users Data Source - Fetch users from Snowflake and generate CSV

# Fetch all test users (filter by UUID prefix)
data "snowflake_users" "test_users" {
  like = "${local.prefix}%"

  depends_on = [
    snowflake_user.person_basic,
    snowflake_user.person_complete,
    snowflake_user.person_rsa,
    snowflake_user.person_special_chars,
    snowflake_user.person_params,
    snowflake_user.person_disabled,
    snowflake_service_user.service_basic,
    snowflake_service_user.service_complete,
    snowflake_service_user.service_rsa,
    snowflake_service_user.service_params,
    snowflake_legacy_service_user.legacy_basic,
    snowflake_legacy_service_user.legacy_complete,
    snowflake_legacy_service_user.legacy_rsa,
    snowflake_user.long_comment,
    snowflake_user.unicode,
  ]
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
  users_csv_header = join(",", [for key in keys(local.users_flattened[0]) : "\"${key}\""])

  # CSV escape function: properly escape special characters for CSV format
  # - Backslashes are escaped first (\ -> \\)
  # - Newlines become literal \n (can be decoded by the migration script)
  # - Double quotes are doubled per RFC 4180 (" -> "")
  csv_escape = {
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
  }

  # Convert each user object to CSV row
  users_csv_rows = [
    for user in local.users_flattened :
    join(",", [
      for key in keys(local.users_flattened[0]) :
      "\"${local.csv_escape[user.name][key]}\""
    ])
  ]

  # Combine header and rows
  users_csv_content = join("\n", concat([local.users_csv_header], local.users_csv_rows))
}

# Write the CSV file
resource "local_file" "users_csv" {
  content  = local.users_csv_content
  filename = "${path.module}/objects.csv"

  # Fail if no users found - this is a test assertion
  lifecycle {
    precondition {
      condition     = length(local.users_flattened) > 0
      error_message = "TEST ASSERTION FAILED: No users found matching ${local.prefix}%. Make sure objects_def.tf resources were created first."
    }
  }
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
