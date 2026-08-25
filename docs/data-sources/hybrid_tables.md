---
page_title: "snowflake_hybrid_tables Data Source - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Data source used to get details of filtered hybrid tables. Filtering is aligned with the current possibilities for SHOW HYBRID TABLES https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables query (like, in, starts_with, limit). The results of SHOW, DESCRIBE, SHOW PARAMETERS, SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, SHOW IMPORTED KEYS, and SHOW INDEXES are encapsulated in one output collection hybrid_tables.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

# snowflake_hybrid_tables (Data Source)

Data source used to get details of filtered hybrid tables. Filtering is aligned with the current possibilities for [SHOW HYBRID TABLES](https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables) query (`like`, `in`, `starts_with`, `limit`). The results of SHOW, DESCRIBE, SHOW PARAMETERS, SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, SHOW IMPORTED KEYS, and SHOW INDEXES are encapsulated in one output collection `hybrid_tables`.

## Example Usage

```terraform
# Simple usage
data "snowflake_hybrid_tables" "simple" {
}

output "simple_output" {
  value = data.snowflake_hybrid_tables.simple.hybrid_tables
}

# Filtering (like)
data "snowflake_hybrid_tables" "like" {
  like = "hybrid-table-name"
}

output "like_output" {
  value = data.snowflake_hybrid_tables.like.hybrid_tables
}

# Filtering by prefix (like)
data "snowflake_hybrid_tables" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_hybrid_tables.like_prefix.hybrid_tables
}

# Filtering (in)
data "snowflake_hybrid_tables" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_hybrid_tables" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "database" : data.snowflake_hybrid_tables.in_database.hybrid_tables,
    "schema" : data.snowflake_hybrid_tables.in_schema.hybrid_tables,
  }
}

# Filtering (starts_with)
data "snowflake_hybrid_tables" "starts_with" {
  starts_with = "prefix"
}

output "starts_with_output" {
  value = data.snowflake_hybrid_tables.starts_with.hybrid_tables
}

# Filtering (limit)
data "snowflake_hybrid_tables" "limit" {
  limit {
    rows = 10
    from = "hybrid-table-"
  }
}

output "limit_output" {
  value = data.snowflake_hybrid_tables.limit.hybrid_tables
}

# Without additional data (to limit the number of calls made for every found hybrid table)
data "snowflake_hybrid_tables" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE TABLE for every hybrid table found and attaches its output to hybrid_tables.*.describe_output field
  with_describe = false
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR TABLE for every hybrid table found and attaches its output to hybrid_tables.*.parameters field
  with_parameters = false
  # with_keys is turned off by default; set to true to call SHOW PRIMARY KEYS / UNIQUE KEYS / IMPORTED KEYS for every hybrid table found and attach the merged constraints (grouped by constraint name, ordered by kind then column names) to hybrid_tables.*.show_keys_output field
  with_keys = false
  # with_indexes is turned off by default; set to true to call SHOW INDEXES for every hybrid table found and attach its output to hybrid_tables.*.show_indexes field
  with_indexes = false
}

output "only_show_output" {
  value = data.snowflake_hybrid_tables.only_show.hybrid_tables
}

# Ensure the number of hybrid tables is equal to at least one element (with the use of postcondition)
data "snowflake_hybrid_tables" "assert_with_postcondition" {
  like = "hybrid-table-name%"
  lifecycle {
    postcondition {
      condition     = length(self.hybrid_tables) > 0
      error_message = "there should be at least one hybrid table"
    }
  }
}

# Ensure the number of hybrid tables is equal to exactly one element (with the use of check block)
check "hybrid_table_check" {
  data "snowflake_hybrid_tables" "assert_with_check_block" {
    like = "hybrid-table-name"
  }

  assert {
    condition     = length(data.snowflake_hybrid_tables.assert_with_check_block.hybrid_tables) == 1
    error_message = "hybrid tables filtered by '${data.snowflake_hybrid_tables.assert_with_check_block.like}' returned ${length(data.snowflake_hybrid_tables.assert_with_check_block.hybrid_tables)} hybrid tables where one was expected"
  }
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `in` (Block List, Max: 1) IN clause to filter the list of objects (see [below for nested schema](#nestedblock--in))
- `like` (String) Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).
- `limit` (Block List, Max: 1) Limits the number of rows returned. If the `limit.from` is set, then the limit will start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`. (see [below for nested schema](#nestedblock--limit))
- `starts_with` (String) Filters the output with **case-sensitive** characters indicating the beginning of the object name.
- `with_describe` (Boolean) (Default: `true`) Runs DESC TABLE for each hybrid table returned by SHOW HYBRID TABLES. The output of describe is saved to the describe_output field. By default this value is set to true.
- `with_indexes` (Boolean) (Default: `false`) Runs SHOW INDEXES for each hybrid table returned by SHOW HYBRID TABLES. The output is saved to the show_indexes field. By default this value is set to false.
- `with_keys` (Boolean) (Default: `false`) Runs SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, and SHOW IMPORTED KEYS for each hybrid table returned by SHOW HYBRID TABLES. The merged constraints are saved to the show_keys_output field, grouped by constraint name and ordered by kind, then by column names. By default this value is set to false.
- `with_parameters` (Boolean) (Default: `true`) Runs SHOW PARAMETERS FOR TABLE for each hybrid table returned by SHOW HYBRID TABLES. The output is saved to the parameters field. By default this value is set to true.

### Read-Only

- `hybrid_tables` (List of Object) Holds the aggregated output of all hybrid table details queries. (see [below for nested schema](#nestedatt--hybrid_tables))
- `id` (String) The ID of this resource.

<a id="nestedblock--in"></a>
### Nested Schema for `in`

Optional:

- `account` (Boolean) Returns records for the entire account.
- `database` (String) Returns records for the current database in use or for a specified database.
- `schema` (String) Returns records for the current schema in use or a specified schema. Use fully qualified name.


<a id="nestedblock--limit"></a>
### Nested Schema for `limit`

Required:

- `rows` (Number) The maximum number of rows to return.

Optional:

- `from` (String) Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.


<a id="nestedatt--hybrid_tables"></a>
### Nested Schema for `hybrid_tables`

Read-Only:

- `describe_output` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--describe_output))
- `parameters` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--parameters))
- `show_indexes` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--show_indexes))
- `show_keys_output` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--show_keys_output))
- `show_output` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--show_output))

<a id="nestedobjatt--hybrid_tables--describe_output"></a>
### Nested Schema for `hybrid_tables.describe_output`

Read-Only:

- `check` (String)
- `collation` (String)
- `comment` (String)
- `default` (String)
- `expression` (String)
- `is_nullable` (Boolean)
- `kind` (String)
- `name` (String)
- `policy_name` (String)
- `primary_key` (Boolean)
- `privacy_domain` (String)
- `schema_evolution_record` (String)
- `type` (String)
- `unique_key` (Boolean)


<a id="nestedobjatt--hybrid_tables--parameters"></a>
### Nested Schema for `hybrid_tables.parameters`

Read-Only:

- `data_retention_time_in_days` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--parameters--data_retention_time_in_days))
- `max_data_extension_time_in_days` (List of Object) (see [below for nested schema](#nestedobjatt--hybrid_tables--parameters--max_data_extension_time_in_days))

<a id="nestedobjatt--hybrid_tables--parameters--data_retention_time_in_days"></a>
### Nested Schema for `hybrid_tables.parameters.data_retention_time_in_days`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--hybrid_tables--parameters--max_data_extension_time_in_days"></a>
### Nested Schema for `hybrid_tables.parameters.max_data_extension_time_in_days`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)



<a id="nestedobjatt--hybrid_tables--show_indexes"></a>
### Nested Schema for `hybrid_tables.show_indexes`

Read-Only:

- `columns` (String)
- `created_on` (String)
- `database_name` (String)
- `included_columns` (String)
- `is_unique` (Boolean)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `schema_name` (String)
- `table_name` (String)


<a id="nestedobjatt--hybrid_tables--show_keys_output"></a>
### Nested Schema for `hybrid_tables.show_keys_output`

Read-Only:

- `columns` (List of String)
- `delete_rule` (String)
- `kind` (String)
- `name` (String)
- `referenced_columns` (List of String)
- `referenced_table` (String)
- `update_rule` (String)


<a id="nestedobjatt--hybrid_tables--show_output"></a>
### Nested Schema for `hybrid_tables.show_output`

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
