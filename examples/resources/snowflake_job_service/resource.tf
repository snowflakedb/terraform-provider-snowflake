# basic resource - from specification file on stage
resource "snowflake_job_service" "basic" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "SERVICE"
  in_compute_pool = "COMPUTE_POOL"
  from_specification {
    stage = "\"DATABASE\".\"SCHEMA\".\"STAGE\""
    file  = "spec.yaml"
  }
}

# basic resource - from specification content
resource "snowflake_job_service" "basic" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "SERVICE"
  in_compute_pool = "COMPUTE_POOL"
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
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "SERVICE"
  in_compute_pool = "COMPUTE_POOL"
  from_specification {
    stage = "\"DATABASE\".\"SCHEMA\".\"STAGE\""
    file  = "spec.yaml"
  }
  external_access_integrations = [
    "INTEGRATION"
  ]
  async           = true
  query_warehouse = "WAREHOUSE"
  comment         = "A service."
}
