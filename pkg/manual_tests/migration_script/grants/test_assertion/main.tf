# =============================================================================
# Test: Verify that the empty grants assertion works
# =============================================================================
# This file tests that the precondition assertion correctly fails when no
# grants are found after filtering. Use this to verify the assertion works.
# Should fail with "TEST ASSERTION FAILED: No grants found"

terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {}

# Query grants to PUBLIC role (always exists, will have some system grants)
data "snowflake_grants" "empty_test" {
  grants_to {
    account_role = "PUBLIC"
  }
}

locals {
  # Filter with an impossible pattern - no grant will match this
  # This simulates what happens when test objects don't exist
  empty_test_grants = [
    for g in data.snowflake_grants.empty_test.grants : g
    if can(regex("THIS_PATTERN_WILL_NEVER_MATCH_ANYTHING_12345", g.grantee_name))
  ]

  # Reuse the same CSV structure
  csv_columns = ["privilege", "granted_on", "grant_on", "name", "granted_to", "grant_to", "grantee_name", "grant_option", "granted_by"]
  csv_header  = join(",", [for col in local.csv_columns : "\"${col}\""])

  csv_escape = {
    for idx, grant in local.empty_test_grants :
    idx => {
      for col in local.csv_columns :
      col => col == "grant_on" || col == "grant_to" ? "" : (
        col == "grant_option" ? tostring(grant[col]) :
        replace(replace(replace(tostring(grant[col]), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      )
    }
  }

  csv_rows = [
    for idx, grant in local.empty_test_grants :
    join(",", [for col in local.csv_columns : "\"${local.csv_escape[idx][col]}\""])
  ]

  csv_content = join("\n", concat([local.csv_header], local.csv_rows))
}

# This resource should FAIL during apply with precondition error
resource "local_file" "empty_test_csv" {
  content  = local.csv_content
  filename = "${path.module}/empty_test.csv"

  lifecycle {
    precondition {
      condition     = length(local.csv_rows) > 0
      error_message = "TEST ASSERTION FAILED: No grants found. This is expected behavior - the assertion is working correctly!"
    }
  }
}

output "grants_count" {
  description = "Should be 0 - testing empty assertion"
  value       = length(local.csv_rows)
}

