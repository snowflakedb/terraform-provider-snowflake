---
page_title: "snowflake_service Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage services. For more information, check services documentation https://docs.snowflake.com/en/sql-reference/sql/create-service. A long-running service is like a web service that does not end automatically. After you create a service, Snowflake manages the running service. For example, if a service container stops, for whatever reason, Snowflake restarts that container so the service runs uninterrupted. See Working with services https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services developer guide for more details.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Managing service state is limited. It is handled by `auto_suspend_secs`, and `auto_resume` fields. The provider does not support managing the state of services in Snowflake with `ALTER ... SUSPEND` and `ALTER ... RESUME`. See [Suspending a service documentation](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services#suspending-a-service) for more details.

-> **Note** During resource deletion and recreation, the provider uses `DROP SERVICE` with `FORCE` option to properly handle services with block storage volumes. Read more in [docs](https://docs.snowflake.com/en/sql-reference/sql/drop-service#force-option).

# snowflake_service (Resource)

Resource used to manage services. For more information, check [services documentation](https://docs.snowflake.com/en/sql-reference/sql/create-service). A long-running service is like a web service that does not end automatically. After you create a service, Snowflake manages the running service. For example, if a service container stops, for whatever reason, Snowflake restarts that container so the service runs uninterrupted. See [Working with services](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services) developer guide for more details.

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

```terraform
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

# basic resource - from specification template file on stage
resource "snowflake_service" "basic" {
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
resource "snowflake_service" "basic" {
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
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `compute_pool` (String) Specifies the name of the compute pool in your account on which to run the service. Identifiers with special or lower-case characters are not supported. This limitation in the provider follows the limitation in Snowflake (see [docs](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool)). Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `database` (String) The database in which to create the service. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `name` (String) Specifies the identifier for the service; must be unique for the schema in which the service is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `schema` (String) The schema in which to create the service. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `auto_resume` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether to automatically resume a service. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `auto_suspend_secs` (Number) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`-1`)) Specifies the number of seconds of inactivity (service is idle) after which Snowflake automatically suspends the service.
- `comment` (String) Specifies a comment for the service.
- `external_access_integrations` (Set of String) Specifies the names of the external access integrations that allow your service to access external sites.
- `from_specification` (Block List, Max: 1) Specifies the service specification to use for the service. Note that external changes on this field and nested fields are not detected. Use correctly formatted YAML files. Watch out for the space/tabs indentation. See [service specification](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/specification-reference#general-guidelines) for more information. (see [below for nested schema](#nestedblock--from_specification))
- `from_specification_template` (Block List, Max: 1) Specifies the service specification template to use for the service. Note that external changes on this field and nested fields are not detected. Use correctly formatted YAML files. Watch out for the space/tabs indentation. See [service specification](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/specification-reference#general-guidelines) for more information. (see [below for nested schema](#nestedblock--from_specification_template))
- `max_instances` (Number) Specifies the maximum number of service instances to run.
- `min_instances` (Number) Specifies the minimum number of service instances to run.
- `min_ready_instances` (Number) Indicates the minimum service instances that must be ready for Snowflake to consider the service is ready to process requests.
- `query_warehouse` (String) Warehouse to use if a service container connects to Snowflake to execute a query but does not explicitly specify a warehouse to use. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE SERVICE` for the given service. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `service_type` (String) Specifies a type for the service. This field is used for checking external changes and recreating the resources if needed.
- `show_output` (List of Object) Outputs the result of `SHOW SERVICES` for the given service. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--from_specification"></a>
### Nested Schema for `from_specification`

Optional:

- `file` (String) The file name of the service specification. Example: `spec.yaml`.
- `path` (String) The path to the service specification file on the given stage. When the path is specified, the `/` character is automatically added as a path prefix. Example: `path/to/spec`.
- `stage` (String) The fully qualified name of the stage containing the service specification file. At symbol (`@`) is added automatically. Example: `"\"<db_name>\".\"<schema_name>\".\"<schema_name>\""`. For more information about this resource, see [docs](./stage).
- `text` (String) The embedded text of the service specification.


<a id="nestedblock--from_specification_template"></a>
### Nested Schema for `from_specification_template`

Required:

- `using` (Block List, Min: 1) List of the specified template variables and the values of those variables. (see [below for nested schema](#nestedblock--from_specification_template--using))

Optional:

- `file` (String) The file name of the service specification template. Example: `spec.yaml`.
- `path` (String) The path to the service specification template file on the given stage. When the path is specified, the `/` character is automatically added as a path prefix. Example: `path/to/spec`.
- `stage` (String) The fully qualified name of the stage containing the service specification template file. At symbol (`@`) is added automatically. Example: `"\"<db_name>\".\"<schema_name>\".\"<schema_name>\""`. For more information about this resource, see [docs](./stage).
- `text` (String) The embedded text of the service specification template.

<a id="nestedblock--from_specification_template--using"></a>
### Nested Schema for `from_specification_template.using`

Required:

- `key` (String) The name of the template variable. The provider wraps it in double quotes by default, so be aware of that while referencing the argument in the spec definition.
- `value` (String) The value to assign to the variable in the template. The provider wraps it in `$$` by default, so be aware of that while referencing the argument in the spec definition. The value must either be alphanumeric or valid JSON.



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

- `auto_resume` (Boolean)
- `auto_suspend_secs` (Number)
- `comment` (String)
- `compute_pool` (String)
- `created_on` (String)
- `current_instances` (Number)
- `database_name` (String)
- `dns_name` (String)
- `external_access_integrations` (Set of String)
- `is_async_job` (Boolean)
- `is_job` (Boolean)
- `is_upgrading` (Boolean)
- `managing_object_domain` (String)
- `managing_object_name` (String)
- `max_instances` (Number)
- `min_instances` (Number)
- `min_ready_instances` (Number)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `query_warehouse` (String)
- `resumed_on` (String)
- `schema_name` (String)
- `spec` (String)
- `spec_digest` (String)
- `status` (String)
- `suspended_on` (String)
- `target_instances` (Number)
- `updated_on` (String)


<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `auto_resume` (Boolean)
- `auto_suspend_secs` (Number)
- `comment` (String)
- `compute_pool` (String)
- `created_on` (String)
- `current_instances` (Number)
- `database_name` (String)
- `dns_name` (String)
- `external_access_integrations` (Set of String)
- `is_async_job` (Boolean)
- `is_job` (Boolean)
- `is_upgrading` (Boolean)
- `managing_object_domain` (String)
- `managing_object_name` (String)
- `max_instances` (Number)
- `min_instances` (Number)
- `min_ready_instances` (Number)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `query_warehouse` (String)
- `resumed_on` (String)
- `schema_name` (String)
- `spec_digest` (String)
- `status` (String)
- `suspended_on` (String)
- `target_instances` (Number)
- `updated_on` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_service.example '"<database_name>"."<schema_name>"."<service_name>"'
```
