# basic resource - from specification file on stage
resource "snowflake_service" "basic" {
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
resource "snowflake_service" "basic" {
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
resource "snowflake_compute_pool" "complete" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "SERVICE"
  in_compute_pool = "COMPUTE_POOL"
  from_specification {
    stage = "\"DATABASE\".\"SCHEMA\".\"STAGE\""
    file  = "spec.yaml"
  }
  auto_suspend_secs = 1200
  external_access_integrations = [
    "INTEGRATION"
  ]
  auto_resume         = true
  min_instances       = 1
  min_ready_instances = 1
  max_instances       = 2
  query_warehouse     = "WAREHOUSE"
  comment             = "A service."
}
