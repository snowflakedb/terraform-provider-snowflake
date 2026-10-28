## Minimal
# from is ForceNew. It is a snow://workspace/... URL, or an internal or temporary stage reference.
resource "snowflake_code_bundle" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "code_bundle_name"
  from     = "@${snowflake_stage.example.fully_qualified_name}"
}

## Complete
resource "snowflake_code_bundle" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "code_bundle_name"
  from     = "@${snowflake_stage.example.fully_qualified_name}"
  comment  = "My code bundle"
}
