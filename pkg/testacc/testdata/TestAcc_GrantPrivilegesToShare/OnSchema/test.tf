resource "snowflake_schema" "test" {
  name     = var.schema
  database = var.database
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.to_share
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = var.to_share
  privileges = var.privileges
  on_schema  = snowflake_schema.test.fully_qualified_name
  depends_on = [snowflake_grant_privileges_to_share.test_setup]
}
