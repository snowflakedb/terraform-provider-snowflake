resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.to_share
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share                = var.to_share
  privileges              = var.privileges
  on_all_tables_in_schema = "\"${var.database}\".\"${var.schema}\""
  depends_on              = [snowflake_grant_privileges_to_share.test_setup]
}
