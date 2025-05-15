data "snowflake_users" "test" {
  like = "non-existing-user"

  lifecycle {
    postcondition {
      condition     = length(self.users) > 0
      error_message = "there should be at least one user"
    }
  }
}
