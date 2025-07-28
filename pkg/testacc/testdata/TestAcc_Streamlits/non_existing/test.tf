data "snowflake_streamlits" "test" {
  like = "non-existing-streamlit"

  lifecycle {
    postcondition {
      condition     = length(self.streamlits) > 0
      error_message = "there should be at least one streamlit"
    }
  }
}
