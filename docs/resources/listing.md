---
page_title: "snowflake_listing Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage listing objects. For more information, check listing documentation https://other-docs.snowflake.com/en/collaboration/collaboration-listings-about.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

!> **Warning** Versioning only works if your listing ever sourced the manifest from stage. This is a Snowflake limitation.

!> **Warning** External changes to the manifest (inlined and staged) won't be detected by the provider automatically. You need to manually trigger updates when manifest content changes.

!> **Warning** This resource isn't suitable for public listings because its review process doesn't align with Terraform's standard method for managing infrastructure resources. The challenge is that the review process often takes time and might need several manual revisions. We need to reconsider how to integrate this process into a resource. Although we plan to support this in the future, it might be added later. Currently, the resource may not function well with public listings because review requests are closely connected to the publish field.

!> **Warning** To use external resources in your manifest (e.g., company logo) you must be sourcing your manifest from a stage. Any references to external resources are relative to the manifest location in the stage.

-> **Note** When using manifest from stage, the change in either stage id, location, or version will create a new listing version that can be seen by calling the [SHOW VERSIONS IN LISTING](https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing) command.

-> **Note** For inlined manifest version, only string is accepted. The manifest structure is not mapped to the resource schema to keep it simple and aligned with other resources that accept similar metadata (e.g., service templates). While it's more recommended to keep your manifest in a stage, the inlined version may be useful for initial setup and testing.

-> **Note** For manifest reference visit [Snowflake's listing manifest reference documentation](https://docs.snowflake.com/en/progaccess/listing-manifest-reference).

# snowflake_listing (Resource)

Resource used to manage listing objects. For more information, check [listing documentation](https://other-docs.snowflake.com/en/collaboration/collaboration-listings-about).

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

```terraform
# basic resource with inlined manifest
resource "snowflake_listing" "basic_inlined" {
  name = "LISTING"
  manifest {
    from_string = <<-EOT
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
EOT
  }
}

# basic resource with manifest in a stage
resource "snowflake_listing" "basic_staged" {
  name = "LISTING"
  manifest {
    from_stage = {
      stage = snowflake_stage.test_stage.fully_qualified_name
    }
  }
}

# complete resource with inlined manifest
resource "snowflake_listing" "basic_inlined" {
  name = "LISTING"
  manifest {
    from_string = <<-EOT
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
EOT
  }

  share = snowflake_share.test_share.fully_qualified_name
  # or
  application_package = "test_application_package"

  publish = true
  comment = "This is a comment for the listing"
}

# complete resource with manifest in a stage
resource "snowflake_listing" "basic_staged" {
  name = "LISTING"
  manifest {
    from_stage = {
      stage           = snowflake_stage.test_stage.fully_qualified_name
      location        = "path/to/manifest"
      version_name    = "v1.0.0"
      version_comment = "Initial version of the manifest"
    }
  }

  share = snowflake_share.test_share.fully_qualified_name
  # or
  application_package = "test_application_package"

  publish = true
  comment = "This is a comment for the listing"
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `manifest` (Block List, Min: 1, Max: 1) Specifies the way manifest is provided for the listing. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint". (see [below for nested schema](#nestedblock--manifest))
- `name` (String) Specifies the listing identifier (name). It must be unique within the organization, regardless of which Snowflake region the account is located in. Must start with an alphabetic character and cannot contain spaces or special characters except for underscores.

### Optional

- `application_package` (String) Specifies the application package attached to the listing.
- `comment` (String) Specifies a comment for the listing.
- `publish` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Determines if the listing should be published.
- `share` (String) Specifies the identifier for the share to attach to the listing.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.
- `show_output` (List of Object) Outputs the result of `SHOW LISTINGS` for the given listing. (see [below for nested schema](#nestedatt--show_output))

<a id="nestedblock--manifest"></a>
### Nested Schema for `manifest`

Optional:

- `from_stage` (Block List, Max: 1) Manifest provided from a given stage. If the manifest file is in the root, only stage needs to be passed. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). A proper YAML indentation (2 spaces) is required. (see [below for nested schema](#nestedblock--manifest--from_stage))
- `from_string` (String) Manifest provided as a string. Wrapping `$$` signs are added by the provider automatically; do not include them. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). Also, the [multiline string syntax](https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings) is a must here. A proper YAML indentation (2 spaces) is required.

<a id="nestedblock--manifest--from_stage"></a>
### Nested Schema for `manifest.from_stage`

Required:

- `stage` (String) Identifier of the stage where the manifest file is located.

Optional:

- `location` (String) Location of the manifest file in the stage. If not specified, the manifest file will be expected to be at the root of the stage.
- `version_comment` (String) Specifies a comment for the listing version. Whenever a new version is created, this comment will be associated with it. The comment on the version will be visible in the [SHOW VERSIONS IN LISTING](https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing) command output.
- `version_name` (String) Represents manifest version name. It's case-sensitive and used in manifest versioning. Version name should be specified or changed whenever any changes in the manifest should be applied to the listing. Later on the versions of the listing can be analyzed by calling the [SHOW VERSIONS IN LISTING](https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing) command. The resource does not track the changes on the specified stage.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--show_output"></a>
### Nested Schema for `show_output`

Read-Only:

- `comment` (String)
- `created_on` (String)
- `detailed_target_accounts` (String)
- `distribution` (String)
- `global_name` (String)
- `is_application` (Boolean)
- `is_by_request` (Boolean)
- `is_limited_trial` (Boolean)
- `is_monetized` (Boolean)
- `is_mountless_queryable` (Boolean)
- `is_targeted` (Boolean)
- `name` (String)
- `organization_profile_name` (String)
- `owner` (String)
- `owner_role_type` (String)
- `profile` (String)
- `published_on` (String)
- `regions` (String)
- `rejected_on` (String)
- `review_state` (String)
- `state` (String)
- `subtitle` (String)
- `target_accounts` (String)
- `title` (String)
- `uniform_listing_locator` (String)
- `updated_on` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_listing.example '"<listing_name>"'
```
