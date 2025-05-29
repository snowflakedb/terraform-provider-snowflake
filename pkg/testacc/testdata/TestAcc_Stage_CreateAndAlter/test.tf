resource "snowflake_stage" "test" {
  name                = var.name
  schema              = var.schema
  database            = var.database
  comment             = var.comment
  url                 = var.url
  storage_integration = var.storage_integration
  credentials         = var.credentials
  encryption          = var.encryption
  file_format         = var.file_format
  copy_options        = var.copy_options
}
