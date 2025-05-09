data "snowflake_image_repositories" "test" {
  like = "non-existing-image-repository"

  lifecycle {
    postcondition {
      condition     = length(self.image_repositories) > 0
      error_message = "there should be at least one image repository"
    }
  }
}
