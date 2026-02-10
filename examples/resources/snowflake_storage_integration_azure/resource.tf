# minimal
resource "snowflake_storage_integration_azure" "minimal" {
  name                      = "example_azure_storage_integration"
  enabled                   = true
  storage_allowed_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/"]
  azure_tenant_id           = "a123b4c5-1234-123a-a12b-1a23b45678c9"
}

# TODO [next PR]: add all fields example
# all fields
