# basic resource
resource "snowflake_catalog_integration_open_catalog" "basic" {
  name    = "example_open_catalog"
  enabled = false

  rest_config {
    catalog_uri  = "https://your-org.snowflakecomputing.com/polaris/api/catalog"
    catalog_name = "my_catalog"
  }

  rest_authentication {
    oauth_client_id      = "your_oauth_client_id"
    oauth_client_secret  = "your_oauth_client_secret"
    oauth_allowed_scopes = ["PRINCIPAL_ROLE:ALL"]
  }
}

# complete resource
resource "snowflake_catalog_integration_open_catalog" "complete" {
  name                     = "example_open_catalog_complete"
  enabled                  = true
  refresh_interval_seconds = 60
  comment                  = "Lorem ipsum"
  catalog_namespace        = "my_namespace"

  rest_config {
    catalog_uri            = "https://your-org.snowflakecomputing.com/polaris/api/catalog"
    catalog_name           = "my_catalog"
    catalog_api_type       = "PUBLIC"
    access_delegation_mode = "EXTERNAL_VOLUME_CREDENTIALS"
  }

  rest_authentication {
    oauth_token_uri      = "https://your-org.snowflakecomputing.com/polaris/api/catalog/v1/oauth/tokens"
    oauth_client_id      = "your_oauth_client_id"
    oauth_client_secret  = "your_oauth_client_secret"
    oauth_allowed_scopes = ["PRINCIPAL_ROLE:ALL"]
  }
}
