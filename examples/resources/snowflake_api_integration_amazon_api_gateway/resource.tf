# basic resource
resource "snowflake_api_integration_amazon_api_gateway" "basic" {
  name                 = "amazon_api_gateway_integration"
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:role/test"
  api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
  enabled              = true
}

# complete resource
resource "snowflake_api_integration_amazon_api_gateway" "complete" {
  name                 = "amazon_api_gateway_integration_complete"
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:role/test"
  api_allowed_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/"]
  api_blocked_prefixes = ["https://123456.execute-api.us-west-2.amazonaws.com/prod/blocked/"]
  api_key              = "my-api-key"
  enabled              = true
  comment              = "Example Amazon API Gateway integration"
}
