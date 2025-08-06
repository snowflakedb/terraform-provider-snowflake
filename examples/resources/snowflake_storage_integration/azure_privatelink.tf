resource "snowflake_storage_integration" "azure_integration" {
  name    = "azure_integration"
  comment = "Azure integration with private link"
  type    = "EXTERNAL_STAGE"

  enabled = true

  storage_provider          = "AZURE"
  azure_tenant_id           = "a123b4c5-1234-123a-a12b-1a23b45678c9"
  storage_allowed_locations = ["azure://myaccount.blob.core.windows.net/mycontainer/path1/"]
  use_private_link_endpoint = true
}
