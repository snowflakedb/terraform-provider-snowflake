## Minimal
resource "snowflake_projection_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "projection_policy_name"
  body     = "PROJECTION_CONSTRAINT(ALLOW => false)"
}

## Complete
resource "snowflake_projection_policy" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "projection_policy_name"
  body     = "PROJECTION_CONSTRAINT(ALLOW => true)"
  comment  = "My projection policy"
}
