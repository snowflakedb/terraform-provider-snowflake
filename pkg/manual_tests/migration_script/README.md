# Migration Script Manual Tests

This directory contains end-to-end tests for the migration script. Each object type has its own folder with test configuration.

## Testing RSA Key Escaping (Users)

To test multi-line RSA public key handling, you need to generate your own keys. The `users/objects_def.tf` has placeholder comments where you can add them.

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
  -import=statement grants < objects.csv > import/main.tf
```

### Step 5: Import in the separate directory

```bash
cd import
terraform init
terraform import .. [command copied from the script's output]
```

This imports the existing Snowflake objects into a fresh state.

### Step 6: Verify plan is empty

```bash
terraform plan -detailed-exitcode
```

This should state that there are `0` fields to change. However, for some objects there are fields that will always show changes. You will need to inspect them manually and make sure that all the fields with changes were expected to have them.

Known fields with expected changes:

- Users
  - `default_secondary_roles_option`
  - `disable_mfa`
  - `disabled`
  - `mins_to_bypass_mfa`
  - `mins_to_unlock`
  - `must_change_password`

### Step 7: Cleanup

```bash
cd ..
terraform destroy
```

## Known Provider Limitations

### Grants: Implicit Grants

Grants with empty `granted_by` are implicit grants (auto-created by Snowflake) and cannot be managed by Terraform. The data source filters these out.
