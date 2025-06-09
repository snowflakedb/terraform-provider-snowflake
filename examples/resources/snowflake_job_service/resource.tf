# basic resource - from specification file on stage
resource "snowflake_job_service" "basic" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    stage = snowflake_stage.test.fully_qualified_name
    file  = "spec.yaml"
  }
}

# basic resource - from specification content
resource "snowflake_job_service" "basic" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    text = <<-EOT
spec:
  containers:
  - name: example-container
    image: /database/schema/image_repository/exampleimage:latest
    EOT
  }
}

# complete resource
resource "snowflake_job_service" "complete" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    stage = snowflake_stage.test.fully_qualified_name
    # or, with explicit stage value
    # stage = "\"DATABASE\".\"SCHEMA\".\"STAGE\""
    path = "path/to/spec"
    file = "spec.yaml"
  }
  external_access_integrations = [
    "INTEGRATION"
  ]
  async           = true
  query_warehouse = "WAREHOUSE"
  comment         = "A service."
}
