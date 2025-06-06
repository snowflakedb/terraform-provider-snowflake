data "snowflake_git_repositories" "test" {
  like = "non-existing-git-repository"

  lifecycle {
    postcondition {
      condition     = length(self.git_repositories) > 0
      error_message = "there should be at least one git repository"
    }
  }
}