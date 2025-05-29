resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.name
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = var.name
  privileges = var.privileges

  on_function = "\"${var.database}\".\"${var.schema}\".\"${var.function_name}\"(${var.argument_type})"
  depends_on  = [snowflake_grant_privileges_to_share.test_setup]
}
