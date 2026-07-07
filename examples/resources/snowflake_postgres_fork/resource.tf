# basic fork - required fields only
resource "snowflake_postgres_fork" "basic" {
  name      = "my_postgres_fork"
  fork_from = "my_source_postgres_instance"
}

# fork with point-in-time recovery using AT timestamp
resource "snowflake_postgres_fork" "pitr_at" {
  name      = "my_postgres_fork_pitr_at"
  fork_from = "my_source_postgres_instance"
  at {
    timestamp = "2025-01-15 12:00:00"
  }
}

# fork with point-in-time recovery using BEFORE timestamp
resource "snowflake_postgres_fork" "pitr_before" {
  name      = "my_postgres_fork_pitr_before"
  fork_from = "my_source_postgres_instance"
  before {
    timestamp = "2025-01-15 12:00:00"
  }
}

# fork with point-in-time recovery using BEFORE offset
resource "snowflake_postgres_fork" "pitr_before_offset" {
  name      = "my_postgres_fork_pitr_before_offset"
  fork_from = "my_source_postgres_instance"
  before {
    offset = "-3600"
  }
}

# complete fork with all options
resource "snowflake_postgres_fork" "complete" {
  name      = "my_postgres_fork_complete"
  fork_from = "my_source_postgres_instance"
  at {
    offset = "-3600"
  }
  compute_family    = "STANDARD_M"
  storage_size_gb   = 20
  high_availability = "true"
  postgres_settings = "{\"max_connections\": 100}"
  comment           = "Forked for testing"
}
