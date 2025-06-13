data "snowflake_services" "test" {
  like = "non-existing-service"

  lifecycle {
    postcondition {
      condition     = length(self.services) > 0
      error_message = "there should be at least one service"
    }
  }
}
