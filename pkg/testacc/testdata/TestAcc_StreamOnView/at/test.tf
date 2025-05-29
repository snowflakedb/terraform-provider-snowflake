resource "snowflake_stream_on_view" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  copy_grants       = var.copy_grants
  view              = var.view
  append_only       = var.append_only
  show_initial_rows = var.show_initial_rows

  at {
    timestamp = try(var.at["timestamp"], null)
    offset    = try(var.at["offset"], null)
    stream    = try(var.at["stream"], null)
    statement = try(var.at["statement"], null)
  }

  comment = var.comment
}
