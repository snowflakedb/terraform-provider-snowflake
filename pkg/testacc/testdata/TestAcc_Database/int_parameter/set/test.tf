resource "snowflake_database" "test" {
  name                        = var.name
  data_retention_time_in_days = var.data_retention_time_in_days
}
