data "snowflake_grants" "test" {
  grants_of {
    database_role = "role"
  }
}
