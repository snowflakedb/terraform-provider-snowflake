data "snowflake_row_access_policies" "test" {
  like = "non-existing-row-access-policy"

  lifecycle {
    postcondition {
      condition     = length(self.row_access_policies) > 0
      error_message = "there should be at least one row access policy"
    }
  }
}
