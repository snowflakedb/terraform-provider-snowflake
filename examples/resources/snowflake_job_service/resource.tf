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

# basic resource - from specification template file on stage
resource "snowflake_job_service" "basic" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification_template {
    stage = snowflake_stage.test.fully_qualified_name
    path  = "path/to/spec"
    file  = "spec.yaml"
    using {
      key   = "tag"
      value = "latest"
    }
  }
}


# basic resource - from specification template content
resource "snowflake_job_service" "basic" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = "SERVICE"
  in_compute_pool = snowflake_compute_pool.test.name
  from_specification {
    text = <<-EOT
spec:
 containers:
 - name: {{ tag }}
   image: /database/schema/image_repository/exampleimage:latest
   EOT
    using {
      key   = "tag"
      value = "latest"
    }
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
  query_warehouse = snowflake_warehouse.test.name
  comment         = "A service."
}
