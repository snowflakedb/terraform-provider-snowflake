# Migration script

<!-- TOC -->
* [Migration script](#migration-script)
    * [Compatibility with Provider Versions](#compatibility-with-provider-versions)
  * [Syntax](#syntax)
  * [Usage](#usage)
    * [Prerequisites](#prerequisites)
    * [Use case: Migrate deprecated resources to new ones](#use-case-migrate-deprecated-resources-to-new-ones)
      * [1. Query Snowflake and save the output](#1-query-snowflake-and-save-the-output)
      * [2. Generate resources and import statements based on the Snowflake output](#2-generate-resources-and-import-statements-based-on-the-snowflake-output)
      * [3. Importing auto-generated resources to the state](#3-importing-auto-generated-resources-to-the-state)
      * [4. Removing old grant resources](#4-removing-old-grant-resources)
      * [5. Update generated resources](#5-update-generated-resources)
    * [Use case: Migrate existing grants to Terraform](#use-case-migrate-existing-grants-to-terraform)
      * [1. Query Snowflake and save the output](#1-query-snowflake-and-save-the-output-1)
      * [2. Generate resources and import statements based on the Snowflake output](#2-generate-resources-and-import-statements-based-on-the-snowflake-output-1)
      * [3. Get the generated resources and import them to the state](#3-get-the-generated-resources-and-import-them-to-the-state)
  * [Limitations](#limitations)
    * [Generated resource names](#generated-resource-names)
    * [No dependencies handling](#no-dependencies-handling)
<!-- TOC -->

This script is designed to assist in migrating existing Snowflake objects into Terraform management
by generating the necessary Terraform resources and import statements based on the Snowflake output.
It can be used for both one-time migrations from deprecated resources to the new ones,
as well as importing existing objects into Terraform state.

The script was provided to give an idea how the migration process can be automated.
It is not officially supported, and we do not prioritize fixes for it.
Feel free to use it as a starting point and modify it to fit your specific needs.
We are open to contributions to enhance its functionality.
If you are planning to contribute, please check out our [contributing guide](./CONTRIBUTING.md).

### Compatibility with Provider Versions

The script is designed to work with the latest version of the provider.
However, if you're using an older version, you can still utilize the script as long as the object types you need are supported and haven't undergone major changes.
For instance, with grants, you can use the script to transition from [old to new grants](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#mapping-from-old-grant-resources-to-the-new-ones)
since they haven't significantly changed from the current provider version, but there may be minor differences like quotes handling in identifiers.

## Syntax

Use the following syntax to run the migration script from your terminal:

```shell
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main [flags] [OBJECT_TYPE] < [INPUT] > [OUTPUT]
```

> **Note**: It's recommended to use the latest version of the script by specifying `@main` at the end of the script path.

where script options are:
- **flags**:
  - `-h`: Displays help information about the script.
  - `-import`: Specifies the import format. Supported values are `statement` (default) and `block`. For example, to generate block imports, specify `-import=block`.
    - `block`: Generates [import blocks](https://developer.hashicorp.com/terraform/language/import) at the bottom of the generated Terraform configuration.
    - `statement`: Generates commented [import commands](https://developer.hashicorp.com/terraform/cli/commands/import) at the bottom of the generated Terraform configuration.
- **OBJECT_TYPES**:
  - `grants`: Generates resources and import statements for Snowflake grants. The expected input is in the form of [`SHOW GRANTS`](https://docs.snowflake.com/en/sql-reference/sql/show-grants) output.

    The allowed SHOW GRANTS commands are:
      - `SHOW GRANTS ON ACCOUNT`
      - `SHOW GRANTS ON <object_type>`
      - `SHOW GRANTS TO ROLE <role_name>`
      - `SHOW GRANTS TO DATABASE ROLE <database_role_name>`

    Supported resources:
      - snowflake_grant_privileges_to_account_role
      - snowflake_grant_privileges_to_database_role
      - snowflake_grant_account_role
      - snowflake_grant_database_role

    Limitations:
      - grants on 'future' or on 'all' objects are not supported
      - all_privileges and always_apply fields are not supported
  - `schemas` which expects a converted CSV output from the snowflake_schemas data source. To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW SCHEMAS output, so the CSV header looks like `"comment","created_on",...,"catalog_value","catalog_level","data_retention_time_in_days_value","data_retention_time_in_days_level",...`
    When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "SCHEMA".
    For more details about using multiple sources, visit the [Multiple sources section](#multiple-sources).

    Supported resources:
      - snowflake_schema
  - `databases` which expects a converted CSV output from the snowflake_databases data source. To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW DATABASES output, so the CSV header looks like `"comment","created_on",...,"catalog_value","catalog_level","data_retention_time_in_days_value","data_retention_time_in_days_level",...`
      When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "DATABASE".
      For more details about using multiple sources, visit the [Multiple sources section](#multiple-sources).

    Supported resources:
      - snowflake_database

  - `warehouses` which expects a converted CSV output from the snowflake_warehouses data source.
      To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW WAREHOUSES output, so the CSV header looks like `"comment","created_on",...,"max_cluster_count","min_cluster_count","name","other",...`
      When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "WAREHOUSE".
      The script always outputs fields that have non-empty default values in Snowflake (they can be removed from the output)

      Caution: Some of the fields are not supported (actives, pendings, failed, suspended, uuid, initially_suspended)

      For more details about using multiple sources, visit [Multiple sources section](#multiple-sources).

    Supported resources:
      - snowflake_warehouse

  - `account_roles` which expects input in the form of [`SHOW ROLES`](https://docs.snowflake.com/en/sql-reference/sql/show-roles) output. Can also be obtained as a converted CSV output from the snowflake_account_roles data source.

    Supported resources:
      - snowflake_account_role

  - `database_roles` which expects input in the form of [`SHOW DATABASE ROLES`](https://docs.snowflake.com/en/sql-reference/sql/show-database-roles) output. Can also be obtained as a converted CSV output from the snowflake_database_roles data source.

    Supported resources:
      - snowflake_database_role

  - `users` which expects a converted CSV output from the snowflake_users data source.
      To support object parameters, one should use the SHOW PARAMETERS output, and combine it with the SHOW USERS output, so the CSV header looks like `"comment","created_on",...,"abort_detached_query_value","abort_detached_query_level","timezone_value","timezone_level",...`
      When the additional columns are present, the resulting resource will have the parameters values, if the parameter level is set to "USER".

      Caution: password parameter is not supported as it is returned in the form of `"***"` from the data source.

      Note: Newlines are allowed only in the `comment`, `rsa_public_key` and `rsa_public_key2` fields, they might cause errors and require manual corrections elsewhere.

      For more details about using multiple sources, visit the [Multiple sources section](#multiple-sources).

    Different user types are mapped to their respective Terraform resources based on the `type` attribute:
      - `PERSON` (or empty) → `snowflake_user` - A human user who can interact with Snowflake
      - `SERVICE` → `snowflake_service_user` - A service or application user without human interaction (cannot use password/SAML authentication, cannot have first_name, last_name, must_change_password)
      - `LEGACY_SERVICE` → `snowflake_legacy_service_user` - Similar to SERVICE but allows password and SAML authentication (cannot have first_name, last_name)

    Supported resources:
      - snowflake_user
      - snowflake_service_user
      - snowflake_legacy_service_user

- **INPUT**:
  - Migration script operates on STDIN input in CSV format. You can redirect the input from a file or pipe it from another command.
- **OUTPUT**:
  - Migration script writes the generated content to STDOUT. You can redirect the output wherever you need to, for example, to a file.
  - **It's user's responsibility to ensure that the output is written securely to a safe location and not to overwrite any important files.**

## Usage

The expected usage use cases for the scripts are either one-time migration from deprecated resources to the new ones,
or importing existing objects into Terraform state. They are both covered in the examples below (in the separate sections for each).
They focus on the grants object_type, but a similar approach could be used for other object types as well.

### Prerequisites

There are a few things needed before you can proceed further:
- To run the script in the recommended way, you will need:
  - Go installed, see https://go.dev/doc/install for more details.
- To test the Terraform configurations locally, you will need:
  - Terraform installed, see https://learn.hashicorp.com/tutorials/terraform/install-cli for more details.
  - Snowflake account with an ability to create sample roles and grant privileges (for simplicity, a user with ACCOUNTADMIN role would be the best choice).
  - Knowledge how to configure Snowflake connection using Terraform Provider, see [Terraforming Snowflake](https://quickstarts.snowflake.com/guide/terraforming_snowflake/index.html#0) guide for more details.

### Use case: Migrate deprecated resources to new ones

In this example, we will focus on migrating from the removed `snowflake_account_grant` resource to the new `snowflake_grant_privileges_to_account_role` resource.
Check [mapping old grant resources to the new ones](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#mapping-from-old-grant-resources-to-the-new-ones)
for more details on how to migrate other grant resources. Our starting configuration will look as follows:

> If you are not sure how to set it up, follow the [Terraforming Snowflake](http://quickstarts.snowflake.com/guide/terraforming_snowflake/index.html#0) guide.

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=0.92.0"
    }
  }
}

provider "snowflake" {
  # Your configuration options
}

resource "snowflake_role" "test_role" {
  name = "TEST_ROLE"
}

resource "snowflake_role" "other_test_role" {
  name = "OTHER_TEST_ROLE"
}

resource "snowflake_account_grant" "grant" {
  roles             = [snowflake_role.test_role.name]
  privilege         = "CREATE ROLE"
  with_grant_option = false
}

resource "snowflake_account_grant" "other_grant" {
  roles             = [snowflake_role.other_test_role.name]
  privilege         = "CREATE DATABASE"
  with_grant_option = true
}
```

#### 1. Query Snowflake and save the output

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

> If you use SnowSight, you can click on the "Download" button and select "CSV" format.
>
> ![Download CSV button in SnowSight](./images/csv_output_download.png)

Whatever way you choose, save the output to a CSV file as `example.csv`.
The file contents should look similar to this:

```csv
"created_on","privilege","granted_on","name","granted_to","grantee_name","grant_option","granted_by"
"2025-09-02 05:42:30.602 -0700","CREATE ROLE","ACCOUNT","VG98132","ROLE","TEST_ROLE","false","USERADMIN"
"2025-09-02 05:41:59.399 -0700","CREATE DATABASE","ACCOUNT","VG98132","ROLE","OTHER_TEST_ROLE","true","SYSADMIN"
```

#### 2. Generate resources and import statements based on the Snowflake output

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

> **It's user's responsibility to ensure that the output is written securely to a safe location and not to overwrite any important files.**

Whichever way you choose, the final configuration should look similar to this:

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=0.92.0"
    }
  }
}

provider "snowflake" {
  # Your configuration options
}

resource "snowflake_role" "test_role" {
  name = "TEST_ROLE"
}

resource "snowflake_role" "other_test_role" {
  name = "OTHER_TEST_ROLE"
}

resource "snowflake_account_grant" "grant" {
  roles             = [snowflake_role.test_role.name]
  privilege         = "CREATE ROLE"
  with_grant_option = false
}

resource "snowflake_account_grant" "other_grant" {
  roles             = [snowflake_role.other_test_role.name]
  privilege         = "CREATE DATABASE"
  with_grant_option = true
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option" {
  # Adjused manually due to planned changes, because of old quotes handling in identifiers (fixed in later versions of the provider)
  account_role_name = "\"TEST_ROLE\""
  on_account = true
  privileges = ["CREATE ROLE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_OTHER_TEST_ROLE_with_grant_option" {
  # Adjused manually due to planned changes, because of old quotes handling in identifiers (fixed in later versions of the provider)
  account_role_name = "\"OTHER_TEST_ROLE\""
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = true
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option
  id = "\"TEST_ROLE\"|false|false|CREATE ROLE|OnAccount"
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_OTHER_TEST_ROLE_with_grant_option
  id = "\"OTHER_TEST_ROLE\"|true|false|CREATE DATABASE|OnAccount"
}
```

#### 3. Importing auto-generated resources to the state

Now that we have the generated resources and import blocks, we can proceed to replace the deprecated resources with the new ones.
First, we need to import the new resources into the state. Run `terraform plan`, to see if there are any issues with the configuration.
The plan should output:

```
Plan: 2 to import, 0 to add, 0 to change, 0 to destroy.
```

Now we are safe to apply the changes with `terraform apply`. You should see the output similar to:

```
Apply complete! Resources: 2 imported, 0 added, 0 changed, 0 destroyed.
```

which means that the resources have been successfully imported into the state.
Remember that, if you chose to use the import block approach, [after importing you can remove the import blocks from the configuration file](https://developer.hashicorp.com/terraform/language/import#plan-and-apply-an-impor).

#### 4. Removing old grant resources

Now that we have the new resources imported into the state, we can proceed to remove the old deprecated resources from the configuration.
To do that, you can either use the [removed blocks](https://support.hashicorp.com/hc/en-us/articles/33229234219411-Terraform-Enterprise-on-Replicated-March-2025-Final-Release) or [state commands](https://developer.hashicorp.com/terraform/cli/commands/state).
In this case, we will proceed with running the following state commands:

```shell
terraform state rm snowflake_account_grant.grant
terraform state rm snowflake_account_grant.other_grant
```

Now, you are safe to remove the old resources from the configuration file.
To confirm that everything is working as expected, run `terraform plan` again.
This time it should output:

```
No changes. Your infrastructure matches the configuration.
```

#### 5. Update generated resources

The last step is optional, but highly recommended. The generated resources have generic names, which are not very user-friendly,
and they do not depend on the existing role resources which they refer to in their configuration.

To rename the resources, you can use the [terraform state mv](https://developer.hashicorp.com/terraform/cli/commands/state/mv) command or [moved block](https://developer.hashicorp.com/terraform/language/moved).

To properly link the resources, you either use the explicit dependency with `depends_on` argument,
or you can use implicit dependency by referring to the existing role resources in the `snowflake_grant_privileges_to_account_role`.

There's also an alternative, which is changing both of those things before importing the resources into the state.
That way, the resources would be imported with the desired names and dependencies from the start.

An example of the updated configuration is shown below:

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=0.92.0"
    }
  }
}

provider "snowflake" {
  # Your configuration options
}

resource "snowflake_role" "test_role" {
  name = "TEST_ROLE"
}

resource "snowflake_role" "other_test_role" {
  name = "OTHER_TEST_ROLE"
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option" {
  # We still need to keep the identifier in quotes, because of old quotes handling in identifiers (fixed in later versions of the provider)
  account_role_name = "\"${snowflake_role.test_role.name}\""
  on_account = true
  privileges = ["CREATE ROLE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_OTHER_TEST_ROLE_with_grant_option" {
  # We still need to keep the identifier in quotes, because of old quotes handling in identifiers (fixed in later versions of the provider)
  account_role_name = "\"${snowflake_role.other_test_role.name}\""
  on_account = true
  privileges = ["CREATE DATABASE"]
  with_grant_option = true
}
```

By following the above steps, you can migrate other deprecated Snowflake resources and start managing the new ones!

### Use case: Migrate existing grants to Terraform

In this case, we will focus on importing existing grants into Terraform state to start managing them.

#### 1. Query Snowflake and save the output

To have some data to work with, we will create two sample roles and grant them some privileges on the account.
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

> If you use SnowSight, you can click on the "Download" button and select "CSV" format.
>
> ![Download CSV button in SnowSight](./images/csv_output_download.png)

Whatever way you choose, save the output to a CSV file as `example.csv`.
The file contents should look similar to this:

```csv
"created_on","privilege","granted_on","name","granted_to","grantee_name","grant_option","granted_by"
"2025-08-29 06:23:25.920 -0700","CREATE DATABASE","ACCOUNT","IYA62698","ROLE","TEST_ROLE","false","SYSADMIN"
"2025-08-29 06:23:27.253 -0700","CREATE ROLE","ACCOUNT","IYA62698","ROLE","TEST_ROLE","false","USERADMIN"
"2025-08-29 06:23:29.902 -0700","CREATE DATABASE","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","SYSADMIN"
"2025-08-29 06:23:31.390 -0700","CREATE ROLE","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","USERADMIN"
"2025-08-29 06:23:32.747 -0700","CREATE USER","ACCOUNT","IYA62698","ROLE","TEST_OTHER_ROLE","false","USERADMIN"
```

#### 2. Generate resources and import statements based on the Snowflake output

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

> **It's user's responsibility to ensure that the output is written securely to a safe location and not to overwrite any important files.**

#### 3. Get the generated resources and import them to the state

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
  on_account        = true
  privileges = ["CREATE DATABASE", "CREATE ROLE", "CREATE USER"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_TEST_ROLE_without_grant_option" {
  account_role_name = "TEST_ROLE"
  on_account        = true
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

Then run `terraform plan`. You should see the following output:

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
Remember that, if you chose to use the import block approach, [after importing you can remove the import blocks from the configuration file](https://developer.hashicorp.com/terraform/language/import#plan-and-apply-an-impor).

By following the above steps, you can migrate other existing Snowflake objects into Terraform and start managing them!

### Multiple sources
Some Snowflake objects (like schemas) have fields returned from more than one SQL command. That's why simply using one `SHOW ...` output will not work. Fields from `DESCRIBE` or `SHOW PARAMETERS` must be also processed.
But outputs from all of these commands must be mapped to the input CSV value of the migration script. To make this easy, we can use a data source output for a given object, which already has the logic for mapping multiple
SQL queries and returning all necessary information.

In general, what we need to do is:
1. Define a data source for the objects you want to import.
1. Use HCL (Terraform's configuration language) to transform the data: merge `show_output` with flattened `parameters` for each object.
1. Write the transformed data to a CSV file using the `local_file` resource.
1. Run the migration script with the generated CSV file.

> **Note:** It's recommended to create a fresh Terraform environment (e.g., a new local directory with a clean state) for this data extraction process to avoid collisions with any existing data sources or state in your main Terraform workspace.

Now, let's look into more details.

As an example, let's import all schemas in a given database. First, we need to define a data source for schemas and use Terraform's HCL to transform the data into CSV format:

```terraform
terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

data "snowflake_schemas" "test" {
  in {
    database = "DATABASE"
  }
}

locals {
  # Transform each schema by merging show_output, describe_output, and flattened parameters
  schemas_flattened = [
    for schema in data.snowflake_schemas.test.schemas : merge(
      schema.show_output[0],
      # Include describe output fields (if describe_output is present)
      length(schema.describe_output) > 0 ? schema.describe_output[0] : {},
      # Flatten parameters: convert each parameter to {param_name}_value and {param_name}_level
      {
        for param_key, param_values in schema.parameters[0] :
        "${param_key}_value" => param_values[0].value
      },
      {
        for param_key, param_values in schema.parameters[0] :
        "${param_key}_level" => param_values[0].level
      }
    )
  ]

  # Get all unique keys from the first schema to create CSV header
  csv_header = join(",", [for key in keys(local.schemas_flattened[0]) : "\"${key}\""])

  # Convert each schema object to CSV row (properly escape quotes and newlines for CSV format)
  csv_escape = {
    for schema in local.schemas_flattened :
    schema.name => {
      for key in keys(local.schemas_flattened[0]) :
      key => replace(
        replace(
          replace(tostring(lookup(schema, key, "")), "\\", "\\\\"),
          "\n", "\\n"
        ),
        "\"", "\"\""
      )
    }
  }

  # Convert each schema object to CSV row
  csv_rows = [
    for schema in local.schemas_flattened :
      join(",", [
        for key in keys(local.schemas_flattened[0]) :
        "\"${local.csv_escape[schema.name][key]}\""
      ])
  ]

  # Combine header and rows
  csv_content = join("\n", concat([local.csv_header], local.csv_rows))
}

resource "local_file" "schemas_csv" {
  content  = local.csv_content
  filename = "${path.module}/schemas.csv"
}
```

After running `terraform apply`, the CSV file will be automatically generated at `schemas.csv`. Now, we can run the migration script like:
```shell
go run github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/migration_script@main -import=block schemas < ./schemas.csv
```

This will output the generated configuration and import blocks for the specified schemas.

## Limitations

### Generated resource names

The resource name generation does not guarantee uniqueness. It bases its parts on the object's state identifier, which is unique,
but because of the [Terraform limitations of the characters in the resource names](https://developer.hashicorp.com/terraform/language/resources/syntax#resource-syntax),
disallowed characters are not included in the name. This means that `!test!` and `@test@` would both be converted to `test`, leading to a name collision.
If you encounter such a situation, you will need to manually rename the resources to ensure uniqueness before proceeding with resource importing.

The exception to this rule are dots which are not removed; they are replaced with underscores, so `test.name` would become `test_name`.
This is only to ensure clarity in the generated names that contain identifiers which are separated by dots, e.g.,
instead of removing the dots in `DATABASE.SCHEMA` (resulting in `DATABASESCHEMA`), we transfer them to `DATABASE_SCHEMA`.

### No dependencies handling

The script does not handle dependencies between resources. If the generated resources depend on other resources,
you will need to manually add the necessary dependencies using `depends_on` argument or implicit dependencies
by referring to the existing resources in the generated resource configuration. It's important to ensure
that all dependent resources are linked to avoid common issues like race conditions (e.g., creating a table on schema that does not exist yet).
To learn more about dependencies, check out the official [Terraform documentation](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies).
