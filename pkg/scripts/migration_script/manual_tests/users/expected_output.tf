# =============================================================================
# Expected Migration Script Output
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block users < users.csv
#
# Use this to compare with actual output:
#   go run .. -import=block users < users.csv > actual_output.tf
#   diff users_expected_output.tf actual_output.tf
#
# NOTE:
# - Snowflake uppercases login_name values
# - Special characters are Unicode-escaped in HCL output
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# =============================================================================

# ------------------------------------------------------------------------------
# LEGACY_SERVICE Users (sorted first alphabetically)
# ------------------------------------------------------------------------------

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_BASIC" {
  name = "MIGRATION_TEST_LEGACY_BASIC"
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_LEGACY_BASIC"
  login_name = "MIGRATION_TEST_LEGACY_BASIC"
}

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_COMPLETE" {
  name = "MIGRATION_TEST_LEGACY_COMPLETE"
  comment = "Complete LEGACY_SERVICE user for migration testing"
  default_namespace = "TEST_DB.TEST_SCHEMA"
  default_role = "PUBLIC"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  display_name = "Migration Test Legacy Service"
  email = "legacy@example.com"
  login_name = "MIGRATION_LEGACY_LOGIN"
}

resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_RSA" {
  name = "MIGRATION_TEST_LEGACY_RSA"
  comment = "LEGACY_SERVICE user with RSA key"
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_LEGACY_RSA"
  login_name = "MIGRATION_LEGACY_RSA_LOGIN"
  rsa_public_key = <<EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz3XLCZnE3iZ969ZvEj/Z
r2xyQrSKxQF8GzRa66FFpOb12INrUviCr5S4ijIU4biu/oZ9kXtmXqcS7BA/k2Kc
8t9w87RlCXBCeGe+sFi8IwFDkdJ+CWu2xTxrkSvp7Po8AqpGws+A7we7+q59DOZ7
uu+/ZIVBIPXfGsCsp/LBD+xllt0M5c/fJt3y+DGmsx+mkdzW1RCj0Nrrp1dX52If
NYtM20ftKBMMiPBCaDJxq3QuYK6JWPmyYrs5bklvSJr5bJiygTjFKgiSUj6h3vnH
pcgO2n75b5INcd9goT/1b6xd/Al7kzGrgalkzXuww9qlJr7R1+qdRH6dEXgPwKM6
pwIDAQAB
EOT
}

# ------------------------------------------------------------------------------
# PERSON Users (snowflake_user)
# ------------------------------------------------------------------------------

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_LONG_COMMENT" {
  name = "MIGRATION_TEST_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_LONG_COMMENT"
  login_name = "MIGRATION_LONG_COMMENT_LOGIN"
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_BASIC" {
  name = "MIGRATION_TEST_PERSON_BASIC"
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_PERSON_BASIC"
  login_name = "MIGRATION_TEST_PERSON_BASIC"
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_COMPLETE" {
  name = "MIGRATION_TEST_PERSON_COMPLETE"
  comment = "Complete PERSON user for migration testing"
  default_namespace = "TEST_DB.TEST_SCHEMA"
  default_role = "PUBLIC"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  display_name = "Migration Test Person"
  email = "john.doe@example.com"
  first_name = "John"
  last_name = "Doe"
  login_name = "MIGRATION_PERSON_LOGIN"
  middle_name = "Michael"
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_DISABLED" {
  name = "MIGRATION_TEST_PERSON_DISABLED"
  comment = "Disabled PERSON user"
  default_secondary_roles_option = "ALL"
  disabled = "true"
  display_name = "MIGRATION_TEST_PERSON_DISABLED"
  login_name = "MIGRATION_DISABLED_LOGIN"
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_PARAMS" {
  name = "MIGRATION_TEST_PERSON_PARAMS"
  abort_detached_query = true
  autocommit = false
  binary_input_format = "HEX"
  binary_output_format = "BASE64"
  client_memory_limit = 2048
  client_metadata_request_use_connection_ctx = true
  client_prefetch_threads = 8
  client_result_chunk_size = 128
  client_result_column_case_insensitive = true
  client_session_keep_alive = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping = "TIMESTAMP_NTZ"
  date_input_format = "YYYY-MM-DD"
  date_output_format = "YYYY-MM-DD"
  default_secondary_roles_option = "ALL"
  display_name = "User with Parameters"
  geography_output_format = "GeoJSON"
  geometry_output_format = "WKT"
  json_indent = 4
  lock_timeout = 43200
  log_level = "INFO"
  login_name = "MIGRATION_PARAMS_LOGIN"
  multi_statement_count = 0
  noorder_sequence_as_default = true
  quoted_identifiers_ignore_case = true
  rows_per_resultset = 100000
  search_path = "$current, $public"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 86400
  strict_json_output = true
  time_input_format = "HH24:MI:SS"
  time_output_format = "HH24:MI:SS"
  timestamp_input_format = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_type_mapping = "TIMESTAMP_NTZ"
  timezone = "America/New_York"
  trace_level = "ON_EVENT"
  transaction_abort_on_error = true
  transaction_default_isolation_level = "READ COMMITTED"
  two_digit_century_start = 1970
  use_cached_result = false
  week_of_year_policy = 1
  week_start = 1
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_RSA" {
  name = "MIGRATION_TEST_PERSON_RSA"
  comment = "User with multi-line RSA key for testing newline escaping"
  default_secondary_roles_option = "ALL"
  display_name = "User with RSA Key"
  email = "rsa.user@example.com"
  login_name = "MIGRATION_RSA_LOGIN"
  rsa_public_key = <<EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt194dL4wn/9BgLPO4LJi
hXB/yFSf6Pja7Ysm1MK4mjdBEFNBDWiD9YIp085zhG01uP8bqLdwy5l3Qc0dgocc
Rt0Ul5nBWfRDqkcT3NPqmAC8v7mZne+f1pcGlOFjnBhnPtE+4bomRwjkF1jWrrzn
2HCwwTELxPHHG10dbJgPtPMql2GUHqm/V2pn2JBlcyR7Kg4oh34DniRhAouTTJDu
5fMj7oBqLXPuya5ck+AY8yzH3bBAtwc879BVduySbNN10QLDfamPHz/G50P4KXWX
dm+i716/OyUEdJN2YzcO8CFl8TBZ7g4uPdjspOuaF7bQu9i7Aia41cHFQZxVBMZu
2wIDAQAB
EOT
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_PERSON_SPECIAL_CHARS" {
  name = "MIGRATION_TEST_PERSON_SPECIAL_CHARS"
  comment = "Comment with special chars: \u003c\u003e\u0026"
  default_secondary_roles_option = "ALL"
  display_name = "User with Special Characters"
  login_name = "MIGRATION_SPECIAL_LOGIN"
}

# ------------------------------------------------------------------------------
# SERVICE Users
# ------------------------------------------------------------------------------

resource "snowflake_service_user" "snowflake_generated_service_user_MIGRATION_TEST_SERVICE_BASIC" {
  name = "MIGRATION_TEST_SERVICE_BASIC"
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_SERVICE_BASIC"
  login_name = "MIGRATION_TEST_SERVICE_BASIC"
}

resource "snowflake_service_user" "snowflake_generated_service_user_MIGRATION_TEST_SERVICE_COMPLETE" {
  name = "MIGRATION_TEST_SERVICE_COMPLETE"
  comment = "Complete SERVICE user for migration testing"
  default_namespace = "TEST_DB.TEST_SCHEMA"
  default_role = "PUBLIC"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  display_name = "Migration Test Service User"
  email = "service@example.com"
  login_name = "MIGRATION_SERVICE_LOGIN"
}

resource "snowflake_service_user" "snowflake_generated_service_user_MIGRATION_TEST_SERVICE_PARAMS" {
  name = "MIGRATION_TEST_SERVICE_PARAMS"
  default_secondary_roles_option = "ALL"
  display_name = "MIGRATION_TEST_SERVICE_PARAMS"
  log_level = "INFO"
  login_name = "MIGRATION_SERVICE_PARAMS_LOGIN"
  statement_timeout_in_seconds = 3600
  trace_level = "ON_EVENT"
}

resource "snowflake_service_user" "snowflake_generated_service_user_MIGRATION_TEST_SERVICE_RSA" {
  name = "MIGRATION_TEST_SERVICE_RSA"
  comment = "SERVICE user with both RSA keys for key rotation testing"
  default_secondary_roles_option = "ALL"
  display_name = "Service User with RSA"
  login_name = "MIGRATION_SERVICE_RSA_LOGIN"
  rsa_public_key = <<EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt194dL4wn/9BgLPO4LJi
hXB/yFSf6Pja7Ysm1MK4mjdBEFNBDWiD9YIp085zhG01uP8bqLdwy5l3Qc0dgocc
Rt0Ul5nBWfRDqkcT3NPqmAC8v7mZne+f1pcGlOFjnBhnPtE+4bomRwjkF1jWrrzn
2HCwwTELxPHHG10dbJgPtPMql2GUHqm/V2pn2JBlcyR7Kg4oh34DniRhAouTTJDu
5fMj7oBqLXPuya5ck+AY8yzH3bBAtwc879BVduySbNN10QLDfamPHz/G50P4KXWX
dm+i716/OyUEdJN2YzcO8CFl8TBZ7g4uPdjspOuaF7bQu9i7Aia41cHFQZxVBMZu
2wIDAQAB
EOT
  rsa_public_key_2 = <<EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz3XLCZnE3iZ969ZvEj/Z
r2xyQrSKxQF8GzRa66FFpOb12INrUviCr5S4ijIU4biu/oZ9kXtmXqcS7BA/k2Kc
8t9w87RlCXBCeGe+sFi8IwFDkdJ+CWu2xTxrkSvp7Po8AqpGws+A7we7+q59DOZ7
uu+/ZIVBIPXfGsCsp/LBD+xllt0M5c/fJt3y+DGmsx+mkdzW1RCj0Nrrp1dX52If
NYtM20ftKBMMiPBCaDJxq3QuYK6JWPmyYrs5bklvSJr5bJiygTjFKgiSUj6h3vnH
pcgO2n75b5INcd9goT/1b6xd/Al7kzGrgalkzXuww9qlJr7R1+qdRH6dEXgPwKM6
pwIDAQAB
EOT
}

resource "snowflake_user" "snowflake_generated_user_MIGRATION_TEST_UNICODE" {
  name = "MIGRATION_TEST_UNICODE"
  comment = "Unicode comment test"
  default_secondary_roles_option = "ALL"
  display_name = "Test User with Unicode"
  first_name = "Jose"
  last_name = "Garcia"
  login_name = "MIGRATION_UNICODE_LOGIN"
}

# ------------------------------------------------------------------------------
# Import Blocks (grouped at the end)
# ------------------------------------------------------------------------------

import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_BASIC
  id = "\"MIGRATION_TEST_LEGACY_BASIC\""
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_COMPLETE
  id = "\"MIGRATION_TEST_LEGACY_COMPLETE\""
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_MIGRATION_TEST_LEGACY_RSA
  id = "\"MIGRATION_TEST_LEGACY_RSA\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_LONG_COMMENT
  id = "\"MIGRATION_TEST_LONG_COMMENT\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_BASIC
  id = "\"MIGRATION_TEST_PERSON_BASIC\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_COMPLETE
  id = "\"MIGRATION_TEST_PERSON_COMPLETE\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_DISABLED
  id = "\"MIGRATION_TEST_PERSON_DISABLED\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_PARAMS
  id = "\"MIGRATION_TEST_PERSON_PARAMS\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_RSA
  id = "\"MIGRATION_TEST_PERSON_RSA\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_PERSON_SPECIAL_CHARS
  id = "\"MIGRATION_TEST_PERSON_SPECIAL_CHARS\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_MIGRATION_TEST_SERVICE_BASIC
  id = "\"MIGRATION_TEST_SERVICE_BASIC\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_MIGRATION_TEST_SERVICE_COMPLETE
  id = "\"MIGRATION_TEST_SERVICE_COMPLETE\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_MIGRATION_TEST_SERVICE_PARAMS
  id = "\"MIGRATION_TEST_SERVICE_PARAMS\""
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_MIGRATION_TEST_SERVICE_RSA
  id = "\"MIGRATION_TEST_SERVICE_RSA\""
}
import {
  to = snowflake_user.snowflake_generated_user_MIGRATION_TEST_UNICODE
  id = "\"MIGRATION_TEST_UNICODE\""
}
