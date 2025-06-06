resource "snowflake_api_integration" "test_api_int" {
  name                 = var.name
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = var.api_allowed_prefixes
  enabled              = true
}

resource "snowflake_external_function" "external_function" {
  name     = var.name
  database = var.database
  schema   = var.schema
  arg {
    name = "ARG1"
    type = "VARCHAR"
  }
  arg {
    name = "ARG2"
    type = "VARCHAR"
  }
  comment                   = var.comment
  return_type               = "VARIANT"
  return_behavior           = "IMMUTABLE"
  api_integration           = snowflake_api_integration.test_api_int.name
  url_of_proxy_and_resource = var.url_of_proxy_and_resource
}
