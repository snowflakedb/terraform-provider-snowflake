#!/bin/bash

# Import existing hybrid table into Terraform state
#
# Hybrid tables are imported using their fully qualified name in the format:
# <database_name>.<schema_name>.<table_name>
#
# Note: Before importing, ensure the preview feature is enabled in your provider:
# provider "snowflake" {
#   preview_features_enabled = ["snowflake_hybrid_table_resource"]
# }

# Example 1: Import a hybrid table with standard identifiers
terraform import snowflake_hybrid_table.orders "MY_DATABASE.MY_SCHEMA.ORDERS"

# Example 2: Import a hybrid table with quoted identifiers (case-sensitive names)
terraform import snowflake_hybrid_table.orders '"MyDatabase"."MySchema"."Orders"'

# Example 3: Import using variables
DATABASE="PROD_DB"
SCHEMA="PUBLIC"
TABLE="CUSTOMER_ORDERS"
terraform import snowflake_hybrid_table.orders "${DATABASE}.${SCHEMA}.${TABLE}"

# After successful import, run terraform plan to see any differences
# between the actual state and your configuration:
# terraform plan

# Important notes:
# 1. The resource name (e.g., "orders") in your Terraform configuration must match
#    the resource name used in the import command
#
# 2. After import, you'll need to define the resource in your .tf files with at least:
#    - database, schema, name (required)
#    - column definitions (at least one column with primary key)
#
# 3. The import command only adds the resource to state - you still need to create
#    the corresponding resource block in your configuration
#
# 4. Indexes, constraints, and other settings will be imported automatically
#    and visible in terraform plan output
#
# 5. Use terraform show to view the imported resource state:
#    terraform show snowflake_hybrid_table.orders

# Example resource block to add after import:
cat <<'EOF'
resource "snowflake_hybrid_table" "orders" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "ORDERS"

  column {
    name        = "order_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name = "customer_name"
    type = "VARCHAR(100)"
  }

  column {
    name = "order_date"
    type = "DATE"
  }

  # Add other columns, indexes, and constraints as needed
  # based on terraform plan output
}
EOF
