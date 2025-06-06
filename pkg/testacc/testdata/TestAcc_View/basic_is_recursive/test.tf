resource "snowflake_view" "test" {
  name         = var.name
  database     = var.database
  schema       = var.schema
  statement    = var.statement
  is_recursive = var.is_recursive

  dynamic "column" {
    for_each = var.column
    content {
      column_name = column.value["column_name"]
    }
  }
}
