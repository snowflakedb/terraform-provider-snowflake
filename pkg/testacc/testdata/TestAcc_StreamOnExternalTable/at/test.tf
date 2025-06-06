resource "snowflake_stream_on_external_table" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema

  copy_grants    = var.copy_grants
  external_table = var.external_table
  insert_only    = var.insert_only

  at {
    timestamp = try(var.at["timestamp"], null)
    offset    = try(var.at["offset"], null)
    stream    = try(var.at["stream"], null)
    statement = try(var.at["statement"], null)
  }

  comment = var.comment
}
