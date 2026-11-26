# basic resource
resource "snowflake_openflow_runtime" "basic" {
  database        = "my_database"
  schema          = "my_schema"
  name            = "my_runtime"
  deployment      = "my_deployment"
  node_type       = "SMALL"
  min_nodes       = 1
  max_nodes       = 1
  execute_as_role = "MY_OPENFLOW_ROLE"
}

# complete resource
resource "snowflake_openflow_runtime" "complete" {
  database        = "my_database"
  schema          = "my_schema"
  name            = "my_runtime_complete"
  deployment      = "my_deployment"
  node_type       = "MEDIUM"
  node_type_tier  = "M6"
  min_nodes       = 1
  max_nodes       = 4
  execute_as_role = "MY_OPENFLOW_ROLE"

  # connectors running in this runtime reach the outside world through these
  external_access_integrations = ["my_external_access_integration"]

  display_name = "My runtime"
  comment      = "Managed by Terraform."

  # provisioning a runtime takes considerably longer than the defaults allow
  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}
