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

## Steps using script

```bash
# List available object types
./run_test.sh --list

# Run a test
./run_test.sh users

# Skip object creation (if objects already exist)
./run_test.sh users --skip-create

# Destroy test resources
./run_test.sh users --destroy
```

## Manual Steps (Without Script)

If you prefer to run the steps manually or need more control:

### Step 1: Navigate to object type folder

```bash
cd users  # or account_roles, etc.
```

### Step 2: Initialize Terraform (first time only)

```bash
terraform init
```

### Step 3: Create test objects on Snowflake

```bash
terraform apply -auto-approve
```

### Step 4: Fetch objects and generate CSV

The same `terraform apply` also runs the data source and generates `objects.csv`.

### Step 5: Run migration script

```bash
cd ..  # back to migration_script folder
go run . -import=block users < manual_tests/users/objects.csv > manual_tests/users/actual_output.tf
```

### Step 6: Compare output

```bash
cd manual_tests/users
diff expected_output.tf actual_output.tf
```

Or for a cleaner comparison (ignoring comments):

```bash
diff <(grep -v '^#' expected_output.tf | grep -v '^$') \
     <(grep -v '^#' actual_output.tf | grep -v '^$')
```

### Step 7: Cleanup (when done)

```bash
terraform destroy -auto-approve
```

## Directory Structure

```
manual_tests/
├── run_test.sh              # Main test runner
├── README.md                # This file
├── users/                   # Users test
│   ├── objects_def.tf       # Creates test users on Snowflake
│   ├── datasource.tf        # Fetches users, generates CSV
│   ├── expected_output.tf   # Expected migration script output
│   ├── objects.csv          # Generated CSV (after terraform apply)
│   └── actual_output.tf     # Generated output (after test run)
└── <new_object_type>/       # Add new object types here
```

## Test Workflow

The test runner (`run_test.sh`) performs the following steps:

1. **Create Objects** (`terraform apply` targeting `objects_def.tf` resources)
   - Creates test objects on Snowflake with various configurations

2. **Fetch Objects** (`terraform apply` targeting `datasource.tf` resources)
   - Fetches objects using the appropriate data source
   - Generates `objects.csv` with proper CSV escaping

3. **Run Migration Script**
   - Runs: `go run .. -import=block <object_type> < objects.csv`
   - Saves output to `actual_output.tf`

4. **Compare Output**
   - Compares `actual_output.tf` with `expected_output.tf`
   - Reports differences

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
cd ..
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

### Step 5: Test

```bash
./run_test.sh warehouses
```

## CSV Format Notes

The CSV files use proper RFC 4180 escaping:

- **Double quotes** are escaped by doubling: `"` → `""`
- **Backslashes** are escaped: `\` → `\\`
- **Newlines** are converted to literal `\n` for multi-line values (like RSA keys)
- All fields are quoted

The migration script's `csvUnescape` function handles decoding these escape sequences.

## Common Issues

### "Object type not found"

Make sure the folder exists and contains `objects_def.tf`, `datasource.tf`, and `expected_output.tf`.

### "Object already exists" errors

If objects already exist on Snowflake, use `--skip-create`:

```bash
./run_test.sh users --skip-create
```

Or manually skip Step 3 and go directly to fetching.

### Terraform state issues

Each object type folder has its own Terraform state. If state gets corrupted:

```bash
cd <object_type>
rm -rf .terraform terraform.tfstate*
terraform init
```

### Differences in expected output

Common reasons for differences:
- **Uppercase**: Snowflake may uppercase some values (e.g., `login_name`)
- **Unicode escaping**: Special characters like `<>&` become `\u003c\u003e\u0026`
- **Ordering**: Resources/imports may be ordered differently
- **Default values**: Snowflake may return additional fields with default values

Run the test, examine the diff, and update `expected_output.tf` if the differences are expected.
