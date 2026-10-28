# from a Snowflake-managed connector definition; the connector carries no configuration yet, so it cannot be
# started until one is supplied
resource "snowflake_openflow_connector" "from_definition" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_connector"
  runtime  = snowflake_openflow_runtime.example.fully_qualified_name

  from {
    definition = "OPENFLOW_POSTGRES_CDC"
  }
}

# from a connector bundle on a stage; the connector arrives already configured. A git repository is a stage,
# so `stage` takes one of those too
resource "snowflake_openflow_connector" "from_stage" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_configured_connector"
  runtime  = snowflake_openflow_runtime.example.fully_qualified_name

  from {
    stage = snowflake_stage.example.fully_qualified_name
    path  = "connectors/postgres"
  }

  # changing this from one non-empty value to another loads the bundle again as a new version of the same
  # connector, rather than recreating it; adding or removing it loads nothing
  version_trigger = "1"
}

# uploading the bundle yourself, then creating the connector from it. Snowflake reads config.json from the
# location and checks that the connector definition it names matches, so that file has to be there; the same
# PUT pattern uploads anything else the bundle needs, such as a driver jar.
#
# PUT runs from the machine executing Terraform, so the local file has to exist there. That rules it out for
# runners that do not have your bundle checked out, where a git repository stage is the better route.
resource "snowflake_stage" "bundles" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_connector_bundles"
}

resource "snowflake_execute" "upload_config" {
  execute = "PUT file:///path/to/bundle/config.json @\"${snowflake_stage.bundles.database}\".\"${snowflake_stage.bundles.schema}\".\"${snowflake_stage.bundles.name}\"/orders/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE"
  revert  = "REMOVE @\"${snowflake_stage.bundles.database}\".\"${snowflake_stage.bundles.schema}\".\"${snowflake_stage.bundles.name}\"/orders/config.json"
}

resource "snowflake_execute" "upload_driver" {
  execute = "PUT file:///path/to/bundle/postgresql.jar @\"${snowflake_stage.bundles.database}\".\"${snowflake_stage.bundles.schema}\".\"${snowflake_stage.bundles.name}\"/orders/ AUTO_COMPRESS = FALSE OVERWRITE = TRUE"
  revert  = "REMOVE @\"${snowflake_stage.bundles.database}\".\"${snowflake_stage.bundles.schema}\".\"${snowflake_stage.bundles.name}\"/orders/postgresql.jar"
}

resource "snowflake_openflow_connector" "from_uploaded_bundle" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_uploaded_connector"
  runtime  = snowflake_openflow_runtime.example.fully_qualified_name

  from {
    stage = snowflake_stage.bundles.fully_qualified_name
    path  = "orders"
  }

  # the bundle has to be on the stage before the connector reads it
  depends_on = [snowflake_execute.upload_config, snowflake_execute.upload_driver]
}

# from a git repository, which is a stage as far as Snowflake is concerned. Preferred when the bundle is
# version controlled, since nothing has to be uploaded from the machine running Terraform.
resource "snowflake_openflow_connector" "from_git" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_git_connector"
  runtime  = snowflake_openflow_runtime.example.fully_qualified_name

  from {
    stage = snowflake_git_repository.example.fully_qualified_name
    path  = "branches/main/connectors/orders"
  }
}

# complete resource
resource "snowflake_openflow_connector" "complete" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_connector_complete"
  runtime  = snowflake_openflow_runtime.example.fully_qualified_name

  from {
    definition = "OPENFLOW_POSTGRES_CDC"
  }

  display_name = "My connector"
  comment      = "Managed by Terraform."

  timeouts {
    create = "45m"
    update = "45m"
    delete = "45m"
  }
}
