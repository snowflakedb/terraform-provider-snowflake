# minimal
resource "snowflake_storage_integration_gcs" "minimal" {
  name                      = "example_gcs_storage_integration"
  enabled                   = true
  storage_allowed_locations = ["gcs://mybucket1/path1"]
}

# all fields
resource "snowflake_storage_integration_gcs" "all" {
  name    = "example_gcs_storage_integration"
  enabled = true
  storage_allowed_locations = [
    "gcs://mybucket1/allowed-location/", "gcs://mybucket1/allowed-location2/"
  ]
  storage_blocked_locations = [
    "gcs://mybucket1/blocked-location/", "gcs://mybucket1/blocked-location2/"
  ]
  comment = "some comment"
}
