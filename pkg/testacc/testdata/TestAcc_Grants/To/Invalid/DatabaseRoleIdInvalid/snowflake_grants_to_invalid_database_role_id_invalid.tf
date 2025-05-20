data "snowflake_grants" "test" {
  grants_to {
    database_role = "role"
  }
}
