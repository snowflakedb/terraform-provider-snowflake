# basic resource
resource "snowflake_postgres_instance" "basic" {
  name                     = "my_postgres_instance"
  compute_family           = "STANDARD_M"
  storage_size_gb          = 10
  authentication_authority = "POSTGRES"
}

# with custom timeouts (postgres operations can take longer than the defaults)
resource "snowflake_postgres_instance" "with_timeouts" {
  name                     = "my_postgres_instance"
  compute_family           = "STANDARD_M"
  storage_size_gb          = 10
  authentication_authority = "POSTGRES"

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
    read   = "10m"
  }
}

# complete resource
resource "snowflake_postgres_instance" "complete" {
  name                     = "my_postgres_instance_complete"
  compute_family           = "STANDARD_M"
  storage_size_gb          = 10
  authentication_authority = "POSTGRES_OR_SNOWFLAKE"
  postgres_version         = 16
  network_policy           = "my_network_policy"
  storage_integration      = "my_storage_integration"
  high_availability        = "true"
  postgres_settings        = "{\"max_connections\": 100}"
  maintenance_window_start = 8
  comment                  = "My Postgres instance"
}
