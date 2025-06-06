data "snowflake_masking_policies" "test" {
  like = "non-existing-masking-policy"

  lifecycle {
    postcondition {
      condition     = length(self.masking_policies) > 0
      error_message = "there should be at least one masking policy"
    }
  }
}
