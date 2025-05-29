data "snowflake_security_integrations" "test" {
  like = "non-existing-security-integration"

  lifecycle {
    postcondition {
      condition     = length(self.security_integrations) > 0
      error_message = "there should be at least one security integration"
    }
  }
}
