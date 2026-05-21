# basic fork
resource "snowflake_postgres_fork" "basic" {
  name      = "my_postgres_fork"
  fork_from = "my_source_postgres_instance"
}

# fork with point-in-time recovery
resource "snowflake_postgres_fork" "pitr" {
  name      = "my_postgres_fork_pitr"
  fork_from = "my_source_postgres_instance"
  at {
    timestamp = "2025-01-15 12:00:00"
  }
}

# fork with all options
resource "snowflake_postgres_fork" "complete" {
  name      = "my_postgres_fork_complete"
  fork_from = "my_source_postgres_instance"
  at {
    offset = "-3600"
  }
  compute_family    = "STANDARD_1"
  storage_size_gb   = 20
  high_availability = true
  comment           = "Forked for testing"
}
