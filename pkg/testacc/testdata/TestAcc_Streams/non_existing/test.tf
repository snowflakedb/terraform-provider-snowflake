data "snowflake_streams" "test" {
  like = "non-existing-stream"

  lifecycle {
    postcondition {
      condition     = length(self.streams) > 0
      error_message = "there should be at least one stream"
    }
  }
}
