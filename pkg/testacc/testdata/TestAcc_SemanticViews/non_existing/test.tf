data "snowflake_semantic_views" "test" {
  like = "non-existing-semantic-view"

  lifecycle {
    postcondition {
      condition     = length(self.semantic_views) > 0
      error_message = "there should be at least one semantic view"
    }
  }
}
