# basic resource
resource "snowflake_openflow_deployment_snowflake_managed" "basic" {
  name = "my_deployment"
}

# with custom timeouts (provisioning a deployment takes considerably longer than the defaults allow)
resource "snowflake_openflow_deployment_snowflake_managed" "with_timeouts" {
  name = "my_deployment"

  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}

# complete resource
resource "snowflake_openflow_deployment_snowflake_managed" "complete" {
  name         = "my_deployment_complete"
  display_name = "My deployment"
  # an existing event table, created with CREATE EVENT TABLE outside Terraform: the provider has no
  # resource for one, and a regular table will not do because Snowflake requires the event schema
  event_table      = "\"<database_name>\".\"<schema_name>\".\"<event_table_name>\""
  use_private_link = "true"
  comment          = "Managed by Terraform."
}
