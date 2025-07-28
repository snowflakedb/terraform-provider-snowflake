---
page_title: "snowflake_job_service Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage job services. For more information, check services documentation https://docs.snowflake.com/en/sql-reference/sql/execute-job-service. Executes a Snowpark Container Services service as a job. A service, created using CREATE SERVICE, is long-running and you must explicitly stop it when it is no longer needed. On the other hand, a job, created using EXECUTE JOB SERVICE (with ASYNC=TRUE in this resource), returns immediately while the job is running. See Working with services https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services developer guide for more details.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

<!-- TODO(SNOW-2129584): address this limitation -->
!> **Caution** Only asynchronous job services are supported. This resource uses `ASYNC=TRUE` during creation. In this case, the command returns immediately while the job is running. Creating a job service automatically executes the job (read more in [EXECUTE JOB SERVICE docs](https://docs.snowflake.com/en/sql-reference/sql/execute-job-service)).

-> **Note** For asynchronous jobs, Snowflake does not perform automatic cleanup after completion. You must either remove the resource or execute the `DROP SERVICE` command to remove the job. If you want to execute the job again, use the [replace flag](https://developer.hashicorp.com/terraform/cli/commands/apply#replace-resource).

-> **Note** During resource deletion and recreation, the provider uses `DROP SERVICE` with `FORCE` option to properly handle services with block storage volumes. Read more in [docs](https://docs.snowflake.com/en/sql-reference/sql/drop-service#force-option).

# snowflake_job_service (Resource)

Resource used to manage job services. For more information, check [services documentation](https://docs.snowflake.com/en/sql-reference/sql/execute-job-service). Executes a Snowpark Container Services service as a job. A service, created using `CREATE SERVICE`, is long-running and you must explicitly stop it when it is no longer needed. On the other hand, a job, created using EXECUTE JOB SERVICE (with `ASYNC=TRUE` in this resource), returns immediately while the job is running. See [Working with services](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services) developer guide for more details.

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

```terraform
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

- `comment` (String) Specifies a comment for the service.
- `external_access_integrations` (Set of String) Specifies the names of the external access integrations that allow your service to access external sites.
- `from_specification` (Block List, Max: 1) Specifies the service specification to use for the service. Note that external changes on this field and nested fields are not detected. Use correctly formatted YAML files. Watch out for the space/tabs indentation. See [service specification](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/specification-reference#general-guidelines) for more information. (see [below for nested schema](#nestedblock--from_specification))
- `from_specification_template` (Block List, Max: 1) Specifies the service specification template to use for the service. Note that external changes on this field and nested fields are not detected. Use correctly formatted YAML files. Watch out for the space/tabs indentation. See [service specification](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/specification-reference#general-guidelines) for more information. (see [below for nested schema](#nestedblock--from_specification_template))
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
- `stage` (String) The fully qualified name of the stage containing the service specification file. At symbol (`@`) is added automatically. Example: `"\"<db_name>\".\"<schema_name>\".\"<stage_name>\""`. For more information about this resource, see [docs](./stage).
- `text` (String) The embedded text of the service specification.


<a id="nestedblock--from_specification_template"></a>
### Nested Schema for `from_specification_template`

Required:

- `using` (Block List, Min: 1) List of the specified template variables and the values of those variables. (see [below for nested schema](#nestedblock--from_specification_template--using))

Optional:

- `file` (String) The file name of the service specification template. Example: `spec.yaml`.
- `path` (String) The path to the service specification template file on the given stage. When the path is specified, the `/` character is automatically added as a path prefix. Example: `path/to/spec`.
- `stage` (String) The fully qualified name of the stage containing the service specification template file. At symbol (`@`) is added automatically. Example: `"\"<db_name>\".\"<schema_name>\".\"<stage_name>\""`. For more information about this resource, see [docs](./stage).
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
terraform import snowflake_job_service.example '"<database_name>"."<schema_name>"."<job_service_name>"'
```
