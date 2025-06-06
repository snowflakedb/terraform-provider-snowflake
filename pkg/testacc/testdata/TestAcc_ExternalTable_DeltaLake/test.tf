resource "snowflake_stage" "test" {
  name        = var.name
  url         = var.location
  database    = var.database
  schema      = var.schema
  credentials = "aws_key_id = '${var.aws_key_id}' aws_secret_key = '${var.aws_secret_key}'"
  file_format = "TYPE = PARQUET NULL_IF = []"
}

resource "snowflake_external_table" "test_table" {
  name         = var.name
  database     = var.database
  schema       = var.schema
  comment      = "Terraform acceptance test"
  table_format = "delta"
  column {
    name = "filename"
    type = "string"
    as   = "metadata$filename"
  }
  column {
    name = "name"
    type = "string"
    as   = "value:name::string"
  }
  partition_by      = ["filename"]
  auto_refresh      = false
  refresh_on_create = false
  file_format       = "TYPE = PARQUET"
  location          = "@\"${var.database}\".\"${var.schema}\".\"${snowflake_stage.test.name}\"/external_tables_test_data/"
}
