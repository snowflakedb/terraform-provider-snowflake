resource "snowflake_storage_integration" "test" {
  storage_allowed_locations = []
  storage_provider          = "S3"
  name                      = "empty_storage_allowed_locations_test"
}