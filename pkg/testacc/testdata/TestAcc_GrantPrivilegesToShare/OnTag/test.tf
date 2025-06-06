resource "snowflake_tag" "test" {
  name     = var.on_tag
  database = var.database
  schema   = var.schema
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.to_share
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = var.to_share
  privileges = var.privileges
  on_tag     = snowflake_tag.test.fully_qualified_name
  depends_on = [snowflake_grant_privileges_to_share.test_setup]
}