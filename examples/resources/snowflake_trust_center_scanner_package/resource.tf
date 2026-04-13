resource "snowflake_trust_center_scanner_package" "security_essentials" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  enabled            = true
  schedule           = "USING CRON 0 2 * * * UTC"

  notification {
    notify_admins      = true
    severity_threshold = "High"
    users              = ["SECURITY_ADMIN"]
  }
}
