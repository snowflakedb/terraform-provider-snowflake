resource "snowflake_api_integration" "test_api_int" {
  name                 = var.name
  api_provider         = "aws_api_gateway"
  api_aws_role_arn     = "arn:aws:iam::000000000001:/role/test"
  api_allowed_prefixes = var.api_allowed_prefixes
  enabled              = true
}

resource "snowflake_function_javascript" "test_func_req_translator" {
  name     = "${var.name}_request_translator"
  database = var.database
  schema   = var.schema
  arguments {
    arg_name      = "EVENT"
    arg_data_type = "OBJECT"
  }
  return_type         = "OBJECT"
  function_definition = <<EOT
		  	let exeprimentName = EVENT.body.data[0][1]
		  	return { "body": { "name": test }}
	  	EOT
}


resource "snowflake_function_javascript" "test_func_res_translator" {
  name     = "${var.name}_response_translator"
  database = var.database
  schema   = var.schema
  arguments {
    arg_name      = "EVENT"
    arg_data_type = "OBJECT"
  }
  return_type         = "OBJECT"
  function_definition = <<EOT
			  return { "body": { "data" :  [[0, EVENT]] } };
		  EOT
}

resource "snowflake_external_function" "external_function" {
  name            = var.name
  database        = var.database
  schema          = var.schema
  comment         = var.comment
  return_type     = "VARIANT"
  return_behavior = "IMMUTABLE"
  api_integration = snowflake_api_integration.test_api_int.name
  header {
    name  = "x-custom-header"
    value = "snowflake"
  }
  max_batch_rows            = 500
  request_translator        = "${var.database}.${var.schema}.${snowflake_function_javascript.test_func_req_translator.name}"
  response_translator       = "${var.database}.${var.schema}.${snowflake_function_javascript.test_func_res_translator.name}"
  url_of_proxy_and_resource = var.url_of_proxy_and_resource
}
