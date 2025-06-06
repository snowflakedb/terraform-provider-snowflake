resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.argument
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body    = var.body
  comment = var.comment
}

data "snowflake_row_access_policies" "test" {
  depends_on = [snowflake_row_access_policy.test]

  like = var.name
  limit {
    rows = 10
    from = snowflake_row_access_policy.test.name
  }
}
