resource "snowflake_stage" "test" {
  name        = var.name
  url         = var.location
  database    = var.database
  schema      = var.schema
  credentials = "aws_key_id = '${var.aws_key_id}' aws_secret_key = '${var.aws_secret_key}'"
  file_format = "TYPE = JSON NULL_IF = []"
}

resource "snowflake_external_table" "test_table" {
  name     = var.name
  database = var.database
  schema   = var.schema
  comment  = "Terraform acceptance test"
  column {
    name = "name"
    type = "string"
    as   = "value:name::string"
  }
  column {
    name = "age"
    type = "number"
    as   = "value:age::number"
  }
  auto_refresh      = false
  refresh_on_create = true
  file_format       = "TYPE = JSON, STRIP_OUTER_ARRAY = TRUE"
  location          = "@\"${var.database}\".\"${var.schema}\".\"${snowflake_stage.test.name}\"/external_tables_test_data/"
}
