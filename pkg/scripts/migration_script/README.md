# Migration script

This script is designed to assist in migrating existing Snowflake objects into Terraform management
by generating the necessary Terraform resources and import statements based on the Snowflake output.
It can be used for both one-time migrations from deprecated resources to the new ones,
as well as importing existing objects into Terraform state.

If there's a need to support more object types, please open an issue or a PR.

## Usage

Example usage will be based on the grant resources, to get more information about other supported object types,
run the command with the `-h` flag.

```shell
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main -h
```

### Perquisites

There are a few things needed before you can proceed further:
- Go installed, see https://go.dev/doc/install for more details.
- Terraform installed, see https://learn.hashicorp.com/tutorials/terraform/install-cli for more details.
- Snowflake account with an ability to create roles and grant privileges (for simplicity, a user with ACCOUNTADMIN role would be the best choice).
- Knowledge how to configure Snowflake connection using Terraform Provider, see [Terraforming Snowflake](https://quickstarts.snowflake.com/guide/terraforming_snowflake/index.html#0) guide for more details.

### 1. Query Snowflake and save the output

For the sake of simplicity, we will only focus on the privileges granted to the account role on the current account.
To set up the necessary objects, run the following commands in Snowflake (e.g., using SnowSight):

```sql
CREATE ROLE TEST_ROLE;
GRANT CREATE DATABASE ON ACCOUNT TO ROLE TEST_ROLE;
GRANT CREATE ROLE ON ACCOUNT TO ROLE TEST_ROLE;

CREATE ROLE TEST_OTHER_ROLE;
GRANT CREATE DATABASE ON ACCOUNT TO ROLE TEST_OTHER_ROLE;
GRANT CREATE ROLE ON ACCOUNT TO ROLE TEST_OTHER_ROLE;
GRANT CREATE USER ON ACCOUNT TO ROLE TEST_OTHER_ROLE;
```

Now, to get the list of grants we are interested in, we can either call

```sql
SHOW GRANTS TO ROLE TEST_ROLE;
SHOW GRANTS TO ROLE TEST_OTHER_ROLE;
```

and combine the outputs, or we can call

```sql
SHOW GRANTS ON ACCOUNT;
```

and filter the output to only include the grants to the roles we are interested in.

Whatever way you choose, save the output to a CSV file, e.g., `example.csv`.
The file contents should look similar to this:
```csv
"created_on","privilege","granted_on","name","granted_to","grantee_name","grant_option","granted_by"
"2025-08-29 06:23:25.920 -0700","CREATE DATABASE","ACCOUNT","IYA62698","ROLE","TEST_ROLE","false","SYSADMIN"
"2025-08-29 06:23:27.253 -0700","CREATE ROLE","ACCOUNT","IYA62698","ROLE","TEST_ROLE","false","USERADMIN"
"2025-08-29 06:23:29.902 -0700","CREATE DATABASE","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","SYSADMIN"
"2025-08-29 06:23:31.390 -0700","CREATE ROLE","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","USERADMIN"
"2025-08-29 06:23:32.747 -0700","CREATE USER","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","USERADMIN"
```

### 2. Generate resources and import statements based on the Snowflake output

To use the script with the saved output, run the following command in the terminal:

```shell
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main -import=block grants < ./example.csv
```

> For more details on the command-line options, run:
> ```shell
> go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main -h
> ```

The script will read the input from the standard input (stdin) and generate the corresponding Terraform resources and import statements as import blocks.
The output is directed to the standard output (stdout), which can be redirected to a file if needed, by running:

```shell
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main -import=block grants < ./example.csv > ./output.tf
```

### 3. Get the generated resources and import them to the state

In this example, we will create a new Terraform project, but you can also use an existing one.
If you are not sure how to do it, follow the [Terraforming Snowflake](http://quickstarts.snowflake.com/guide/terraforming_snowflake/index.html#0) guide.
Create a new directory and navigate to it, create a new file named `main.tf` with the following content:

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=2.6.0"
    }
  }
}

provider "snowflake" {
  # Your configuration options (remember to connect to the same account where you created the roles)
}
```

and run the `terraform init` command.

After project initialization, copy the generated resources and import blocks from the previous step into the main file.
The file should look similar to this:

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=2.6.0"
    }
  }
}

provider "snowflake" {
  # Your configuration options (remember to connect to the same account where you created the roles)
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_OTHER_ROLE_without_grant_option" {
  account_role_name = "TEST_OTHER_ROLE"
  on_account = true
  privileges = ["CREATE DATABASE", "CREATE ROLE", "CREATE USER"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account = true
  privileges = ["CREATE DATABASE", "CREATE ROLE"]
  with_grant_option = false
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_OTHER_ROLE_without_grant_option
  id = "\"TEST_OTHER_ROLE\"|false|false|CREATE DATABASE,CREATE ROLE,CREATE USER|OnAccount"
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option
  id = "\"TEST_ROLE\"|false|false|CREATE DATABASE,CREATE ROLE|OnAccount"
}
```

Then run `terraform plan`. You should the following output:

```
  # snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_OTHER_ROLE_without_grant_option will be imported
    resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_OTHER_ROLE_without_grant_option" {
        account_role_name = "\"TEST_OTHER_ROLE\""
        all_privileges    = false
        always_apply      = false
        id                = "\"TEST_OTHER_ROLE\"|false|false|CREATE DATABASE,CREATE ROLE,CREATE USER|OnAccount"
        on_account        = true
        privileges        = [
            "CREATE DATABASE",
            "CREATE ROLE",
            "CREATE USER",
        ]
        with_grant_option = false
    }

  # snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option will be imported
    resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option" {
        account_role_name = "\"TEST_ROLE\""
        all_privileges    = false
        always_apply      = false
        id                = "\"TEST_ROLE\"|false|false|CREATE DATABASE,CREATE ROLE|OnAccount"
        on_account        = true
        privileges        = [
            "CREATE DATABASE",
            "CREATE ROLE",
        ]
        with_grant_option = false
    }

Plan: 2 to import, 0 to add, 0 to change, 0 to destroy.
```

which indicates that the resources are ready to be imported.

Finally, run the following command to import the resources into the state, by running `terraform apply`.
At the end of the command, you should see an output similar to this:

```
Apply complete! Resources: 2 imported, 0 added, 0 changed, 0 destroyed.
```

which means that the resources have been successfully imported into the state.

By following the above steps, you can migrate other existing Snowflake objects into Terraform and start manging them!