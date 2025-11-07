data "snowflake_authentication_policies" "test" {
  like = "non-existing-authentication-policy"

  lifecycle {
    postcondition {
      condition     = length(self.authentication_policies) > 0
      error_message = "there should be at least one authentication policy"
    }
  }
}
