terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
      version = "= 2.12.0"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "snowflake" {
  # Uses default configuration from ~/.snowflake/config or environment variables
}

# ------------------------------------------------------------------------------
# UUID for unique naming
# ------------------------------------------------------------------------------
resource "random_uuid" "test_run" {}

locals {
  # Use first 8 chars of UUID for shorter names
  test_id = upper(substr(random_uuid.test_run.result, 0, 8))
  prefix  = "MIGRATION_TEST_${local.test_id}"
}

# Output the prefix for use in datasource
output "test_prefix" {
  description = "The unique prefix used for all test objects"
  value       = local.prefix
}

# ------------------------------------------------------------------------------
# PERSON Users (type = 'PERSON' or empty)
# ------------------------------------------------------------------------------

# 1. Basic PERSON user - minimal attributes
resource "snowflake_user" "person_basic" {
  name = "${local.prefix}_PERSON_BASIC"
}

# 2. PERSON user with all basic attributes
resource "snowflake_user" "person_complete" {
  name                            = "${local.prefix}_PERSON_COMPLETE"
  login_name                      = "${local.test_id}_person_login"
  display_name                    = "Migration Test Person"
  first_name                      = "John"
  middle_name                     = "Michael"
  last_name                       = "Doe"
  email                           = "john.doe@example.com"
  comment                         = "Complete PERSON user for migration testing"
  default_warehouse               = "COMPUTE_WH"
  default_namespace               = "TEST_DB.TEST_SCHEMA"
  default_role                    = "PUBLIC"
  default_secondary_roles_option  = "ALL"
  disabled                        = false
  must_change_password            = false
}

# 3. PERSON user with multi-line RSA public key (tests newline escaping)
# NOTE: To test RSA key escaping, generate a key and uncomment rsa_public_key below.
#       See README.md for instructions on generating RSA keys.
resource "snowflake_user" "person_rsa" {
  name         = "${local.prefix}_PERSON_RSA"
  login_name   = "${local.test_id}_rsa_login"
  display_name = "User with RSA Key"
  email        = "rsa.user@example.com"
  comment      = "User with multi-line RSA key for testing newline escaping"
  # Uncomment and paste your generated public key here:
  # rsa_public_key = <<-EOT
  # <PASTE_RSA_PUBLIC_KEY_1_HERE>
  # EOT
}

# 4. PERSON user with special characters in comment (tests quote escaping)
resource "snowflake_user" "person_special_chars" {
  name         = "${local.prefix}_PERSON_SPECIAL_CHARS"
  login_name   = "${local.test_id}_special_login"
  display_name = "User with Special Characters"
  comment      = "Comment with special chars: <>&"
}

# 5. PERSON user with all (possible) session parameters set
resource "snowflake_user" "person_params" {
  name         = "${local.prefix}_PERSON_PARAMS"
  login_name   = "${local.test_id}_params_login"
  display_name = "User with Parameters"

  # Session parameters - ALL parameters supported by the migration script
  abort_detached_query                          = true
  autocommit                                    = false
  binary_input_format                           = "HEX"
  binary_output_format                          = "BASE64"
  client_memory_limit                           = 2048
  client_metadata_request_use_connection_ctx    = true
  client_prefetch_threads                       = 8
  client_result_chunk_size                      = 128
  client_result_column_case_insensitive         = true
  client_session_keep_alive                     = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping                 = "TIMESTAMP_NTZ"
  date_input_format                             = "YYYY-MM-DD"
  date_output_format                            = "YYYY-MM-DD"
  enable_unload_physical_type_optimization      = true
  enable_unredacted_query_syntax_error          = true
  error_on_nondeterministic_merge               = false
  error_on_nondeterministic_update              = false
  geography_output_format                       = "GeoJSON"
  geometry_output_format                        = "WKT"
  jdbc_treat_decimal_as_int                     = true
  jdbc_treat_timestamp_ntz_as_utc               = true
  jdbc_use_session_timezone                     = true
  json_indent                                   = 4
  lock_timeout                                  = 43200
  log_level                                     = "INFO"
  multi_statement_count                         = 0
  noorder_sequence_as_default                   = true
  odbc_treat_decimal_as_int                     = true
  prevent_unload_to_internal_stages             = true
  query_tag                                     = "migration_test"
  quoted_identifiers_ignore_case                = true
  rows_per_resultset                            = 100000
  search_path                                   = "$current, $public"
  statement_queued_timeout_in_seconds           = 300
  statement_timeout_in_seconds                  = 86400
  strict_json_output                            = true
  time_input_format                             = "HH24:MI:SS"
  time_output_format                            = "HH24:MI:SS"
  timestamp_day_is_always_24h                   = true
  timestamp_input_format                        = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ltz_output_format                   = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_tz_output_format                    = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_type_mapping                        = "TIMESTAMP_NTZ"
  timezone                                      = "America/New_York"
  trace_level                                   = "ON_EVENT"
  transaction_abort_on_error                    = true
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1970
  unsupported_ddl_action                        = "IGNORE"
  use_cached_result                             = false
  week_of_year_policy                           = 1
  week_start                                    = 1

  # Note: The following parameters require special setup and are not tested here:
  # - network_policy: requires a network policy to exist
  # - active_python_profiler: requires specific profiler setup
  # - python_profiler_modules: requires specific profiler setup
  # - python_profiler_target_stage: requires a stage to exist
  # - s3_stage_vpce_dns_name: requires AWS VPC endpoint setup
  # - simulated_data_sharing_consumer: requires data sharing setup

  # There are also a few unsupported ones:
  # - client_enable_log_info_statement_parameters
  # - client_metadata_use_session_database
  # - csv_timestamp_format
  # - hybrid_table_lock_timeout
  # - js_treat_integer_as_big_int
}

# 6. PERSON user with disabled = true
resource "snowflake_user" "person_disabled" {
  name       = "${local.prefix}_PERSON_DISABLED"
  login_name = "${local.test_id}_disabled_login"
  disabled   = true
  comment    = "Disabled PERSON user"
}

# ------------------------------------------------------------------------------
# SERVICE Users (type = 'SERVICE')
# Note: SERVICE users CANNOT have: first_name, middle_name, last_name,
#       password, must_change_password, mins_to_bypass_mfa
# ------------------------------------------------------------------------------

# 7. Basic SERVICE user - minimal attributes
resource "snowflake_service_user" "service_basic" {
  name = "${local.prefix}_SERVICE_BASIC"
}

# 8. SERVICE user with all allowed attributes
resource "snowflake_service_user" "service_complete" {
  name                            = "${local.prefix}_SERVICE_COMPLETE"
  login_name                      = "${local.test_id}_service_login"
  display_name                    = "Migration Test Service User"
  email                           = "service@example.com"
  comment                         = "Complete SERVICE user for migration testing"
  default_warehouse               = "COMPUTE_WH"
  default_namespace               = "TEST_DB.TEST_SCHEMA"
  default_role                    = "PUBLIC"
  default_secondary_roles_option  = "ALL"
  disabled                        = false
}

# 9. SERVICE user with RSA keys (common for service accounts)
# NOTE: To test RSA key escaping with key rotation, generate two keys and uncomment below.
#       See README.md for instructions on generating RSA keys.
resource "snowflake_service_user" "service_rsa" {
  name         = "${local.prefix}_SERVICE_RSA"
  login_name   = "${local.test_id}_service_rsa_login"
  display_name = "Service User with RSA"
  comment      = "SERVICE user with both RSA keys for key rotation testing"
  # Uncomment and paste your generated public keys here:
  # rsa_public_key = <<-EOT
  # <PASTE_RSA_PUBLIC_KEY_1_HERE>
  # EOT
  # rsa_public_key_2 = <<-EOT
  # <PASTE_RSA_PUBLIC_KEY_2_HERE>
  # EOT
}

# 10. SERVICE user with all (possible) session parameters set
resource "snowflake_service_user" "service_params" {
  name       = "${local.prefix}_SERVICE_PARAMS"
  login_name = "${local.test_id}_service_params_login"

  # Session parameters - ALL parameters supported for SERVICE users
  abort_detached_query                          = true
  autocommit                                    = false
  binary_input_format                           = "HEX"
  binary_output_format                          = "BASE64"
  client_memory_limit                           = 2048
  client_metadata_request_use_connection_ctx    = true
  client_prefetch_threads                       = 8
  client_result_chunk_size                      = 128
  client_result_column_case_insensitive         = true
  client_session_keep_alive                     = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping                 = "TIMESTAMP_NTZ"
  date_input_format                             = "YYYY-MM-DD"
  date_output_format                            = "YYYY-MM-DD"
  enable_unload_physical_type_optimization      = true
  enable_unredacted_query_syntax_error          = true
  error_on_nondeterministic_merge               = false
  error_on_nondeterministic_update              = false
  geography_output_format                       = "GeoJSON"
  geometry_output_format                        = "WKT"
  jdbc_treat_decimal_as_int                     = true
  jdbc_treat_timestamp_ntz_as_utc               = true
  jdbc_use_session_timezone                     = true
  json_indent                                   = 4
  lock_timeout                                  = 43200
  log_level                                     = "INFO"
  multi_statement_count                         = 0
  noorder_sequence_as_default                   = true
  odbc_treat_decimal_as_int                     = true
  prevent_unload_to_internal_stages             = true
  query_tag                                     = "service_migration_test"
  quoted_identifiers_ignore_case                = true
  rows_per_resultset                            = 100000
  search_path                                   = "$current, $public"
  statement_queued_timeout_in_seconds           = 300
  statement_timeout_in_seconds                  = 86400
  strict_json_output                            = true
  time_input_format                             = "HH24:MI:SS"
  time_output_format                            = "HH24:MI:SS"
  timestamp_day_is_always_24h                   = true
  timestamp_input_format                        = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ltz_output_format                   = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_ntz_output_format                   = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format                       = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_tz_output_format                    = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_type_mapping                        = "TIMESTAMP_NTZ"
  timezone                                      = "America/New_York"
  trace_level                                   = "ON_EVENT"
  transaction_abort_on_error                    = true
  transaction_default_isolation_level           = "READ COMMITTED"
  two_digit_century_start                       = 1970
  unsupported_ddl_action                        = "IGNORE"
  use_cached_result                             = false
  week_of_year_policy                           = 1
  week_start                                    = 1
}

# ------------------------------------------------------------------------------
# LEGACY_SERVICE Users (type = 'LEGACY_SERVICE')
# Note: LEGACY_SERVICE users CANNOT have: first_name, middle_name, last_name,
#       mins_to_bypass_mfa
# BUT CAN have: password, must_change_password
# ------------------------------------------------------------------------------

# 11. Basic LEGACY_SERVICE user - minimal attributes
resource "snowflake_legacy_service_user" "legacy_basic" {
  name = "${local.prefix}_LEGACY_BASIC"
}

# 12. LEGACY_SERVICE user with all allowed attributes
resource "snowflake_legacy_service_user" "legacy_complete" {
  name                            = "${local.prefix}_LEGACY_COMPLETE"
  login_name                      = "${local.test_id}_legacy_login"
  display_name                    = "Migration Test Legacy Service"
  email                           = "legacy@example.com"
  comment                         = "Complete LEGACY_SERVICE user for migration testing"
  default_warehouse               = "COMPUTE_WH"
  default_namespace               = "TEST_DB.TEST_SCHEMA"
  default_role                    = "PUBLIC"
  default_secondary_roles_option  = "ALL"
  disabled                        = false
  must_change_password            = false
}

# 13. LEGACY_SERVICE user with RSA keys
# NOTE: To test RSA key escaping, generate a key and uncomment rsa_public_key below.
#       See README.md for instructions on generating RSA keys.
resource "snowflake_legacy_service_user" "legacy_rsa" {
  name       = "${local.prefix}_LEGACY_RSA"
  login_name = "${local.test_id}_legacy_rsa_login"
  comment    = "LEGACY_SERVICE user with RSA key"
  # Uncomment and paste your generated public key here:
  # rsa_public_key = <<-EOT
  # <PASTE_RSA_PUBLIC_KEY_3_HERE>
  # EOT
}

# ------------------------------------------------------------------------------
# Edge Cases
# ------------------------------------------------------------------------------

# 14. User with very long comment (tests field handling)
resource "snowflake_user" "long_comment" {
  name       = "${local.prefix}_LONG_COMMENT"
  login_name = "${local.test_id}_long_comment_login"
  comment    = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# 15. User with unicode characters in names
resource "snowflake_user" "unicode" {
  name         = "${local.prefix}_UNICODE"
  login_name   = "${local.test_id}_unicode_login"
  display_name = "Tëst Üsér wïth Ünïcödé"
  first_name   = "José"
  last_name    = "García"
  comment      = "Unicode comment test: émojis 🎉, accents àéîõü, symbols €£¥, Chinese 中文, Japanese 日本語"
}
