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
  body = var.body
}
