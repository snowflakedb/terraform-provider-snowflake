# Migration Script Manual Tests

This directory contains end-to-end tests for the migration script. Each object type has its own folder with test configuration.

## Quick Start

### Step 0: Set account information in environment variables:
- SNOWFLAKE_ACCOUNT
- SNOWFLAKE_ORGANIZATION
- SNOWFLAKE_HOST
- SNOWFLAKE_USER
- SNOWFLAKE_PASSWORD
- SNOWFLAKE_ROLE

### Step 1: Navigate to object type folder

```bash
cd users  # or account_roles, etc.
```

### Step 2: Clean up any previous state

```bash
rm -rf .terraform terraform.tfstate terraform.tfstate.backup actual_output.tf
```

### Step 3: Move expected_output.tf out of the way (it has import blocks)

```bash
mv expected_output.tf expected_output.tf.bak
```

### Step 4: Initialize Terraform

```bash
terraform init
```

### Step 5: Create test objects on Snowflake

```bash
terraform apply -auto-approve
```

### Step 6: Restore expected_output.tf

```bash
mv expected_output.tf.bak expected_output.tf
```

### Step 7: Fetch objects via data source and generate CSV

```bash
terraform apply -auto-approve -target=data.snowflake_users.test_users -target=local_file.users_csv
```

### Step 8: Run migration script

```bash
cd ../..
go run . -import=block users < manual_tests/users/objects.csv > manual_tests/users/actual_output.tf
```

### Step 9: Compare output

```bash
cd manual_tests/users
diff <(grep -v '^#' expected_output.tf | grep -v '^$') \
     <(grep -v '^#' actual_output.tf | grep -v '^$')
```

### Step 10: Cleanup

```bash
mv expected_output.tf expected_output.tf.bak
terraform destroy -auto-approve
mv expected_output.tf.bak expected_output.tf
```

## Directory Structure

```
manual_tests/
├── README.md                # This file
├── users/                   # Users test
│   ├── objects_def.tf       # Creates test users on Snowflake
│   ├── datasource.tf        # Fetches users, generates CSV
│   ├── expected_output.tf   # Expected migration script output
│   ├── objects.csv          # Generated CSV (after terraform apply)
│   └── actual_output.tf     # Generated output (after test run)
└── <new_object_type>/       # Add new object types here
```

## Adding a New Object Type

To add tests for a new object type (e.g., `warehouses`):

### Step 1: Create the folder

```bash
mkdir -p warehouses
```

### Step 2: Create `objects_def.tf`

Create test objects with various configurations:

```hcl
terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {}

# Basic warehouse
resource "snowflake_warehouse" "basic" {
  name = "MIGRATION_TEST_WH_BASIC"
}

# Warehouse with all parameters
resource "snowflake_warehouse" "complete" {
  name           = "MIGRATION_TEST_WH_COMPLETE"
  comment        = "Test warehouse for migration"
  warehouse_size = "XSMALL"
  # ... more parameters
}
```

**Important naming convention:** Use `MIGRATION_TEST_` prefix for all test objects.

### Step 3: Create `datasource.tf`

Fetch the objects and generate CSV:

```hcl
# Fetch test warehouses
data "snowflake_warehouses" "test_warehouses" {
  like = "MIGRATION_TEST_WH_%"
}

locals {
  # Flatten the data source output
  warehouses_flattened = [
    for wh in data.snowflake_warehouses.test_warehouses.warehouses :
    wh.show_output[0]
  ]

  # Create CSV header
  csv_header = length(local.warehouses_flattened) > 0 ? join(",", [
    for key in keys(local.warehouses_flattened[0]) : "\"${key}\""
  ]) : ""

  # CSV escape function
  csv_escape = length(local.warehouses_flattened) > 0 ? {
    for wh in local.warehouses_flattened :
    wh.name => {
      for key in keys(local.warehouses_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(wh, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  } : {}

  # Create CSV rows
  csv_rows = length(local.warehouses_flattened) > 0 ? [
    for wh in local.warehouses_flattened :
      join(",", [
        for key in keys(local.warehouses_flattened[0]) :
        "\"${local.csv_escape[wh.name][key]}\""
      ])
  ] : []

  csv_content = join("\n", concat([local.csv_header], local.csv_rows))
}

# Write CSV file
resource "local_file" "csv" {
  content  = local.csv_content
  filename = "${path.module}/objects.csv"
}

# Debug outputs
output "objects_found" {
  value = length(local.warehouses_flattened)
}
```

### Step 4: Create `expected_output.tf`

Run the test manually first to generate expected output:

```bash
# Create objects
cd warehouses
terraform init
terraform apply

# Run migration script to see output
cd ../..
go run . -import=block warehouses < warehouses/objects.csv

# Copy the output to expected_output.tf and add comments
```

**Expected output format:**

```hcl
# Expected Migration Script Output
# Run: go run .. -import=block warehouses < objects.csv

resource "snowflake_warehouse" "snowflake_generated_warehouse_MIGRATION_TEST_WH_BASIC" {
  name = "MIGRATION_TEST_WH_BASIC"
}

# ... more resources ...

import {
  to = snowflake_warehouse.snowflake_generated_warehouse_MIGRATION_TEST_WH_BASIC
  id = "\"MIGRATION_TEST_WH_BASIC\""
}
```

## CSV Format Notes

The CSV files use proper RFC 4180 escaping:

- **Double quotes** are escaped by doubling: `"` → `""`
- **Backslashes** are escaped: `\` → `\\`
- **Newlines** are converted to literal `\n` for multi-line values (like RSA keys)
- All fields are quoted

The migration script's `csvUnescape` function handles decoding these escape sequences.

## Common Issues

### Differences in expected output

Common reasons for differences:
- **Uppercase**: Snowflake may uppercase some values (e.g., `login_name`)
- **Unicode escaping**: Special characters like `<>&` become `\u003c\u003e\u0026`
- **Ordering**: Resources/imports may be ordered differently
- **Default values**: Snowflake may return additional fields with default values

Run the test, examine the diff, and update `expected_output.tf` if the differences are expected.
