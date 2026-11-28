---
page_title: "snowflake_object_parameter Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

!> **Warning** This resource shouldn't be used with `snowflake_current_account` resource in the same configuration to set parameters on the current account, as it may lead to unexpected behavior. Unless this resource is only used to manage the following parameters that are not supported by `snowflake_current_account`. More details in the [snowflake_current_account](./current_account) resource documentation.

-> **Note** This resource does not support all account parameters. The supported ones are listed below. This feature gap will be addressed in future releases.

~> **Note** The `OBJECT_PARAMETER_PARSE_ON_SET` experiment is enabled by default. Covered parameters are parsed with the same rules as the provider's typed parameter fields, instead of sending `value` as a raw SQL token. String values are quoted by the provider, so do not wrap them in quotes, and a value that does not match the parameter fails in the provider. Covered today: account-level (`on_account = true`) parameters and Database, Schema, Warehouse, User, Task, Service, Hybrid Table, and Iceberg Table. Anything not covered yet keeps the previous behavior. This includes enum values added by Snowflake before the provider catalog recognizes them; the provider logs a warning before falling back. To keep sending `value` as a raw SQL token, add `OBJECT_PARAMETER_PARSE_ON_SET` to the [`experimental_features_disabled`](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#experimental_features_disabled-1) provider field. See the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#new-feature-account_parameter_catalog_write_path-and-object_parameter_parse_on_set-enabled-by-default) for details.

~> **Note** Forward compatibility: parameters and object types that the provider does not support yet can be set with `value_type` (e.g. `key = "SOME_NEW_PARAMETER"`, `value = "10"`, `value_type = "INTEGER"`). `object_type` accepts any object type that Snowflake supports: for the object types that the provider does not support (e.g. `TABLE` or `PIPE`), every parameter is sent as is, and it is rendered by `value_type` when the field is set. The supported parameters are listed per object type in the `key` field description below (and the supported object types in `object_type`), and `value_type` is ignored for them (a warning is logged, see the `value_type` description). `value_type` requires the `OBJECT_PARAMETER_PARSE_ON_SET` experiment (enabled by default), because the rules for the value follow its behavior. When a later provider version starts to support the parameter or the object type, it is sent through the typed path, and the configuration does not have to change. The transition does not need configuration changes only if `value_type` matches the type of the parameter in Snowflake. Some types are easy to confuse: for example, an identifier (like a network policy name) is also accepted by Snowflake as a `STRING`, but once the provider supports the parameter, it is sent as a quoted, case-sensitive identifier, which may show a plan difference or fail. Check the parameter type in the [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/parameters) before choosing `value_type`. See the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#new-feature-forward-compatible-parameters-and-object-types-in-snowflake_object_parameter) for details.

# snowflake_object_parameter (Resource)



## Example Usage

```terraform
resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_object_parameter" "o" {
  key         = "SUSPEND_TASK_AFTER_NUM_FAILURES"
  value       = "33"
  object_type = "DATABASE"
  object_identifier {
    name = snowflake_database.d.name
  }
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_object_parameter" "o2" {
  key         = "USER_TASK_TIMEOUT_MS"
  value       = "500"
  object_type = "SCHEMA"
  object_identifier {
    database = snowflake_database.d.name
    name     = snowflake_schema.s.name
  }
}

resource "snowflake_table" "t" {
  name     = "TEST_TABLE"
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  column {
    name = "id"
    type = "NUMBER"
  }
}

resource "snowflake_object_parameter" "o3" {
  key         = "DATA_RETENTION_TIME_IN_DAYS"
  value       = "89"
  object_type = "TABLE"
  object_identifier {
    database = snowflake_database.d.name
    schema   = snowflake_schema.s.name
    name     = snowflake_table.t.name
  }
}

// Setting object parameter at account level
resource "snowflake_object_parameter" "o4" {
  key        = "DATA_RETENTION_TIME_IN_DAYS"
  value      = "89"
  on_account = true
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) Name of object parameter. Valid values are those in [object parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#object-parameters). With the `OBJECT_PARAMETER_PARSE_ON_SET` experiment (enabled by default), the provider supports the following parameters, and it parses and quotes their values itself: on the account level (`on_account = true`): `ABORT_DETACHED_QUERY` | `ACTIVE_PYTHON_PROFILER` | `ALLOW_BIND_VALUES_ACCESS` | `ALLOW_CLIENT_MFA_CACHING` | `ALLOW_ID_TOKEN` | `ALLOWED_SPCS_WORKLOAD_TYPES` | `AUTOCOMMIT` | `BASE_LOCATION_PREFIX` | `BINARY_INPUT_FORMAT` | `BINARY_OUTPUT_FORMAT` | `CATALOG` | `CATALOG_SYNC` | `CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS` | `CLIENT_ENCRYPTION_KEY_SIZE` | `CLIENT_MEMORY_LIMIT` | `CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX` | `CLIENT_METADATA_USE_SESSION_DATABASE` | `CLIENT_PREFETCH_THREADS` | `CLIENT_RESULT_CHUNK_SIZE` | `CLIENT_RESULT_COLUMN_CASE_INSENSITIVE` | `CLIENT_SESSION_KEEP_ALIVE` | `CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY` | `CLIENT_TIMESTAMP_TYPE_MAPPING` | `CORTEX_CODE_CLI_DAILY_EST_CREDIT_LIMIT_PER_USER` | `CORTEX_CODE_DESKTOP_DAILY_EST_CREDIT_LIMIT_PER_USER` | `CORTEX_CODE_SNOWSIGHT_DAILY_EST_CREDIT_LIMIT_PER_USER` | `CORTEX_ENABLED_CROSS_REGION` | `CORTEX_MODELS_ALLOWLIST` | `CSV_TIMESTAMP_FORMAT` | `DATA_METRIC_SCHEDULE` | `DATA_RETENTION_TIME_IN_DAYS` | `DATE_INPUT_FORMAT` | `DATE_OUTPUT_FORMAT` | `DEFAULT_DBT_VERSION` | `DEFAULT_DDL_COLLATION` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU` | `DEFAULT_NULL_ORDERING` | `DEFAULT_STREAMLIT_COMPUTE_POOL` | `DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE` | `DISABLE_UI_DOWNLOAD_BUTTON` | `DISABLE_USER_PRIVILEGE_GRANTS` | `DISALLOWED_SPCS_WORKLOAD_TYPES` | `ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG` | `ENABLE_BUDGET_EVENT_LOGGING` | `ENABLE_CONSOLE_OUTPUT` | `ENABLE_CORTEX_ANALYST` | `ENABLE_DATA_COMPACTION` | `ENABLE_EGRESS_COST_OPTIMIZER` | `ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS` | `ENABLE_ICEBERG_MERGE_ON_READ` | `ENABLE_IDENTIFIER_FIRST_LOGIN` | `ENABLE_INTERNAL_STAGES_PRIVATELINK` | `ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB` | `ENABLE_PER_ACCOUNT_APP_SERVICE_PRIVATELINK_URL` | `ENABLE_PERSONAL_DATABASE` | `ENABLE_SPCS_BLOCK_STORAGE_SNOWFLAKE_FULL_ENCRYPTION_ENFORCEMENT` | `ENABLE_TAG_PROPAGATION_EVENT_LOGGING` | `ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY` | `ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE` | `ENABLE_UNHANDLED_EXCEPTIONS_REPORTING` | `ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION` | `ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR` | `ENABLE_UNREDACTED_SECURE_OBJECT_ERROR` | `ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES` | `ERROR_ON_NONDETERMINISTIC_MERGE` | `ERROR_ON_NONDETERMINISTIC_UPDATE` | `EVENT_TABLE` | `EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST` | `EXTERNAL_VOLUME` | `GEOGRAPHY_OUTPUT_FORMAT` | `GEOMETRY_OUTPUT_FORMAT` | `HYBRID_TABLE_LOCK_TIMEOUT` | `ICEBERG_VERSION_DEFAULT` | `INITIAL_REPLICATION_SIZE_LIMIT_IN_TB` | `JDBC_TREAT_DECIMAL_AS_INT` | `JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC` | `JDBC_USE_SESSION_TIMEZONE` | `JS_TREAT_INTEGER_AS_BIGINT` | `JSON_INDENT` | `LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE` | `LOCK_TIMEOUT` | `LOG_EVENT_LEVEL` | `LOG_LEVEL` | `MAX_CONCURRENCY_LEVEL` | `MAX_DATA_EXTENSION_TIME_IN_DAYS` | `METRIC_LEVEL` | `MIN_DATA_RETENTION_TIME_IN_DAYS` | `MULTI_STATEMENT_COUNT` | `NETWORK_POLICY` | `NOORDER_SEQUENCE_AS_DEFAULT` | `OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST` | `ODBC_TREAT_DECIMAL_AS_INT` | `PERIODIC_DATA_REKEYING` | `PIPE_EXECUTION_PAUSED` | `PREVENT_LOAD_FROM_INLINE_URL` | `PREVENT_UNLOAD_TO_INLINE_URL` | `PREVENT_UNLOAD_TO_INTERNAL_STAGES` | `PYTHON_PROFILER_MODULES` | `PYTHON_PROFILER_TARGET_STAGE` | `QUERY_TAG` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `READ_CONSISTENCY_MODE` | `REPLACE_INVALID_CHARACTERS` | `REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION` | `REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION` | `ROW_TIMESTAMP_DEFAULT` | `ROWS_PER_RESULTSET` | `S3_STAGE_VPCE_DNS_NAME` | `SEARCH_PATH` | `SERVERLESS_TASK_MAX_STATEMENT_SIZE` | `SERVERLESS_TASK_MIN_STATEMENT_SIZE` | `SHARE_RESTRICTIONS` | `SIMULATED_DATA_SHARING_CONSUMER` | `SQL_TRACE_QUERY_TEXT` | `SSO_LOGIN_PAGE` | `STATEMENT_QUEUED_TIMEOUT_IN_SECONDS` | `STATEMENT_TIMEOUT_IN_SECONDS` | `STORAGE_SERIALIZATION_POLICY` | `STRICT_JSON_OUTPUT` | `SUSPEND_TASK_AFTER_NUM_FAILURES` | `TASK_AUTO_RETRY_ATTEMPTS` | `TIME_INPUT_FORMAT` | `TIME_OUTPUT_FORMAT` | `TIMESTAMP_DAY_IS_ALWAYS_24H` | `TIMESTAMP_INPUT_FORMAT` | `TIMESTAMP_LTZ_OUTPUT_FORMAT` | `TIMESTAMP_NTZ_OUTPUT_FORMAT` | `TIMESTAMP_OUTPUT_FORMAT` | `TIMESTAMP_TYPE_MAPPING` | `TIMESTAMP_TZ_OUTPUT_FORMAT` | `TIMEZONE` | `TRACE_LEVEL` | `TRANSACTION_ABORT_ON_ERROR` | `TRANSACTION_DEFAULT_ISOLATION_LEVEL` | `TWO_DIGIT_CENTURY_START` | `UNSUPPORTED_DDL_ACTION` | `USE_CACHED_RESULT` | `USE_WORKSPACES_FOR_SQL` | `USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE` | `USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS` | `USER_TASK_TIMEOUT_MS` | `WEEK_OF_YEAR_POLICY` | `WEEK_START`; `DATABASE`: `CATALOG` | `DATA_RETENTION_TIME_IN_DAYS` | `DEFAULT_DDL_COLLATION` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU` | `ENABLE_CONSOLE_OUTPUT` | `EXTERNAL_VOLUME` | `LOG_EVENT_LEVEL` | `LOG_LEVEL` | `MAX_DATA_EXTENSION_TIME_IN_DAYS` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `REPLACE_INVALID_CHARACTERS` | `STORAGE_SERIALIZATION_POLICY` | `SUSPEND_TASK_AFTER_NUM_FAILURES` | `TASK_AUTO_RETRY_ATTEMPTS` | `TRACE_LEVEL` | `USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE` | `USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS` | `USER_TASK_TIMEOUT_MS`; `SCHEMA`: `CATALOG` | `DATA_RETENTION_TIME_IN_DAYS` | `DEFAULT_DDL_COLLATION` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU` | `DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU` | `ENABLE_CONSOLE_OUTPUT` | `EXTERNAL_VOLUME` | `LOG_EVENT_LEVEL` | `LOG_LEVEL` | `MAX_DATA_EXTENSION_TIME_IN_DAYS` | `PIPE_EXECUTION_PAUSED` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `REPLACE_INVALID_CHARACTERS` | `STORAGE_SERIALIZATION_POLICY` | `SUSPEND_TASK_AFTER_NUM_FAILURES` | `TASK_AUTO_RETRY_ATTEMPTS` | `TRACE_LEVEL` | `USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE` | `USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS` | `USER_TASK_TIMEOUT_MS`; `WAREHOUSE`: `MAX_CONCURRENCY_LEVEL` | `STATEMENT_QUEUED_TIMEOUT_IN_SECONDS` | `STATEMENT_TIMEOUT_IN_SECONDS`; `USER`: `ABORT_DETACHED_QUERY` | `ACTIVE_PYTHON_PROFILER` | `AUTOCOMMIT` | `BINARY_INPUT_FORMAT` | `BINARY_OUTPUT_FORMAT` | `CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS` | `CLIENT_MEMORY_LIMIT` | `CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX` | `CLIENT_METADATA_USE_SESSION_DATABASE` | `CLIENT_PREFETCH_THREADS` | `CLIENT_RESULT_CHUNK_SIZE` | `CLIENT_RESULT_COLUMN_CASE_INSENSITIVE` | `CLIENT_SESSION_KEEP_ALIVE` | `CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY` | `CLIENT_TIMESTAMP_TYPE_MAPPING` | `CSV_TIMESTAMP_FORMAT` | `DATE_INPUT_FORMAT` | `DATE_OUTPUT_FORMAT` | `ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS` | `ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION` | `ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR` | `ERROR_ON_NONDETERMINISTIC_MERGE` | `ERROR_ON_NONDETERMINISTIC_UPDATE` | `GEOGRAPHY_OUTPUT_FORMAT` | `GEOMETRY_OUTPUT_FORMAT` | `HYBRID_TABLE_LOCK_TIMEOUT` | `JDBC_TREAT_DECIMAL_AS_INT` | `JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC` | `JDBC_USE_SESSION_TIMEZONE` | `JS_TREAT_INTEGER_AS_BIGINT` | `JSON_INDENT` | `LOCK_TIMEOUT` | `LOG_EVENT_LEVEL` | `LOG_LEVEL` | `MULTI_STATEMENT_COUNT` | `NETWORK_POLICY` | `NOORDER_SEQUENCE_AS_DEFAULT` | `ODBC_TREAT_DECIMAL_AS_INT` | `PREVENT_UNLOAD_TO_INTERNAL_STAGES` | `PYTHON_PROFILER_MODULES` | `PYTHON_PROFILER_TARGET_STAGE` | `QUERY_TAG` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `ROWS_PER_RESULTSET` | `S3_STAGE_VPCE_DNS_NAME` | `SEARCH_PATH` | `SIMULATED_DATA_SHARING_CONSUMER` | `STATEMENT_QUEUED_TIMEOUT_IN_SECONDS` | `STATEMENT_TIMEOUT_IN_SECONDS` | `STRICT_JSON_OUTPUT` | `TIME_INPUT_FORMAT` | `TIME_OUTPUT_FORMAT` | `TIMESTAMP_DAY_IS_ALWAYS_24H` | `TIMESTAMP_INPUT_FORMAT` | `TIMESTAMP_LTZ_OUTPUT_FORMAT` | `TIMESTAMP_NTZ_OUTPUT_FORMAT` | `TIMESTAMP_OUTPUT_FORMAT` | `TIMESTAMP_TYPE_MAPPING` | `TIMESTAMP_TZ_OUTPUT_FORMAT` | `TIMEZONE` | `TRACE_LEVEL` | `TRANSACTION_ABORT_ON_ERROR` | `TRANSACTION_DEFAULT_ISOLATION_LEVEL` | `TWO_DIGIT_CENTURY_START` | `UNSUPPORTED_DDL_ACTION` | `USE_CACHED_RESULT` | `WEEK_OF_YEAR_POLICY` | `WEEK_START`; `TASK`: `ABORT_DETACHED_QUERY` | `AUTOCOMMIT` | `BINARY_INPUT_FORMAT` | `BINARY_OUTPUT_FORMAT` | `CLIENT_MEMORY_LIMIT` | `CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX` | `CLIENT_PREFETCH_THREADS` | `CLIENT_RESULT_CHUNK_SIZE` | `CLIENT_RESULT_COLUMN_CASE_INSENSITIVE` | `CLIENT_SESSION_KEEP_ALIVE` | `CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY` | `CLIENT_TIMESTAMP_TYPE_MAPPING` | `DATE_INPUT_FORMAT` | `DATE_OUTPUT_FORMAT` | `ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION` | `ERROR_ON_NONDETERMINISTIC_MERGE` | `ERROR_ON_NONDETERMINISTIC_UPDATE` | `GEOGRAPHY_OUTPUT_FORMAT` | `GEOMETRY_OUTPUT_FORMAT` | `JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC` | `JDBC_USE_SESSION_TIMEZONE` | `JSON_INDENT` | `LOCK_TIMEOUT` | `LOG_EVENT_LEVEL` | `LOG_LEVEL` | `MULTI_STATEMENT_COUNT` | `NOORDER_SEQUENCE_AS_DEFAULT` | `ODBC_TREAT_DECIMAL_AS_INT` | `QUERY_TAG` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `ROWS_PER_RESULTSET` | `S3_STAGE_VPCE_DNS_NAME` | `SERVERLESS_TASK_MAX_STATEMENT_SIZE` | `SERVERLESS_TASK_MIN_STATEMENT_SIZE` | `STATEMENT_QUEUED_TIMEOUT_IN_SECONDS` | `STATEMENT_TIMEOUT_IN_SECONDS` | `STRICT_JSON_OUTPUT` | `SUSPEND_TASK_AFTER_NUM_FAILURES` | `TASK_AUTO_RETRY_ATTEMPTS` | `TIME_INPUT_FORMAT` | `TIME_OUTPUT_FORMAT` | `TIMESTAMP_DAY_IS_ALWAYS_24H` | `TIMESTAMP_INPUT_FORMAT` | `TIMESTAMP_LTZ_OUTPUT_FORMAT` | `TIMESTAMP_NTZ_OUTPUT_FORMAT` | `TIMESTAMP_OUTPUT_FORMAT` | `TIMESTAMP_TYPE_MAPPING` | `TIMESTAMP_TZ_OUTPUT_FORMAT` | `TIMEZONE` | `TRACE_LEVEL` | `TRANSACTION_ABORT_ON_ERROR` | `TRANSACTION_DEFAULT_ISOLATION_LEVEL` | `TWO_DIGIT_CENTURY_START` | `UNSUPPORTED_DDL_ACTION` | `USE_CACHED_RESULT` | `USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE` | `USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS` | `USER_TASK_TIMEOUT_MS` | `WEEK_OF_YEAR_POLICY` | `WEEK_START`; `SERVICE`: `SERVICE_CALLER_TOKEN_VALIDITY_SECS`; `HYBRID TABLE`: `DATA_RETENTION_TIME_IN_DAYS` | `MAX_DATA_EXTENSION_TIME_IN_DAYS`; `ICEBERG TABLE`: `ALLOW_ROW_TIMESTAMP` | `CATALOG` | `CATALOG_SYNC` | `DATA_METRIC_SCHEDULE` | `DATA_RETENTION_TIME_IN_DAYS` | `DEFAULT_DDL_COLLATION` | `ENABLE_DATA_COMPACTION` | `ENABLE_ICEBERG_MERGE_ON_READ` | `EXTERNAL_VOLUME` | `ICEBERG_MERGE_ON_READ_BEHAVIOR` | `LOG_EVENT_LEVEL` | `MAX_DATA_EXTENSION_TIME_IN_DAYS` | `OPTIMIZE_DATA_LAYOUT` | `QUOTED_IDENTIFIERS_IGNORE_CASE` | `REPLACE_INVALID_CHARACTERS` | `STORAGE_SERIALIZATION_POLICY` | `TARGET_FILE_SIZE`. Any other parameter, and any parameter on any other object type, is sent to Snowflake as is: with `value_type`, the value is rendered by it, and without it, the value is sent as written.
- `value` (String) Value of object parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.

### Optional

- `object_identifier` (Block List) Specifies the object identifier for the object parameter. If no value is provided, then the resource will default to setting the object parameter at account level. (see [below for nested schema](#nestedblock--object_identifier))
- `object_type` (String) Type of object to which the parameter applies. Valid values are those in [object types](https://docs.snowflake.com/en/sql-reference/parameters.html#object-types). If no value is provided, then the resource will default to setting the object parameter at account level. The provider supports the parameters of the following object types: `DATABASE` | `SCHEMA` | `WAREHOUSE` | `USER` | `TASK` | `SERVICE` | `HYBRID TABLE` | `ICEBERG TABLE` (see `key`). Other object types are sent as is, with the value rendered by `value_type`.
- `on_account` (Boolean) (Default: `false`) If true, the object parameter will be set on the account level.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `value_type` (String) Type of `value`, used only for parameters that the provider does not support yet (forward compatibility). Valid values are (case-insensitive): `STRING` | `IDENTIFIER` | `INTEGER` | `FLOAT` | `BOOLEAN`. `value` is rendered the same way as for the parameters that the provider supports: `STRING` (also use it for enum parameters) is sent as a single-quoted string, `IDENTIFIER` as a quoted identifier (it can be fully qualified, e.g. `DB.SCHEMA.NAME`), and the others as a number or a boolean. The field is required for a parameter that the provider does not support, and it is ignored for the supported parameters (they are listed in the `key` field description, per object type), so it can stay in the configuration when a later provider version starts to support the parameter. Changes of this field are ignored for the supported parameters, but setting it for a supported parameter logs a warning. The field works only with the OBJECT_PARAMETER_PARSE_ON_SET experiment (enabled by default). The transition does not need configuration changes only if `value_type` matches the type of the parameter in Snowflake. Some types are easy to confuse: for example, an identifier (like a network policy name) is also accepted by Snowflake as a `STRING`, but once the provider supports the parameter, it is sent as a quoted, case-sensitive identifier, which may show a plan difference or fail. Check the parameter type in the [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/parameters) before choosing `value_type`.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--object_identifier"></a>
### Nested Schema for `object_identifier`

Required:

- `name` (String) Name of the object to set the parameter for.

Optional:

- `database` (String) Name of the database that the object was created in.
- `schema` (String) Name of the schema that the object was created in.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

## Import

Import is supported using the following syntax:

```shell
# The ID is <key>|<object_type>|<object_identifier>, where <object_identifier> is the fully qualified name of the object
# with every part double-quoted: "<name>" for an account object (e.g. a user or a warehouse),
# "<database_name>"."<name>" for a schema, "<database_name>"."<schema_name>"."<name>" for a schema object.
# Quote the whole ID so the shell does not read | as a pipe.
terraform import snowflake_object_parameter.example 'DATA_RETENTION_TIME_IN_DAYS|DATABASE|"<database_name>"'
terraform import snowflake_object_parameter.example 'DATA_RETENTION_TIME_IN_DAYS|TABLE|"<database_name>"."<schema_name>"."<table_name>"'
```
