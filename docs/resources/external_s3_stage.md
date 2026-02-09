---
page_title: "snowflake_external_s3_stage Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage external S3 stages. For more information, check external stage documentation https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Temporary stages are not supported because they result in per-session objects.

-> **Note** External changes detection on `credentials`, and `encryption` fields are not supported because Snowflake does not return such settings in DESCRIBE or SHOW STAGE output.

-> **Note** Due to Snowflake limitations, when `directory.auto_refresh` is set to a new value in the configuration, the resource is recreated. When it is unset, the provider alters the whole `directory` field with the `enable` value from the configuration.

-> **Note** Integration based stages are not allowed to be altered to use privatelink endpoint. You must either alter the storage integration itself, or first unset the storage integration from the stage instead.

# snowflake_external_s3_stage (Resource)

Resource used to manage external S3 stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).

```terraform
# Basic resource with storage integration
resource "snowflake_external_s3_stage" "basic" {
  name     = "my_s3_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"
}

# Complete resource with all options
resource "snowflake_external_s3_stage" "complete" {
  name                 = "complete_s3_stage"
  database             = "my_database"
  schema               = "my_schema"
  url                  = "s3://mybucket/mypath/"
  storage_integration  = snowflake_storage_integration.s3.name
  aws_access_point_arn = "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"

  encryption {
    aws_cse {
      master_key = var.s3_master_key
    }
  }

  directory {
    enable            = true
    refresh_on_create = true
    auto_refresh      = false
  }

  comment = "Fully configured S3 external stage"
}

# Resource with AWS key credentials instead of storage integration
resource "snowflake_external_s3_stage" "with_key_credentials" {
  name     = "s3_stage_with_keys"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_key_id     = var.aws_access_key_id
    aws_secret_key = var.aws_secret_access_key
    aws_token      = var.aws_token
  }
}

# Resource with AWS IAM role credentials
resource "snowflake_external_s3_stage" "with_role_credentials" {
  name     = "s3_stage_with_role"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_role = var.aws_role_arn
  }
}

# Resource with SSE-S3 encryption
resource "snowflake_external_s3_stage" "sse_s3" {
  name                = "s3_stage_sse_s3"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    aws_sse_s3 {}
  }
}

# Resource with SSE-KMS encryption
resource "snowflake_external_s3_stage" "sse_kms" {
  name                = "s3_stage_sse_kms"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    aws_sse_kms {
      kms_key_id = var.kms_key_id
    }
  }
}

# Resource with encryption set to none
resource "snowflake_external_s3_stage" "no_encryption" {
  name                = "s3_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    none {}
  }
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database in which to create the stage. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `name` (String) Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `schema` (String) The schema in which to create the stage. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `url` (String) Specifies the URL for the S3 bucket (e.g., 's3://bucket-name/path/').

### Optional

- `aws_access_point_arn` (String) Specifies the ARN for an AWS S3 Access Point to use for data transfer.
- `comment` (String) Specifies a comment for the stage.
- `credentials` (Block List, Max: 1) Specifies the AWS credentials for the external stage. (see [below for nested schema](#nestedblock--credentials))
- `directory` (Block List, Max: 1) Directory tables store a catalog of staged files in cloud storage. (see [below for nested schema](#nestedblock--directory))
- `encryption` (Block List, Max: 1) Specifies the encryption settings for the S3 external stage. (see [below for nested schema](#nestedblock--encryption))
- `storage_integration` (String) Specifies the name of the storage integration used to delegate authentication responsibility to a Snowflake identity. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `use_privatelink_endpoint` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether to use a private link endpoint for S3 storage.

### Read-Only

- `cloud` (String) Specifies a cloud provider for the stage. This field is used for checking external changes and recreating the resources if needed.
- `describe_output` (List of Object) Outputs the result of `DESCRIBE STAGE` for the given stage. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `show_output` (List of Object) Outputs the result of `SHOW STAGES` for the given stage. (see [below for nested schema](#nestedatt--show_output))
- `stage_type` (String) Specifies a type for the stage. This field is used for checking external changes and recreating the resources if needed.

<a id="nestedblock--credentials"></a>
### Nested Schema for `credentials`

Optional:

- `aws_key_id` (String, Sensitive) Specifies the AWS access key ID.
- `aws_role` (String) Specifies the AWS IAM role ARN to use for accessing the bucket.
- `aws_secret_key` (String, Sensitive) Specifies the AWS secret access key.
- `aws_token` (String, Sensitive) Specifies the AWS session token for temporary credentials.


<a id="nestedblock--directory"></a>
### Nested Schema for `directory`

Required:

- `enable` (Boolean) Specifies whether to enable a directory table on the external stage.

Optional:

- `auto_refresh` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether Snowflake should enable triggering automatic refreshes of the directory table metadata.
- `refresh_on_create` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Specifies whether to automatically refresh the directory table metadata once, immediately after the stage is created.This field is used only when creating the object. Changes on this field are ignored after creation.


<a id="nestedblock--encryption"></a>
### Nested Schema for `encryption`

Optional:

- `aws_cse` (Block List, Max: 1) AWS client-side encryption using a master key. (see [below for nested schema](#nestedblock--encryption--aws_cse))
- `aws_sse_kms` (Block List, Max: 1) AWS server-side encryption using KMS-managed keys. (see [below for nested schema](#nestedblock--encryption--aws_sse_kms))
- `aws_sse_s3` (Block List, Max: 1) AWS server-side encryption using S3-managed keys. (see [below for nested schema](#nestedblock--encryption--aws_sse_s3))
- `none` (Block List, Max: 1) No encryption. (see [below for nested schema](#nestedblock--encryption--none))

<a id="nestedblock--encryption--aws_cse"></a>
### Nested Schema for `encryption.aws_cse`

Required:

- `master_key` (String, Sensitive) Specifies the 128-bit or 256-bit client-side master key.


<a id="nestedblock--encryption--aws_sse_kms"></a>
### Nested Schema for `encryption.aws_sse_kms`

Optional:

- `kms_key_id` (String) Specifies the KMS-managed key ID.


<a id="nestedblock--encryption--aws_sse_s3"></a>
### Nested Schema for `encryption.aws_sse_s3`


<a id="nestedblock--encryption--none"></a>
### Nested Schema for `encryption.none`



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

- `directory_table` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--directory_table))
- `file_format` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format))
- `location` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--location))
- `privatelink` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--privatelink))

<a id="nestedobjatt--describe_output--directory_table"></a>
### Nested Schema for `describe_output.directory_table`

Read-Only:

- `auto_refresh` (Boolean)
- `enable` (Boolean)


<a id="nestedobjatt--describe_output--file_format"></a>
### Nested Schema for `describe_output.file_format`

Read-Only:

- `csv` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--csv))
- `format_name` (String)

<a id="nestedobjatt--describe_output--file_format--csv"></a>
### Nested Schema for `describe_output.file_format.csv`

Read-Only:

- `binary_format` (String)
- `compression` (String)
- `date_format` (String)
- `empty_field_as_null` (Boolean)
- `encoding` (String)
- `error_on_column_count_mismatch` (Boolean)
- `escape` (String)
- `escape_unenclosed_field` (String)
- `field_delimiter` (String)
- `field_optionally_enclosed_by` (String)
- `file_extension` (String)
- `multi_line` (Boolean)
- `null_if` (List of String)
- `parse_header` (Boolean)
- `record_delimiter` (String)
- `replace_invalid_characters` (Boolean)
- `skip_blank_lines` (Boolean)
- `skip_byte_order_mark` (Boolean)
- `skip_header` (Number)
- `time_format` (String)
- `timestamp_format` (String)
- `trim_space` (Boolean)
- `type` (String)
- `validate_utf8` (Boolean)



<a id="nestedobjatt--describe_output--location"></a>
### Nested Schema for `describe_output.location`

Read-Only:

- `aws_access_point_arn` (String)
- `url` (String)


<a id="nestedobjatt--describe_output--privatelink"></a>
### Nested Schema for `describe_output.privatelink`

Read-Only:

- `use_privatelink_endpoint` (Boolean)



<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `cloud` (String)
- `comment` (String)
- `created_on` (String)
- `database_name` (String)
- `directory_enabled` (Boolean)
- `endpoint` (String)
- `has_credentials` (Boolean)
- `has_encryption_key` (Boolean)
- `name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `region` (String)
- `schema_name` (String)
- `storage_integration` (String)
- `type` (String)
- `url` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_external_s3_stage.example '"<database_name>"."<schema_name>"."<stage_name>"'
```
