resource "snowflake_view" "test" {
  name        = var.name
  database    = var.database
  schema      = var.schema
  statement   = var.statement
  copy_grants = var.copy_grants
  is_secure   = var.is_secure

  dynamic "column" {
    for_each = var.column
    content {
      column_name = column.value["column_name"]
    }
  }
}
