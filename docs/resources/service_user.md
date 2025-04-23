---
page_title: "snowflake_service_user Resource - terraform-provider-snowflake"
subcategory: "Stable"
description: |-
  Resource used to manage service user objects. For more information, check user documentation https://docs.snowflake.com/en/sql-reference/commands-user-role#user-management.
---

!> **Caution** Use `network_policy` attribute instead of the [`snowflake_network_policy_attachment`](./network_policy_attachment) resource. `snowflake_network_policy_attachment` will be reworked in the following versions of the provider which may still affect this resource.

!> **Sensitive values** This resource's `display_name`, `show_output.display_name`, `show_output.email`, `show_output.login_name`, `show_output.first_name` and `show_output.last_name` fields are not marked as sensitive in the provider. Ensure that no personal data, sensitive data, export-controlled data, or other regulated data is entered as metadata when using the provider. If you use one of these fields, they may be present in logs, so ensure that the provider logs are properly restricted. For more information, see [Sensitive values limitations](../#sensitive-values-limitations) and [Metadata fields in Snowflake](https://docs.snowflake.com/en/sql-reference/metadata).

-> **Note** `snowflake_user_password_policy_attachment` will be reworked in the following versions of the provider which may still affect this resource.

-> **Note** Attaching user policies will be handled in the following versions of the provider which may still affect this resource.

-> **Note** Other two user types are handled in separate resources: `snowflake_legacy_service_user` for user type `legacy_service` and `snowflake_user` for user type `person`.

-> **Note** External changes to `days_to_expiry` and `mins_to_unlock` are not currently handled by the provider (because the value changes continuously on Snowflake side after setting it).

# snowflake_service_user (Resource)

Resource used to manage service user objects. For more information, check [user documentation](https://docs.snowflake.com/en/sql-reference/commands-user-role#user-management).

## Example Usage

```terraform
# minimal
resource "snowflake_service_user" "minimal" {
  name = "Snowflake Service User - minimal"
}

# with all attributes set
resource "snowflake_service_user" "service_user" {
  name         = "Snowflake Service User"
  login_name   = var.login_name
  comment      = "A service user of snowflake."
  disabled     = "false"
  display_name = "Snowflake Service User"
  email        = var.email

  default_warehouse              = "warehouse"
  default_secondary_roles_option = "ALL"
  default_role                   = "role1"
  default_namespace              = "some.namespace"

  mins_to_unlock = 9
  days_to_expiry = 8

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}

# all parameters set on the resource level
resource "snowflake_service_user" "u" {
  name = "Snowflake Service User with all parameters"

  abort_detached_query                          = true
  autocommit                                    = false
  binary_input_format                           = "UTF8"
  binary_output_format                          = "BASE64"
  client_memory_limit                           = 1024
  client_metadata_request_use_connection_ctx    = true
  client_prefetch_threads                       = 2
  client_result_chunk_size                      = 48
  client_result_column_case_insensitive         = true
  client_session_keep_alive                     = true
  client_session_keep_alive_heartbeat_frequency = 2400
  client_timestamp_type_mapping                 = "TIMESTAMP_NTZ"
  date_input_format                             = "YYYY-MM-DD"
  date_output_format                            = "YY-MM-DD"
  enable_unload_physical_type_optimization      = false
  enable_unredacted_query_syntax_error          = true
  error_on_nondeterministic_merge               = false
  error_on_nondeterministic_update              = true
  geography_output_format                       = "WKB"
  geometry_output_format                        = "WKB"
  jdbc_treat_decimal_as_int                     = false
  jdbc_treat_timestamp_ntz_as_utc               = true
  jdbc_use_session_timezone                     = false
  json_indent                                   = 4
  lock_timeout                                  = 21222
  log_level                                     = "ERROR"
  multi_statement_count                         = 0
  network_policy                                = "BVYDGRAT_0D5E3DD1_F644_03DE_318A_1179886518A7"
  noorder_sequence_as_default                   = false
  odbc_treat_decimal_as_int                     = true
  prevent_unload_to_internal_stages             = true
  query_tag                                     = "some_tag"
  quoted_identifiers_ignore_case                = true
  rows_per_resultset                            = 2
  search_path                                   = "$public, $current"
  simulated_data_sharing_consumer               = "some_consumer"
  statement_queued_timeout_in_seconds           = 10
  statement_timeout_in_seconds                  = 10
  strict_json_output                            = true
  s3_stage_vpce_dns_name                        = "vpce-id.s3.region.vpce.amazonaws.com"
  time_input_format                             = "HH24:MI"
  time_output_format                            = "HH24:MI"
  timestamp_day_is_always_24h                   = true
  timestamp_input_format                        = "YYYY-MM-DD"
  timestamp_ltz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS"
  timestamp_type_mapping                        = "TIMESTAMP_LTZ"
  timestamp_tz_output_format                    = "YYYY-MM-DD HH24:MI:SS"
  timezone                                      = "Europe/Warsaw"
  trace_level                                   = "ON_EVENT"
  transaction_abort_on_error                    = true
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1980
  unsupported_ddl_action                        = "FAIL"
  use_cached_result                             = false
  week_of_year_policy                           = 1
  week_start                                    = 1
}

variable "email" {
  type      = string
  sensitive = true
}

variable "login_name" {
  type      = string
  sensitive = true
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the user. Note that if you do not supply login_name this will be used as login_name. Check the [docs](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters). Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `abort_detached_query` (Boolean) Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption). For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query).
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
- `comment` (String) Specifies a comment for the user.
- `date_input_format` (String) Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [DATE_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#date-input-format).
- `date_output_format` (String) Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [DATE_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#date-output-format).
- `days_to_expiry` (Number) Specifies the number of days after which the user status is set to `Expired` and the user is no longer allowed to log in. This is useful for defining temporary users (i.e. users who should only have access to Snowflake for a limited time period). In general, you should not set this property for [account administrators](https://docs.snowflake.com/en/user-guide/security-access-control-considerations.html#label-accountadmin-users) (i.e. users with the `ACCOUNTADMIN` role) because Snowflake locks them out when they become `Expired`. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `default_namespace` (String) Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the namespace exists.
- `default_role` (String) Specifies the role that is active by default for the user’s session upon login. Note that specifying a default role for a user does **not** grant the role to the user. The role must be granted explicitly to the user using the [GRANT ROLE](https://docs.snowflake.com/en/sql-reference/sql/grant-role) command. In addition, the CREATE USER operation does not verify that the role exists. For more information about this resource, see [docs](./account_role).
- `default_secondary_roles_option` (String) (Default: `DEFAULT`) Specifies the secondary roles that are active for the user’s session upon login. Valid values are (case-insensitive): `DEFAULT` | `NONE` | `ALL`. More information can be found in [doc](https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties).
- `default_warehouse` (String) Specifies the virtual warehouse that is active by default for the user’s session upon login. Note that the CREATE USER operation does not verify that the warehouse exists. For more information about this resource, see [docs](./warehouse).
- `disabled` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether the user is disabled, which prevents logging in and aborts all the currently-running queries for the user. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `display_name` (String) Name displayed for the user in the Snowflake web interface.
- `email` (String, Sensitive) Email address for the user.
- `enable_unload_physical_type_optimization` (Boolean) Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table). For more information, check [ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-unload-physical-type-optimization).
- `enable_unredacted_query_syntax_error` (Boolean) Controls whether query text is redacted if a SQL query fails due to a syntax or parsing error. If `FALSE`, the content of a failed query is redacted in the views, pages, and functions that provide a query history. Only users with a role that is granted or inherits the AUDIT privilege can set the ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR parameter. When using the ALTER USER command to set the parameter to `TRUE` for a particular user, modify the user that you want to see the query text, not the user who executed the query (if those are different users). For more information, check [ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-unredacted-query-syntax-error).
- `error_on_nondeterministic_merge` (Boolean) Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row. For more information, check [ERROR_ON_NONDETERMINISTIC_MERGE docs](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-merge).
- `error_on_nondeterministic_update` (Boolean) Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row. For more information, check [ERROR_ON_NONDETERMINISTIC_UPDATE docs](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-update).
- `geography_output_format` (String) Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography). For more information, check [GEOGRAPHY_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#geography-output-format).
- `geometry_output_format` (String) Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry). For more information, check [GEOMETRY_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#geometry-output-format).
- `jdbc_treat_decimal_as_int` (Boolean) Specifies how JDBC processes columns that have a scale of zero (0). For more information, check [JDBC_TREAT_DECIMAL_AS_INT docs](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-decimal-as-int).
- `jdbc_treat_timestamp_ntz_as_utc` (Boolean) Specifies how JDBC processes TIMESTAMP_NTZ values. For more information, check [JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC docs](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-timestamp-ntz-as-utc).
- `jdbc_use_session_timezone` (Boolean) Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class. For more information, check [JDBC_USE_SESSION_TIMEZONE docs](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-use-session-timezone).
- `json_indent` (Number) Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element. For more information, check [JSON_INDENT docs](https://docs.snowflake.com/en/sql-reference/parameters#json-indent).
- `lock_timeout` (Number) Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement. For more information, check [LOCK_TIMEOUT docs](https://docs.snowflake.com/en/sql-reference/parameters#lock-timeout).
- `log_level` (String) Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested. For more information about log levels, see [Setting log level](https://docs.snowflake.com/en/developer-guide/logging-tracing/logging-log-level). For more information, check [LOG_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#log-level).
- `login_name` (String, Sensitive) The name users use to log in. If not supplied, snowflake will use name instead. Login names are always case-insensitive.
- `mins_to_unlock` (Number) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`-1`)) Specifies the number of minutes until the temporary lock on the user login is cleared. To protect against unauthorized user login, Snowflake places a temporary lock on a user after five consecutive unsuccessful login attempts. When creating a user, this property can be set to prevent them from logging in until the specified amount of time passes. To remove a lock immediately for a user, specify a value of 0 for this parameter. **Note** because this value changes continuously after setting it, the provider is currently NOT handling the external changes to it. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `multi_statement_count` (Number) Number of statements to execute when using the multi-statement capability. For more information, check [MULTI_STATEMENT_COUNT docs](https://docs.snowflake.com/en/sql-reference/parameters#multi-statement-count).
- `network_policy` (String) Specifies the network policy to enforce for your account. Network policies enable restricting access to your account based on users’ IP address. For more details, see [Controlling network traffic with network policies](https://docs.snowflake.com/en/user-guide/network-policies). Any existing network policy (created using [CREATE NETWORK POLICY](https://docs.snowflake.com/en/sql-reference/sql/create-network-policy)). For more information, check [NETWORK_POLICY docs](https://docs.snowflake.com/en/sql-reference/parameters#network-policy).
- `noorder_sequence_as_default` (Boolean) Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values). For more information, check [NOORDER_SEQUENCE_AS_DEFAULT docs](https://docs.snowflake.com/en/sql-reference/parameters#noorder-sequence-as-default).
- `odbc_treat_decimal_as_int` (Boolean) Specifies how ODBC processes columns that have a scale of zero (0). For more information, check [ODBC_TREAT_DECIMAL_AS_INT docs](https://docs.snowflake.com/en/sql-reference/parameters#odbc-treat-decimal-as-int).
- `prevent_unload_to_internal_stages` (Boolean) Specifies whether to prevent data unload operations to internal (Snowflake) stages using [COPY INTO <location>](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location) statements. For more information, check [PREVENT_UNLOAD_TO_INTERNAL_STAGES docs](https://docs.snowflake.com/en/sql-reference/parameters#prevent-unload-to-internal-stages).
- `query_tag` (String) Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions. For more information, check [QUERY_TAG docs](https://docs.snowflake.com/en/sql-reference/parameters#query-tag).
- `quoted_identifiers_ignore_case` (Boolean) Specifies whether letters in double-quoted object identifiers are stored and resolved as uppercase letters. By default, Snowflake preserves the case of alphabetic characters when storing and resolving double-quoted identifiers (see [Identifier resolution](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing)). You can use this parameter in situations in which [third-party applications always use double quotes around identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing-parameter). For more information, check [QUOTED_IDENTIFIERS_IGNORE_CASE docs](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case).
- `rows_per_resultset` (Number) Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum. For more information, check [ROWS_PER_RESULTSET docs](https://docs.snowflake.com/en/sql-reference/parameters#rows-per-resultset).
- `rsa_public_key` (String) Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.
- `rsa_public_key_2` (String) Specifies the user’s second RSA public key; used to rotate the public and private keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.
- `s3_stage_vpce_dns_name` (String) Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation). For more information, check [S3_STAGE_VPCE_DNS_NAME docs](https://docs.snowflake.com/en/sql-reference/parameters#s3-stage-vpce-dns-name).
- `search_path` (String) Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name. For more information, check [SEARCH_PATH docs](https://docs.snowflake.com/en/sql-reference/parameters#search-path).
- `simulated_data_sharing_consumer` (String) Specifies the name of a consumer account to simulate for testing/validating shared data, particularly shared secure views. When this parameter is set in a session, shared views return rows as if executed in the specified consumer account rather than the provider account. For more information, see [Introduction to Secure Data Sharing](https://docs.snowflake.com/en/user-guide/data-sharing-intro) and [Working with shares](https://docs.snowflake.com/en/user-guide/data-sharing-provider). For more information, check [SIMULATED_DATA_SHARING_CONSUMER docs](https://docs.snowflake.com/en/sql-reference/parameters#simulated-data-sharing-consumer).
- `statement_queued_timeout_in_seconds` (Number) Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged. For more information, check [STATEMENT_QUEUED_TIMEOUT_IN_SECONDS docs](https://docs.snowflake.com/en/sql-reference/parameters#statement-queued-timeout-in-seconds).
- `statement_timeout_in_seconds` (Number) Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system. For more information, check [STATEMENT_TIMEOUT_IN_SECONDS docs](https://docs.snowflake.com/en/sql-reference/parameters#statement-timeout-in-seconds).
- `strict_json_output` (Boolean) This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON. For more information, check [STRICT_JSON_OUTPUT docs](https://docs.snowflake.com/en/sql-reference/parameters#strict-json-output).
- `time_input_format` (String) Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session). For more information, check [TIME_INPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#time-input-format).
- `time_output_format` (String) Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). For more information, check [TIME_OUTPUT_FORMAT docs](https://docs.snowflake.com/en/sql-reference/parameters#time-output-format).
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
- `week_of_year_policy` (Number) Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year. For more information, check [WEEK_OF_YEAR_POLICY docs](https://docs.snowflake.com/en/sql-reference/parameters#week-of-year-policy).
- `week_start` (Number) Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week. For more information, check [WEEK_START docs](https://docs.snowflake.com/en/sql-reference/parameters#week-start).

### Read-Only

- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `parameters` (List of Object) Outputs the result of `SHOW PARAMETERS IN USER` for the given user. (see [below for nested schema](#nestedatt--parameters))
- `show_output` (List of Object) Outputs the result of `SHOW USER` for the given user. (see [below for nested schema](#nestedatt--show_output))
- `user_type` (String) Specifies a type for the user.

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
- `enable_unredacted_query_syntax_error` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--enable_unredacted_query_syntax_error))
- `error_on_nondeterministic_merge` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--error_on_nondeterministic_merge))
- `error_on_nondeterministic_update` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--error_on_nondeterministic_update))
- `geography_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--geography_output_format))
- `geometry_output_format` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--geometry_output_format))
- `jdbc_treat_decimal_as_int` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--jdbc_treat_decimal_as_int))
- `jdbc_treat_timestamp_ntz_as_utc` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--jdbc_treat_timestamp_ntz_as_utc))
- `jdbc_use_session_timezone` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--jdbc_use_session_timezone))
- `json_indent` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--json_indent))
- `lock_timeout` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--lock_timeout))
- `log_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--log_level))
- `multi_statement_count` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--multi_statement_count))
- `network_policy` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--network_policy))
- `noorder_sequence_as_default` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--noorder_sequence_as_default))
- `odbc_treat_decimal_as_int` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--odbc_treat_decimal_as_int))
- `prevent_unload_to_internal_stages` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--prevent_unload_to_internal_stages))
- `query_tag` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--query_tag))
- `quoted_identifiers_ignore_case` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--quoted_identifiers_ignore_case))
- `rows_per_resultset` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--rows_per_resultset))
- `s3_stage_vpce_dns_name` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--s3_stage_vpce_dns_name))
- `search_path` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--search_path))
- `simulated_data_sharing_consumer` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--simulated_data_sharing_consumer))
- `statement_queued_timeout_in_seconds` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--statement_queued_timeout_in_seconds))
- `statement_timeout_in_seconds` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--statement_timeout_in_seconds))
- `strict_json_output` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--strict_json_output))
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


<a id="nestedobjatt--parameters--enable_unredacted_query_syntax_error"></a>
### Nested Schema for `parameters.enable_unredacted_query_syntax_error`

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


<a id="nestedobjatt--parameters--jdbc_treat_decimal_as_int"></a>
### Nested Schema for `parameters.jdbc_treat_decimal_as_int`

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


<a id="nestedobjatt--parameters--network_policy"></a>
### Nested Schema for `parameters.network_policy`

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


<a id="nestedobjatt--parameters--prevent_unload_to_internal_stages"></a>
### Nested Schema for `parameters.prevent_unload_to_internal_stages`

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


<a id="nestedobjatt--parameters--simulated_data_sharing_consumer"></a>
### Nested Schema for `parameters.simulated_data_sharing_consumer`

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

- `comment` (String)
- `created_on` (String)
- `days_to_expiry` (String)
- `default_namespace` (String)
- `default_role` (String)
- `default_secondary_roles` (String)
- `default_warehouse` (String)
- `disabled` (Boolean)
- `display_name` (String)
- `email` (String)
- `expires_at_time` (String)
- `ext_authn_duo` (Boolean)
- `ext_authn_uid` (String)
- `first_name` (String)
- `has_mfa` (Boolean)
- `has_password` (Boolean)
- `has_rsa_public_key` (Boolean)
- `last_name` (String)
- `last_success_login` (String)
- `locked_until_time` (String)
- `login_name` (String)
- `mins_to_bypass_mfa` (String)
- `mins_to_unlock` (String)
- `must_change_password` (Boolean)
- `name` (String)
- `owner` (String)
- `snowflake_lock` (Boolean)
- `type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_service_user.example '"<user_name>"'
```
