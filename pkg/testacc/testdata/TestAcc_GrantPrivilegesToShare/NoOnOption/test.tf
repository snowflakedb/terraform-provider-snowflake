resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = "some_share"
  privileges = ["USAGE"]
}
