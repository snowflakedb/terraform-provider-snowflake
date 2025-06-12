---
page_title: "snowflake_services Data Source - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Data source used to get details of filtered services. Filtering is aligned with the current possibilities for SHOW SERVICES https://docs.snowflake.com/en/sql-reference/sql/show-services query. The results of SHOW and DESCRIBE are encapsulated in one output collection services. By default, the results includes both services and job services. If you want to filter only services or job service, set mode with a relevant option.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

# snowflake_services (Data Source)

Data source used to get details of filtered services. Filtering is aligned with the current possibilities for [SHOW SERVICES](https://docs.snowflake.com/en/sql-reference/sql/show-services) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `services`. By default, the results includes both services and job services. If you want to filter only services or job service, set `mode` with a relevant option.

## Example Usage

```terraform
# Simple usage
data "snowflake_services" "simple" {
}

output "simple_output" {
  value = data.snowflake_services.simple.services
}

# Filtering (like)
data "snowflake_services" "like" {
  like = "service-name"
}

output "like_output" {
  value = data.snowflake_services.like.services
}

# Filtering by prefix (like)
data "snowflake_services" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_services.like_prefix.services
}

# Filtering (starts_with)
data "snowflake_services" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_services.starts_with.services
}

# Filtering (in)
data "snowflake_services" "in_account" {
  in {
    account = true
  }
}

data "snowflake_services" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_services" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_services" "in_compute_pool" {
  in {
    compute_pool = "<compute_pool_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_services.in_account.services,
    "database" : data.snowflake_services.in_database.services,
    "schema" : data.snowflake_services.in_schema.services,
    "compute_pool" : data.snowflake_services.in_compute_pool.services,
  }
}

# Filtering (limit)
data "snowflake_services" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_services.limit.services
}

# Filtering (jobs only)
data "snowflake_services" "jobs_only" {
  mode = "JOBS_ONLY"
}

output "jobs_only_output" {
  value = data.snowflake_services.jobs_only.services
}

# Filtering (exclude jobs)
data "snowflake_services" "exclude_jobs" {
  mode = "SERVICES_ONLY"
}

output "exclude_jobs_output" {
  value = data.snowflake_services.exclude_jobs.services
}

# Without additional data (to limit the number of calls make for every found service)
data "snowflake_services" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SERVICE for every service found and attaches its output to services.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_services.only_show.services
}

# Ensure the number of services is equal to at least one element (with the use of postcondition)
data "snowflake_services" "assert_with_postcondition" {
  like = "service-name%"
  lifecycle {
    postcondition {
      condition     = length(self.services) > 0
      error_message = "there should be at least one service"
    }
  }
}

# Ensure the number of services is equal to exactly one element (with the use of check block)
check "service_check" {
  data "snowflake_services" "assert_with_check_block" {
    like = "service-name"
  }

  assert {
    condition     = length(data.snowflake_services.assert_with_check_block.services) == 1
    error_message = "services filtered by '${data.snowflake_services.assert_with_check_block.like}' returned ${length(data.snowflake_services.assert_with_check_block.services)} services where one was expected"
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
- `mode` (String) (Default: `ALL`) The type filtering of `SHOW SERVICES` results. `ALL` returns both services and job services. `JOBS_ONLY` returns only job services (`JOB` option in SQL). `SERVICES_ONLY` returns only services (`EXCLUDE_JOBS` option in SQL).
- `starts_with` (String) Filters the output with **case-sensitive** characters indicating the beginning of the object name.
- `with_describe` (Boolean) (Default: `true`) Runs DESC SERVICE for each service returned by SHOW SERVICES. The output of describe is saved to the description field. By default this value is set to true.

### Read-Only

- `id` (String) The ID of this resource.
- `services` (List of Object) Holds the aggregated output of all services details queries. (see [below for nested schema](#nestedatt--services))

<a id="nestedblock--in"></a>
### Nested Schema for `in`

Optional:

- `account` (Boolean) Returns records for the entire account.
- `compute_pool` (String) Returns records for the specified compute pool.
- `database` (String) Returns records for the current database in use or for a specified database.
- `schema` (String) Returns records for the current schema in use or a specified schema. Use fully qualified name.


<a id="nestedblock--limit"></a>
### Nested Schema for `limit`

Required:

- `rows` (Number) The maximum number of rows to return.

Optional:

- `from` (String) Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.


<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `describe_output` (List of Object) (see [below for nested schema](#nestedobjatt--services--describe_output))
- `show_output` (List of Object) (see [below for nested schema](#nestedobjatt--services--show_output))

<a id="nestedobjatt--services--describe_output"></a>
### Nested Schema for `services.describe_output`

Read-Only:

- `auto_resume` (Boolean)
- `auto_suspend_secs` (Number)
- `comment` (String)
- `compute_pool` (String)
- `created_on` (String)
- `current_instances` (Number)
- `database_name` (String)
- `dns_name` (String)
- `external_access_integrations` (Set of String)
- `is_async_job` (Boolean)
- `is_job` (Boolean)
- `is_upgrading` (Boolean)
- `managing_object_domain` (String)
- `managing_object_name` (String)
- `max_instances` (Number)
- `min_instances` (Number)
- `min_ready_instances` (Number)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `query_warehouse` (String)
- `resumed_on` (String)
- `schema_name` (String)
- `spec` (String)
- `spec_digest` (String)
- `status` (String)
- `suspended_on` (String)
- `target_instances` (Number)
- `updated_on` (String)


<a id="nestedobjatt--services--show_output"></a>
### Nested Schema for `services.show_output`

Read-Only:

- `auto_resume` (Boolean)
- `auto_suspend_secs` (Number)
- `comment` (String)
- `compute_pool` (String)
- `created_on` (String)
- `current_instances` (Number)
- `database_name` (String)
- `dns_name` (String)
- `external_access_integrations` (Set of String)
- `is_async_job` (Boolean)
- `is_job` (Boolean)
- `is_upgrading` (Boolean)
- `managing_object_domain` (String)
- `managing_object_name` (String)
- `max_instances` (Number)
- `min_instances` (Number)
- `min_ready_instances` (Number)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `query_warehouse` (String)
- `resumed_on` (String)
- `schema_name` (String)
- `spec_digest` (String)
- `status` (String)
- `suspended_on` (String)
- `target_instances` (Number)
- `updated_on` (String)
