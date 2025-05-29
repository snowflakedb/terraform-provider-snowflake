resource "snowflake_table" "test" {
  name     = var.on_table
  database = var.database
  schema   = var.schema
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.to_share
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = var.to_share
  privileges = var.privileges
  on_table   = snowflake_table.test.fully_qualified_name
  depends_on = [snowflake_grant_privileges_to_share.test_setup]
}
