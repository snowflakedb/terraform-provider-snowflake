# minimal
resource "snowflake_storage_integration_azure" "minimal" {
  name                      = "example_azure_storage_integration"
  enabled                   = true
  storage_allowed_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/"]
  azure_tenant_id           = "a123b4c5-1234-123a-a12b-1a23b45678c9"
}

# all fields
resource "snowflake_storage_integration_azure" "all" {
  name    = "example_azure_storage_integration"
  enabled = true
  storage_allowed_locations = [
    "azure://myaccount.blob.core.windows.net/mycontainer/allowed-location/",
    "azure://myaccount.blob.core.windows.net/mycontainer/allowed-location2/"
  ]
  storage_blocked_locations = [
    "azure://myaccount.blob.core.windows.net/mycontainer/blocked-location/",
    "azure://myaccount.blob.core.windows.net/mycontainer/blocked-location2/"
  ]
  use_privatelink_endpoint = "true"
  comment                  = "some comment"

  azure_tenant_id = "a123b4c5-1234-123a-a12b-1a23b45678c9"
}
