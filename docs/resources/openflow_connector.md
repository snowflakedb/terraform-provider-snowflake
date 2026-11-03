---
page_title: "snowflake_openflow_connector Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage Openflow connectors, which run inside an Openflow runtime. Every mutating statement is asynchronous, so create and update return once the connector settles. Starting and stopping are operational actions and are not exposed here, although the provider may execute such commands when altering the object; the connector's state is reported in show_output. For more information, check Openflow connector documentation https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Every mutating statement is asynchronous and the provider waits for the connector to settle before returning, so applies take minutes rather than seconds. If you encounter timeout errors, use a [`timeouts` block](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts) to set higher limits for your environment.

-> **Note** `from` is create-only and is not read back. Snowflake resolves a connector definition whichever source the connector was created from, so a value read from `SHOW` could not be told apart from a configured one; `show_output.connector_definition` reports what Snowflake resolved. External changes to the block are therefore not detected, and after `terraform import` the block is absent from state, so the first plan asks to replace the connector.

-> **Note** Starting and stopping are operational actions rather than desired state, so they are not exposed here, although the provider may execute such commands when altering the object. The connector's current status is reported in `show_output`.

# snowflake_openflow_connector (Resource)

Resource used to manage Openflow connectors, which run inside an Openflow runtime. Every mutating statement is asynchronous, so create and update return once the connector settles. Starting and stopping are operational actions and are not exposed here, although the provider may execute such commands when altering the object; the connector's state is reported in `show_output`. For more information, check [Openflow connector documentation](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector).

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

```terraform
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
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database in which to create the Openflow connector. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `from` (Block List, Min: 1, Max: 1) Specifies what the connector is created from. Snowflake has no ALTER for it, so changing it recreates the connector. Note that external changes on this field and nested fields are not detected: Snowflake resolves a definition for a connector created from a stage too, so a value read from SHOW could not be told apart from a configured one. `show_output.connector_definition` reports what Snowflake resolved. (see [below for nested schema](#nestedblock--from))
- `name` (String) Specifies the identifier for the Openflow connector; must be unique for the schema in which the connector is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `runtime` (String) Specifies the fully qualified name of the Openflow runtime the connector runs in. The connector is created in the runtime's schema, so `database` and `schema` must match it. Snowflake has no ALTER for it, so changing it recreates the connector.
- `schema` (String) The schema in which to create the Openflow connector. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `comment` (String) Specifies a comment for the Openflow connector.
- `display_name` (String) A free-text alias for the connector. Shown in the Openflow UI in place of the connector's identifier when set.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `version_trigger` (String) An arbitrary value, such as a hash of the bundle's files, that, if and only if changed from a non-empty to a different non-empty value (or known after apply), loads the bundle at `from` again as a new configuration version of the same connector, rather than recreating it. When you add this field to the configuration, or remove it, nothing is loaded, so add it before the bundle changes. A running connector is stopped for the load and started again afterwards; a stopped one stays stopped. When a version is loaded, `show_output` and `describe_output` are marked as computed. Only for a connector created from a stage.

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE OPENFLOW CONNECTOR` for the given connector. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `show_output` (List of Object) Outputs the result of `SHOW OPENFLOW CONNECTORS` for the given connector. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--from"></a>
### Nested Schema for `from`

Optional:

- `definition` (String) Catalog definition ID for the connector type, for example `OPENFLOW_POSTGRES_CDC`. List the available IDs with the `snowflake_openflow_connector_definitions` data source. A connector created this way is a draft: it settles on STOPPED and stays there until a configuration version is committed, which this resource does not do.
- `path` (String) Path to the bundle within the stage. The bundle's root is used when omitted.
- `stage` (String) Identifier of a stage holding a complete configuration bundle, which is how a connector arrives already configured and able to start without a commit. A git repository stage works here too.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--describe_output"></a>
### Nested Schema for `describe_output`

Read-Only:

- `comment` (String)
- `connector_definition` (String)
- `connector_url` (String)
- `default_version` (String)
- `default_version_alias` (String)
- `default_version_git_commit_hash` (String)
- `default_version_location_uri` (String)
- `default_version_name` (String)
- `default_version_source_location_uri` (String)
- `display_name` (String)
- `last_version_alias` (String)
- `last_version_git_commit_hash` (String)
- `last_version_location_uri` (String)
- `last_version_name` (String)
- `last_version_source_location_uri` (String)
- `live_version_location_uri` (String)
- `name` (String)
- `owner` (String)
- `runtime` (String)
- `status` (String)


<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `comment` (String)
- `connector_definition` (String)
- `connector_url` (String)
- `created_on` (String)
- `database_name` (String)
- `default_version` (String)
- `default_version_alias` (String)
- `default_version_location_uri` (String)
- `default_version_name` (String)
- `default_version_source_location_uri` (String)
- `display_name` (String)
- `live_version_location_uri` (String)
- `name` (String)
- `owner` (String)
- `runtime` (String)
- `schema_name` (String)
- `status` (String)
- `updated_on` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_openflow_connector.example '"<database_name>"."<schema_name>"."<connector_name>"'
```
