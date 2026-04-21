# Basic + OAuth
resource "snowflake_catalog_integration_iceberg_rest" "oauth_basic" {
  name    = "example_iceberg_rest_oauth"
  enabled = false

  rest_config {
    catalog_uri = "https://your-iceberg-rest.example.com"
  }

  oauth_rest_authentication {
    oauth_client_id      = "your_oauth_client_id"
    oauth_client_secret  = "your_oauth_client_secret"
    oauth_allowed_scopes = ["PRINCIPAL_ROLE:ALL"]
  }
}

# Complete + OAuth
resource "snowflake_catalog_integration_iceberg_rest" "oauth_complete" {
  name                     = "example_iceberg_rest_oauth_complete"
  enabled                  = true
  refresh_interval_seconds = 60
  comment                  = "Lorem ipsum"
  catalog_namespace        = "my_namespace"

  rest_config {
    catalog_uri            = "https://your-iceberg-rest.example.com"
    prefix                 = "/api/v1"
    catalog_name           = "my_catalog"
    catalog_api_type       = "PUBLIC"
    access_delegation_mode = "EXTERNAL_VOLUME_CREDENTIALS"
  }

  oauth_rest_authentication {
    oauth_token_uri      = "https://idp.example.com/oauth/token"
    oauth_client_id      = "your_oauth_client_id"
    oauth_client_secret  = "your_oauth_client_secret"
    oauth_allowed_scopes = ["PRINCIPAL_ROLE:ALL"]
  }
}

# Bearer token
resource "snowflake_catalog_integration_iceberg_rest" "bearer" {
  name    = "example_iceberg_rest_bearer"
  enabled = false

  rest_config {
    catalog_uri = "https://your-iceberg-rest.example.com"
  }

  bearer_rest_authentication {
    bearer_token = "your_static_bearer_token"
  }
}

# SigV4
resource "snowflake_catalog_integration_iceberg_rest" "sigv4" {
  name    = "example_iceberg_rest_sigv4"
  enabled = false

  rest_config {
    catalog_uri = "https://your-iceberg-rest.example.com"
  }

  sigv4_rest_authentication {
    sigv4_iam_role       = "arn:aws:iam::123456789012:role/YourRole"
    sigv4_signing_region = "us-east-1"
    sigv4_external_id    = "optional-external-id"
  }
}
