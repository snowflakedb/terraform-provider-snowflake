# minimal
resource "snowflake_storage_integration_gcs" "minimal" {
  name                      = "example_gcs_storage_integration"
  enabled                   = true
  storage_allowed_locations = ["gcs://mybucket1/path1"]
}

# TODO [next PR]: add all fields example
# all fields
