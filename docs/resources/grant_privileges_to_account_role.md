---
page_title: "snowflake_grant_privileges_to_account_role Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  
---


!> **Warning** Be careful when using `always_apply` field. It will always produce a plan (even when no changes were made) and can be harmful in some setups. For more details why we decided to introduce it to go our document explaining those design decisions (coming soon).

~> **Note** Manage grants on `HYBRID TABLE` by specifying `TABLE` or `TABLES` in `object_type` field. This applies to a single object, all objects, or future objects. This reflects the current behavior in Snowflake.

~> **Note** When granting privileges on applications (for example, the default "SNOWFLAKE" application) use `on_account_object.object_type = "DATABASE"` instead.

~> **Note** When using `IMPORTED PRIVILEGES` privilege, the `with_grant_option` field is not supported. Additionally, when the `IMPORTED PRIVILEGES` privilege is not set in the config, and it is granted externally, this change is not detected because of Snowflake limitations. Also, granting individual privileges on imported database is not allowed, this is a Snowflake limitation. Use `IMPORTED PRIVILEGES` instead.

~> **Note** Please, follow the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/security-access-control-considerations) for best practices on access control. The provider does not enforce any specific methodology, so it is essential for users to choose the appropriate strategy for seamless privilege management. Additionally, refer to [this link](https://docs.snowflake.com/en/user-guide/security-access-control-privileges) for a list of all available privileges in Snowflake.

# snowflake_grant_privileges_to_account_role (Resource)



## Example Usage

```terraform
resource "snowflake_database" "db" {
  name = "database"
}

resource "snowflake_schema" "my_schema" {
  database = snowflake_database.db.name
  name     = "my_schema"
}

resource "snowflake_account_role" "db_role" {
  name = "role_name"
}

##################################
### on account privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["CREATE DATABASE", "CREATE USER"]
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
}

## ID: "\"role_name\"|false|false|CREATE DATABASE,CREATE USER|OnAccount"

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnAccount"

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|true|ALL|OnAccount"

##################################
### on account object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["CREATE SCHEMA", "CREATE DATABASE ROLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|CREATE SCHEMA,CREATE DATABASE ROLE|OnAccountObject|DATABASE|\"database\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnAccountObject|DATABASE|\"database\""

# grant IMPORTED PRIVILEGES on SNOWFLAKE application
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  privileges        = ["IMPORTED PRIVILEGES"]
  on_account_object {
    object_type = "DATABASE" # All applications should be using DATABASE object_type
    object_name = "SNOWFLAKE"
  }
}

## ID: "\"role_name\"|false|false|IMPORTED PRIVILEGES|OnAccountObject|DATABASE|\"SNOWFLAKE\""

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|true|ALL|OnAccountObject|DATABASE|\"database\""

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnSchema|\"database\".\"my_schema\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|MODIFY,CREATE TABLE|OnSchema|OnSchema|\"database\".\"my_schema\""

# all schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    all_schemas_in_database = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnAllSchemasInDatabase|\"database\""

# future schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    future_schemas_in_database = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnFutureSchemasInDatabase|\"database\""

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "REFERENCES"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
}

## ID: "\"role_name\"|false|false|SELECT,REFERENCES|OnSchemaObject|VIEW|\"database\".\"my_schema\".\"my_view\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnSchemaObject|OnObject|VIEW|\"database\".\"my_schema\".\"my_view\""

# all in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnAll|TABLES|InDatabase|\"database\""

# all in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnAll|TABLES|InSchema|\"database\".\"my_schema\""

# future in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnFuture|TABLES|InDatabase|\"database\""

# future in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnFuture|TABLES|InSchema|\"database\".\"my_schema\""
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_role_name` (String) The fully qualified name of the account role to which privileges will be granted. For more information about this resource, see [docs](./account_role).

### Optional

- `all_privileges` (Boolean) (Default: `false`) Grant all privileges on the account role. When all privileges cannot be granted, the provider returns a warning, which is aligned with the Snowsight behavior.
- `always_apply` (Boolean) (Default: `false`) If true, the resource will always produce a “plan” and on “apply” it will re-grant defined privileges. It is supposed to be used only in “grant privileges on all X’s in database / schema Y” or “grant all privileges to X” scenarios to make sure that every new object in a given database / schema is granted by the account role and every new privilege is granted to the database role. Important note: this flag is not compliant with the Terraform assumptions of the config being eventually convergent (producing an empty plan).
- `always_apply_trigger` (String) (Default: ``) This is a helper field and should not be set. Its main purpose is to help to achieve the functionality described by the always_apply field.
- `on_account` (Boolean) (Default: `false`) If true, the privileges will be granted on the account.
- `on_account_object` (Block List, Max: 1) Specifies the account object on which privileges will be granted (see [below for nested schema](#nestedblock--on_account_object))
- `on_schema` (Block List, Max: 1) Specifies the schema on which privileges will be granted. (see [below for nested schema](#nestedblock--on_schema))
- `on_schema_object` (Block List, Max: 1) Specifies the schema object on which privileges will be granted. (see [below for nested schema](#nestedblock--on_schema_object))
- `privileges` (Set of String) The privileges to grant on the account role. This field is case-sensitive; use only upper-case privileges.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `with_grant_option` (Boolean) (Default: `false`) Specifies whether the grantee can grant the privileges to other users.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--on_account_object"></a>
### Nested Schema for `on_account_object`

Required:

- `object_name` (String) The fully qualified name of the object on which privileges will be granted.
- `object_type` (String) The object type of the account object on which privileges will be granted. Valid values are: USER | RESOURCE MONITOR | WAREHOUSE | COMPUTE POOL | DATABASE | INTEGRATION | FAILOVER GROUP | REPLICATION GROUP | EXTERNAL VOLUME


<a id="nestedblock--on_schema"></a>
### Nested Schema for `on_schema`

Optional:

- `all_schemas_in_database` (String) The fully qualified name of the database.
- `future_schemas_in_database` (String) The fully qualified name of the database.
- `schema_name` (String) The fully qualified name of the schema.


<a id="nestedblock--on_schema_object"></a>
### Nested Schema for `on_schema_object`

Optional:

- `all` (Block List, Max: 1) Configures the privilege to be granted on all objects in either a database or schema. (see [below for nested schema](#nestedblock--on_schema_object--all))
- `future` (Block List, Max: 1) Configures the privilege to be granted on future objects in either a database or schema. (see [below for nested schema](#nestedblock--on_schema_object--future))
- `object_name` (String) The fully qualified name of the object on which privileges will be granted.
- `object_type` (String) The object type of the schema object on which privileges will be granted. Valid values are: AGGREGATION POLICY | ALERT | AUTHENTICATION POLICY | CORTEX SEARCH SERVICE | DATA METRIC FUNCTION | DYNAMIC TABLE | EVENT TABLE | EXTERNAL TABLE | FILE FORMAT | FUNCTION | GIT REPOSITORY | HYBRID TABLE | IMAGE REPOSITORY | ICEBERG TABLE | MASKING POLICY | MATERIALIZED VIEW | MODEL | NETWORK RULE | NOTEBOOK | PACKAGES POLICY | PASSWORD POLICY | PIPE | PROCEDURE | PROJECTION POLICY | ROW ACCESS POLICY | SECRET | SERVICE | SESSION POLICY | SEQUENCE | SNAPSHOT | STAGE | STREAM | TABLE | TAG | TASK | VIEW | STREAMLIT | DATASET

<a id="nestedblock--on_schema_object--all"></a>
### Nested Schema for `on_schema_object.all`

Required:

- `object_type_plural` (String) The plural object type of the schema object on which privileges will be granted. Valid values are: AGGREGATION POLICIES | ALERTS | AUTHENTICATION POLICIES | CORTEX SEARCH SERVICES | DATA METRIC FUNCTIONS | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | FILE FORMATS | FUNCTIONS | GIT REPOSITORIES | HYBRID TABLES | IMAGE REPOSITORIES | ICEBERG TABLES | MASKING POLICIES | MATERIALIZED VIEWS | MODELS | NETWORK RULES | NOTEBOOKS | PACKAGES POLICIES | PASSWORD POLICIES | PIPES | PROCEDURES | PROJECTION POLICIES | ROW ACCESS POLICIES | SECRETS | SERVICES | SESSION POLICIES | SEQUENCES | SNAPSHOTS | STAGES | STREAMS | TABLES | TAGS | TASKS | VIEWS | STREAMLITS | DATASETS.

Optional:

- `in_database` (String)
- `in_schema` (String)


<a id="nestedblock--on_schema_object--future"></a>
### Nested Schema for `on_schema_object.future`

Required:

- `object_type_plural` (String) The plural object type of the schema object on which privileges will be granted. Valid values are: ALERTS | AUTHENTICATION POLICIES | CORTEX SEARCH SERVICES | DATA METRIC FUNCTIONS | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | FILE FORMATS | FUNCTIONS | GIT REPOSITORIES | HYBRID TABLES | ICEBERG TABLES | MATERIALIZED VIEWS | MODELS | NETWORK RULES | NOTEBOOKS | PASSWORD POLICIES | PIPES | PROCEDURES | SECRETS | SERVICES | SEQUENCES | SNAPSHOTS | STAGES | STREAMS | TABLES | TASKS | VIEWS | STREAMLITS | DATASETS.

Optional:

- `in_database` (String)
- `in_schema` (String)



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

## Known limitations
- Setting the `CREATE SNOWFLAKE.ML.ANOMALY_DETECTION` or `CREATE SNOWFLAKE.ML.FORECAST` privileges on schema results in a permadiff because of the probably incorrect Snowflake's behavior of `SHOW GRANTS ON <object_type> <object_name>`. More in the [comment](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2651#issuecomment-2022634952).

## Import

~> **Note** All the ..._name parts should be fully qualified names (where every part is quoted), e.g. for schema object it is `"<database_name>"."<schema_name>"."<object_name>"`
~> **Note** To import all_privileges write ALL or ALL PRIVILEGES in place of `<privileges>`

Import is supported using the following syntax:

`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>'`

where:
- account_role_name - fully qualified identifier
- with_grant_option - boolean
- always_apply - boolean
- privileges - list of privileges, comma separated; to import all_privileges write "ALL" or "ALL PRIVILEGES"
- grant_type - enum
- grant_data - enum data

It has varying number of parts, depending on grant_type. All the possible types are:

### OnAccount
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnAccount'`

### OnAccountObject
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnAccountObject|<object_type>|<object_name>'`

### OnSchema

On schema contains inner types for all options.

#### OnSchema
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnSchema|<schema_name>'`

#### OnAllSchemasInDatabase
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnAllSchemasInDatabase|<database_name>'`

#### OnFutureSchemasInDatabase
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnFutureSchemasInDatabase|<database_name>'`

### OnSchemaObject

On schema object contains inner types for all options.

#### OnObject
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>'`

#### OnAll

On all contains inner types for all options.

##### InDatabase
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnAll|<object_type_plural>|InDatabase|<identifier>'`

##### InSchema
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnAll|<object_type_plural>|InSchema|<identifier>'`

#### OnFuture

On future contains inner types for all options.

##### InDatabase
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnFuture|<object_type_plural>|InDatabase|<identifier>'`

##### InSchema
`terraform import snowflake_grant_privileges_to_account_role.example '<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnFuture|<object_type_plural>|InSchema|<identifier>'`

### Import examples

#### Grant all privileges OnAccountObject (Database)
`terraform import snowflake_grant_privileges_to_account_role.example '"test_db_role"|false|false|ALL|OnAccountObject|DATABASE|"test_db"'`

#### Grant list of privileges OnAllSchemasInDatabase
`terraform import snowflake_grant_privileges_to_account_role.example '"test_db_role"|false|false|CREATE TAG,CREATE TABLE|OnSchema|OnAllSchemasInDatabase|"test_db"'`

#### Grant list of privileges on table
`terraform import snowflake_grant_privileges_to_account_role.example '"test_db_role"|false|false|SELECT,DELETE,INSERT|OnSchemaObject|OnObject|TABLE|"test_db"."test_schema"."test_table"'`

#### Grant list of privileges OnAll tables in schema
`terraform import snowflake_grant_privileges_to_account_role.example '"test_db_role"|false|false|SELECT,DELETE,INSERT|OnSchemaObject|OnAll|TABLES|InSchema|"test_db"."test_schema"'`

