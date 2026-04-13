data "snowflake_trust_center_scanner_packages" "all" {}

data "snowflake_trust_center_scanner_packages" "security" {
  like = "SECURITY%"
}
