# basic resource
resource "snowflake_postgres_instance" "basic" {
  name                     = "my_postgres_instance"
  compute_family           = "STANDARD_1"
  storage_size_gb          = 10
  authentication_authority = "POSTGRES"
}

# complete resource
resource "snowflake_postgres_instance" "complete" {
  name                     = "my_postgres_instance"
  compute_family           = "STANDARD_1"
  storage_size_gb          = 10
  authentication_authority = "POSTGRES"
  postgres_version         = 16
  high_availability        = true
  comment                  = "My Postgres instance"
}
