# basic resource (without oauth_scopes — scopes are inherited from the security integration)
resource "snowflake_secret_with_client_credentials" "basic" {
  name               = "EXAMPLE_SECRET"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  api_authentication = snowflake_api_authentication_integration_with_client_credentials.example.fully_qualified_name
}

# resource with explicit oauth_scopes
resource "snowflake_secret_with_client_credentials" "with_scopes" {
  name               = "EXAMPLE_SECRET"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  api_authentication = snowflake_api_authentication_integration_with_client_credentials.example.fully_qualified_name
  oauth_scopes       = ["useraccount", "testscope"]
}

# resource with all fields set
resource "snowflake_secret_with_client_credentials" "complete" {
  name               = "EXAMPLE_SECRET"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  api_authentication = snowflake_api_authentication_integration_with_client_credentials.example.fully_qualified_name
  oauth_scopes       = ["useraccount", "testscope"]
  comment            = "EXAMPLE_COMMENT"
}
