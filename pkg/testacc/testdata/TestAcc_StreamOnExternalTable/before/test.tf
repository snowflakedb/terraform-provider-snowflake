resource "snowflake_stream_on_external_table" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  copy_grants    = var.copy_grants
  external_table = var.external_table
  insert_only    = var.insert_only

  before {
    timestamp = try(var.before["timestamp"], null)
    offset    = try(var.before["offset"], null)
    stream    = try(var.before["stream"], null)
    statement = try(var.before["statement"], null)
  }

  comment = var.comment
}
