# Minimal
resource "snowflake_function_python" "minimal" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "my_function_function"
  runtime_version = "3.9"
  arguments {
    arg_data_type = "NUMBER(36, 2)"
    arg_name      = "x"
  }
  return_type         = "NUMBER(36, 2)"
  handler             = "some_function"
  function_definition = <<EOT
def some_function(x):
  result = ''
  for a in range(5):
    result += x
  return result
EOT
}
