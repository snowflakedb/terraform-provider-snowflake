# basic resource
resource "snowflake_catalog_integration_aws_glue" "basic" {
  name              = "example"
  enabled           = false
  glue_aws_role_arn = "arn:aws:iam::123456789012:role/testRole"
  glue_catalog_id   = "123456789012"
}

# complete resource
resource "snowflake_catalog_integration_aws_glue" "complete" {
  name                     = "example_complete"
  enabled                  = true
  refresh_interval_seconds = 60
  comment                  = "Lorem ipsum"
  glue_aws_role_arn        = "arn:aws:iam::123456789012:role/testRole"
  glue_catalog_id          = "123456789012"
  glue_region              = "us-east-1"
  catalog_namespace        = "myNamespace"
}
