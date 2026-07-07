# basic resource - no authentication secrets (none allowed)
resource "snowflake_api_integration_git_repository_token" "basic" {
  name                              = "git_repository_token_integration"
  no_allowed_authentication_secrets = true
  api_allowed_prefixes              = ["https://github.com/my-org/"]
  enabled                           = true
}

# with specific allowed secrets
resource "snowflake_api_integration_git_repository_token" "with_secrets" {
  name                           = "git_repository_token_integration_secrets"
  allowed_authentication_secrets = ["my_db.my_schema.my_secret"]
  api_allowed_prefixes           = ["https://github.com/my-org/"]
  enabled                        = true
}

# complete resource - all secrets allowed
resource "snowflake_api_integration_git_repository_token" "complete" {
  name                               = "git_repository_token_integration_complete"
  all_allowed_authentication_secrets = true
  api_allowed_prefixes               = ["https://github.com/my-org/"]
  api_blocked_prefixes               = ["https://github.com/my-org/private-repo/"]
  enabled                            = true
  comment                            = "Example Git Repository Token integration"
}
