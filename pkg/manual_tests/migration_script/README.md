# Migration Script Manual Tests

This directory contains end-to-end tests for the migration script. Each object type has its own folder with test configuration.

## Testing RSA Key Escaping (Users)

To test multi-line RSA public key handling, you need to generate your own keys. The `users/objects_def.tf` has placeholder comments where you can add them.

### Step 1: Generate RSA key pairs

```bash
# Generate 3 key pairs (for person_rsa, service_rsa key 1 & 2, legacy_rsa)
for i in 1 2 3; do
  openssl genrsa -out rsa_key_$i.pem 2048
  openssl rsa -in rsa_key_$i.pem -pubout -out rsa_key_$i.pub
done
```

### Step 2: Extract the public key body (without headers)

```bash
# View key 1 (for person_rsa and service_rsa key 1)
sed -n '2,$p' rsa_key_1.pub | sed '$d'

# View key 2 (for service_rsa key 2)
sed -n '2,$p' rsa_key_2.pub | sed '$d'

# View key 3 (for legacy_rsa)
sed -n '2,$p' rsa_key_3.pub | sed '$d'
```

### Step 3: Add keys to objects_def.tf

Edit `users/objects_def.tf` and find the commented `rsa_public_key` sections. Uncomment them and paste the key bodies:

```hcl
# Example for person_rsa:
rsa_public_key = <<-EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
...your key content...
...AQAB
EOT
```

### Step 4: Cleanup generated keys

```bash
rm -f rsa_key_*.pem rsa_key_*.pub
```

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
```

This imports the existing Snowflake objects into a fresh state.

### Step 6: Verify plan is empty

```bash
terraform plan -detailed-exitcode
```

This should X to import, 0 to add, 0 to change, 0 to destroy.

### Step 7: Cleanup

```bash
cd ..
terraform destroy
```

## Known Provider Limitations

### Grants: Implicit Grants

Grants with empty `granted_by` are implicit grants (auto-created by Snowflake) and cannot be managed by Terraform. The datasource filters these out.
