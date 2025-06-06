resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = "some_share"
  on_database = "some_database"
}
