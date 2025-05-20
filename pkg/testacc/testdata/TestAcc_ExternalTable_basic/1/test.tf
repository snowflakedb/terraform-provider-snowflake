resource "snowflake_stage" "test" {
  name        = var.name
  url         = var.location
  database    = var.database
  schema      = var.schema
  credentials = "aws_key_id = '${var.aws_key_id}' aws_secret_key = '${var.aws_secret_key}'"
  file_format = "TYPE = JSON NULL_IF = []"
}
