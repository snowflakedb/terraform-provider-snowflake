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
terraform apply
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
