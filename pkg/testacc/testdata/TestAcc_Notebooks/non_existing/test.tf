data "snowflake_notebooks" "test" {
  like = "non-existing-notebook"

  lifecycle {
    postcondition {
      condition     = length(self.notebooks) > 0
      error_message = "there should be at least one notebook"
    }
  }
}
