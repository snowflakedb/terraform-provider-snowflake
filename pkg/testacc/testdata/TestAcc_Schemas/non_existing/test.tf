data "snowflake_schemas" "test" {
  like = "non-existing-schema"

  lifecycle {
    postcondition {
      condition     = length(self.schemas) > 0
      error_message = "there should be at least one schema"
    }
  }
}
