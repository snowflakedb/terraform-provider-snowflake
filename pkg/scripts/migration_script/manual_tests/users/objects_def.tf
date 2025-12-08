# =============================================================================
# User Objects for Migration Script Testing
# =============================================================================
# This Terraform file creates users covering all edge cases for testing the
# migration script.
#
# Usage:
#   1. Configure Snowflake provider credentials (env vars or config)
#   2. Run: terraform init && terraform apply
#   3. Use main.tf to generate CSV: cd .. && terraform apply
#   4. Test migration: go run . -import=block users < users.csv
#   5. Cleanup: terraform destroy
# =============================================================================

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {
  # Configure via environment variables:
  # SNOWFLAKE_ACCOUNT, SNOWFLAKE_USER, SNOWFLAKE_PASSWORD, etc.
}

# ------------------------------------------------------------------------------
# PERSON Users (type = 'PERSON' or empty)
# ------------------------------------------------------------------------------

# 1. Basic PERSON user - minimal attributes
resource "snowflake_user" "person_basic" {
  name = "MIGRATION_TEST_PERSON_BASIC"
}

# 2. PERSON user with all basic attributes
resource "snowflake_user" "person_complete" {
  name                         = "MIGRATION_TEST_PERSON_COMPLETE"
  login_name                   = "migration_person_login"
  display_name                 = "Migration Test Person"
  first_name                   = "John"
  middle_name                  = "Michael"
  last_name                    = "Doe"
  email                        = "john.doe@example.com"
  comment                      = "Complete PERSON user for migration testing"
  default_warehouse            = "COMPUTE_WH"
  default_namespace            = "TEST_DB.TEST_SCHEMA"
  default_role                 = "PUBLIC"
  default_secondary_roles_option = "ALL"
  disabled                     = false
  must_change_password         = false
}

# 3. PERSON user with multi-line RSA public key (tests newline escaping)
resource "snowflake_user" "person_rsa" {
  name         = "MIGRATION_TEST_PERSON_RSA"
  login_name   = "migration_rsa_login"
  display_name = "User with RSA Key"
  email        = "rsa.user@example.com"
  comment      = "User with multi-line RSA key for testing newline escaping"
  rsa_public_key = <<-EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt194dL4wn/9BgLPO4LJi
hXB/yFSf6Pja7Ysm1MK4mjdBEFNBDWiD9YIp085zhG01uP8bqLdwy5l3Qc0dgocc
Rt0Ul5nBWfRDqkcT3NPqmAC8v7mZne+f1pcGlOFjnBhnPtE+4bomRwjkF1jWrrzn
2HCwwTELxPHHG10dbJgPtPMql2GUHqm/V2pn2JBlcyR7Kg4oh34DniRhAouTTJDu
5fMj7oBqLXPuya5ck+AY8yzH3bBAtwc879BVduySbNN10QLDfamPHz/G50P4KXWX
dm+i716/OyUEdJN2YzcO8CFl8TBZ7g4uPdjspOuaF7bQu9i7Aia41cHFQZxVBMZu
2wIDAQAB
EOT
}

# 4. PERSON user with special characters in comment (tests quote escaping)
resource "snowflake_user" "person_special_chars" {
  name         = "MIGRATION_TEST_PERSON_SPECIAL_CHARS"
  login_name   = "migration_special_login"
  display_name = "User with Special Characters"
  comment      = "Comment with special chars: <>&"
}

# 5. PERSON user with session parameters set
resource "snowflake_user" "person_params" {
  name         = "MIGRATION_TEST_PERSON_PARAMS"
  login_name   = "migration_params_login"
  display_name = "User with Parameters"

  # Session parameters
  abort_detached_query                       = true
  autocommit                                 = false
  binary_input_format                        = "HEX"
  binary_output_format                       = "BASE64"
  client_memory_limit                        = 2048
  client_metadata_request_use_connection_ctx = true
  client_prefetch_threads                    = 8
  client_result_chunk_size                   = 128
  client_result_column_case_insensitive      = true
  client_session_keep_alive                  = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping              = "TIMESTAMP_NTZ"
  date_input_format                          = "YYYY-MM-DD"
  date_output_format                         = "YYYY-MM-DD"
  geography_output_format                    = "GeoJSON"
  geometry_output_format                     = "WKT"
  json_indent                                = 4
  lock_timeout                               = 43200
  log_level                                  = "INFO"
  multi_statement_count                      = 0
  noorder_sequence_as_default                = true
  quoted_identifiers_ignore_case             = true
  rows_per_resultset                         = 100000
  search_path                                = "$current, $public"
  statement_queued_timeout_in_seconds        = 300
  statement_timeout_in_seconds               = 86400
  strict_json_output                         = true
  time_input_format                          = "HH24:MI:SS"
  time_output_format                         = "HH24:MI:SS"
  timestamp_day_is_always_24h                = true
  timestamp_input_format                     = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ltz_output_format                = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_ntz_output_format                = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format                    = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_tz_output_format                 = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_type_mapping                     = "TIMESTAMP_NTZ"
  timezone                                   = "America/New_York"
  trace_level                                = "ON_EVENT"
  transaction_abort_on_error                 = true
  transaction_default_isolation_level        = "READ COMMITTED"
  two_digit_century_start                    = 1970
  unsupported_ddl_action                     = "IGNORE"
  use_cached_result                          = false
  week_of_year_policy                        = 1
  week_start                                 = 1
}

# 6. PERSON user with disabled = true
resource "snowflake_user" "person_disabled" {
  name       = "MIGRATION_TEST_PERSON_DISABLED"
  login_name = "migration_disabled_login"
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
  name = "MIGRATION_TEST_SERVICE_BASIC"
}

# 8. SERVICE user with all allowed attributes
resource "snowflake_service_user" "service_complete" {
  name                         = "MIGRATION_TEST_SERVICE_COMPLETE"
  login_name                   = "migration_service_login"
  display_name                 = "Migration Test Service User"
  email                        = "service@example.com"
  comment                      = "Complete SERVICE user for migration testing"
  default_warehouse            = "COMPUTE_WH"
  default_namespace            = "TEST_DB.TEST_SCHEMA"
  default_role                 = "PUBLIC"
  default_secondary_roles_option = "ALL"
  disabled                     = false
}

# 9. SERVICE user with RSA keys (common for service accounts)
resource "snowflake_service_user" "service_rsa" {
  name         = "MIGRATION_TEST_SERVICE_RSA"
  login_name   = "migration_service_rsa_login"
  display_name = "Service User with RSA"
  comment      = "SERVICE user with both RSA keys for key rotation testing"
  rsa_public_key = <<-EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt194dL4wn/9BgLPO4LJi
hXB/yFSf6Pja7Ysm1MK4mjdBEFNBDWiD9YIp085zhG01uP8bqLdwy5l3Qc0dgocc
Rt0Ul5nBWfRDqkcT3NPqmAC8v7mZne+f1pcGlOFjnBhnPtE+4bomRwjkF1jWrrzn
2HCwwTELxPHHG10dbJgPtPMql2GUHqm/V2pn2JBlcyR7Kg4oh34DniRhAouTTJDu
5fMj7oBqLXPuya5ck+AY8yzH3bBAtwc879BVduySbNN10QLDfamPHz/G50P4KXWX
dm+i716/OyUEdJN2YzcO8CFl8TBZ7g4uPdjspOuaF7bQu9i7Aia41cHFQZxVBMZu
2wIDAQAB
EOT
  rsa_public_key_2 = <<-EOT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz3XLCZnE3iZ969ZvEj/Z
r2xyQrSKxQF8GzRa66FFpOb12INrUviCr5S4ijIU4biu/oZ9kXtmXqcS7BA/k2Kc
8t9w87RlCXBCeGe+sFi8IwFDkdJ+CWu2xTxrkSvp7Po8AqpGws+A7we7+q59DOZ7
uu+/ZIVBIPXfGsCsp/LBD+xllt0M5c/fJt3y+DGmsx+mkdzW1RCj0Nrrp1dX52If
NYtM20ftKBMMiPBCaDJxq3QuYK6JWPmyYrs5bklvSJr5bJiygTjFKgiSUj6h3vnH
pcgO2n75b5INcd9goT/1b6xd/Al7kzGrgalkzXuww9qlJr7R1+qdRH6dEXgPwKM6
pwIDAQAB
EOT
}

# 10. SERVICE user with parameters
resource "snowflake_service_user" "service_params" {
  name       = "MIGRATION_TEST_SERVICE_PARAMS"
  login_name = "migration_service_params_login"

  statement_timeout_in_seconds = 3600
  trace_level                  = "ON_EVENT"
  log_level                    = "INFO"
}

# ------------------------------------------------------------------------------
# LEGACY_SERVICE Users (type = 'LEGACY_SERVICE')
# Note: LEGACY_SERVICE users CANNOT have: first_name, middle_name, last_name,
#       mins_to_bypass_mfa
# BUT CAN have: password, must_change_password
# ------------------------------------------------------------------------------

# 11. Basic LEGACY_SERVICE user - minimal attributes
resource "snowflake_legacy_service_user" "legacy_basic" {
  name = "MIGRATION_TEST_LEGACY_BASIC"
}

# 12. LEGACY_SERVICE user with all allowed attributes
resource "snowflake_legacy_service_user" "legacy_complete" {
  name                         = "MIGRATION_TEST_LEGACY_COMPLETE"
  login_name                   = "migration_legacy_login"
  display_name                 = "Migration Test Legacy Service"
  email                        = "legacy@example.com"
  comment                      = "Complete LEGACY_SERVICE user for migration testing"
  default_warehouse            = "COMPUTE_WH"
  default_namespace            = "TEST_DB.TEST_SCHEMA"
  default_role                 = "PUBLIC"
  default_secondary_roles_option = "ALL"
  disabled                     = false
  must_change_password         = false
}

# 13. LEGACY_SERVICE user with RSA keys
resource "snowflake_legacy_service_user" "legacy_rsa" {
  name       = "MIGRATION_TEST_LEGACY_RSA"
  login_name = "migration_legacy_rsa_login"
  comment    = "LEGACY_SERVICE user with RSA key"
  rsa_public_key = <<-EOT
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
# Edge Cases
# ------------------------------------------------------------------------------

# 14. User with very long comment (tests field handling)
resource "snowflake_user" "long_comment" {
  name       = "MIGRATION_TEST_LONG_COMMENT"
  login_name = "migration_long_comment_login"
  comment    = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# 15. User with unicode characters in names
resource "snowflake_user" "unicode" {
  name         = "MIGRATION_TEST_UNICODE"
  login_name   = "migration_unicode_login"
  display_name = "Test User with Unicode"
  first_name   = "Jose"
  last_name    = "Garcia"
  comment      = "Unicode comment test"
}

