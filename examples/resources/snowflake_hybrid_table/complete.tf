# Complete hybrid table example demonstrating all features
#
# This example shows advanced configurations including:
# - Identity/autoincrement columns
# - Default values
# - Out-of-line primary key constraint
# - Unique constraints
# - Multiple indexes
# - Collation
# - Comments on all elements
#
# Note: Hybrid tables require preview feature enablement in the provider:
# provider "snowflake" {
#   preview_features_enabled = ["snowflake_hybrid_table_resource"]
# }

resource "snowflake_hybrid_table" "customers" {
  database   = "MY_DATABASE"
  schema     = "MY_SCHEMA"
  name       = "CUSTOMERS"
  comment    = "Customer master table with comprehensive structure"
  or_replace = false

  # Identity column with autoincrement
  column {
    name     = "customer_id"
    type     = "NUMBER(38,0)"
    nullable = false
    comment  = "Auto-incrementing customer identifier"

    identity {
      start_num = 1
      step_num  = 1
    }
  }

  # Column with default expression
  column {
    name     = "created_at"
    type     = "TIMESTAMP_NTZ"
    nullable = false
    comment  = "Timestamp when record was created"

    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  # String column with collation
  column {
    name     = "customer_name"
    type     = "VARCHAR(200)"
    nullable = false
    collate  = "en-ci"
    comment  = "Customer full name with case-insensitive collation"
  }

  # Email column with unique constraint (inline)
  column {
    name     = "email"
    type     = "VARCHAR(255)"
    nullable = false
    unique   = true
    comment  = "Customer email address (unique)"
  }

  column {
    name    = "phone"
    type    = "VARCHAR(20)"
    comment = "Customer phone number"
  }

  # Column with default value
  column {
    name     = "status"
    type     = "VARCHAR(20)"
    nullable = false
    comment  = "Customer status"

    default {
      expression = "'ACTIVE'"
    }
  }

  column {
    name = "loyalty_points"
    type = "NUMBER(10,0)"

    default {
      expression = "0"
    }
  }

  # Nullable timestamp column
  column {
    name    = "last_purchase_date"
    type    = "DATE"
    comment = "Date of most recent purchase"
  }

  # Out-of-line primary key constraint (composite key example)
  # Note: Can also use single column if needed
  primary_key {
    name    = "pk_customers"
    columns = ["customer_id"]
  }

  # Additional unique constraint on multiple columns
  unique_constraint {
    name    = "uk_customer_email_phone"
    columns = ["email", "phone"]
  }

  # Index on frequently queried column
  index {
    name    = "idx_customer_name"
    columns = ["customer_name"]
  }

  # Composite index for common query patterns
  index {
    name    = "idx_status_created"
    columns = ["status", "created_at"]
  }

  # Index on date column for time-based queries
  index {
    name    = "idx_last_purchase"
    columns = ["last_purchase_date"]
  }

  # Data retention configuration
  data_retention_time_in_days = 30
}

# Example outputs to access computed attributes
output "customer_table_fqn" {
  description = "Fully qualified name of the hybrid table"
  value       = snowflake_hybrid_table.customers.fully_qualified_name
}

output "customer_table_show_output" {
  description = "Output from SHOW HYBRID TABLES command"
  value       = snowflake_hybrid_table.customers.show_output
}

output "customer_table_describe_output" {
  description = "Output from DESCRIBE HYBRID TABLE command (column details)"
  value       = snowflake_hybrid_table.customers.describe_output
}
