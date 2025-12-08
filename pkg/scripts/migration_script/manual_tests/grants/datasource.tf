# =============================================================================
# Grants Data Source - Fetch grants from Snowflake and generate CSV
# =============================================================================
# This file fetches the test grants created by objects_def.tf and
# generates a CSV file that can be used to test the migration script.
#
# Usage:
#   1. First create grants: terraform apply (with objects_def.tf)
#   2. Generate CSV: terraform apply (this will create objects.csv)
#   3. Test migration: go run .. -import=block grants < objects.csv
#
# Note: Grants data source returns different fields than other data sources.
# We need to fetch grants TO specific roles and grants OF specific roles.
# =============================================================================

# ==============================================================================
# GRANTS TO ACCOUNT ROLES
# ==============================================================================

# Fetch all grants TO the privilege test role
data "snowflake_grants" "to_priv_role" {
  grants_to {
    account_role = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  }
}

# Fetch all grants TO the parent role (to get role grants)
data "snowflake_grants" "to_parent_role" {
  grants_to {
    account_role = "MIGRATION_TEST_GRANT_PARENT_ROLE"
  }
}

# ==============================================================================
# GRANTS OF ACCOUNT ROLES (role-to-role and role-to-user grants)
# ==============================================================================

# Fetch grants OF the child role (shows who it's granted to)
data "snowflake_grants" "of_child_role" {
  grants_of {
    account_role = "MIGRATION_TEST_GRANT_CHILD_ROLE"
  }
}

# ==============================================================================
# GRANTS TO DATABASE ROLES
# ==============================================================================

# Fetch all grants TO the privilege database role
data "snowflake_grants" "to_priv_db_role" {
  grants_to {
    database_role = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  }
}

# Fetch all grants TO the parent database role
data "snowflake_grants" "to_parent_db_role" {
  grants_to {
    database_role = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PARENT_DBROLE\""
  }
}

# ==============================================================================
# GRANTS OF DATABASE ROLES
# ==============================================================================

# Fetch grants OF the child database role
data "snowflake_grants" "of_child_db_role" {
  grants_of {
    database_role = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_CHILD_DBROLE\""
  }
}

# Fetch grants OF the priv database role (granted to account role)
data "snowflake_grants" "of_priv_db_role" {
  grants_of {
    database_role = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  }
}

locals {
  # Combine all grants from different sources
  all_grants = concat(
    [for g in data.snowflake_grants.to_priv_role.grants : g],
    [for g in data.snowflake_grants.of_child_role.grants : g],
    [for g in data.snowflake_grants.to_priv_db_role.grants : g],
    [for g in data.snowflake_grants.to_parent_db_role.grants : g],
    [for g in data.snowflake_grants.of_child_db_role.grants : g],
    [for g in data.snowflake_grants.of_priv_db_role.grants : g]
  )

  # Filter to only our test grants (by grantee_name containing MIGRATION_TEST)
  # and remove duplicates
  test_grants = [
    for g in local.all_grants : g
    if can(regex("MIGRATION_TEST", g.grantee_name)) || can(regex("MIGRATION_TEST", g.name))
  ]

  # CSV header - matches the GrantCsvRow struct
  grants_csv_header = "\"privilege\",\"granted_on\",\"grant_on\",\"name\",\"granted_to\",\"grant_to\",\"grantee_name\",\"grant_option\",\"granted_by\""

  # CSV escape function
  csv_escape_grant = {
    for idx, grant in local.test_grants :
    idx => {
      privilege   = replace(replace(replace(tostring(grant.privilege), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      granted_on  = replace(replace(replace(tostring(grant.granted_on), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      grant_on    = ""
      name        = replace(replace(replace(tostring(grant.name), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      granted_to  = replace(replace(replace(tostring(grant.granted_to), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      grant_to    = ""
      grantee_name = replace(replace(replace(tostring(grant.grantee_name), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      grant_option = tostring(grant.grant_option)
      granted_by  = replace(replace(replace(tostring(grant.granted_by), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
    }
  }

  # Convert each grant to CSV row
  grants_csv_rows = [
    for idx, grant in local.test_grants :
    "\"${local.csv_escape_grant[idx].privilege}\",\"${local.csv_escape_grant[idx].granted_on}\",\"${local.csv_escape_grant[idx].grant_on}\",\"${local.csv_escape_grant[idx].name}\",\"${local.csv_escape_grant[idx].granted_to}\",\"${local.csv_escape_grant[idx].grant_to}\",\"${local.csv_escape_grant[idx].grantee_name}\",\"${local.csv_escape_grant[idx].grant_option}\",\"${local.csv_escape_grant[idx].granted_by}\""
  ]

  # Remove duplicates by converting to set and back
  grants_csv_rows_unique = tolist(toset(local.grants_csv_rows))

  # Combine header and rows
  grants_csv_content = length(local.grants_csv_rows_unique) > 0 ? join("\n", concat([local.grants_csv_header], local.grants_csv_rows_unique)) : "# No grants found"
}

# Write the CSV file
resource "local_file" "grants_csv" {
  content  = local.grants_csv_content
  filename = "${path.module}/objects.csv"
}

# Output for debugging
output "grants_found" {
  description = "Number of test grants found"
  value       = length(local.grants_csv_rows_unique)
}

output "grants_to_priv_role_count" {
  description = "Grants to priv role"
  value       = length(data.snowflake_grants.to_priv_role.grants)
}

output "grants_of_child_role_count" {
  description = "Grants of child role"
  value       = length(data.snowflake_grants.of_child_role.grants)
}

output "grants_to_priv_db_role_count" {
  description = "Grants to priv db role"
  value       = length(data.snowflake_grants.to_priv_db_role.grants)
}

