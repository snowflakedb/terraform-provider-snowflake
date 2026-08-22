resource "snowflake_trust_center_scanner" "mfa_check" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  scanner_id         = "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"
  enabled            = true
  schedule           = "USING CRON 0 0 * * * UTC"

  notification {
    notify_admins      = true
    severity_threshold = "Critical"
  }
}
