data "snowflake_trust_center_scanners" "all" {}

data "snowflake_trust_center_scanners" "security_essentials" {
  scanner_package_id = "SECURITY_ESSENTIALS"
}
