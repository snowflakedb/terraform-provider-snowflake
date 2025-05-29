resource "snowflake_stream_on_table" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  copy_grants       = var.copy_grants
  table             = var.table
  append_only       = var.append_only
  show_initial_rows = var.show_initial_rows

  before {
    timestamp = try(var.before["timestamp"], null)
    offset    = try(var.before["offset"], null)
    stream    = try(var.before["stream"], null)
    statement = try(var.before["statement"], null)
  }

  comment = var.comment
}
