---
page_title: "snowflake_grant_ownership Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  
---

~> **Note** For more details about granting ownership, please visit [`GRANT OWNERSHIP` Snowflake documentation page](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership), and our [grant ownership resource overview](../guides/grant_ownership_resource_overview).

~> **Note** Manage grants on `HYBRID TABLE` by specifying `TABLE` or `TABLES` in `object_type` field. This applies to a single object, all objects, or future objects. This reflects the current behavior in Snowflake.

~> **Note** Changes introduced in [2025_01 Snowflake BCR Bundle](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_01/bcr-1870) alter how future ownership grant hierarchy works. After the change schema-level future ownership grant(s) will override those defined on the database-level. Previously, the database-level ones took the precedence.

~> **Warning** Be careful when using this resource on managed schema or any object within it, as in Snowflake managed schemas follow different ownership and privilege model. For more details, please visit [Schema documentation](https://docs.snowflake.com/en/sql-reference/sql/create-schema#optional-parameters) (scroll to `WITH MANAGED ACCESS` option).

!> **Warning** Grant ownership resource still has some limitations. Delete operation is not implemented for on_future grants (you have to remove the config and then revoke ownership grant on future X manually).

# snowflake_grant_ownership (Resource)



## Example Usage

For more examples, head over to our usage guide where we present how to use the grant_ownership resource in [common use cases](../guides/grant_ownership_common_use_cases).

```terraform
##################################
### on object to account role
##################################

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name   = snowflake_role.test.name
  outbound_privileges = "COPY"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on object to database role
##################################

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_database_role" "test" {
  name     = "test_database_role"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  database_role_name  = snowflake_database_role.test.fully_qualified_name
  outbound_privileges = "REVOKE"
  on {
    object_type = "SCHEMA"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
  }
}

##################################
### on all tables in database to account role
##################################

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.test.name
    }
  }
}

##################################
### on all tables in schema to account role
##################################

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    all {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

##################################
### on future tables in database to account role
##################################

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.test.name
    }
  }
}

##################################
### on future tables in schema to account role
##################################

resource "snowflake_account_role" "test" {
  name = "test_role"
}

resource "snowflake_database" "test" {
  name = "test_database"
}

resource "snowflake_schema" "test" {
  name     = "test_schema"
  database = snowflake_database.test.name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    future {
      object_type_plural = "TABLES"
      in_schema          = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\""
    }
  }
}

##################################
### RoleBasedAccessControl (RBAC example)
##################################

resource "snowflake_account_role" "test" {
  name = "role"
}

resource "snowflake_database" "test" {
  name = "database"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}

resource "snowflake_grant_account_role" "test" {
  role_name = snowflake_account_role.test.name
  user_name = "username"
}

provider "snowflake" {
  profile = "default"
  alias   = "secondary"
  role    = snowflake_account_role.test.name
}

## With ownership on the database, the secondary provider is able to create schema on it without any additional privileges.
resource "snowflake_schema" "test" {
  depends_on = [snowflake_grant_ownership.test, snowflake_grant_account_role.test]
  provider   = snowflake.secondary
  database   = snowflake_database.test.name
  name       = "schema"
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

## Granting ownership on pipes
To transfer ownership of a pipe, there must be additional conditions met. Otherwise, additional manual work
will be needed afterward or in some cases, the ownership won't be transferred (resulting in error).

To transfer ownership of a pipe(s) **fully automatically**, one of the following conditions has to be met:
- OPERATE and MONITOR privileges are granted to the current role on the pipe(s) and `outbound_privileges` field is set to `COPY`.
- The pipe(s) running status is paused (additional privileges and fields set are needed to pause and resume the pipe before and after ownership transfer. If it's already paused, nothing additional is needed and the pipe will remain paused after the ownership transfer).

To transfer ownership of a pipe(s) **semi-automatically** you have to:
1. Pause the pipe(s) you want to transfer ownership of (using [ALTER PIPE](https://docs.snowflake.com/en/sql-reference/sql/alter-pipe#syntax); see PIPE_EXECUTION_PAUSED).
2. Create Terraform configuration with the `snowflake_grant_ownership` resource and perform ownership transfer with the `terraform apply`.
3. To resume the pipe(s) after ownership transfer use [PIPE_FORCE_RESUME system function](https://docs.snowflake.com/en/sql-reference/functions/system_pipe_force_resume).

## Granting ownership on task
Granting ownership on single task requires:
- Either OWNERSHIP or OPERATE privilege to suspend the task (and its root)
- Role that will be granted ownership has to have USAGE granted on the warehouse assigned to the task, as well as EXECUTE TASK granted globally
- The outbound privileges set to `outbound_privileges = "COPY"` if you want to move grants automatically to the owner (also enables the provider to resume the task automatically)
If originally the first owner won't be granted with OPERATE, USAGE (on the warehouse), EXECUTE TASK (on the account), and outbound privileges won't be set to `COPY`, then you have to resume suspended tasks manually.

## Granting ownership on all tasks in database/schema
Granting ownership on all tasks requires less privileges than granting ownership on one task, because it does a little bit less and requires additional work to be done after.
The only thing you have to take care of is to resume tasks after grant ownership transfer. If all of your tasks are managed by the Snowflake Terraform Plugin, this should
be as simple as running `terraform apply` second time (assuming the currently used role is privileged enough to be able to resume the tasks).
If your tasks are not managed by the Snowflake Terraform Plugin, you should resume them yourself manually.

## Granting ownership on external tables
Transferring ownership on an external table or its parent database blocks automatic refreshes of the table metadata by setting the `AUTO_REFRESH` property to `FALSE`.
Right now, there's no way to check the `AUTO_REFRESH` state of the external table and because of that, a manual step is required after ownership transfer.
To set the `AUTO_REFRESH` property back to `TRUE` (after you transfer ownership), use the [ALTER EXTERNAL TABLE](https://docs.snowflake.com/en/sql-reference/sql/alter-external-table) command.

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `on` (Block List, Min: 1, Max: 1) Configures which object(s) should transfer their ownership to the specified role. (see [below for nested schema](#nestedblock--on))

### Optional

- `account_role_name` (String) The fully qualified name of the account role to which privileges will be granted. For more information about this resource, see [docs](./account_role).
- `database_role_name` (String) The fully qualified name of the database role to which privileges will be granted. For more information about this resource, see [docs](./database_role).
- `outbound_privileges` (String) Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role. Available options are: REVOKE for removing existing privileges and COPY to transfer them with ownership. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#optional-parameters).
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--on"></a>
### Nested Schema for `on`

Optional:

- `all` (Block List, Max: 1) Configures the privilege to be granted on all objects in either a database or schema. (see [below for nested schema](#nestedblock--on--all))
- `future` (Block List, Max: 1) Configures the privilege to be granted on all objects in either a database or schema. (see [below for nested schema](#nestedblock--on--future))
- `object_name` (String) Specifies the identifier for the object on which you are transferring ownership.
- `object_type` (String) Specifies the type of object on which you are transferring ownership. Available values are: AGGREGATION POLICY | ALERT | AUTHENTICATION POLICY | COMPUTE POOL | DATA METRIC FUNCTION | DATABASE | DATABASE ROLE | DYNAMIC TABLE | EVENT TABLE | EXTERNAL TABLE | EXTERNAL VOLUME | FAILOVER GROUP | FILE FORMAT | FUNCTION | GIT REPOSITORY | HYBRID TABLE | ICEBERG TABLE | IMAGE REPOSITORY | INTEGRATION | MATERIALIZED VIEW | NETWORK POLICY | NETWORK RULE | PACKAGES POLICY | PIPE | PROCEDURE | MASKING POLICY | PASSWORD POLICY | PROJECTION POLICY | REPLICATION GROUP | RESOURCE MONITOR | ROLE | ROW ACCESS POLICY | SCHEMA | SESSION POLICY | SECRET | SEQUENCE | STAGE | STREAM | TABLE | TAG | TASK | USER | VIEW | WAREHOUSE

<a id="nestedblock--on--all"></a>
### Nested Schema for `on.all`

Required:

- `object_type_plural` (String) Specifies the type of object in plural form on which you are transferring ownership. Available values are: AGGREGATION POLICIES | ALERTS | AUTHENTICATION POLICIES | COMPUTE POOLS | DATA METRIC FUNCTIONS | DATABASES | DATABASE ROLES | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | EXTERNAL VOLUMES | FAILOVER GROUPS | FILE FORMATS | FUNCTIONS | GIT REPOSITORIES | HYBRID TABLES | ICEBERG TABLES | IMAGE REPOSITORIES | INTEGRATIONS | MATERIALIZED VIEWS | NETWORK POLICIES | NETWORK RULES | PACKAGES POLICIES | PIPES | PROCEDURES | MASKING POLICIES | PASSWORD POLICIES | PROJECTION POLICIES | REPLICATION GROUPS | RESOURCE MONITORS | ROLES | ROW ACCESS POLICIES | SCHEMAS | SESSION POLICIES | SECRETS | SEQUENCES | STAGES | STREAMS | TABLES | TAGS | TASKS | USERS | VIEWS | WAREHOUSES. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters).

Optional:

- `in_database` (String) The fully qualified name of the database. For more information about this resource, see [docs](./database).
- `in_schema` (String) The fully qualified name of the schema. For more information about this resource, see [docs](./schema).


<a id="nestedblock--on--future"></a>
### Nested Schema for `on.future`

Required:

- `object_type_plural` (String) Specifies the type of object in plural form on which you are transferring ownership. Available values are: AGGREGATION POLICIES | ALERTS | AUTHENTICATION POLICIES | COMPUTE POOLS | DATA METRIC FUNCTIONS | DATABASES | DATABASE ROLES | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | EXTERNAL VOLUMES | FAILOVER GROUPS | FILE FORMATS | FUNCTIONS | GIT REPOSITORIES | HYBRID TABLES | ICEBERG TABLES | IMAGE REPOSITORIES | INTEGRATIONS | MATERIALIZED VIEWS | NETWORK POLICIES | NETWORK RULES | PACKAGES POLICIES | PIPES | PROCEDURES | MASKING POLICIES | PASSWORD POLICIES | PROJECTION POLICIES | REPLICATION GROUPS | RESOURCE MONITORS | ROLES | ROW ACCESS POLICIES | SCHEMAS | SESSION POLICIES | SECRETS | SEQUENCES | STAGES | STREAMS | TABLES | TAGS | TASKS | USERS | VIEWS | WAREHOUSES. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters).

Optional:

- `in_database` (String) The fully qualified name of the database. For more information about this resource, see [docs](./database).
- `in_schema` (String) The fully qualified name of the schema. For more information about this resource, see [docs](./schema).



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

## Import

~> **Note** All the ..._name parts should be fully qualified names (where every part is quoted), e.g. for schema object it is `"<database_name>"."<schema_name>"."<object_name>"`

Import is supported using the following syntax:

`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|<grant_type>|<grant_data>'`

where:
- role_type - string - type of granted role (either ToAccountRole or ToDatabaseRole)
- role_name - string - fully qualified identifier for either account role or database role (depending on the role_type)
- outbound_privileges_behavior - string - behavior specified for existing roles (can be either COPY or REVOKE)
- grant_type - enum
- grant_data - data dependent on grant_type

It has varying number of parts, depending on grant_type. All the possible types are:

### OnObject
`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|OnObject|<object_type>|<object_name>'`

### OnAll (contains inner types: InDatabase | InSchema)

#### InDatabase
`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|OnAll|<object_type_plural>|InDatabase|<database_name>'`

#### InSchema
`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|OnAll|<object_type_plural>|InSchema|<schema_name>'`

### OnFuture (contains inner types: InDatabase | InSchema)

#### InDatabase
`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|OnFuture|<object_type_plural>|InDatabase|<database_name>'`

#### InSchema
`terraform import snowflake_grant_ownership.example '<role_type>|<role_identifier>|<outbound_privileges_behavior>|OnFuture|<object_type_plural>|InSchema|<schema_name>'`

### Import examples

#### OnObject on Schema ToAccountRole
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"|COPY|OnObject|SCHEMA|"database_name"."schema_name"'`

#### OnObject on Schema ToDatabaseRole
`terraform import snowflake_grant_ownership.example 'ToDatabaseRole|"database_name"."database_role_name"|COPY|OnObject|SCHEMA|"database_name"."schema_name"'`

#### OnObject on Table
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"|COPY|OnObject|TABLE|"database_name"."schema_name"."table_name"'`

#### OnAll InDatabase
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"|REVOKE|OnAll|TABLES|InDatabase|"database_name"'`

#### OnAll InSchema
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"||OnAll|TABLES|InSchema|"database_name"."schema_name"'`

#### OnFuture InDatabase
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"||OnFuture|TABLES|InDatabase|"database_name"'`

#### OnFuture InSchema
`terraform import snowflake_grant_ownership.example 'ToAccountRole|"account_role"|COPY|OnFuture|TABLES|InSchema|"database_name"."schema_name"'`
