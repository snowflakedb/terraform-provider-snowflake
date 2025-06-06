data "snowflake_tags" "test" {
  like = "non-existing-tag"

  lifecycle {
    postcondition {
      condition     = length(self.tags) > 0
      error_message = "there should be at least one tag"
    }
  }
}
