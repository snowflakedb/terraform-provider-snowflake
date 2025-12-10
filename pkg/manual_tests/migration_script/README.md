# Migration Script Manual Tests

This directory contains end-to-end tests for the migration script. Each object type has its own folder with test configuration.

## Quick Start

### Step 1: Navigate to object type folder

```bash
cd grants  # or schemas, warehouses, users, etc.
```

### Step 2: Clean up any previous state

```bash
rm -rf .terraform terraform.tfstate terraform.tfstate.backup import/.terraform import/terraform.tfstate import/terraform.tfstate.backup import/main.tf
```

### Step 3: Initialize and create test objects

```bash
terraform init
terraform apply -auto-approve
```

This creates objects from `objects_def.tf` and generates `objects.csv` via `datasource.tf`.

### Step 4: Run migration script (output to import directory)

```bash
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@dev \
  -import=block grants < objects.csv > import/main.tf
```

### Step 5: Import in the separate directory

```bash
cd import
terraform init
terraform apply
```

This imports the existing Snowflake objects into a fresh state.

### Step 6: Verify plan is empty

```bash
terraform plan
```

### Step 7: Cleanup

```bash
cd ..
terraform destroy
```

## Test Assertions

The `datasource.tf` includes a **precondition assertion** that fails if no grants are found. This prevents generating an empty CSV that would cause silent failures.

### How it works

```hcl
resource "local_file" "grants_csv" {
  # ...
  lifecycle {
    precondition {
      condition     = length(local.grants_csv_rows_unique) > 0
      error_message = "TEST ASSERTION FAILED: No grants found. Make sure objects_def.tf resources were created first."
    }
  }
}
```

### Testing the assertion

To verify the assertion works correctly, use the `test_assertion/` subdirectory:

```bash
cd test_assertion
terraform init
terraform apply
# Expected error: "TEST ASSERTION FAILED: No grants found..."
```
