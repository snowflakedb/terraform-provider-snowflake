---
page_title: "snowflake_hybrid_table Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage hybrid table objects. For more information, check hybrid table documentation https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

# snowflake_hybrid_table (Resource)

Resource used to manage hybrid table objects. For more information, check [hybrid table documentation](https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table).

## Example Usage

```terraform
## Minimal - Single column with primary key
resource "snowflake_hybrid_table" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "hybrid_table_name"

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }
}

## With multiple columns and inline constraints
resource "snowflake_hybrid_table" "with_constraints" {
  database = "database_name"
  schema   = "schema_name"
  name     = "users"

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
    comment     = "User ID"
  }

  column {
    name     = "email"
    type     = "VARCHAR(255)"
    nullable = false
    unique   = true
    comment  = "User email address"
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  comment = "User data hybrid table"
}

## With composite primary key and indexes
resource "snowflake_hybrid_table" "with_composite_key" {
  database = "database_name"
  schema   = "schema_name"
  name     = "order_items"

  column {
    name     = "order_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "item_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "quantity"
    type = "NUMBER(10,0)"
  }

  column {
    name = "price"
    type = "NUMBER(10,2)"
  }

  # Out-of-line primary key constraint
  primary_key {
    name    = "pk_order_items"
    columns = ["order_id", "item_id"]
  }

  # Secondary index for faster queries
  index {
    name    = "idx_order_id"
    columns = ["order_id"]
  }

  data_retention_time_in_days = 7
  comment                     = "Order items with composite primary key"
}

## With foreign key constraint
resource "snowflake_hybrid_table" "orders" {
  database = "database_name"
  schema   = "schema_name"
  name     = "orders"

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name = "customer_id"
    type = "NUMBER(38,0)"
    foreign_key {
      table_name  = "${var.database}.${var.schema}.customers"
      column_name = "id"
    }
  }

  column {
    name = "order_date"
    type = "DATE"
  }

  comment = "Orders table with foreign key to customers"
}

## With identity column (auto-increment)
resource "snowflake_hybrid_table" "with_identity" {
  database = "database_name"
  schema   = "schema_name"
  name     = "products"

  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = false
    identity {
      start_num = 1
      step_num  = 1
    }
    primary_key = true
  }

  column {
    name = "product_name"
    type = "VARCHAR(200)"
  }

  column {
    name = "price"
    type = "NUMBER(10,2)"
  }
}

## Complete example with all features
resource "snowflake_hybrid_table" "complete" {
  database  = "database_name"
  schema    = "schema_name"
  name      = "complete_example"
  or_replace = false

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    comment     = "Primary identifier"
  }

  column {
    name     = "code"
    type     = "VARCHAR(50)"
    nullable = false
    collate  = "en-ci"
  }

  column {
    name = "data"
    type = "VARIANT"
  }

  column {
    name = "updated_at"
    type = "TIMESTAMP_NTZ"
    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  primary_key {
    name    = "pk_complete"
    columns = ["id"]
  }

  unique_constraint {
    name    = "uq_code"
    columns = ["code"]
  }

  index {
    name    = "idx_updated_at"
    columns = ["updated_at"]
  }

  data_retention_time_in_days = 30
  comment                     = "Complete hybrid table example"
}
```

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).

-> **Note** Hybrid tables require at least one column with a primary key constraint. This can be specified either as an inline constraint on a column or as an out-of-line primary key constraint.

-> **Note** Changes to column definitions (name, type, constraints) typically require table recreation. Plan carefully when defining your hybrid table schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `column` (Block List, Min: 1) Definitions of columns for the hybrid table. (see [below for nested schema](#nestedblock--column))
- `database` (String) The database in which to create the hybrid table. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `name` (String) Specifies the identifier for the hybrid table; must be unique for the schema in which the hybrid table is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `schema` (String) The schema in which to create the hybrid table. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `comment` (String) Specifies a comment for the hybrid table.
- `data_retention_time_in_days` (Number) Specifies the retention period for the table in days. Valid values are between 0 and 90.
- `foreign_key` (Block Set) Out-of-line foreign key constraint definitions. (see [below for nested schema](#nestedblock--foreign_key))
- `index` (Block Set) Definitions of indexes for the hybrid table. (see [below for nested schema](#nestedblock--index))
- `or_replace` (Boolean) Specifies whether to replace the hybrid table if it already exists. Default: `false`.
- `primary_key` (Block List, Max: 1) Out-of-line primary key constraint definition. (see [below for nested schema](#nestedblock--primary_key))
- `unique_constraint` (Block Set) Out-of-line unique constraint definitions. (see [below for nested schema](#nestedblock--unique_constraint))

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE HYBRID TABLE` for the given hybrid table. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `show_output` (List of Object) Outputs the result of `SHOW HYBRID TABLES` for the given hybrid table. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--column"></a>
### Nested Schema for `column`

Required:

- `name` (String) Column name.
- `type` (String) Column data type.

Optional:

- `collate` (String) Collation specification for string column.
- `comment` (String) Column comment.
- `default` (Block List, Max: 1) Defines the default value for the column. (see [below for nested schema](#nestedblock--column--default))
- `foreign_key` (Block List, Max: 1) Inline foreign key constraint for the column. (see [below for nested schema](#nestedblock--column--foreign_key))
- `identity` (Block List, Max: 1) Defines the identity/autoincrement configuration for the column. (see [below for nested schema](#nestedblock--column--identity))
- `nullable` (Boolean) Specifies whether the column can contain NULL values. Default is true (nullable). Default: `true`.
- `primary_key` (Boolean) Specifies whether the column is a primary key (inline constraint). Default: `false`.
- `unique` (Boolean) Specifies whether the column has a unique constraint (inline constraint). Default: `false`.

<a id="nestedblock--column--default"></a>
### Nested Schema for `column.default`

Optional:

- `expression` (String) Default value expression.
- `sequence` (String) Fully qualified name of sequence for default value.


<a id="nestedblock--column--foreign_key"></a>
### Nested Schema for `column.foreign_key`

Required:

- `column_name` (String) Column name in the referenced table.
- `table_name` (String) Name of the table being referenced.


<a id="nestedblock--column--identity"></a>
### Nested Schema for `column.identity`

Optional:

- `start_num` (Number) Starting value for the identity column. Default: `1`.
- `step_num` (Number) Step/increment value for the identity column. Default: `1`.



<a id="nestedblock--foreign_key"></a>
### Nested Schema for `foreign_key`

Required:

- `columns` (List of String) List of column names forming the foreign key.
- `references_columns` (List of String) List of column names in the referenced table.
- `references_table` (String) Name of the table being referenced.

Optional:

- `name` (String) Name of the foreign key constraint.


<a id="nestedblock--index"></a>
### Nested Schema for `index`

Required:

- `columns` (List of String) List of column names to include in the index.
- `name` (String) Index name.


<a id="nestedblock--primary_key"></a>
### Nested Schema for `primary_key`

Required:

- `columns` (List of String) List of column names forming the primary key.

Optional:

- `name` (String) Name of the primary key constraint.


<a id="nestedblock--unique_constraint"></a>
### Nested Schema for `unique_constraint`

Required:

- `columns` (List of String) List of column names forming the unique constraint.

Optional:

- `name` (String) Name of the unique constraint.


<a id="nestedatt--describe_output"></a>
### Nested Schema for `describe_output`

Read-Only:

- `check` (String)
- `comment` (String)
- `default` (String)
- `expression` (String)
- `is_nullable` (String)
- `kind` (String)
- `name` (String)
- `policy_name` (String)
- `primary_key` (String)
- `privacy_domain` (String)
- `schema_evolution_record` (String)
- `type` (String)
- `unique_key` (String)


<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `bytes` (Number)
- `comment` (String)
- `created_on` (String)
- `database_name` (String)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `rows` (Number)
- `schema_name` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_hybrid_table.example '"<database_name>"."<schema_name>"."<hybrid_table_name>"'
```
