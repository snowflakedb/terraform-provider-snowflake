---
page_title: "snowflake_grant_privileges_to_database_role Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  
---


!> **Warning** Be careful when using `always_apply` field. It will always produce a plan (even when no changes were made) and can be harmful in some setups. For more details why we decided to introduce it to go our document explaining those design decisions (coming soon).

~> **Note** Manage grants on `HYBRID TABLE` by specifying `TABLE` or `TABLES` in `object_type` field. This applies to a single object, all objects, or future objects. This reflects the current behavior in Snowflake.

~> **Note** Please, follow the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/security-access-control-considerations) for best practices on access control. The provider does not enforce any specific methodology, so it is essential for users to choose the appropriate strategy for seamless privilege management. Additionally, refer to [this link](https://docs.snowflake.com/en/user-guide/security-access-control-privileges) for a list of all available privileges in Snowflake.

# snowflake_grant_privileges_to_database_role (Resource)



## Example Usage

```terraform
resource "snowflake_database" "db" {
  name = "database"
}

resource "snowflake_schema" "my_schema" {
  database = snowflake_database.db.name
  name     = "my_schema"
}

resource "snowflake_database_role" "db_role" {
  database = snowflake_database.db.name
  name     = "db_role_name"
}

##################################
### on database privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["CREATE", "MONITOR"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
  all_privileges     = true
  with_grant_option  = true
}

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
  always_apply       = true
  all_privileges     = true
  with_grant_option  = true
}

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all schemas in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    all_schemas_in_database = snowflake_database_role.db_role.database
  }
}

# future schemas in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    future_schemas_in_database = snowflake_database_role.db_role.database
  }
}

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "REFERENCES"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database_role.db_role.database
    }
  }
}

# all in schema
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

# future in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database_role.db_role.database
    }
  }
}

# future in schema
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database_role_name` (String) The fully qualified name of the database role to which privileges will be granted. For more information about this resource, see [docs](./database_role).

### Optional

- `all_privileges` (Boolean) (Default: `false`) Grant all privileges on the database role.
- `always_apply` (Boolean) (Default: `false`) If true, the resource will always produce a “plan” and on “apply” it will re-grant defined privileges. It is supposed to be used only in “grant privileges on all X’s in database / schema Y” or “grant all privileges to X” scenarios to make sure that every new object in a given database / schema is granted by the account role and every new privilege is granted to the database role. Important note: this flag is not compliant with the Terraform assumptions of the config being eventually convergent (producing an empty plan).
- `always_apply_trigger` (String) (Default: ``) This is a helper field and should not be set. Its main purpose is to help to achieve the functionality described by the always_apply field.
- `on_database` (String) The fully qualified name of the database on which privileges will be granted. For more information about this resource, see [docs](./database).
- `on_schema` (Block List, Max: 1) Specifies the schema on which privileges will be granted. (see [below for nested schema](#nestedblock--on_schema))
- `on_schema_object` (Block List, Max: 1) Specifies the schema object on which privileges will be granted. (see [below for nested schema](#nestedblock--on_schema_object))
- `privileges` (Set of String) The privileges to grant on the database role.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `with_grant_option` (Boolean) (Default: `false`) If specified, allows the recipient role to grant the privileges to other roles.

### Read-Only

- `id` (String) The ID of this resource.

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

- `in_database` (String) The fully qualified name of the database.
- `in_schema` (String) The fully qualified name of the schema.


<a id="nestedblock--on_schema_object--future"></a>
### Nested Schema for `on_schema_object.future`

Required:

- `object_type_plural` (String) The plural object type of the schema object on which privileges will be granted. Valid values are: ALERTS | AUTHENTICATION POLICIES | DATA METRIC FUNCTIONS | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | FILE FORMATS | FUNCTIONS | GIT REPOSITORIES | HYBRID TABLES | ICEBERG TABLES | MATERIALIZED VIEWS | MODELS | NETWORK RULES | NOTEBOOKS | PASSWORD POLICIES | PIPES | PROCEDURES | SECRETS | SERVICES | SEQUENCES | SNAPSHOTS | STAGES | STREAMS | TABLES | TASKS | VIEWS | STREAMLITS | DATASETS.

Optional:

- `in_database` (String) The fully qualified name of the database.
- `in_schema` (String) The fully qualified name of the schema.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

## Import

~> **Note** All the ..._name parts should be fully qualified names (where every part is quoted), e.g. for database object it is `"<database_name>"."<object_name>"`
~> **Note** To import all_privileges write ALL or ALL PRIVILEGES in place of `<privileges>`

Import is supported using the following syntax:

`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>'`

where:
- database_role_name - fully qualified identifier
- with_grant_option - boolean
- always_apply - boolean
- privileges - list of privileges, comma separated; to import all_privileges write "ALL" or "ALL PRIVILEGES"
- grant_type - enum
- grant_data - enum data

It has varying number of parts, depending on grant_type. All the possible types are:

### OnDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnDatabase|<database_name>'`

### OnSchema

On schema contains inner types for all options.

#### OnSchema
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnSchema|<schema_name>'`

#### OnAllSchemasInDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnAllSchemasInDatabase|<database_name>'`

#### OnFutureSchemasInDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|OnFutureSchemasInDatabase|<database_name>'`

### OnSchemaObject

On schema object contains inner types for all options.

#### OnObject
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>'`

#### OnAll

On all contains inner types for all options.

##### InDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnAll|<object_type_plural>|InDatabase|<identifier>'`

##### InSchema
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnAll|<object_type_plural>|InSchema|<identifier>'`

#### OnFuture

On future contains inner types for all options.

##### InDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnFuture|<object_type_plural>|InDatabase|<identifier>'`

##### InSchema
`terraform import snowflake_grant_privileges_to_database_role.example '<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnFuture|<object_type_plural>|InSchema|<identifier>'`

### Import examples

#### Grant all privileges OnDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '"test_db"."test_db_role"|false|false|ALL|OnDatabase|"test_db"'`

#### Grant list of privileges OnAllSchemasInDatabase
`terraform import snowflake_grant_privileges_to_database_role.example '"test_db"."test_db_role"|false|false|CREATE TAG,CREATE TABLE|OnSchema|OnAllSchemasInDatabase|"test_db"'`

#### Grant list of privileges on table
`terraform import snowflake_grant_privileges_to_database_role.example '"test_db"."test_db_role"|false|false|SELECT,DELETE,INSERT|OnSchemaObject|OnObject|TABLE|"test_db"."test_schema"."test_table"'`

#### Grant list of privileges OnAll tables in schema
`terraform import snowflake_grant_privileges_to_database_role.example '"test_db"."test_db_role"|false|false|SELECT,DELETE,INSERT|OnSchemaObject|OnAll|TABLES|InSchema|"test_db"."test_schema"'`

