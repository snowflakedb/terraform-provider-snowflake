## Minimal
resource "snowflake_aggregation_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "aggregation_policy_name"
  body     = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"
}

## Complete
resource "snowflake_aggregation_policy" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "aggregation_policy_name"
  body     = "CASE WHEN CURRENT_ROLE() = 'ADMIN' THEN NO_AGGREGATION_CONSTRAINT() ELSE AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5) END"
  comment  = "My aggregation policy"
}
