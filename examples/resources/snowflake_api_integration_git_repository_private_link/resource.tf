# basic resource - no authentication secrets (none allowed)
resource "snowflake_api_integration_git_repository_private_link" "basic" {
  name                               = "git_repository_private_link_integration"
  use_privatelink_endpoint           = true
  no_allowed_authentication_secrets  = true
  api_allowed_prefixes               = ["https://github.example.internal/my-org/"]
  enabled                            = true
}

# with specific allowed secrets
resource "snowflake_api_integration_git_repository_private_link" "with_secrets" {
  name                            = "git_repository_private_link_integration_secrets"
  use_privatelink_endpoint        = true
  allowed_authentication_secrets  = ["my_db.my_schema.my_secret"]
  api_allowed_prefixes            = ["https://github.example.internal/my-org/"]
  enabled                         = true
}

# complete resource - all secrets allowed, with TLS certificates
resource "snowflake_api_integration_git_repository_private_link" "complete" {
  name                              = "git_repository_private_link_integration_complete"
  use_privatelink_endpoint          = true
  all_allowed_authentication_secrets = true
  tls_trusted_certificates          = ["my_db.my_schema.my_cert_secret"]
  api_allowed_prefixes              = ["https://github.example.internal/my-org/"]
  api_blocked_prefixes              = ["https://github.example.internal/my-org/private-repo/"]
  enabled                           = true
  comment                           = "Example Git Repository Private Link integration"
}
