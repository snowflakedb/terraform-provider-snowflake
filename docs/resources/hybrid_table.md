---
page_title: "snowflake_hybrid_table Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  Resource for managing Snowflake Hybrid Tables
---

# snowflake_hybrid_table (Resource)

A resource for managing hybrid tables in Snowflake. Hybrid tables combine the transactional capabilities of OLTP systems with the analytics capabilities of Snowflake's data platform. They support primary keys, foreign keys, unique constraints, and secondary indexes, all of which must be enforced.

**Important:** Hybrid tables are available in AWS and Azure commercial regions only and require Snowflake Enterprise Edition or higher.

## Example Usage

### Basic Hybrid Table

```terraform
resource "snowflake_hybrid_table" "example" {
  database = "MYDB"
  schema   = "MYSCHEMA"
  name     = "ORDERS"
  comment  = "Orders hybrid table"

  column {
    name     = "order_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "order_date"
    type = "TIMESTAMP_NTZ"
  }

  column {
    name = "amount"
    type = "NUMBER(10,2)"
  }

  constraint {
    name    = "pk_orders"
    type    = "PRIMARY KEY"
    columns = ["order_id"]
  }
}
```

### With Secondary Indexes

```terraform
resource "snowflake_hybrid_table" "with_indexes" {
  database = "MYDB"
  schema   = "MYSCHEMA"
  name     = "CUSTOMERS"

  column {
    name     = "customer_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "email"
    type = "VARCHAR(100)"
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
  }

  constraint {
    name    = "pk_customers"
    type    = "PRIMARY KEY"
    columns = ["customer_id"]
  }

  constraint {
    name    = "uk_email"
    type    = "UNIQUE"
    columns = ["email"]
  }

  index {
    name    = "idx_created_at"
    columns = ["created_at"]
  }

  index {
    name    = "idx_email_created"
    columns = ["email", "created_at"]
  }
}
```

### With Foreign Key Constraint

```terraform
resource "snowflake_hybrid_table" "customers" {
  database = "ECOMMERCE"
  schema   = "PUBLIC"
  name     = "CUSTOMERS"

  column {
    name     = "customer_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "name"
    type = "VARCHAR(200)"
  }

  constraint {
    name    = "pk_customers"
    type    = "PRIMARY KEY"
    columns = ["customer_id"]
  }
}

resource "snowflake_hybrid_table" "orders" {
  database = "ECOMMERCE"
  schema   = "PUBLIC"
  name     = "ORDERS"

  column {
    name     = "order_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "customer_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "total_amount"
    type = "NUMBER(10,2)"
  }

  constraint {
    name    = "pk_orders"
    type    = "PRIMARY KEY"
    columns = ["order_id"]
  }

  constraint {
    name    = "fk_customer"
    type    = "FOREIGN KEY"
    columns = ["customer_id"]
    foreign_key {
      table_id  = snowflake_hybrid_table.customers.fully_qualified_name
      columns   = ["customer_id"]
      on_delete = "CASCADE"
      on_update = "RESTRICT"
    }
  }

  index {
    name    = "idx_customer"
    columns = ["customer_id"]
  }

  depends_on = [snowflake_hybrid_table.customers]
}
```

## Schema

### Required

- `column` (Block List, Min: 1) Definitions of columns to create in the hybrid table. Minimum one required. (see [below for nested schema](#nestedblock--column))
- `constraint` (Block List, Min: 1) Definitions of constraints for the hybrid table. At least one PRIMARY KEY constraint is required. Constraints cannot be modified after table creation. (see [below for nested schema](#nestedblock--constraint))
- `database` (String) The database in which to create the hybrid table.
- `name` (String) Specifies the identifier for the hybrid table; must be unique for the database and schema in which the table is created.
- `schema` (String) The schema in which to create the hybrid table.

### Optional

- `comment` (String) Specifies a comment for the hybrid table
- `index` (Block List) Definitions of secondary indexes for the hybrid table (see [below for nested schema](#nestedblock--index))

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE TABLE COLUMNS` for the given hybrid table. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `show_output` (List of Object) Outputs the result of `SHOW TABLES` for the given hybrid table. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--column"></a>
### Nested Schema for `column`

Required:

- `name` (String) Column name
- `type` (String) Column type, e.g. NUMBER, VARCHAR. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).

Optional:

- `comment` (String) Column comment
- `nullable` (Boolean) Whether this column can contain null values. Set to false for columns used in primary key constraints. Default: `true`

<a id="nestedblock--constraint"></a>
### Nested Schema for `constraint`

Required:

- `columns` (List of String) Columns to use in the constraint
- `name` (String) Name of the constraint
- `type` (String) Type of constraint: PRIMARY KEY, FOREIGN KEY, or UNIQUE. All constraints must be ENFORCED.

Optional:

- `foreign_key` (Block List, Max: 1) Foreign key reference details. Required when type is FOREIGN KEY. (see [below for nested schema](#nestedblock--constraint--foreign_key))

<a id="nestedblock--constraint--foreign_key"></a>
### Nested Schema for `constraint.foreign_key`

Required:

- `columns` (List of String) Columns in the referenced table
- `table_id` (String) Identifier of the referenced table in the format database.schema.table

Optional:

- `match` (String) The match type for the foreign key: FULL, SIMPLE, or PARTIAL. Note: MATCH is not supported for hybrid tables and will be ignored.
- `on_delete` (String) Action to perform when the primary/unique key is deleted: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION. Note: not supported for hybrid tables and will be ignored.
- `on_update` (String) Action to perform when the primary/unique key is updated: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION. Note: not supported for hybrid tables and will be ignored.

<a id="nestedblock--index"></a>
### Nested Schema for `index`

Required:

- `columns` (List of String) Columns to include in the index
- `name` (String) Name of the index

<a id="nestedatt--describe_output"></a>
### Nested Schema for `describe_output`

Read-Only:

- `check` (String)
- `collation` (String)
- `comment` (String)
- `default` (String)
- `expression` (String)
- `is_nullable` (Boolean)
- `is_primary` (Boolean)
- `is_unique` (Boolean)
- `kind` (String)
- `name` (String)
- `policy_name` (String)
- `schema_evolution_record` (String)
- `type` (String)

<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `automatic_clustering` (Boolean)
- `budget` (String)
- `bytes` (Number)
- `change_tracking` (Boolean)
- `cluster_by` (String)
- `comment` (String)
- `created_on` (String)
- `database_name` (String)
- `dropped_on` (String)
- `enable_schema_evolution` (Boolean)
- `is_event` (Boolean)
- `is_external` (Boolean)
- `kind` (String)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `retention_time` (Number)
- `rows` (Number)
- `schema_name` (String)
- `search_optimization` (Boolean)
- `search_optimization_bytes` (Number)
- `search_optimization_progress` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_hybrid_table.example "database|schema|name"
```

## Important Notes

- **Primary Key Required**: At least one PRIMARY KEY constraint must be defined. This is validated at the SDK level.
- **Enforced Constraints Only**: All constraints (PRIMARY KEY, FOREIGN KEY, UNIQUE) must be ENFORCED. NOT ENFORCED constraints are not supported for hybrid tables.
- **Immutable Constraints**: Constraints cannot be modified after table creation. Any changes to constraints will force recreation of the table.
- **Secondary Indexes**: Indexes can be added or removed dynamically without recreating the table.
- **Regional Availability**: Hybrid tables are only available in AWS and Azure commercial regions.
- **Edition Requirement**: Requires Snowflake Enterprise Edition or higher.
- **Limited ALTER Support**: Some ALTER TABLE operations available for standard tables are not supported for hybrid tables (e.g., changing PRIMARY KEY, modifying constraints).

## References

- [Snowflake Hybrid Tables Documentation](https://docs.snowflake.com/en/user-guide/tables-hybrid)
- [CREATE HYBRID TABLE](https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table)
