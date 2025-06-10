# basic resource - from specification file on stage
resource "snowflake_service" "basic" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    stage = snowflake_stage.basic.fully_qualified_name
    file  = "spec.yaml"
  }
}

# basic resource - from specification content
resource "snowflake_service" "basic" {
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
resource "snowflake_compute_pool" "complete" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    stage = snowflake_stage.complete.fully_qualified_name
    # or, with explicit stage value
    # stage = "\"DATABASE\".\"SCHEMA\".\"STAGE\""
    path = "path/to/spec"
    file = "spec.yaml"
  }
  auto_suspend_secs = 1200
  external_access_integrations = [
    "INTEGRATION"
  ]
  auto_resume         = true
  min_instances       = 1
  min_ready_instances = 1
  max_instances       = 2
  query_warehouse     = snowflake_warehouse.test.name
  comment             = "A service."
}
