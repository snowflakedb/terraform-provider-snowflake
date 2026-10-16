# Basic
resource "snowflake_catalog_linked_database" "basic" {
  name = "example_catalog_linked_database"

  linked_catalog {
    catalog = snowflake_catalog_integration_iceberg_rest.example.name
  }

  external_volume = snowflake_external_volume.example.name
}

# Complete
resource "snowflake_catalog_linked_database" "complete" {
  name = "example_catalog_linked_database_complete"

  linked_catalog {
    catalog                     = snowflake_catalog_integration_iceberg_rest.example.name
    allowed_namespaces          = ["ns1", "ns2"]
    blocked_namespaces          = ["ns3"]
    allowed_write_operations    = "ALL"
    namespace_mode              = "FLATTEN_NESTED_NAMESPACE"
    namespace_flatten_delimiter = "_"
    sync_interval_seconds       = 60
  }

  external_volume          = snowflake_external_volume.example.name
  catalog_case_sensitivity = "CASE_INSENSITIVE"
  comment                  = "synced from an external Iceberg REST catalog"
}
