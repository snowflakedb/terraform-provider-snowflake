---
page_title: "snowflake_stage_external_s3 Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage external S3 stages. For more information, check external stage documentation https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Temporary stages are not supported because they result in per-session objects.

-> **Note** External changes detection on `credentials`, and `encryption` fields are not supported because Snowflake does not return such settings in DESCRIBE or SHOW STAGE output.

-> **Note** Due to Snowflake limitations, when `directory.auto_refresh` is set to a new value in the configuration, the resource is recreated. When it is unset, the provider alters the whole `directory` field with the `enable` value from the configuration.

-> **Note** Integration based stages are not allowed to be altered to use privatelink endpoint. You must either alter the storage integration itself, or first unset the storage integration from the stage instead.

-> **Note** This resource is meant only for S3 stages, not S3-compatible stages. For S3-compatible stages, use the `snowflake_stage_external_s3_compatible` resource instead. Do not use this resource with `s3compat://` URLs.

# snowflake_stage_external_s3 (Resource)

Resource used to manage external S3 stages. For more information, check [external stage documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stage#external-stage-parameters-externalstageparams).

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).

```terraform
# Basic resource with storage integration
resource "snowflake_stage_external_s3" "basic" {
  name     = "my_s3_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"
}

# Complete resource with all options
resource "snowflake_stage_external_s3" "complete" {
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
resource "snowflake_stage_external_s3" "with_key_credentials" {
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
resource "snowflake_stage_external_s3" "with_role_credentials" {
  name     = "s3_stage_with_role"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  credentials {
    aws_role = var.aws_role_arn
  }
}

# Resource with SSE-S3 encryption
resource "snowflake_stage_external_s3" "sse_s3" {
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
resource "snowflake_stage_external_s3" "sse_kms" {
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
resource "snowflake_stage_external_s3" "no_encryption" {
  name                = "s3_stage_no_encryption"
  database            = "my_database"
  schema              = "my_schema"
  url                 = "s3://mybucket/mypath/"
  storage_integration = snowflake_storage_integration.s3.name

  encryption {
    none {}
  }
}

# resource with inline CSV file format
resource "snowflake_stage_external_s3" "with_csv_format" {
  name     = "s3_csv_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    csv {
      compression                    = "GZIP"
      record_delimiter               = "\n"
      field_delimiter                = "|"
      multi_line                     = "false"
      file_extension                 = ".csv"
      skip_header                    = 1 # or parse_header = true
      skip_blank_lines               = "true"
      date_format                    = "AUTO"
      time_format                    = "AUTO"
      timestamp_format               = "AUTO"
      binary_format                  = "HEX"
      escape                         = "\\"
      escape_unenclosed_field        = "\\"
      trim_space                     = "false"
      field_optionally_enclosed_by   = "\""
      null_if                        = ["NULL", ""]
      error_on_column_count_mismatch = "true"
      replace_invalid_characters     = "false"
      empty_field_as_null            = "true"
      skip_byte_order_mark           = "true"
      encoding                       = "UTF8"
    }
  }
}

# resource with inline JSON file format
resource "snowflake_stage_external_s3" "with_json_format" {
  name     = "s3_json_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    json {
      compression                = "AUTO"
      date_format                = "AUTO"
      time_format                = "AUTO"
      timestamp_format           = "AUTO"
      binary_format              = "HEX"
      trim_space                 = "false"
      multi_line                 = "false"
      null_if                    = ["NULL", ""]
      file_extension             = ".json"
      enable_octal               = "false"
      allow_duplicate            = "false"
      strip_outer_array          = "false"
      strip_null_values          = "false"
      replace_invalid_characters = "false" # or ignore_utf8_errors = true
      skip_byte_order_mark       = "false"
    }
  }
}

# resource with inline AVRO file format
resource "snowflake_stage_external_s3" "with_avro_format" {
  name     = "s3_avro_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    avro {
      compression                = "GZIP"
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline ORC file format
resource "snowflake_stage_external_s3" "with_orc_format" {
  name     = "s3_orc_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    orc {
      trim_space                 = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline Parquet file format
resource "snowflake_stage_external_s3" "with_parquet_format" {
  name     = "s3_parquet_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    parquet {
      compression                = "SNAPPY"
      binary_as_text             = "true"
      use_logical_type           = "true"
      trim_space                 = "false"
      use_vectorized_scanner     = "false"
      replace_invalid_characters = "false"
      null_if                    = ["NULL", ""]
    }
  }
}

# resource with inline XML file format
resource "snowflake_stage_external_s3" "with_xml_format" {
  name     = "s3_xml_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    xml {
      compression                = "AUTO"
      preserve_space             = "false"
      strip_outer_element        = "false"
      disable_auto_convert       = "false"
      replace_invalid_characters = "false" # or ignore_utf8_errors = true
      skip_byte_order_mark       = "false"
    }
  }
}

# resource with named file format
resource "snowflake_stage_external_s3" "with_named_format" {
  name     = "s3_named_format_stage"
  database = "my_database"
  schema   = "my_schema"
  url      = "s3://mybucket/mypath/"

  file_format {
    format_name = snowflake_file_format.test.fully_qualified_name
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
- `file_format` (Block List, Max: 1) Specifies the file format for the stage. (see [below for nested schema](#nestedblock--file_format))
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



<a id="nestedblock--file_format"></a>
### Nested Schema for `file_format`

Optional:

- `avro` (Block List, Max: 1) AVRO file format options. (see [below for nested schema](#nestedblock--file_format--avro))
- `csv` (Block List, Max: 1) CSV file format options. (see [below for nested schema](#nestedblock--file_format--csv))
- `format_name` (String) Fully qualified name of the file format (e.g., 'database.schema.format_name').
- `json` (Block List, Max: 1) JSON file format options. (see [below for nested schema](#nestedblock--file_format--json))
- `orc` (Block List, Max: 1) ORC file format options. (see [below for nested schema](#nestedblock--file_format--orc))
- `parquet` (Block List, Max: 1) Parquet file format options. (see [below for nested schema](#nestedblock--file_format--parquet))
- `xml` (Block List, Max: 1) XML file format options. (see [below for nested schema](#nestedblock--file_format--xml))

<a id="nestedblock--file_format--avro"></a>
### Nested Schema for `file_format.avro`

Optional:

- `compression` (String) Specifies the compression format. Valid values: `AUTO` | `GZIP` | `BROTLI` | `ZSTD` | `DEFLATE` | `RAW_DEFLATE` | `NONE`.
- `null_if` (List of String) String used to convert to and from SQL NULL.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `trim_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to remove white space from fields. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--file_format--csv"></a>
### Nested Schema for `file_format.csv`

Optional:

- `binary_format` (String) Defines the encoding format for binary input or output. Valid values: `HEX` | `BASE64` | `UTF8`.
- `compression` (String) Specifies the compression format. Valid values: `AUTO` | `GZIP` | `BZ2` | `BROTLI` | `ZSTD` | `DEFLATE` | `RAW_DEFLATE` | `NONE`.
- `date_format` (String) Defines the format of date values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `empty_field_as_null` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to insert SQL NULL for empty fields in an input file. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `encoding` (String) Specifies the character set of the source data when loading data into a table. Valid values: `BIG5` | `EUCJP` | `EUCKR` | `GB18030` | `IBM420` | `IBM424` | `ISO2022CN` | `ISO2022JP` | `ISO2022KR` | `ISO88591` | `ISO88592` | `ISO88595` | `ISO88596` | `ISO88597` | `ISO88598` | `ISO88599` | `ISO885915` | `KOI8R` | `SHIFTJIS` | `UTF8` | `UTF16` | `UTF16BE` | `UTF16LE` | `UTF32` | `UTF32BE` | `UTF32LE` | `WINDOWS1250` | `WINDOWS1251` | `WINDOWS1252` | `WINDOWS1253` | `WINDOWS1254` | `WINDOWS1255` | `WINDOWS1256`.
- `error_on_column_count_mismatch` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to generate a parsing error if the number of delimited columns in an input file does not match the number of columns in the corresponding table. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `escape` (String) Single character string used as the escape character for field values. Use `NONE` to specify no escape character. NOTE: This value may be not imported properly from Snowflake. Snowflake returns escaped values.
- `escape_unenclosed_field` (String) Single character string used as the escape character for unenclosed field values only. Use `NONE` to specify no escape character. NOTE: This value may be not imported properly from Snowflake. Snowflake returns escaped values.
- `field_delimiter` (String) One or more singlebyte or multibyte characters that separate fields in an input file. Use `NONE` to specify no delimiter.
- `field_optionally_enclosed_by` (String) Character used to enclose strings. Use `NONE` to specify no enclosure character.
- `file_extension` (String) Specifies the extension for files unloaded to a stage.
- `multi_line` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to parse CSV files containing multiple records on a single line. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `null_if` (List of String) String used to convert to and from SQL NULL.
- `parse_header` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to use the first row headers in the data files to determine column names. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `record_delimiter` (String) One or more singlebyte or multibyte characters that separate records in an input file. Use `NONE` to specify no delimiter.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `skip_blank_lines` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies to skip any blank lines encountered in the data files. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `skip_byte_order_mark` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to skip the BOM (byte order mark) if present in a data file. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `skip_header` (Number) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`-1`)) Number of lines at the start of the file to skip.
- `time_format` (String) Defines the format of time values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `timestamp_format` (String) Defines the format of timestamp values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `trim_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to remove white space from fields. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--file_format--json"></a>
### Nested Schema for `file_format.json`

Optional:

- `allow_duplicate` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to allow duplicate object field names (only the last one will be preserved). Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `binary_format` (String) Defines the encoding format for binary input or output. Valid values: `HEX` | `BASE64` | `UTF8`.
- `compression` (String) Specifies the compression format. Valid values: `AUTO` | `GZIP` | `BZ2` | `BROTLI` | `ZSTD` | `DEFLATE` | `RAW_DEFLATE` | `NONE`.
- `date_format` (String) Defines the format of date values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `enable_octal` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that enables parsing of octal numbers. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `file_extension` (String) Specifies the extension for files unloaded to a stage.
- `ignore_utf8_errors` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether UTF-8 encoding errors produce error conditions. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `multi_line` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to allow multiple records on a single line. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `null_if` (List of String) String used to convert to and from SQL NULL.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `skip_byte_order_mark` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to skip the BOM (byte order mark) if present in a data file. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `strip_null_values` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that instructs the JSON parser to remove object fields or array elements containing null values. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `strip_outer_array` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that instructs the JSON parser to remove outer brackets. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `time_format` (String) Defines the format of time values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `timestamp_format` (String) Defines the format of timestamp values in the data files. Use `AUTO` to have Snowflake auto-detect the format.
- `trim_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to remove white space from fields. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--file_format--orc"></a>
### Nested Schema for `file_format.orc`

Optional:

- `null_if` (List of String) String used to convert to and from SQL NULL.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `trim_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to remove white space from fields. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--file_format--parquet"></a>
### Nested Schema for `file_format.parquet`

Optional:

- `binary_as_text` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to interpret columns with no defined logical data type as UTF-8 text. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `compression` (String) Specifies the compression format. Valid values: `AUTO` | `LZO` | `SNAPPY` | `NONE`.
- `null_if` (List of String) String used to convert to and from SQL NULL.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `trim_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to remove white space from fields. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `use_logical_type` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to use Parquet logical types when loading data. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `use_vectorized_scanner` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to use a vectorized scanner for loading Parquet files. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.


<a id="nestedblock--file_format--xml"></a>
### Nested Schema for `file_format.xml`

Optional:

- `compression` (String) Specifies the compression format. Valid values: `AUTO` | `GZIP` | `BZ2` | `BROTLI` | `ZSTD` | `DEFLATE` | `RAW_DEFLATE` | `NONE`.
- `disable_auto_convert` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether the XML parser disables automatic conversion of numeric and Boolean values from text to native representation. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `ignore_utf8_errors` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether UTF-8 encoding errors produce error conditions. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `preserve_space` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether the XML parser preserves leading and trailing spaces in element content. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `replace_invalid_characters` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `skip_byte_order_mark` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether to skip the BOM (byte order mark) if present in a data file. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `strip_outer_element` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Boolean that specifies whether the XML parser strips out the outer XML element, exposing 2nd level elements as separate documents. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.



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
- `last_refreshed_on` (String)


<a id="nestedobjatt--describe_output--file_format"></a>
### Nested Schema for `describe_output.file_format`

Read-Only:

- `avro` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--avro))
- `csv` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--csv))
- `format_name` (String)
- `json` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--json))
- `orc` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--orc))
- `parquet` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--parquet))
- `xml` (List of Object) (see [below for nested schema](#nestedobjatt--describe_output--file_format--xml))

<a id="nestedobjatt--describe_output--file_format--avro"></a>
### Nested Schema for `describe_output.file_format.avro`

Read-Only:

- `compression` (String)
- `null_if` (List of String)
- `replace_invalid_characters` (Boolean)
- `trim_space` (Boolean)
- `type` (String)


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


<a id="nestedobjatt--describe_output--file_format--json"></a>
### Nested Schema for `describe_output.file_format.json`

Read-Only:

- `allow_duplicate` (Boolean)
- `binary_format` (String)
- `compression` (String)
- `date_format` (String)
- `enable_octal` (Boolean)
- `file_extension` (String)
- `ignore_utf8_errors` (Boolean)
- `multi_line` (Boolean)
- `null_if` (List of String)
- `replace_invalid_characters` (Boolean)
- `skip_byte_order_mark` (Boolean)
- `strip_null_values` (Boolean)
- `strip_outer_array` (Boolean)
- `time_format` (String)
- `timestamp_format` (String)
- `trim_space` (Boolean)
- `type` (String)


<a id="nestedobjatt--describe_output--file_format--orc"></a>
### Nested Schema for `describe_output.file_format.orc`

Read-Only:

- `null_if` (List of String)
- `replace_invalid_characters` (Boolean)
- `trim_space` (Boolean)
- `type` (String)


<a id="nestedobjatt--describe_output--file_format--parquet"></a>
### Nested Schema for `describe_output.file_format.parquet`

Read-Only:

- `binary_as_text` (Boolean)
- `compression` (String)
- `null_if` (List of String)
- `replace_invalid_characters` (Boolean)
- `trim_space` (Boolean)
- `type` (String)
- `use_logical_type` (Boolean)
- `use_vectorized_scanner` (Boolean)


<a id="nestedobjatt--describe_output--file_format--xml"></a>
### Nested Schema for `describe_output.file_format.xml`

Read-Only:

- `compression` (String)
- `disable_auto_convert` (Boolean)
- `ignore_utf8_errors` (Boolean)
- `preserve_space` (Boolean)
- `replace_invalid_characters` (Boolean)
- `skip_byte_order_mark` (Boolean)
- `strip_outer_element` (Boolean)
- `type` (String)



<a id="nestedobjatt--describe_output--location"></a>
### Nested Schema for `describe_output.location`

Read-Only:

- `aws_access_point_arn` (String)
- `url` (List of String)


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
terraform import snowflake_stage_external_s3.example '"<database_name>"."<schema_name>"."<stage_name>"'
```
