# Grants Data Source - Fetch grants from Snowflake and generate CSV

# ==============================================================================
# GRANTS TO ACCOUNT ROLES
# ==============================================================================

# Fetch all grants TO the privilege test role
data "snowflake_grants" "to_priv_role" {
  grants_to {
    account_role = snowflake_account_role.priv_role.name
  }

  depends_on = [
    snowflake_grant_privileges_to_account_role.on_account,
    snowflake_grant_privileges_to_account_role.on_database,
    snowflake_grant_privileges_to_account_role.on_warehouse,
    snowflake_grant_privileges_to_account_role.on_schema,
    snowflake_grant_privileges_to_account_role.on_table,
    snowflake_grant_privileges_to_account_role.with_grant_option,
    snowflake_grant_database_role.dbrole_to_role,
  ]
}

# Fetch all grants TO the parent role (to get role grants)
data "snowflake_grants" "to_parent_role" {
  grants_to {
    account_role = snowflake_account_role.parent_role.name
  }

  depends_on = [
    snowflake_grant_account_role.role_to_role,
  ]
}

# ==============================================================================
# GRANTS OF ACCOUNT ROLES (role-to-role and role-to-user grants)
# ==============================================================================

# Fetch grants OF the child role (shows who it's granted to)
data "snowflake_grants" "of_child_role" {
  grants_of {
    account_role = snowflake_account_role.child_role.name
  }

  depends_on = [
    snowflake_grant_account_role.role_to_role,
  ]
}

# ==============================================================================
# GRANTS TO DATABASE ROLES
# ==============================================================================

# Fetch all grants TO the privilege database role
data "snowflake_grants" "to_priv_db_role" {
  grants_to {
    database_role = snowflake_database_role.priv_db_role.fully_qualified_name
  }

  depends_on = [
    snowflake_grant_privileges_to_database_role.on_database,
    snowflake_grant_privileges_to_database_role.on_schema,
    snowflake_grant_privileges_to_database_role.on_table,
    snowflake_grant_privileges_to_database_role.with_grant_option,
  ]
}

# Fetch all grants TO the parent database role
data "snowflake_grants" "to_parent_db_role" {
  grants_to {
    database_role = snowflake_database_role.parent_db_role.fully_qualified_name
  }

  depends_on = [
    snowflake_grant_database_role.dbrole_to_dbrole,
  ]
}

# ==============================================================================
# GRANTS OF DATABASE ROLES
# ==============================================================================

# Fetch grants OF the child database role
data "snowflake_grants" "of_child_db_role" {
  grants_of {
    database_role = snowflake_database_role.child_db_role.fully_qualified_name
  }

  depends_on = [
    snowflake_grant_database_role.dbrole_to_dbrole,
  ]
}

# Fetch grants OF the priv database role (granted to account role)
data "snowflake_grants" "of_priv_db_role" {
  grants_of {
    database_role = snowflake_database_role.priv_db_role.fully_qualified_name
  }

  depends_on = [
    snowflake_grant_database_role.dbrole_to_role,
  ]
}

locals {
  # Combine all grants from different sources
  all_grants = concat(
    [for g in data.snowflake_grants.to_priv_role.grants : g],
    [for g in data.snowflake_grants.to_parent_role.grants : g],
    [for g in data.snowflake_grants.of_child_role.grants : g],
    [for g in data.snowflake_grants.to_priv_db_role.grants : g],
    [for g in data.snowflake_grants.to_parent_db_role.grants : g],
    [for g in data.snowflake_grants.of_child_db_role.grants : g],
    [for g in data.snowflake_grants.of_priv_db_role.grants : g]
  )

  # Filter grants:
  # - Must contain test prefix in grantee_name or name (test objects only)
  # - Must have a non-empty privilege (grants_of returns role membership rows with
  #   empty privilege - the migration script can't handle these)
  # - Must have a non-empty granted_by (implicit grants cannot be managed by Terraform)
  test_grants = [
    for g in local.all_grants : g
    if (can(regex(local.prefix, g.grantee_name)) || can(regex(local.prefix, g.name))) &&
       g.privilege != "" &&
       g.granted_by != ""
  ]

  # CSV header - matches the GrantCsvRow struct
  grants_csv_header = join(",", [for col in keys(local.all_grants[0]) : "\"${col}\""])

  # CSV escape helper - escapes special chars for CSV format
  csv_escape = {
    for idx, grant in local.test_grants :
    idx => {
      for col in keys(local.all_grants[0]) :
      col => (
        col == "grant_option" ? tostring(grant[col]) :
        replace(replace(replace(tostring(grant[col]), "\\", "\\\\"), "\n", "\\n"), "\"", "\"\"")
      )
    }
  }

  # Convert each grant to CSV row
  grants_csv_rows = [
    for idx, grant in local.test_grants :
    join(",", [for col in keys(local.all_grants[0]) : "\"${local.csv_escape[idx][col]}\""])
  ]

  # Remove duplicates by converting to set and back
  grants_csv_rows_unique = tolist(toset(local.grants_csv_rows))

  # Combine header and rows
  grants_csv_content = join("\n", concat([local.grants_csv_header], local.grants_csv_rows_unique))
}

# Write the CSV file
resource "local_file" "grants_csv" {
  content  = local.grants_csv_content
  filename = "${path.module}/objects.csv"

  # Fail if no grants found - this is a test assertion
  lifecycle {
    precondition {
      condition     = length(local.grants_csv_rows_unique) > 0
      error_message = "TEST ASSERTION FAILED: No grants found matching ${local.prefix}. Make sure objects_def.tf resources were created first."
    }
  }
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
