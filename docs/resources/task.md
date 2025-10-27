---
page_title: "snowflake_task Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  Resource used to manage task objects. For more information, check task documentation https://docs.snowflake.com/en/user-guide/tasks-intro.
---

!> **Sensitive values** This resource's `config`, `show_output.config` and `show_output.definition` fields are not marked as sensitive in the provider. Ensure that no personal data, sensitive data, export-controlled data, or other regulated data is entered as metadata when using the provider. For more information, see [Sensitive values limitations](../#sensitive-values-limitations) and [Metadata fields in Snowflake](https://docs.snowflake.com/en/sql-reference/metadata).

# snowflake_task (Resource)

Resource used to manage task objects. For more information, check [task documentation](https://docs.snowflake.com/en/user-guide/tasks-intro).

## Example Usage

```terraform
# Basic standalone task
resource "snowflake_task" "task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  schedule {
    minutes = 5
  }
  sql_statement = "select 1"
}

# Basic serverless task
resource "snowflake_task" "serverless_task" {
  database                                 = "database"
  schema                                   = "schema"
  name                                     = "task"
  user_task_managed_initial_warehouse_size = "XSMALL"
  started                                  = true
  schedule {
    minutes = 5
  }
  sql_statement = "select 1"
}

# Basic child task
resource "snowflake_task" "child_task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  # You can do it by referring to task by computed fully_qualified_name field or write the task name in manually if it's not managed by Terraform
  after         = [snowflake_task.root_task.fully_qualified_name, "<database_name>.<schema_name>.<root_task_name>"]
  sql_statement = "select 1"
}

# Basic finalizer task
resource "snowflake_task" "child_task" {
  database  = "database"
  schema    = "schema"
  name      = "task"
  warehouse = "warehouse"
  started   = true
  # You can do it by referring to task by computed fully_qualified_name field or write the task name in manually if it's not managed by Terraform
  finalize      = snowflake_task.root_task.fully_qualified_name
  sql_statement = "select 1"
}

# Complete child task
resource "snowflake_task" "complete_child_task" {
  database      = "database"
  schema        = "schema"
  name          = "complete_child_task"
  started       = true
  sql_statement = "CALL my_stored_procedure()"

  after                       = [snowflake_task.root_task.fully_qualified_name]
  config                      = "{\"output_dir\": \"/temp/\", \"learning_rate\": 0.1}"
  allow_overlapping_execution = true
  error_integration           = snowflake_notification_integration.example.fully_qualified_name
  when                        = "SYSTEM$STREAM_HAS_DATA('my_stream')"
  comment                     = "complete child task with all parameters"
  target_completion_interval  = "10 MINUTES"

  # Task Parameters
  suspend_task_after_num_failures               = 5
  task_auto_retry_attempts                      = 3
  user_task_managed_initial_warehouse_size      = "MEDIUM"
  user_task_minimum_trigger_interval_in_seconds = 60
  user_task_timeout_ms                          = 7200000
  serverless_task_min_statement_size            = "SMALL"
  serverless_task_max_statement_size            = "LARGE"

  # Session Parameters
  abort_detached_query         = false
  autocommit                   = true
  binary_input_format          = "HEX"
  binary_output_format         = "HEX"
  json_indent                  = 4
  lock_timeout                 = 21600
  log_level                    = "INFO"
  query_tag                    = "child_task_query"
  statement_timeout_in_seconds = 86400
  timezone                     = "America/New_York"
}

# Complete standalone task
resource "snowflake_task" "test" {
  database      = "database"
  schema        = "schema"
  name          = "task"
  warehouse     = snowflake_warehouse.example.fully_qualified_name
  started       = true
  sql_statement = "select 1"

  config                      = "{\"key\":\"value\"}"
  allow_overlapping_execution = true
  error_integration           = snowflake_notification_integration.example.fully_qualified_name
  when                        = "SYSTEM$STREAM_HAS_DATA('<stream_name>')"
  comment                     = "complete task"

  schedule {
    minutes = 10
  }

  # Session Parameters
  suspend_task_after_num_failures               = 10
  task_auto_retry_attempts                      = 0
  user_task_managed_initial_warehouse_size      = "Medium"
  user_task_minimum_trigger_interval_in_seconds = 30
  user_task_timeout_ms                          = 3600000
  abort_detached_query                          = false
  autocommit                                    = true
  binary_input_format                           = "HEX"
  binary_output_format                          = "HEX"
  client_memory_limit                           = 1536
  client_metadata_request_use_connection_ctx    = false
  client_prefetch_threads                       = 4
  client_result_chunk_size                      = 160
  client_result_column_case_insensitive         = false
  client_session_keep_alive                     = false
  client_session_keep_alive_heartbeat_frequency = 3600
  client_timestamp_type_mapping                 = "TIMESTAMP_LTZ"
  date_input_format                             = "AUTO"
  date_output_format                            = "YYYY-MM-DD"
  enable_unload_physical_type_optimization      = true
  error_on_nondeterministic_merge               = true
  error_on_nondeterministic_update              = false
  geography_output_format                       = "GeoJSON"
  geometry_output_format                        = "GeoJSON"
  jdbc_use_session_timezone                     = true
  json_indent                                   = 2
  lock_timeout                                  = 43200
  log_level                                     = "OFF"
  multi_statement_count                         = 1
  noorder_sequence_as_default                   = true
  odbc_treat_decimal_as_int                     = false
  query_tag                                     = ""
  quoted_identifiers_ignore_case                = false
  rows_per_resultset                            = 0
  s3_stage_vpce_dns_name                        = ""
  search_path                                   = "$current, $public"
  statement_queued_timeout_in_seconds           = 0
  statement_timeout_in_seconds                  = 172800
  strict_json_output                            = false
  timestamp_day_is_always_24h                   = false
  timestamp_input_format                        = "AUTO"
  timestamp_ltz_output_format                   = ""
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM"
  timestamp_type_mapping                        = "TIMESTAMP_NTZ"
  timestamp_tz_output_format                    = ""
  timezone                                      = "America/Los_Angeles"
  time_input_format                             = "AUTO"
  time_output_format                            = "HH24:MI:SS"
  trace_level                                   = "OFF"
  transaction_abort_on_error                    = false
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1970
  unsupported_ddl_action                        = "ignore"
  use_cached_result                             = true
  week_of_year_policy                           = 0
  week_start                                    = 0
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database in which to create the task. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `name` (String) Specifies the identifier for the task; must be unique for the database and schema in which the task is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `schema` (String) The schema in which to create the task. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `sql_statement` (String) Any single SQL statement, or a call to a stored procedure, executed when the task runs.
- `started` (Boolean) Specifies if the task should be started or suspended.

### Optional

- `abort_detached_query` (Boolean) Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption). For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query).
- `after` (Set of String) Specifies one or more predecessor tasks for the current task. Use this option to [create a DAG](https://docs.snowflake.com/en/user-guide/tasks-graphs.html#label-task-dag) of tasks or add this task to an existing DAG. A DAG is a series of tasks that starts with a scheduled root task and is linked together by dependencies. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `allow_overlapping_execution` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) By default, Snowflake ensures that only one instance of a particular DAG is allowed to run at a time, setting the parameter value to TRUE permits DAG runs to overlap. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `autocommit` (Boolean) Specifies whether autocommit is enabled for the session. Autocommit determines whether a DML statement, when executed without an active transaction, is automatically committed after the statement successfully completes. For more information, see [Transactions](https://docs.snowflake.com/en/sql-reference/transactions). For more information, check [AUTOCOMMIT docs](https://docs.snowflake.com/en/sql-reference/parameters#autocommit).
- `binary_input_format` (String) The format of VARCHAR values passed as input to VARCHAR-to-BINARY conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output). For more information, check [BINARY_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#binary-input-format).
- `binary_output_format` (String) The format for VARCHAR values returned as output by BINARY-to-VARCHAR conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output). For more information, check [BINARY_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#binary-output-format).
- `client_memory_limit` (Number) Parameter that specifies the maximum amount of memory the JDBC driver or ODBC driver should use for the result set from queries (in MB). For more information, check [CLIENT_MEMORY_LIMIT docs](https://docs.snowflake.com/en/sql-reference/parameters#client-memory-limit).
- `client_metadata_request_use_connection_ctx` (Boolean) For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases/schemas to the current database/schema. The narrower search typically returns fewer rows and executes more quickly. For more information, check [CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX docs](https://docs.snowflake.com/en/sql-reference/parameters#client-metadata-request-use-connection-ctx).
- `client_prefetch_threads` (Number) Parameter that specifies the number of threads used by the client to pre-fetch large result sets. The driver will attempt to honor the parameter value, but defines the minimum and maximum values (depending on your system’s resources) to improve performance. For more information, check [CLIENT_PREFETCH_THREADS docs](https://docs.snowflake.com/en/sql-reference/parameters#client-prefetch-threads).
- `client_result_chunk_size` (Number) Parameter that specifies the maximum size of each set (or chunk) of query results to download (in MB). The JDBC driver downloads query results in chunks. For more information, check [CLIENT_RESULT_CHUNK_SIZE docs](https://docs.snowflake.com/en/sql-reference/parameters#client-result-chunk-size).
- `client_result_column_case_insensitive` (Boolean) Parameter that indicates whether to match column name case-insensitively in ResultSet.get* methods in JDBC. For more information, check [CLIENT_RESULT_COLUMN_CASE_INSENSITIVE docs](https://docs.snowflake.com/en/sql-reference/parameters#client-result-column-case-insensitive).
- `client_session_keep_alive` (Boolean) Parameter that indicates whether to force a user to log in again after a period of inactivity in the session. For more information, check [CLIENT_SESSION_KEEP_ALIVE docs](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive).
- `client_session_keep_alive_heartbeat_frequency` (Number) Number of seconds in-between client attempts to update the token for the session. For more information, check [CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY docs](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive-heartbeat-frequency).
- `client_timestamp_type_mapping` (String) Specifies the [TIMESTAMP_* variation](https://docs.snowflake.com/en/sql-reference/data-types-datetime.html#label-datatypes-timestamp-variations) to use when binding timestamp variables for JDBC or ODBC applications that use the bind API to load data. For more information, check [CLIENT_TIMESTAMP_TYPE_MAPPING docs](https://docs.snowflake.com/en/sql-reference/parameters#client-timestamp-type-mapping).
- `comment` (String) Specifies a comment for the task.
- `config` (String) Specifies a string representation of key value pairs that can be accessed by all tasks in the task graph. Must be in JSON format.
- `date_input_format` (String) Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [DATE_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#date-input-format).
- `date_output_format` (String) Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [DATE_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#date-output-format).
- `enable_unload_physical_type_optimization` (Boolean) Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table). For more information, check [ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-unload-physical-type-optimization).
- `error_integration` (String) Specifies the name of the notification integration used for error notifications. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`. For more information about this resource, see [docs](./notification_integration).
- `error_on_nondeterministic_merge` (Boolean) Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row. For more information, check [ERROR_ON_NONDETERMINISTIC_MERGE docs](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-merge).
- `error_on_nondeterministic_update` (Boolean) Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row. For more information, check [ERROR_ON_NONDETERMINISTIC_UPDATE docs](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-update).
- `finalize` (String) Specifies the name of a root task that the finalizer task is associated with. Finalizer tasks run after all other tasks in the task graph run to completion. You can define the SQL of a finalizer task to handle notifications and the release and cleanup of resources that a task graph uses. For more information, see [Release and cleanup of task graphs](https://docs.snowflake.com/en/user-guide/tasks-graphs.html#label-finalizer-task). Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `geography_output_format` (String) Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography). For more information, check [GEOGRAPHY_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#geography-output-format).
- `geometry_output_format` (String) Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry). For more information, check [GEOMETRY_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#geometry-output-format).
- `jdbc_treat_timestamp_ntz_as_utc` (Boolean) Specifies how JDBC processes TIMESTAMP_NTZ values. For more information, check [JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC docs](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-timestamp-ntz-as-utc).
- `jdbc_use_session_timezone` (Boolean) Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class. For more information, check [JDBC_USE_SESSION_TIMEZONE docs](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-use-session-timezone).
- `json_indent` (Number) Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element. For more information, check [JSON_INDENT docs](https://docs.snowflake.com/en/sql-reference/parameters#json-indent).
- `lock_timeout` (Number) Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement. For more information, check [LOCK_TIMEOUT docs](https://docs.snowflake.com/en/sql-reference/parameters#lock-timeout).
- `log_level` (String) Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested. For more information about log levels, see [Setting log level](https://docs.snowflake.com/en/developer-guide/logging-tracing/logging-log-level). For more information, check [LOG_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#log-level).
- `multi_statement_count` (Number) Number of statements to execute when using the multi-statement capability. For more information, check [MULTI_STATEMENT_COUNT docs](https://docs.snowflake.com/en/sql-reference/parameters#multi-statement-count).
- `noorder_sequence_as_default` (Boolean) Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values). For more information, check [NOORDER_SEQUENCE_AS_DEFAULT docs](https://docs.snowflake.com/en/sql-reference/parameters#noorder-sequence-as-default).
- `odbc_treat_decimal_as_int` (Boolean) Specifies how ODBC processes columns that have a scale of zero (0). For more information, check [ODBC_TREAT_DECIMAL_AS_INT docs](https://docs.snowflake.com/en/sql-reference/parameters#odbc-treat-decimal-as-int).
- `query_tag` (String) Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions. For more information, check [QUERY_TAG docs](https://docs.snowflake.com/en/sql-reference/parameters#query-tag).
- `quoted_identifiers_ignore_case` (Boolean) Specifies whether letters in double-quoted object identifiers are stored and resolved as uppercase letters. By default, Snowflake preserves the case of alphabetic characters when storing and resolving double-quoted identifiers (see [Identifier resolution](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing)). You can use this parameter in situations in which [third-party applications always use double quotes around identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing-parameter). For more information, check [QUOTED_IDENTIFIERS_IGNORE_CASE docs](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case).
- `rows_per_resultset` (Number) Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum. For more information, check [ROWS_PER_RESULTSET docs](https://docs.snowflake.com/en/sql-reference/parameters#rows-per-resultset).
- `s3_stage_vpce_dns_name` (String) Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation). For more information, check [S3_STAGE_VPCE_DNS_NAME docs](https://docs.snowflake.com/en/sql-reference/parameters#s3-stage-vpce-dns-name).
- `schedule` (Block List, Max: 1) The schedule for periodically running the task. This can be a cron or interval in minutes. (Conflicts with finalize and after; when set, one of the sub-fields `minutes` or `using_cron` should be set) (see [below for nested schema](#nestedblock--schedule))
- `search_path` (String) Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name. For more information, check [SEARCH_PATH docs](https://docs.snowflake.com/en/sql-reference/parameters#search-path).
- `serverless_task_max_statement_size` (String) Specifies the maximum warehouse size for serverless tasks. Valid values are (case-insensitive): `XSMALL` | `X-SMALL` | `SMALL` | `MEDIUM` | `LARGE` | `XLARGE` | `X-LARGE` | `XXLARGE` | `X2LARGE` | `2X-LARGE` | `XXXLARGE` | `X3LARGE` | `3X-LARGE` | `X4LARGE` | `4X-LARGE` | `X5LARGE` | `5X-LARGE` | `X6LARGE` | `6X-LARGE`. For more information, check [SERVERLESS_TASK_MAX_STATEMENT_SIZE docs](https://docs.snowflake.com/en/sql-reference/parameters#serverless-task-max-statement-size).
- `serverless_task_min_statement_size` (String) Specifies the minimum warehouse size for serverless tasks. Valid values are (case-insensitive): `XSMALL` | `X-SMALL` | `SMALL` | `MEDIUM` | `LARGE` | `XLARGE` | `X-LARGE` | `XXLARGE` | `X2LARGE` | `2X-LARGE` | `XXXLARGE` | `X3LARGE` | `3X-LARGE` | `X4LARGE` | `4X-LARGE` | `X5LARGE` | `5X-LARGE` | `X6LARGE` | `6X-LARGE`. For more information, check [SERVERLESS_TASK_MIN_STATEMENT_SIZE docs](https://docs.snowflake.com/en/sql-reference/parameters#serverless-task-min-statement-size).
- `statement_queued_timeout_in_seconds` (Number) Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged. For more information, check [STATEMENT_QUEUED_TIMEOUT_IN_SECONDS docs](https://docs.snowflake.com/en/sql-reference/parameters#statement-queued-timeout-in-seconds).
- `statement_timeout_in_seconds` (Number) Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system. For more information, check [STATEMENT_TIMEOUT_IN_SECONDS docs](https://docs.snowflake.com/en/sql-reference/parameters#statement-timeout-in-seconds).
- `strict_json_output` (Boolean) This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON. For more information, check [STRICT_JSON_OUTPUT docs](https://docs.snowflake.com/en/sql-reference/parameters#strict-json-output).
- `suspend_task_after_num_failures` (Number) Specifies the number of consecutive failed task runs after which the current task is suspended automatically. The default is 0 (no automatic suspension). For more information, check [SUSPEND_TASK_AFTER_NUM_FAILURES docs](https://docs.snowflake.com/en/sql-reference/parameters#suspend-task-after-num-failures).
- `target_completion_interval` (String) Specifies the target completion interval for tasks. Format: '<num> { HOURS | MINUTES | SECONDS }' (e.g., '10 MINUTES'). For more information, check [TARGET_COMPLETION_INTERVAL docs](https://docs.snowflake.com/en/sql-reference/parameters#target-completion-interval).
- `task_auto_retry_attempts` (Number) Specifies the number of automatic task graph retry attempts. If any task graphs complete in a FAILED state, Snowflake can automatically retry the task graphs from the last task in the graph that failed. For more information, check [TASK_AUTO_RETRY_ATTEMPTS docs](https://docs.snowflake.com/en/sql-reference/parameters#task-auto-retry-attempts).
- `time_input_format` (String) Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session). For more information, check [TIME_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#time-input-format).
- `time_output_format` (String) Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [TIME_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#time-output-format).
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `timestamp_day_is_always_24h` (Boolean) Specifies whether the [DATEADD](https://docs.snowflake.com/en/sql-reference/functions/dateadd) function (and its aliases) always consider a day to be exactly 24 hours for expressions that span multiple days. For more information, check [TIMESTAMP_DAY_IS_ALWAYS_24H docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-day-is-always-24h).
- `timestamp_input_format` (String) Specifies the input format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported timestamp format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of timestamps stored in the system during the session). For more information, check [TIMESTAMP_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-input-format).
- `timestamp_ltz_output_format` (String) Specifies the display format for the TIMESTAMP_LTZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [TIMESTAMP_LTZ_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ltz-output-format).
- `timestamp_ntz_output_format` (String) Specifies the display format for the TIMESTAMP_NTZ data type. For more information, check [TIMESTAMP_NTZ_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ntz-output-format).
- `timestamp_output_format` (String) Specifies the display format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [TIMESTAMP_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-output-format).
- `timestamp_type_mapping` (String) Specifies the TIMESTAMP_* variation that the TIMESTAMP data type alias maps to. For more information, check [TIMESTAMP_TYPE_MAPPING docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-type-mapping).
- `timestamp_tz_output_format` (String) Specifies the display format for the TIMESTAMP_TZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [TIMESTAMP_TZ_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-tz-output-format).
- `timezone` (String) Specifies the time zone for the session. You can specify a [time zone name](https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab) or a [link name](https://data.iana.org/time-zones/tzdb-2021a/backward) from release 2021a of the [IANA Time Zone Database](https://www.iana.org/time-zones) (e.g. America/Los_Angeles, Europe/London, UTC, Etc/GMT, etc.). For more information, check [TIMEZONE docs](https://docs.snowflake.com/en/sql-reference/parameters#timezone).
- `trace_level` (String) Controls how trace events are ingested into the event table. For more information about trace levels, see [Setting trace level](https://docs.snowflake.com/en/developer-guide/logging-tracing/tracing-trace-level). For more information, check [TRACE_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#trace-level).
- `transaction_abort_on_error` (Boolean) Specifies the action to perform when a statement issued within a non-autocommit transaction returns with an error. For more information, check [TRANSACTION_ABORT_ON_ERROR docs](https://docs.snowflake.com/en/sql-reference/parameters#transaction-abort-on-error).
- `transaction_default_isolation_level` (String) Specifies the isolation level for transactions in the user session. For more information, check [TRANSACTION_DEFAULT_ISOLATION_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#transaction-default-isolation-level).
- `two_digit_century_start` (Number) Specifies the “century start” year for 2-digit years (i.e. the earliest year such dates can represent). This parameter prevents ambiguous dates when importing or converting data with the `YY` date format component (i.e. years represented as 2 digits). For more information, check [TWO_DIGIT_CENTURY_START docs](https://docs.snowflake.com/en/sql-reference/parameters#two-digit-century-start).
- `unsupported_ddl_action` (String) Determines if an unsupported (i.e. non-default) value specified for a constraint property returns an error. For more information, check [UNSUPPORTED_DDL_ACTION docs](https://docs.snowflake.com/en/sql-reference/parameters#unsupported-ddl-action).
- `use_cached_result` (Boolean) Specifies whether to reuse persisted query results, if available, when a matching query is submitted. For more information, check [USE_CACHED_RESULT docs](https://docs.snowflake.com/en/sql-reference/parameters#use-cached-result).
- `user_task_managed_initial_warehouse_size` (String) Specifies the size of the compute resources to provision for the first run of the task, before a task history is available for Snowflake to determine an ideal size. Once a task has successfully completed a few runs, Snowflake ignores this parameter setting. Valid values are (case-insensitive): %s. (Conflicts with warehouse). For more information about warehouses, see [docs](./warehouse). For more information, check [USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE docs](https://docs.snowflake.com/en/sql-reference/parameters#user-task-managed-initial-warehouse-size).
- `user_task_minimum_trigger_interval_in_seconds` (Number) Minimum amount of time between Triggered Task executions in seconds For more information, check [USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS docs](https://docs.snowflake.com/en/sql-reference/parameters#user-task-minimum-trigger-interval-in-seconds).
- `user_task_timeout_ms` (Number) Specifies the time limit on a single run of the task before it times out (in milliseconds). For more information, check [USER_TASK_TIMEOUT_MS docs](https://docs.snowflake.com/en/sql-reference/parameters#user-task-timeout-ms).
- `warehouse` (String) The warehouse the task will use. Omit this parameter to use Snowflake-managed compute resources for runs of this task. Due to Snowflake limitations warehouse identifier can consist of only upper-cased letters. (Conflicts with user_task_managed_initial_warehouse_size) For more information about this resource, see [docs](./warehouse).
- `week_of_year_policy` (Number) Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year. For more information, check [WEEK_OF_YEAR_POLICY docs](https://docs.snowflake.com/en/sql-reference/parameters#week-of-year-policy).
- `week_start` (Number) Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week. For more information, check [WEEK_START docs](https://docs.snowflake.com/en/sql-reference/parameters#week-start).
- `when` (String) Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported. When a task is triggered (based on its SCHEDULE or AFTER setting), it validates the conditions of the expression to determine whether to execute. If the conditions of the expression are not met, then the task skips the current run. Any tasks that identify this task as a predecessor also don’t run.

### Read-Only

- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `parameters` (List of Object) Outputs the result of `SHOW PARAMETERS IN TASK` for the given task. (see [below for nested schema](#nestedatt--parameters))
- `show_output` (List of Object) Outputs the result of `SHOW TASKS` for the given task. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--schedule"></a>
### Nested Schema for `schedule`

Optional:

- `minutes` (Number) Specifies an interval (in minutes) of wait time inserted between runs of the task. Accepts positive integers only. (conflicts with `using_cron`)
- `using_cron` (String) Specifies a cron expression and time zone for periodically running the task. Supports a subset of standard cron utility syntax. (conflicts with `minutes`)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--parameters"></a>
### Nested Schema for `parameters`

Read-Only:

- `abort_detached_query` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--abort_detached_query))
- `autocommit` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--autocommit))
- `binary_input_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--binary_input_format))
- `binary_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--binary_output_format))
- `client_memory_limit` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_memory_limit))
- `client_metadata_request_use_connection_ctx` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_metadata_request_use_connection_ctx))
- `client_prefetch_threads` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_prefetch_threads))
- `client_result_chunk_size` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_result_chunk_size))
- `client_result_column_case_insensitive` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_result_column_case_insensitive))
- `client_session_keep_alive` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_session_keep_alive))
- `client_session_keep_alive_heartbeat_frequency` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_session_keep_alive_heartbeat_frequency))
- `client_timestamp_type_mapping` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--client_timestamp_type_mapping))
- `date_input_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--date_input_format))
- `date_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--date_output_format))
- `enable_unload_physical_type_optimization` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--enable_unload_physical_type_optimization))
- `error_on_nondeterministic_merge` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--error_on_nondeterministic_merge))
- `error_on_nondeterministic_update` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--error_on_nondeterministic_update))
- `geography_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--geography_output_format))
- `geometry_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--geometry_output_format))
- `jdbc_treat_timestamp_ntz_as_utc` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--jdbc_treat_timestamp_ntz_as_utc))
- `jdbc_use_session_timezone` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--jdbc_use_session_timezone))
- `json_indent` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--json_indent))
- `lock_timeout` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--lock_timeout))
- `log_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--log_level))
- `multi_statement_count` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--multi_statement_count))
- `noorder_sequence_as_default` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--noorder_sequence_as_default))
- `odbc_treat_decimal_as_int` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--odbc_treat_decimal_as_int))
- `query_tag` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--query_tag))
- `quoted_identifiers_ignore_case` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--quoted_identifiers_ignore_case))
- `rows_per_resultset` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--rows_per_resultset))
- `s3_stage_vpce_dns_name` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--s3_stage_vpce_dns_name))
- `search_path` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--search_path))
- `statement_queued_timeout_in_seconds` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--statement_queued_timeout_in_seconds))
- `statement_timeout_in_seconds` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--statement_timeout_in_seconds))
- `strict_json_output` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--strict_json_output))
- `suspend_task_after_num_failures` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--suspend_task_after_num_failures))
- `task_auto_retry_attempts` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--task_auto_retry_attempts))
- `time_input_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--time_input_format))
- `time_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--time_output_format))
- `timestamp_day_is_always_24h` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_day_is_always_24h))
- `timestamp_input_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_input_format))
- `timestamp_ltz_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_ltz_output_format))
- `timestamp_ntz_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_ntz_output_format))
- `timestamp_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_output_format))
- `timestamp_type_mapping` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_type_mapping))
- `timestamp_tz_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timestamp_tz_output_format))
- `timezone` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--timezone))
- `trace_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--trace_level))
- `transaction_abort_on_error` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--transaction_abort_on_error))
- `transaction_default_isolation_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--transaction_default_isolation_level))
- `two_digit_century_start` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--two_digit_century_start))
- `unsupported_ddl_action` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--unsupported_ddl_action))
- `use_cached_result` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--use_cached_result))
- `user_task_managed_initial_warehouse_size` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--user_task_managed_initial_warehouse_size))
- `user_task_minimum_trigger_interval_in_seconds` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--user_task_minimum_trigger_interval_in_seconds))
- `user_task_timeout_ms` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--user_task_timeout_ms))
- `week_of_year_policy` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--week_of_year_policy))
- `week_start` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--week_start))

<a id="nestedobjatt--parameters--abort_detached_query"></a>
### Nested Schema for `parameters.abort_detached_query`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--autocommit"></a>
### Nested Schema for `parameters.autocommit`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--binary_input_format"></a>
### Nested Schema for `parameters.binary_input_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--binary_output_format"></a>
### Nested Schema for `parameters.binary_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_memory_limit"></a>
### Nested Schema for `parameters.client_memory_limit`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_metadata_request_use_connection_ctx"></a>
### Nested Schema for `parameters.client_metadata_request_use_connection_ctx`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_prefetch_threads"></a>
### Nested Schema for `parameters.client_prefetch_threads`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_result_chunk_size"></a>
### Nested Schema for `parameters.client_result_chunk_size`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_result_column_case_insensitive"></a>
### Nested Schema for `parameters.client_result_column_case_insensitive`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_session_keep_alive"></a>
### Nested Schema for `parameters.client_session_keep_alive`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_session_keep_alive_heartbeat_frequency"></a>
### Nested Schema for `parameters.client_session_keep_alive_heartbeat_frequency`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--client_timestamp_type_mapping"></a>
### Nested Schema for `parameters.client_timestamp_type_mapping`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--date_input_format"></a>
### Nested Schema for `parameters.date_input_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--date_output_format"></a>
### Nested Schema for `parameters.date_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--enable_unload_physical_type_optimization"></a>
### Nested Schema for `parameters.enable_unload_physical_type_optimization`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--error_on_nondeterministic_merge"></a>
### Nested Schema for `parameters.error_on_nondeterministic_merge`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--error_on_nondeterministic_update"></a>
### Nested Schema for `parameters.error_on_nondeterministic_update`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--geography_output_format"></a>
### Nested Schema for `parameters.geography_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--geometry_output_format"></a>
### Nested Schema for `parameters.geometry_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--jdbc_treat_timestamp_ntz_as_utc"></a>
### Nested Schema for `parameters.jdbc_treat_timestamp_ntz_as_utc`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--jdbc_use_session_timezone"></a>
### Nested Schema for `parameters.jdbc_use_session_timezone`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--json_indent"></a>
### Nested Schema for `parameters.json_indent`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--lock_timeout"></a>
### Nested Schema for `parameters.lock_timeout`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--log_level"></a>
### Nested Schema for `parameters.log_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--multi_statement_count"></a>
### Nested Schema for `parameters.multi_statement_count`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--noorder_sequence_as_default"></a>
### Nested Schema for `parameters.noorder_sequence_as_default`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--odbc_treat_decimal_as_int"></a>
### Nested Schema for `parameters.odbc_treat_decimal_as_int`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--query_tag"></a>
### Nested Schema for `parameters.query_tag`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--quoted_identifiers_ignore_case"></a>
### Nested Schema for `parameters.quoted_identifiers_ignore_case`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--rows_per_resultset"></a>
### Nested Schema for `parameters.rows_per_resultset`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--s3_stage_vpce_dns_name"></a>
### Nested Schema for `parameters.s3_stage_vpce_dns_name`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--search_path"></a>
### Nested Schema for `parameters.search_path`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--statement_queued_timeout_in_seconds"></a>
### Nested Schema for `parameters.statement_queued_timeout_in_seconds`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--statement_timeout_in_seconds"></a>
### Nested Schema for `parameters.statement_timeout_in_seconds`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--strict_json_output"></a>
### Nested Schema for `parameters.strict_json_output`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--suspend_task_after_num_failures"></a>
### Nested Schema for `parameters.suspend_task_after_num_failures`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--task_auto_retry_attempts"></a>
### Nested Schema for `parameters.task_auto_retry_attempts`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--time_input_format"></a>
### Nested Schema for `parameters.time_input_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--time_output_format"></a>
### Nested Schema for `parameters.time_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_day_is_always_24h"></a>
### Nested Schema for `parameters.timestamp_day_is_always_24h`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_input_format"></a>
### Nested Schema for `parameters.timestamp_input_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_ltz_output_format"></a>
### Nested Schema for `parameters.timestamp_ltz_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_ntz_output_format"></a>
### Nested Schema for `parameters.timestamp_ntz_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_output_format"></a>
### Nested Schema for `parameters.timestamp_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_type_mapping"></a>
### Nested Schema for `parameters.timestamp_type_mapping`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timestamp_tz_output_format"></a>
### Nested Schema for `parameters.timestamp_tz_output_format`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--timezone"></a>
### Nested Schema for `parameters.timezone`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--trace_level"></a>
### Nested Schema for `parameters.trace_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--transaction_abort_on_error"></a>
### Nested Schema for `parameters.transaction_abort_on_error`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--transaction_default_isolation_level"></a>
### Nested Schema for `parameters.transaction_default_isolation_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--two_digit_century_start"></a>
### Nested Schema for `parameters.two_digit_century_start`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--unsupported_ddl_action"></a>
### Nested Schema for `parameters.unsupported_ddl_action`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--use_cached_result"></a>
### Nested Schema for `parameters.use_cached_result`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--user_task_managed_initial_warehouse_size"></a>
### Nested Schema for `parameters.user_task_managed_initial_warehouse_size`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--user_task_minimum_trigger_interval_in_seconds"></a>
### Nested Schema for `parameters.user_task_minimum_trigger_interval_in_seconds`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--user_task_timeout_ms"></a>
### Nested Schema for `parameters.user_task_timeout_ms`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--week_of_year_policy"></a>
### Nested Schema for `parameters.week_of_year_policy`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--week_start"></a>
### Nested Schema for `parameters.week_start`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)



<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `allow_overlapping_execution` (Boolean)
- `budget` (String)
- `comment` (String)
- `condition` (String)
- `config` (String)
- `created_on` (String)
- `database_name` (String)
- `definition` (String)
- `error_integration` (String)
- `id` (String)
- `last_committed_on` (String)
- `last_suspended_on` (String)
- `last_suspended_reason` (String)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `predecessors` (Set of String)
- `schedule` (String)
- `schema_name` (String)
- `state` (String)
- `task_relations` (List of Object) (see [below for nested schema](#nestedobjatt--show_output--task_relations))
- `warehouse` (String)

<a id="nestedobjatt--show_output--task_relations"></a>
### Nested Schema for `show_output.task_relations`

Read-Only:

- `finalized_root_task` (String)
- `finalizer` (String)
- `predecessors` (List of String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_task.example '"<database_name>"."<schema_name>"."<task_name>"'
```
