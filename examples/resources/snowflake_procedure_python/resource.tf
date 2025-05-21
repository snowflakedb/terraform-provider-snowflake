resource "snowflake_procedure_python" "w" {
  database = "Database"
  schema   = "Schema"
  name     = "Name"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "echoVarchar"
  procedure_definition = <<EOT
  def echoVarchar(x):
  result = ''
  for a in range(5):
    result += x
  return result
EOT
  runtime_version      = "3.9"
  snowpark_package     = "1.14.0"
}
