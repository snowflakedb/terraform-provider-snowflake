---
page_title: "snowflake_iceberg_table Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage a Snowflake-managed Iceberg table. For more information, check the official documentation https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

~> **Note** Any change to the `column` block (adding, removing, renaming, retyping, or reordering a column) recreates the whole table, because column definitions can currently only be set at creation time (Snowflake `ALTER ICEBERG TABLE` column operations are not yet used by this resource). This will be addressed in a future release.

~> **Note** `primary_key_constraint`, `unique_constraint`, `foreign_key_constraint`, and `check_constraint` can only be set at creation time; changing or removing them recreates the whole table. They also are not read back from Snowflake, so external changes to these constraints (e.g. added, dropped, or altered outside Terraform) are not detected, and after importing the resource, the first `terraform plan` may show a diff for these fields even without a config change.

~> **Note** `path_layout`, `error_logging`, and `change_tracking` are not returned by `SHOW`/`DESCRIBE ICEBERG TABLE`, so external changes to these fields are not detected. `cluster_by` is not read back either, because Snowflake does not expose the original clustering key expression for Iceberg tables.

# snowflake_iceberg_table (Resource)

Resource used to manage a Snowflake-managed Iceberg table. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake).

## Example Usage

```terraform
# Basic - only required fields
resource "snowflake_iceberg_table" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }
  column {
    name = "NAME"
    type = "VARCHAR(16777216)"
  }
}

# Complete - every field set (except cluster_by, which conflicts with partition_by - see below)
resource "snowflake_iceberg_table" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"
  comment  = "COMMENT"

  # Snowflake-managed parameters
  external_volume                 = "EXTERNAL_VOLUME"
  catalog                         = "SNOWFLAKE"
  catalog_sync                    = "CATALOG_INTEGRATION"
  target_file_size                = "64MB"
  storage_serialization_policy    = "OPTIMIZED"
  data_retention_time_in_days     = 5
  max_data_extension_time_in_days = 10
  enable_data_compaction          = true
  enable_iceberg_merge_on_read    = true

  # Fields that cannot be changed after creation
  base_location   = "iceberg_table"
  path_layout     = "FLAT"
  change_tracking = "true"
  iceberg_version = 2

  # Fields that can be altered in place
  error_logging = "true"

  column {
    name     = "ID"
    type     = "NUMBER(38,0)"
    not_null = "true"
    comment  = "Primary identifier"
  }
  column {
    name    = "NAME"
    type    = "VARCHAR(16777216)"
    comment = "Name of the entity"
    masking_policy {
      policy_name = "MASKING_POLICY"
      using       = ["NAME"]
    }
  }
  column {
    name = "REGION"
    type = "VARCHAR(16777216)"
    projection_policy {
      policy_name = "PROJECTION_POLICY"
    }
  }
  column {
    name = "STATUS"
    type = "VARCHAR(16777216)"
    default {
      expression = "'ACTIVE'"
    }
  }
  column {
    name = "CATEGORY"
    type = "VARCHAR(16777216)"
    masking_policy {
      policy_name = "CONDITIONAL_MASKING_POLICY"
      using       = ["CATEGORY", "STATUS"]
    }
  }
  column {
    name = "CREATED_AT"
    type = "TIMESTAMP_NTZ(9)"
    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }
  column {
    name = "REF_ID"
    type = "NUMBER(38,0)"
  }

  primary_key_constraint {
    name               = "PK"
    column             = ["ID"]
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Primary key constraint"
  }

  unique_constraint {
    name               = "NAME_UQ"
    column             = ["NAME"]
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Unique constraint on name"
  }

  foreign_key_constraint {
    name               = "FK"
    column             = ["REF_ID"]
    table_name         = "OTHER_DATABASE.OTHER_SCHEMA.OTHER_TABLE"
    ref_column         = ["ID"]
    match              = "SIMPLE"
    on_update          = "CASCADE"
    on_delete          = "SET NULL"
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Foreign key constraint"
  }

  check_constraint {
    name       = "CHK"
    expression = "ID > 0"
    validate   = "true"
  }

  row_access_policy {
    policy_name = "ROW_ACCESS_POLICY"
    on          = ["ID"]
  }

  aggregation_policy {
    policy_name = "AGGREGATION_POLICY"
    entity_key  = ["ID"]
  }

  partition_by {
    identity = "REGION"
  }
  partition_by {
    bucket {
      num_buckets = 4
      column      = "ID"
    }
  }
  partition_by {
    truncate {
      width  = 10
      column = "NAME"
    }
  }
  partition_by {
    year = "CREATED_AT"
  }
  partition_by {
    month = "CREATED_AT"
  }
  partition_by {
    day = "CREATED_AT"
  }
  partition_by {
    hour = "CREATED_AT"
  }
}

# cluster_by conflicts with partition_by, so it is shown on a separate resource.
resource "snowflake_iceberg_table" "complete_with_cluster_by" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }
  column {
    name = "NAME"
    type = "VARCHAR(16777216)"
  }

  cluster_by = ["ID", "NAME"]
}
```

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `column` (Block List, Min: 1) Definitions of the columns to create in the table. Minimum one required. (see [below for nested schema](#nestedblock--column))
- `database` (String) The database in which to create the Iceberg table. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `name` (String) Specifies the identifier for the Iceberg table; must be unique for the schema in which the Iceberg table is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `schema` (String) The schema in which to create the Iceberg table. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `aggregation_policy` (Block List, Max: 1) Specifies the aggregation policy to set on a Iceberg table. (see [below for nested schema](#nestedblock--aggregation_policy))
- `base_location` (String) The path to a directory where Snowflake can write data and metadata files for the Iceberg table. Specify a relative path from the table's `EXTERNAL_VOLUME` location.
- `catalog` (String) Specifies the identifier for the catalog integration to use for the Iceberg table. If not specified, the account-level default is used.
- `catalog_sync` (String) Specifies the name of the catalog integration that Snowflake uses to automatically synchronize the Iceberg table with an external catalog. For more information, check [CATALOG_SYNC docs](https://docs.snowflake.com/en/sql-reference/parameters#catalog-sync).
- `change_tracking` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether to enable change tracking on the Iceberg table. Cannot be changed after creation. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `check_constraint` (Block List) Defines a table-level CHECK constraint. (see [below for nested schema](#nestedblock--check_constraint))
- `cluster_by` (List of String) A list of one or more table columns/expressions to be used as clustering key(s) for the table. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `comment` (String) Specifies a comment for the Iceberg table.
- `data_retention_time_in_days` (Number) Specifies the retention period for the Iceberg table so that Time Travel actions can be performed on historical data. For more information, check [DATA_RETENTION_TIME_IN_DAYS docs](https://docs.snowflake.com/en/sql-reference/parameters#data-retention-time-in-days).
- `enable_data_compaction` (Boolean) Specifies whether automatic background data compaction is enabled for the Iceberg table. For more information, check [ENABLE_DATA_COMPACTION docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-data-compaction).
- `enable_iceberg_merge_on_read` (Boolean) Specifies whether merge-on-read is enabled for the Iceberg table. For more information, check [ENABLE_ICEBERG_MERGE_ON_READ docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-iceberg-merge-on-read).
- `error_logging` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether error logging is enabled for the Iceberg table. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `external_volume` (String) Specifies the identifier for the external volume where the Iceberg table stores its metadata files and data in Parquet format. If not specified, the account-level default is used.
- `foreign_key_constraint` (Block List) Defines a table-level FOREIGN KEY constraint. (see [below for nested schema](#nestedblock--foreign_key_constraint))
- `iceberg_version` (Number) Specifies the Iceberg table format version.
- `max_data_extension_time_in_days` (Number) Specifies the maximum number of days for which Snowflake can extend the data retention period for the Iceberg table to prevent streams on the table from becoming stale. For more information, check [MAX_DATA_EXTENSION_TIME_IN_DAYS docs](https://docs.snowflake.com/en/sql-reference/parameters#max-data-extension-time-in-days).
- `partition_by` (Block List) Defines the partitioning for the Iceberg table. Cannot be changed after creation. Exactly one of identity, bucket, truncate, year, month, day, or hour must be set for each entry. Cannot be used together with `cluster_by`. (see [below for nested schema](#nestedblock--partition_by))
- `path_layout` (String) Specifies the storage layout for the Iceberg table's Parquet files. Valid values are: [FLAT HIERARCHICAL]. Cannot be changed after creation. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `primary_key_constraint` (Block List, Max: 1) Defines a table-level PRIMARY KEY constraint. (see [below for nested schema](#nestedblock--primary_key_constraint))
- `row_access_policy` (Block List, Max: 1) Specifies the row access policy to set on a Iceberg table. (see [below for nested schema](#nestedblock--row_access_policy))
- `storage_serialization_policy` (String) Specifies the storage serialization policy for the Iceberg table. Valid values are: [COMPATIBLE OPTIMIZED]. Cannot be changed after creation. For more information, check [STORAGE_SERIALIZATION_POLICY docs](https://docs.snowflake.com/en/sql-reference/parameters#storage-serialization-policy).
- `target_file_size` (String) Specifies the target file size (in bytes) used when writing the Iceberg table's Parquet files. Valid values are: [AUTO 16MB 32MB 64MB 128MB]. For more information, check [TARGET_FILE_SIZE docs](https://docs.snowflake.com/en/sql-reference/parameters#target-file-size).
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `unique_constraint` (Block List) Defines a table-level UNIQUE constraint. (see [below for nested schema](#nestedblock--unique_constraint))

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE ICEBERG TABLE` for the given Iceberg table. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `parameters` (List of Object) Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table. (see [below for nested schema](#nestedatt--parameters))
- `show_output` (List of Object) Outputs the result of `SHOW ICEBERG TABLES` for the given Iceberg table. Note that this value will be only recomputed whenever values of fields affecting the output change. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--column"></a>
### Nested Schema for `column`

Required:

- `name` (String) Column name.
- `type` (String) Column type, e.g. VARIANT. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).

Optional:

- `comment` (String) Column comment.
- `default` (Block List, Max: 1) Defines the column default value. (see [below for nested schema](#nestedblock--column--default))
- `masking_policy` (Block List, Max: 1) Specifies the masking policy to set on a column. For more information about this resource, see [docs](./masking_policy). (see [below for nested schema](#nestedblock--column--masking_policy))
- `not_null` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether to restrict the column to NOT NULL values.
- `projection_policy` (Block List, Max: 1) Specifies the projection policy to set on a column. (see [below for nested schema](#nestedblock--column--projection_policy))

<a id="nestedblock--column--default"></a>
### Nested Schema for `column.default`

Required:

- `expression` (String) The default expression value for the column.


<a id="nestedblock--column--masking_policy"></a>
### Nested Schema for `column.masking_policy`

Required:

- `policy_name` (String) Masking policy name. For more information about this resource, see [docs](./masking_policy).

Optional:

- `using` (List of String) Specifies the arguments to pass into the conditional masking policy SQL expression, in order. The first column in the list specifies the column for the policy conditions to mask or tokenize the data and must match the column to which the masking policy is set. The additional columns specify the columns to evaluate to determine whether to mask or tokenize the data in each row of the query result when a query is made on the first column. If the USING clause is omitted, Snowflake treats the conditional masking policy as a normal masking policy.


<a id="nestedblock--column--projection_policy"></a>
### Nested Schema for `column.projection_policy`

Required:

- `policy_name` (String) Projection policy name.



<a id="nestedblock--aggregation_policy"></a>
### Nested Schema for `aggregation_policy`

Required:

- `policy_name` (String) Aggregation policy name.

Optional:

- `entity_key` (Set of String) Defines which columns uniquely identify an entity within the Iceberg table.


<a id="nestedblock--check_constraint"></a>
### Nested Schema for `check_constraint`

Required:

- `expression` (String) The CHECK constraint expression.

Optional:

- `name` (String) Name of the constraint.
- `validate` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether existing data is validated against the constraint (`true`, `ENABLE VALIDATE`) or not (`false`, `ENABLE NOVALIDATE`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--foreign_key_constraint"></a>
### Nested Schema for `foreign_key_constraint`

Required:

- `column` (List of String) The local column(s) the foreign key is defined on.
- `table_name` (String) The table that the foreign key references.

Optional:

- `comment` (String) Constraint comment.
- `deferrable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is deferrable (`true`) or not deferrable (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enabled (`true`) or disabled (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enforced` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enforced (`true`) or not enforced (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `initially_deferred` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is initially deferred (`true`) or initially immediate (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `match` (String) The match type for the foreign key. Valid values are: [FULL SIMPLE PARTIAL].
- `name` (String) Name of the constraint.
- `on_delete` (String) Specifies the action to perform when the referenced primary/unique key is deleted. Valid values are: [CASCADE SET NULL SET DEFAULT RESTRICT NO ACTION].
- `on_update` (String) Specifies the action to perform when the referenced primary/unique key is updated. Valid values are: [CASCADE SET NULL SET DEFAULT RESTRICT NO ACTION].
- `ref_column` (List of String) The column(s) in the referenced table that the foreign key references.
- `rely` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether a constraint in NOVALIDATE mode is taken into account (`true`) or not (`false`) during query rewrite. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `validate` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether to validate existing data on the table when the constraint is created (`true`) or skip validation (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--partition_by"></a>
### Nested Schema for `partition_by`

Optional:

- `bucket` (Block List, Max: 1) Partitions the table by hashing the column into a fixed number of buckets. (see [below for nested schema](#nestedblock--partition_by--bucket))
- `day` (String) Partitions the table by the day component of the column.
- `hour` (String) Partitions the table by the hour component of the column.
- `identity` (String) Name of the column to use as-is for partitioning.
- `month` (String) Partitions the table by the month component of the column.
- `truncate` (Block List, Max: 1) Partitions the table by truncating the column value to a fixed width. (see [below for nested schema](#nestedblock--partition_by--truncate))
- `year` (String) Partitions the table by the year component of the column.

<a id="nestedblock--partition_by--bucket"></a>
### Nested Schema for `partition_by.bucket`

Required:

- `column` (String) Name of the column to bucket.
- `num_buckets` (Number) Number of buckets to hash the column values into.


<a id="nestedblock--partition_by--truncate"></a>
### Nested Schema for `partition_by.truncate`

Required:

- `column` (String) Name of the column to truncate.
- `width` (Number) Width to truncate the column value to.



<a id="nestedblock--primary_key_constraint"></a>
### Nested Schema for `primary_key_constraint`

Required:

- `column` (List of String) The column(s) the constraint applies to.

Optional:

- `comment` (String) Constraint comment.
- `deferrable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is deferrable (`true`) or not deferrable (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enabled (`true`) or disabled (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enforced` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enforced (`true`) or not enforced (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `initially_deferred` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is initially deferred (`true`) or initially immediate (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `name` (String) Name of the constraint.
- `rely` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether a constraint in NOVALIDATE mode is taken into account (`true`) or not (`false`) during query rewrite. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `validate` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether to validate existing data on the table when the constraint is created (`true`) or skip validation (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--row_access_policy"></a>
### Nested Schema for `row_access_policy`

Required:

- `on` (Set of String) Defines which columns are affected by the policy.
- `policy_name` (String) Row access policy name. For more information about this resource, see [docs](./row_access_policy).


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedblock--unique_constraint"></a>
### Nested Schema for `unique_constraint`

Required:

- `column` (List of String) The column(s) the constraint applies to.

Optional:

- `comment` (String) Constraint comment.
- `deferrable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is deferrable (`true`) or not deferrable (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enable` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enabled (`true`) or disabled (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `enforced` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is enforced (`true`) or not enforced (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `initially_deferred` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether the constraint is initially deferred (`true`) or initially immediate (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `name` (String) Name of the constraint.
- `rely` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether a constraint in NOVALIDATE mode is taken into account (`true`) or not (`false`) during query rewrite. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `validate` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Whether to validate existing data on the table when the constraint is created (`true`) or skip validation (`false`). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedatt--describe_output"></a>
### Nested Schema for `describe_output`

Read-Only:

- `check` (String)
- `comment` (String)
- `default` (String)
- `expression` (String)
- `is_nullable` (Boolean)
- `kind` (String)
- `name` (String)
- `name_mapping` (String)
- `policy_name` (String)
- `primary_key` (Boolean)
- `privacy_domain` (String)
- `source_iceberg_type` (String)
- `type` (String)
- `unique_key` (Boolean)
- `write_default` (String)


<a id="nestedatt--parameters"></a>
### Nested Schema for `parameters`

Read-Only:

- `catalog` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--catalog))
- `catalog_sync` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--catalog_sync))
- `data_retention_time_in_days` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--data_retention_time_in_days))
- `enable_data_compaction` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--enable_data_compaction))
- `enable_iceberg_merge_on_read` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--enable_iceberg_merge_on_read))
- `external_volume` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--external_volume))
- `max_data_extension_time_in_days` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--max_data_extension_time_in_days))
- `storage_serialization_policy` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--storage_serialization_policy))
- `target_file_size` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--target_file_size))

<a id="nestedobjatt--parameters--catalog"></a>
### Nested Schema for `parameters.catalog`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--catalog_sync"></a>
### Nested Schema for `parameters.catalog_sync`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--data_retention_time_in_days"></a>
### Nested Schema for `parameters.data_retention_time_in_days`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--enable_data_compaction"></a>
### Nested Schema for `parameters.enable_data_compaction`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--enable_iceberg_merge_on_read"></a>
### Nested Schema for `parameters.enable_iceberg_merge_on_read`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--external_volume"></a>
### Nested Schema for `parameters.external_volume`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--max_data_extension_time_in_days"></a>
### Nested Schema for `parameters.max_data_extension_time_in_days`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--storage_serialization_policy"></a>
### Nested Schema for `parameters.storage_serialization_policy`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--target_file_size"></a>
### Nested Schema for `parameters.target_file_size`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)



<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `auto_refresh_status` (List of Object) (see [below for nested schema](#nestedobjatt--show_output--auto_refresh_status))
- `base_location` (String)
- `can_write_metadata` (Boolean)
- `catalog_name` (String)
- `catalog_namespace` (String)
- `catalog_sync_name` (String)
- `catalog_table_name` (String)
- `comment` (String)
- `created_on` (String)
- `current_partition_spec_id` (Number)
- `database_name` (String)
- `external_volume_name` (String)
- `iceberg_table_format_version` (Number)
- `iceberg_table_type` (String)
- `name` (String)
- `name_mapping` (String)
- `owner` (String)
- `owner_role_type` (String)
- `partition_specs` (List of Object) (see [below for nested schema](#nestedobjatt--show_output--partition_specs))
- `schema_name` (String)

<a id="nestedobjatt--show_output--auto_refresh_status"></a>
### Nested Schema for `show_output.auto_refresh_status`

Read-Only:

- `current_snapshot_id` (Number)
- `execution_state` (String)
- `last_snapshot_time` (String)
- `last_updated_time` (String)
- `pending_snapshot_count` (Number)


<a id="nestedobjatt--show_output--partition_specs"></a>
### Nested Schema for `show_output.partition_specs`

Read-Only:

- `fields` (List of Object) (see [below for nested schema](#nestedobjatt--show_output--partition_specs--fields))
- `spec_id` (Number)

<a id="nestedobjatt--show_output--partition_specs--fields"></a>
### Nested Schema for `show_output.partition_specs.fields`

Read-Only:

- `field_id` (Number)
- `name` (String)
- `source_id` (Number)
- `transform` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_iceberg_table.example '"<database_name>"."<schema_name>"."<table_name>"'
```
