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
    name = "filename"
    type = "string"
    as   = "metadata$filename"
  }
  column {
    name = "name"
    type = "varchar(200)"
    as   = "value:name::string"
  }
  column {
    name = "age"
    type = "number(2, 2)"
    as   = "value:age::number"
  }
  partition_by      = ["filename"]
  auto_refresh      = false
  refresh_on_create = true
  file_format       = "TYPE = JSON, STRIP_OUTER_ARRAY = TRUE"
  location          = "@\"${var.database}\".\"${var.schema}\".\"${snowflake_stage.test.name}\"/external_tables_test_data/"
}
