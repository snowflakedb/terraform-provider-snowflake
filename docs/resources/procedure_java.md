---
page_title: "snowflake_procedure_java Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage java procedure objects. For more information, check procedure documentation https://docs.snowflake.com/en/sql-reference/sql/create-procedure.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

!> **Sensitive values** This resource's `procedure_definition` and `show_output.arguments_raw` fields are not marked as sensitive in the provider. Ensure that no personal data, sensitive data, export-controlled data, or other regulated data is entered as metadata when using the provider. If you use one of these fields, they may be present in logs, so ensure that the provider logs are properly restricted. For more information, see [Sensitive values limitations](../#sensitive-values-limitations) and [Metadata fields in Snowflake](https://docs.snowflake.com/en/sql-reference/metadata).

-> **Note** External changes to `is_secure` and `null_input_behavior` are not currently supported. They will be handled in the following versions of the provider which may still affect this resource.

-> **Note** `COPY GRANTS` and `OR REPLACE` are not currently supported.

-> **Note** `RETURN... [[ NOT ] NULL]` is not currently supported. It will be improved in the following versions of the provider which may still affect this resource.

-> **Note** Use of return type `TABLE` is currently limited. It will be improved in the following versions of the provider which may still affect this resource.

-> **Note** Snowflake is not returning full data type information for arguments which may lead to unexpected plan outputs. Diff suppression for such cases will be improved.

-> **Note** Snowflake is not returning the default values for arguments so argument's `arg_default_value` external changes cannot be tracked.

-> **Note** Limit the use of special characters (`.`, `'`, `/`, `"`, `(`, `)`, `[`, `]`, `{`, `}`, ` `) in argument names, stage ids, and secret ids. It's best to limit to only alphanumeric and underscores. There is a lot of parsing of SHOW/DESCRIBE outputs involved and using special characters may limit the possibility to achieve the correct results.

-> **Note** Diff suppression for `import.stage_location` is currently not implemented. Make sure you are using the fully qualified stage name (with all the double quotes), e.g. `"\"database_name\".\"schema_name\".\"stage_name\""`.

~> **Required warehouse** This resource may require active warehouse. Please, make sure you have either set a DEFAULT_WAREHOUSE for the user, or specified a warehouse in the provider configuration.

# snowflake_procedure_java (Resource)

Resource used to manage java procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).

## Example Usage

```terraform
# basic example
resource "snowflake_procedure_java" "basic" {
  database = "Database"
  schema   = "Schema"
  name     = "ProcedureName"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "TestFunc.echoVarchar"
  procedure_definition = <<EOT
  import com.snowflake.snowpark_java.*;
  class TestFunc {
    public static String echoVarchar(Session session, String x) {
      return x;
    }
  }
EOT
  runtime_version      = "11"
  snowpark_package     = "1.14.0"
}

# full example
resource "snowflake_procedure_java" "full" {
  database = "Database"
  schema   = "Schema"
  name     = "ProcedureName"
  arguments {
    arg_data_type = "VARCHAR(100)"
    arg_name      = "x"
  }
  return_type          = "VARCHAR(100)"
  handler              = "TestFunc.echoVarchar"
  procedure_definition = <<EOT
    import com.snowflake.snowpark_java.*;
  class TestFunc {
    public static String echoVarchar(Session session, String x) {
      return x;
    }
  }
EOT
  runtime_version      = "11"
  snowpark_package     = "1.14.0"

  comment    = "some comment"
  execute_as = "CALLER"
  target_path {
    path_on_stage  = "tf-1734028493-OkoTf.jar"
    stage_location = snowflake_stage.example.fully_qualified_name
  }
  packages = ["com.snowflake:telemetry:0.1.0"]
  imports {
    path_on_stage  = "tf-1734028486-OLJpF.jar"
    stage_location = "~"
  }
  imports {
    path_on_stage  = "tf-1734028491-EMoDC.jar"
    stage_location = "~"
  }
  is_secure           = "false"
  null_input_behavior = "CALLED ON NULL INPUT"
  external_access_integrations = [
    "INTEGRATION_1", "INTEGRATION_2"
  ]
  secrets {
    secret_id            = snowflake_secret_with_generic_string.example1.fully_qualified_name
    secret_variable_name = "abc"
  }
  secrets {
    secret_id            = snowflake_secret_with_generic_string.example2.fully_qualified_name
    secret_variable_name = "def"
  }
}
```
-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database in which to create the procedure. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `handler` (String) Use the fully qualified name of the method or function for the stored procedure. This is typically in the following form `com.my_company.my_package.MyClass.myMethod` where `com.my_company.my_package` corresponds to the package containing the object or class: `package com.my_company.my_package;`.
- `name` (String) The name of the procedure; the identifier does not need to be unique for the schema in which the procedure is created because stored procedures are [identified and resolved by the combination of the name and argument types](https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-naming-conventions.html#label-procedure-function-name-overloading). Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `return_type` (String) Specifies the type of the result returned by the stored procedure. For `<result_data_type>`, use the Snowflake data type that corresponds to the type of the language that you are using (see [SQL-Java Data Type Mappings](https://docs.snowflake.com/en/developer-guide/udf-stored-procedure-data-type-mapping.html#label-sql-java-data-type-mappings)). For `RETURNS TABLE ( [ col_name col_data_type [ , ... ] ] )`, if you know the Snowflake data types of the columns in the returned table, specify the column names and types. Otherwise (e.g. if you are determining the column types during run time), you can omit the column names and types (i.e. `TABLE ()`).
- `runtime_version` (String) The language runtime version to use. Currently, the supported versions are: 11.
- `schema` (String) The schema in which to create the procedure. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `snowpark_package` (String) The Snowpark package is required for stored procedures, so it must always be present. For more information about Snowpark, see [Snowpark API](https://docs.snowflake.com/en/developer-guide/snowpark/index).

### Optional

- `arguments` (Block List) List of the arguments for the procedure. Consult the [docs](https://docs.snowflake.com/en/sql-reference/sql/create-procedure#all-languages) for more details. (see [below for nested schema](#nestedblock--arguments))
- `comment` (String) (Default: `user-defined procedure`) Specifies a comment for the procedure.
- `enable_console_output` (Boolean) Enable stdout/stderr fast path logging for anonyous stored procs. This is a public parameter (similar to LOG_LEVEL). For more information, check [ENABLE_CONSOLE_OUTPUT docs](https://docs.snowflake.com/en/sql-reference/parameters#enable-console-output).
- `execute_as` (String) Specifies whether the stored procedure executes with the privileges of the owner (an “owner’s rights” stored procedure) or with the privileges of the caller (a “caller’s rights” stored procedure). If you execute the statement CREATE PROCEDURE … EXECUTE AS CALLER, then in the future the procedure will execute as a caller’s rights procedure. If you execute CREATE PROCEDURE … EXECUTE AS OWNER, then the procedure will execute as an owner’s rights procedure. For more information, see [Understanding caller’s rights and owner’s rights stored procedures](https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-rights). Valid values are (case-insensitive): `CALLER` | `OWNER`.
- `external_access_integrations` (Set of String) The names of [external access integrations](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) needed in order for this procedure’s handler code to access external networks. An external access integration specifies [network rules](https://docs.snowflake.com/en/sql-reference/sql/create-network-rule) and [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) that specify external locations and credentials (if any) allowed for use by handler code when making requests of an external network, such as an external REST API.
- `imports` (Block Set) The location (stage), path, and name of the file(s) to import. You must set the IMPORTS clause to include any files that your stored procedure depends on. If you are writing an in-line stored procedure, you can omit this clause, unless your code depends on classes defined outside the stored procedure or resource files. If you are writing a stored procedure with a staged handler, you must also include a path to the JAR file containing the stored procedure’s handler code. The IMPORTS definition cannot reference variables from arguments that are passed into the stored procedure. Each file in the IMPORTS clause must have a unique name, even if the files are in different subdirectories or different stages. (see [below for nested schema](#nestedblock--imports))
- `is_secure` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies that the procedure is secure. For more information about secure procedures, see [Protecting Sensitive Information with Secure UDFs and Stored Procedures](https://docs.snowflake.com/en/developer-guide/secure-udf-procedure). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `log_level` (String) LOG_LEVEL to use when filtering events For more information, check [LOG_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#log-level).
- `metric_level` (String) METRIC_LEVEL value to control whether to emit metrics to Event Table For more information, check [METRIC_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#metric-level).
- `null_input_behavior` (String) Specifies the behavior of the procedure when called with null inputs. Valid values are (case-insensitive): `CALLED ON NULL INPUT` | `RETURNS NULL ON NULL INPUT`.
- `packages` (Set of String) List of the names of packages deployed in Snowflake that should be included in the handler code’s execution environment. The Snowpark package is required for stored procedures, but is specified in the `snowpark_package` attribute. For more information about Snowpark, see [Snowpark API](https://docs.snowflake.com/en/developer-guide/snowpark/index).
- `procedure_definition` (String) Defines the code executed by the stored procedure. The definition can consist of any valid code. Wrapping `$$` signs are added by the provider automatically; do not include them. The `procedure_definition` value must be Java source code. For more information, see [Java (using Snowpark)](https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java). To mitigate permadiff on this field, the provider replaces blank characters with a space. This can lead to false positives in cases where a change in case or run of whitespace is semantically significant.
- `secrets` (Block Set) Assigns the names of [secrets](https://docs.snowflake.com/en/sql-reference/sql/create-secret) to variables so that you can use the variables to reference the secrets when retrieving information from secrets in handler code. Secrets you specify here must be allowed by the [external access integration](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration) specified as a value of this CREATE FUNCTION command’s EXTERNAL_ACCESS_INTEGRATIONS parameter. (see [below for nested schema](#nestedblock--secrets))
- `target_path` (Block Set, Max: 1) Use the fully qualified name of the method or function for the stored procedure. This is typically in the following form `com.my_company.my_package.MyClass.myMethod` where `com.my_company.my_package` corresponds to the package containing the object or class: `package com.my_company.my_package;`. (see [below for nested schema](#nestedblock--target_path))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `trace_level` (String) Trace level value to use when generating/filtering trace events For more information, check [TRACE_LEVEL docs](https://docs.snowflake.com/en/sql-reference/parameters#trace-level).

### Read-Only

- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `parameters` (List of Object) Outputs the result of `SHOW PARAMETERS IN PROCEDURE` for the given procedure. (see [below for nested schema](#nestedatt--parameters))
- `procedure_language` (String) Specifies language for the procedure. Used to detect external changes.
- `show_output` (List of Object) Outputs the result of `SHOW PROCEDURE` for the given procedure. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--arguments"></a>
### Nested Schema for `arguments`

Required:

- `arg_data_type` (String) The argument type.
- `arg_name` (String) The argument name. The provider wraps it in double quotes by default, so be aware of that while referencing the argument in the procedure definition.

Optional:

- `arg_default_value` (String) Optional default value for the argument. For text values use single quotes. Numeric values can be unquoted. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".


<a id="nestedblock--imports"></a>
### Nested Schema for `imports`

Required:

- `path_on_stage` (String) Path for import on stage, without the leading `/`.
- `stage_location` (String) Stage location without leading `@`. To use your user's stage set this to `~`, otherwise pass fully qualified name of the stage (with every part contained in double quotes or use `snowflake_stage.<your stage's resource name>.fully_qualified_name` if you manage this stage through terraform).


<a id="nestedblock--secrets"></a>
### Nested Schema for `secrets`

Required:

- `secret_id` (String) Fully qualified name of the allowed [secret](https://docs.snowflake.com/en/sql-reference/sql/create-secret). You will receive an error if you specify a SECRETS value whose secret isn’t also included in an integration specified by the EXTERNAL_ACCESS_INTEGRATIONS parameter.
- `secret_variable_name` (String) The variable that will be used in handler code when retrieving information from the secret.


<a id="nestedblock--target_path"></a>
### Nested Schema for `target_path`

Required:

- `path_on_stage` (String) Path for import on stage, without the leading `/`.
- `stage_location` (String) Stage location without leading `@`. To use your user's stage set this to `~`, otherwise pass fully qualified name of the stage (with every part contained in double quotes or use `snowflake_stage.<your stage's resource name>.fully_qualified_name` if you manage this stage through terraform).


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--parameters"></a>
### Nested Schema for `parameters`

Read-Only:

- `enable_console_output` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--enable_console_output))
- `log_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--log_level))
- `metric_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--metric_level))
- `trace_level` (List of Object) (see [below for nested schema](#nestedobjatt--parameters--trace_level))

<a id="nestedobjatt--parameters--enable_console_output"></a>
### Nested Schema for `parameters.enable_console_output`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--log_level"></a>
### Nested Schema for `parameters.log_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--metric_level"></a>
### Nested Schema for `parameters.metric_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)


<a id="nestedobjatt--parameters--trace_level"></a>
### Nested Schema for `parameters.trace_level`

Read-Only:

- `default` (String)
- `description` (String)
- `key` (String)
- `level` (String)
- `value` (String)



<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `arguments_raw` (String)
- `catalog_name` (String)
- `created_on` (String)
- `description` (String)
- `external_access_integrations` (String)
- `is_aggregate` (Boolean)
- `is_ansi` (Boolean)
- `is_builtin` (Boolean)
- `is_secure` (Boolean)
- `is_table_function` (Boolean)
- `max_num_arguments` (Number)
- `min_num_arguments` (Number)
- `name` (String)
- `schema_name` (String)
- `secrets` (String)
- `valid_for_clustering` (Boolean)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_procedure_java.example '"<database_name>"."<schema_name>"."<function_name>"(varchar, varchar, varchar)'
```

Note: Snowflake is not returning all information needed to populate the state correctly after import (e.g. data types with attributes like NUMBER(32, 10) are returned as NUMBER, default values for arguments are not returned at all).
Also, `ALTER` for functions is very limited so most of the attributes on this resource are marked as force new. Because of that, in multiple situations plan won't be empty after importing and manual state operations may be required.
