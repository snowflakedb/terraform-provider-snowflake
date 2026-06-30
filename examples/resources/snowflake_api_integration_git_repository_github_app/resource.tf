# basic resource
resource "snowflake_api_integration_git_repository_github_app" "basic" {
  name                 = "git_repository_github_app_integration"
  api_allowed_prefixes = ["https://github.com/my-org/"]
  enabled              = true
}

# complete resource
resource "snowflake_api_integration_git_repository_github_app" "complete" {
  name                 = "git_repository_github_app_integration_complete"
  api_allowed_prefixes = ["https://github.com/my-org/"]
  api_blocked_prefixes = ["https://github.com/my-org/private-repo/"]
  enabled              = true
  comment              = "Example Git Repository GitHub App integration"
}
