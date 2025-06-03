# basic resource
resource "snowflake_git_repository" "basic" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "GIT_REPOSITORY"
  origin          = "https://github/com/user/repo"
  api_integration = "API_INTEGRATION"
}

# complete resource
resource "snowflake_git_repository" "complete" {
  name            = "GIT_REPOSITORY"
  database        = "DATABASE"
  schema          = "SCHEMA"
  origin          = "https://github/com/user/repo"
  api_integration = "API_INTEGRATION"
  git_credentials = "GIT_CREDENTIALS"
  comment         = "comment"
}