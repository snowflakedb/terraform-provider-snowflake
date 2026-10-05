---
page_title: "Getting Started with Openflow"
subcategory: ""
description: |-

---

# Getting started with Openflow

[Snowflake Openflow](https://docs.snowflake.com/en/user-guide/data-integration/openflow/about) is an integration service that connects any data source and any destination with hundreds of processors supporting structured and unstructured text, images, audio, video, and sensor data. This guide shows how to stand up a working Openflow Connector on a Snowflake-managed Deployment with Terraform, including the objects it needs to reach a source and land data into Snowflake.

There are existing guides that this one assumes you are familiar with:

- The [Openflow gen 2 quickstart](https://docs.snowflake.com/en/user-guide/data-integration/openflow/gen2/quickstart) explains what the Openflow gen 2 objects are and walks through the same setup using Snowsight and SQL. Read it first if Openflow gen 2 itself is new to you; this guide assumes you know what a deployment, a runtime and a connector are.
- Each resource's own documentation lists every field and its behavior: [snowflake_openflow_deployment_snowflake_managed](../resources/openflow_deployment_snowflake_managed), [snowflake_openflow_deployment_byoc](../resources/openflow_deployment_byoc), [snowflake_openflow_runtime](../resources/openflow_runtime), [snowflake_openflow_connector](../resources/openflow_connector).
- This guide will focus on creating a Snowflake-managed deployment using [snowflake_openflow_deployment_snowflake_managed](../resources/openflow_deployment_snowflake_managed).

!> **Caution: Preview Feature** At the time of writing this guide, all Openflow resources and data sources are preview features. Breaking changes are expected, even without bumping the major version. If they are still preview features, each one has to be enabled by name in `preview_features_enabled`, as shown below.

## Table of contents

* [How the objects fit together](#how-the-objects-fit-together)
* [Before you start](#before-you-start)
* [Provider configuration](#provider-configuration)
* [Step 1: where the data lands](#step-1-where-the-data-lands)
* [Step 2: reaching the source](#step-2-reaching-the-source)
* [Step 3: the deployment](#step-3-the-deployment)
* [Step 4: the runtime](#step-4-the-runtime)
* [Step 5: grants for the runtime's role](#step-5-grants-for-the-runtimes-role)
* [Step 6: the connector](#step-6-the-connector)
* [After the apply](#after-the-apply)
* [What Terraform does not manage](#what-terraform-does-not-manage)
* [Things worth knowing](#things-worth-knowing)
* [Troubleshooting](#troubleshooting)

## How the objects fit together

Three Openflow objects nest inside one another:

| Object | Scope | Holds |
|---|---|---|
| Deployment | account | the compute environment runtimes run in |
| Runtime | schema | the nodes connectors execute on, and the egress they are allowed |
| Connector | schema | one data flow, from a source into Snowflake |

A runtime names the deployment it runs in. A connector names the runtime it runs in, and is created in that runtime's schema, so **a connector's `database` and `schema` must match its runtime's**.

Around them sit ordinary Snowflake objects, all of which have stable resources in this provider:

| Resource | Why Openflow needs it |
|---|---|
| [snowflake_database](../resources/database), [snowflake_schema](../resources/schema) | the runtime and its connectors live here |
| [snowflake_warehouse](../resources/warehouse) | connectors merge changes into the destination with it |
| [snowflake_secret_with_generic_string](../resources/secret_with_generic_string) | holds source credentials, so they never appear in a connector's configuration file |
| [snowflake_network_rule](../resources/network_rule) + [snowflake_external_access_integration](../resources/external_access_integration) | Snowflake blocks outbound traffic by default; without these the runtime cannot reach the source |
| [snowflake_stage_internal](../resources/stage_internal) or [snowflake_git_repository](../resources/git_repository) | holds a connector's configuration bundle |
| [snowflake_grant_privileges_to_account_role](../resources/grant_privileges_to_account_role) | the role the runtime executes as needs access to the secret, the destination and the warehouse |

## Before you start

You need a role that can create Openflow objects. The [Snowflake Deployment setup guide](https://docs.snowflake.com/en/user-guide/data-integration/openflow/setup-openflow-spcs) has you create one with `CREATE ROLE IF NOT EXISTS openflow_admin`.

-> **Note** Snowflake folds unquoted identifiers to upper case, so `create role openflow_admin` stores the role as `OPENFLOW_ADMIN`. The provider quotes identifiers it is given, so passing `openflow_admin` to `execute_as_role` fails with `Role openflow_admin does not exist or is not accessible`. Write role names in the case Snowflake stored them.

**Account-level privileges have to exist before the first apply.** If Terraform connects as the Openflow role, it cannot grant these to itself, so run them once as a role that can grant on the account:

```sql
USE ROLE ACCOUNTADMIN;

GRANT CREATE DATABASE                    ON ACCOUNT TO ROLE OPENFLOW_ADMIN;
GRANT CREATE WAREHOUSE                   ON ACCOUNT TO ROLE OPENFLOW_ADMIN;
GRANT CREATE EXTERNAL ACCESS INTEGRATION ON ACCOUNT TO ROLE OPENFLOW_ADMIN;
GRANT CREATE OPENFLOW DEPLOYMENT         ON ACCOUNT TO ROLE OPENFLOW_ADMIN;
GRANT CREATE COMPUTE POOL                ON ACCOUNT TO ROLE OPENFLOW_ADMIN;
```

`CREATE OPENFLOW RUNTIME`, `CREATE OPENFLOW CONNECTOR`, `CREATE SECRET`, and `CREATE NETWORK RULE` are privileges on a schema rather than on the account. These are the ones this guide needs, not every schema privilege. If Terraform creates the schema, the role owns it and they come with ownership. If the schema already exists and is owned by someone else, grant them the same way, on the schema.

These grants can also be managed in Terraform with [snowflake_grant_privileges_to_account_role](../resources/grant_privileges_to_account_role) — the `privileges` list is passed through to SQL, so Openflow privileges can be granted the same way as any other privilege. Doing so needs a [second provider configuration aliased](https://developer.hashicorp.com/terraform/language/block/provider#alias) to a role that can grant, since the granting role and the Openflow role are usually different.

## Provider configuration

Openflow resources are preview features, so list the ones you use. `snowflake_external_access_integration` is a preview feature too; the database, schema, warehouse, secret, network rule, stage, grant and `snowflake_execute` resources are stable and need no entry.

```terraform
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "~> 2.21"
    }
    # only needed if you render a connector's config.json, as shown in step 6
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

provider "snowflake" {
  preview_features_enabled = [
    "snowflake_openflow_deployment_snowflake_managed_resource",
    "snowflake_openflow_runtime_resource",
    "snowflake_openflow_connector_resource",
    "snowflake_external_access_integration_resource",
    "snowflake_openflow_deployments_datasource",
    "snowflake_openflow_connector_definitions_datasource",
  ]
}
```

Authentication is configured the same way as for any other use of this provider; see the [authentication methods guide](./authentication_methods).

## Step 1: where the data lands

Nothing here is Openflow-specific, but the schema matters: the runtime and its connectors have to share one.

```terraform
resource "snowflake_database" "openflow" {
  name    = "OPENFLOW_DB"
  comment = "Holds the runtime and its connectors."
}

resource "snowflake_schema" "openflow" {
  database = snowflake_database.openflow.name
  name     = "OPENFLOW_SCH"
}

resource "snowflake_warehouse" "openflow" {
  name           = "OPENFLOW_WH"
  warehouse_size = "XSMALL"
  auto_suspend   = 60
}
```

## Step 2: reaching the source

A connector needs two things to talk to a source outside Snowflake: credentials it can resolve at run time, and permission to leave the account.

```terraform
resource "snowflake_secret_with_generic_string" "source_password" {
  database      = snowflake_database.openflow.name
  schema        = snowflake_schema.openflow.name
  name          = "SOURCE_PASSWORD"
  secret_string = var.source_password
}

resource "snowflake_network_rule" "network_rule_1" {
  database   = snowflake_database.openflow.name
  schema     = snowflake_schema.openflow.name
  name       = "NETWORK_RULE_1"
  type       = "HOST_PORT"
  mode       = "EGRESS"
  value_list = ["${var.source_host}:${var.source_port}"]
}

resource "snowflake_network_rule" "network_rule_2" {
  database   = snowflake_database.openflow.name
  schema     = snowflake_schema.openflow.name
  name       = "NETWORK_RULE_2"
  type       = "HOST_PORT"
  mode       = "EGRESS"
  value_list = ["${var.network_rule_2_host}:${var.source_port}"]
}

resource "snowflake_external_access_integration" "source" {
  name    = "OPENFLOW_EAI"
  enabled = true
  allowed_network_rules = [
    snowflake_network_rule.network_rule_1.fully_qualified_name,
    snowflake_network_rule.network_rule_2.fully_qualified_name,
  ]

  # without this the runtime can reach the host but cannot read the password
  allowed_authentication_secrets {
    secrets = [snowflake_secret_with_generic_string.source_password.fully_qualified_name]
  }
}
```

The connector's configuration references the secret by fully qualified name, so the password itself never ends up in a configuration file or on a stage. Keep it out of the state file's plain sight as well; see [protecting secret values](./authentication_methods#protecting-secret-values).

## Step 3: the deployment

A deployment is the compute environment runtimes run in. This walkthrough uses a [Snowflake-managed deployment](../resources/openflow_deployment_snowflake_managed), which runs on Snowflake-managed compute.

A [BYOC deployment](../resources/openflow_deployment_byoc) runs in your own cloud account instead. Terraform can take it as far as `INACTIVE`; it becomes `ACTIVE` only after you apply the deployment's CloudFormation template in your cloud account, and that template is only downloadable from the Openflow UI. The rest of this guide assumes a Snowflake-managed deployment.

```terraform
resource "snowflake_openflow_deployment_snowflake_managed" "main" {
  name    = "OPENFLOW_DEPLOYMENT"
  comment = "Managed by Terraform."

  # provisioning a deployment is the slowest step in an Openflow configuration
  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}
```

Creating one takes a few minutes. If you already have a deployment, reuse it rather than creating another: read it with the [snowflake_openflow_deployments](../data-sources/openflow_deployments) data source and pass its name to the runtime.

```terraform
data "snowflake_openflow_deployments" "existing" {
  like = "EXISTING_DEPLOYMENT_NAME"
}

locals {
  deployment_name = data.snowflake_openflow_deployments.existing.openflow_deployments[0].show_output[0].name
}
```

A deployment read through a data source is never managed, so `terraform destroy` leaves it alone — which is usually what you want for a deployment shared by several configurations.

## Step 4: the runtime

The runtime is the compute connectors execute on. It carries the external access integration, so every connector inside it inherits that egress, and it names the role connectors execute as.

```terraform
resource "snowflake_openflow_runtime" "main" {
  database        = snowflake_database.openflow.name
  schema          = snowflake_schema.openflow.name
  name            = "OPENFLOW_RUNTIME"
  deployment      = snowflake_openflow_deployment_snowflake_managed.main.name
  node_type       = "MEDIUM"
  min_nodes       = 1
  max_nodes       = 1
  execute_as_role = "OPENFLOW_RUNTIME_EXECUTE_AS_RL"

  external_access_integrations = [snowflake_external_access_integration.source.name]

  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}
```

`deployment` and `node_type` are create-only: Snowflake has no `ALTER` for them, so changing either replaces the runtime, and replacing a runtime destroys the connectors in it. `min_nodes` and `max_nodes` can be changed in place.

The runtime's Openflow canvas URL is in `describe_output[0].server_url`. Surface it with an [output](https://developer.hashicorp.com/terraform/language/values/outputs); the outputs used in this guide are collected in [After the apply](#after-the-apply).

## Step 5: grants for the runtime's role

The connector runs as the runtime's `execute_as_role`, not as the role Terraform connects with. Without these grants the connector is created but fails as soon as it tries to read its secret or write its destination.

```terraform
resource "snowflake_grant_privileges_to_account_role" "secret" {
  account_role_name = "OPENFLOW_RUNTIME_EXECUTE_AS_RL"
  privileges        = ["USAGE"]

  on_schema_object {
    object_type = "SECRET"
    object_name = snowflake_secret_with_generic_string.source_password.fully_qualified_name
  }
}

resource "snowflake_grant_privileges_to_account_role" "destination" {
  account_role_name = "OPENFLOW_RUNTIME_EXECUTE_AS_RL"
  privileges        = ["USAGE", "CREATE SCHEMA"]

  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.openflow.name
  }
}

resource "snowflake_grant_privileges_to_account_role" "warehouse" {
  account_role_name = "OPENFLOW_RUNTIME_EXECUTE_AS_RL"
  privileges        = ["USAGE", "OPERATE"]

  on_account_object {
    object_type = "WAREHOUSE"
    object_name = snowflake_warehouse.openflow.name
  }
}
```

## Step 6: the connector

A connector is the one Openflow object whose configuration is not expressed as SQL arguments. `CREATE OPENFLOW CONNECTOR` either names a Snowflake-managed definition or reads a configuration bundle from a stage. The `from` block covers both, and it is create-only: changing it replaces the connector.

### From a definition

The simplest form, and a draft: the connector is created with no configuration, settles on `STOPPED`, and stays there until a version is committed.

```terraform
resource "snowflake_openflow_connector" "draft" {
  database = snowflake_database.openflow.name
  schema   = snowflake_schema.openflow.name
  name     = "OPENFLOW_CONNECTOR"
  runtime  = snowflake_openflow_runtime.main.fully_qualified_name

  from {
    definition = "OPENFLOW_POSTGRES_CDC"
  }
}
```

List the available definition IDs with the [snowflake_openflow_connector_definitions](../data-sources/openflow_connector_definitions) data source.

### From a configuration bundle

To have a connector arrive fully configured and startable, point it at a stage path holding a configuration bundle. A bundle is a `config.json` plus every file that file references — for a Postgres CDC connector, that means the Postgres JDBC driver jar.

The examples below use `OPENFLOW_POSTGRES_CDC`. Other connectors differ only in which properties `config.json` carries and which assets it references; the Terraform mechanics are identical.

#### What goes in config.json

The property names and allowed values for a given connector are documented with that connector, and the [gen 2 quickstart](https://docs.snowflake.com/en/user-guide/data-integration/openflow/gen2/quickstart) walks through filling them in. The part that matters for Terraform is the three ways a property carries a value:

- `STRING_LITERAL` — a plain value, such as a JDBC URL or the destination database name.
- `SECRET_REFERENCE` — the **fully qualified name** of a Snowflake secret. The connector resolves it at run time, so no credential is ever written into the bundle.
- `ASSET_REFERENCE` — a **bare file name** of another file in the same bundle, such as a driver jar.

An abbreviated `config.json.tftpl` for Postgres CDC, showing all three:

```json
{
  "configFormatVersion": 1,
  "connectorDefinitionId": "OPENFLOW_POSTGRES_CDC",
  "configuration": [
    {
      "name": "Source",
      "properties": {
        "Source Database Connection URL": {
          "valueType": "STRING_LITERAL",
          "value": "jdbc:postgresql://${source_host}:${source_port}/${source_database}?sslmode=require"
        },
        "Source Database Password": {
          "valueType": "SECRET_REFERENCE",
          "fullyQualifiedSecretName": "${secret_fully_qualified_name}"
        },
        "Source Database Driver": {
          "valueType": "ASSET_REFERENCE",
          "assetIds": ["${jdbc_driver_file}"]
        }
      }
    },
    {
      "name": "Destination details",
      "properties": {
        "Snowflake Destination Database": {
          "valueType": "STRING_LITERAL",
          "value": "${destination_database}"
        },
        "Snowflake Warehouse": {
          "valueType": "STRING_LITERAL",
          "value": "${destination_warehouse}"
        }
      }
    }
  ]
}
```

#### Rendering it

Render the file rather than committing a finished one. Its values name objects Terraform creates — the secret's fully qualified name, the destination database, the warehouse — and rendering is what keeps the file and those objects from drifting apart when a name changes.

```terraform
locals {
  jdbc_driver_file = "postgresql-42.7.7.jar"
  jdbc_driver_path = abspath("${path.module}/assets/${local.jdbc_driver_file}")
  bundle_path      = "example"
}

resource "local_file" "config" {
  filename = "${path.module}/generated/config.json"

  content = templatefile("${path.module}/templates/config.json.tftpl", {
    source_host                 = var.source_host
    source_port                 = var.source_port
    source_database             = var.source_database
    secret_fully_qualified_name = snowflake_secret_with_generic_string.source_password.fully_qualified_name
    jdbc_driver_file            = local.jdbc_driver_file
    destination_database        = snowflake_database.openflow.name
    destination_warehouse       = snowflake_warehouse.openflow.name
  })
}
```

#### Uploading the bundle

There is no resource for putting a file on a stage, so use [snowflake_execute](../resources/execute) to run `PUT`. Every file in the bundle goes to the **same stage path**, `config.json` and assets alike.

```terraform
resource "snowflake_stage_internal" "bundles" {
  database = snowflake_database.openflow.name
  schema   = snowflake_schema.openflow.name
  name     = "CONNECTOR_BUNDLES"
}

resource "snowflake_execute" "upload_config" {
  execute = "PUT file://${local_file.config.filename} @${snowflake_stage_internal.bundles.fully_qualified_name}/${local.bundle_path}/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE"
  revert  = "REMOVE @${snowflake_stage_internal.bundles.fully_qualified_name}/${local.bundle_path}/config.json"
}

resource "snowflake_execute" "upload_driver" {
  execute = "PUT file://${local.jdbc_driver_path} @${snowflake_stage_internal.bundles.fully_qualified_name}/${local.bundle_path}/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE"
  revert  = "REMOVE @${snowflake_stage_internal.bundles.fully_qualified_name}/${local.bundle_path}/${local.jdbc_driver_file}"

  lifecycle {
    precondition {
      condition     = fileexists(local.jdbc_driver_path)
      error_message = "Download the Postgres JDBC driver into assets/ before applying."
    }
  }
}
```

`AUTO_COMPRESS = FALSE` is not optional: with the default, Snowflake gzips the file and the connector finds `config.json.gz` instead of the file it is looking for. `OVERWRITE = TRUE` makes re-uploads idempotent. `revert` is what `terraform destroy` runs, so the stage is left as it was found.

Assets are referenced by bare file name, which is why they must share a path with `config.json`. Downloading a driver jar is outside Terraform's remit — fetch it in a setup script, or vendor it — but a `precondition` turns a missing file into a clear plan-time error rather than a confusing failure during create.

-> **Note** `PUT` uploads from the machine running Terraform, so every local file in the bundle has to exist there. On a CI runner without them, use a [git repository stage](#from-a-git-repository) instead.

#### Creating the connector from it

```terraform
resource "snowflake_openflow_connector" "example" {
  database = snowflake_database.openflow.name
  schema   = snowflake_schema.openflow.name
  name     = "OPENFLOW_CONNECTOR"
  runtime  = snowflake_openflow_runtime.main.fully_qualified_name

  from {
    stage = snowflake_stage_internal.bundles.fully_qualified_name
    path  = local.bundle_path
  }

  timeouts {
    create = "45m"
    update = "45m"
    delete = "45m"
  }

  # the whole bundle has to be on the stage, and the grants in place, before the connector reads it
  depends_on = [
    snowflake_execute.upload_config,
    snowflake_execute.upload_driver,
    snowflake_grant_privileges_to_account_role.secret,
    snowflake_grant_privileges_to_account_role.destination,
    snowflake_grant_privileges_to_account_role.warehouse,
  ]
}
```

The connector reads the bundle **once, at create time**. Changing `config.json` afterwards and re-applying uploads the new file but leaves the connector as it was: `from` is create-only, and the contents behind it are not tracked. To pick up a changed configuration, either replace the connector (`terraform apply -replace=snowflake_openflow_connector.example`) or commit a new version.

Confirm Snowflake actually read the file from the `connector_definition` [output](https://developer.hashicorp.com/terraform/language/values/outputs) after the apply — it reports the definition Snowflake resolved, so `OPENFLOW_POSTGRES_CDC` there means your `config.json` was found and parsed rather than silently skipped.

### From a git repository

A [git repository](../resources/git_repository) is a stage as far as Snowflake is concerned, so a version-controlled bundle needs no upload step at all. This is the better route once a configuration is past the experimenting stage.

```terraform
resource "snowflake_openflow_connector" "example" {
  database = snowflake_database.openflow.name
  schema   = snowflake_schema.openflow.name
  name     = "OPENFLOW_CONNECTOR"
  runtime  = snowflake_openflow_runtime.main.fully_qualified_name

  from {
    stage = snowflake_git_repository.bundles.fully_qualified_name
    path  = "branches/main/connectors/example"
  }

  depends_on = [
    snowflake_grant_privileges_to_account_role.secret,
    snowflake_grant_privileges_to_account_role.destination,
    snowflake_grant_privileges_to_account_role.warehouse,
  ]
}
```

## After the apply

A first apply that creates its own deployment takes minutes; reusing an existing deployment is faster. Every Openflow statement is asynchronous and the provider waits for the object to settle, so a long apply is expected rather than a hang.

A connector created from a bundle is configured but **stopped**. Surface Snowflake-reported values with [outputs](https://developer.hashicorp.com/terraform/language/values/outputs):

```terraform
output "runtime_url" {
  value = snowflake_openflow_runtime.main.describe_output[0].server_url
}

output "connector_status" {
  value = snowflake_openflow_connector.example.show_output[0].status
}

output "connector_url" {
  value = snowflake_openflow_connector.example.show_output[0].connector_url
}

output "connector_definition" {
  value = snowflake_openflow_connector.example.show_output[0].connector_definition
}
```

After the apply, read them with `terraform output`:

```
terraform output runtime_url
terraform output connector_status
terraform output connector_url
terraform output connector_definition
```

`connector_definition` is what Snowflake resolved from your bundle. Open `connector_url` to start the connector in the Openflow canvas.

## What Terraform does not manage

Openflow has a large operational surface that is deliberately not modeled as desired state:

- **Starting and stopping connectors.** Whether a connector is running is an operational decision, not configuration. The current status is reported in `show_output`.
- **Connector version management.** Committing a configuration version, upgrading a connector, and promoting versions are not modeled here.
- **The BYOC cloud infrastructure.** If you later use a BYOC deployment, the CloudFormation template is downloaded from the Openflow UI and applied in your own cloud account. That path is outside this walkthrough.

## Things worth knowing

**Applies are slow and asynchronous.** Set generous `timeouts` blocks on all three Openflow resources. The defaults are sized for typical cases, and provisioning that exceeds them fails an apply partway.

**Several fields are create-only.** A runtime's `deployment` and `node_type`, and a connector's `from` and `runtime`, have no `ALTER` in Snowflake, so changing them replaces the object. Replacing a deployment destroys every runtime and connector in it; read a plan carefully before applying one that shows a replacement.

**`from` is not read back.** Snowflake resolves a definition for connectors created from a stage too, so a value read from `SHOW` could not be distinguished from a configured one. External changes to `from` are therefore not detected, and after a `terraform import` the block is absent from state, so the first plan asks to replace the connector.

**Destroy is slow too**, for the same reason as create: a connector has to terminate before it can be dropped, and a deployment is again the longest part.

**Import** works for all three objects with the syntax shown on each resource page — useful for bringing a deployment created in Snowsight under management, or recovering after an apply times out mid-create.

## Troubleshooting

| Symptom | Cause |
|---|---|
| `Role openflow_admin does not exist or is not accessible` | identifier case: Snowflake stored the role as `OPENFLOW_ADMIN` and the provider quotes what you pass |
| `003001 (42501): Insufficient privileges to operate on account` | one of the account-level grants in [Before you start](#before-you-start) is missing; the message names the privilege |
| Connector created but immediately fails or reports `START_FAILED` | usually a missing grant to `execute_as_role`, or a file the configuration references missing from the stage path |
| Connector stays `STOPPED` and the canvas shows no configuration | it was created `from { definition = ... }`, so it is a draft until a version is committed |
| `show_output.connector_definition` is not what you expected | Snowflake did not read your `config.json`; check it is at the exact `path` on the stage |
| Apply fails with a timeout | raise the `timeouts` block, then import the object if it was created before the timeout fired |

## Further reading

- [Openflow gen 2 quickstart](https://docs.snowflake.com/en/user-guide/data-integration/openflow/gen2/quickstart)
- [Set up Openflow - Snowflake Deployment](https://docs.snowflake.com/en/user-guide/data-integration/openflow/setup-openflow-spcs)
- [Set up Openflow - BYOC](https://docs.snowflake.com/en/user-guide/data-integration/openflow/setup-openflow)
- [CREATE OPENFLOW DEPLOYMENT](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-deployment), [CREATE OPENFLOW RUNTIME](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-runtime), [CREATE OPENFLOW CONNECTOR](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector)
- [snowflake_openflow_connectors](../data-sources/openflow_connectors), [snowflake_openflow_runtimes](../data-sources/openflow_runtimes), [snowflake_openflow_deployments](../data-sources/openflow_deployments), [snowflake_openflow_connector_definitions](../data-sources/openflow_connector_definitions) data sources
