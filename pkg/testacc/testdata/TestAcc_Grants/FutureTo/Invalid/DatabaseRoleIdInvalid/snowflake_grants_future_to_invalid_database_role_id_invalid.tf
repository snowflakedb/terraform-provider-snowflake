data "snowflake_grants" "test" {
  future_grants_to {
    database_role = "role"
  }
}
