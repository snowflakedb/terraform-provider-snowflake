---
page_title: "snowflake_user_programmatic_access_token Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage user programmatic access tokens. For more information, check user programmatic access tokens documentation https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token. A programmatic access token is a token that can be used to authenticate to an endpoint. See Using programmatic access tokens for authentication https://docs.snowflake.com/en/user-guide/programmatic-access-tokens user guide for more details.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** External changes to `mins_to_bypass_network_policy_requirement` are not handled by the provider because the value changes continuously on Snowflake side after setting it.

-> **Note** External changes to `days_to_expiry` are not handled by the provider because Snowflake returns `expires_at` which is the token expiration date. Also, the provider does not handle expired tokens automatically. Please change the value of `days_to_expiry` to force a new expiration date.

-> **Note** External changes to `token` are not handled by the provider because the data in this field can be updated only when the token is created or rotated.

-> **Note** Rotating a token can be done by changing the value of `keepers` field. See an example below.

-> **Note** In order to authenticate with PAT with role restriction, you need to grant the role to the user. You can use the [snowflake_grant_account_role](./grant_account_role) resource to do this.

# snowflake_user_programmatic_access_token (Resource)

Resource used to manage user programmatic access tokens. For more information, check [user programmatic access tokens documentation](https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token). A programmatic access token is a token that can be used to authenticate to an endpoint. See [Using programmatic access tokens for authentication](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens) user guide for more details.

## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

```terraform
# basic resource
resource "snowflake_user_programmatic_access_token" "basic" {
  user = "USER"
  name = "TOKEN"
}

# complete resource
resource "snowflake_user_programmatic_access_token" "complete" {
  user                                      = "USER"
  name                                      = "TOKEN"
  role_restriction                          = "ROLE"
  days_to_expiry                            = 30
  mins_to_bypass_network_policy_requirement = 10
  disabled                                  = false
  comment                                   = "COMMENT"
}

# Set up dependencies and reference them from the token resource.
resource "snowflake_account_role" "role" {
  name = "ROLE"
}

resource "snowflake_user" "user" {
  name = "USER"
}

# Grant the role to the user. This is required to authenticate with PAT with role restriction.
resource "snowflake_grant_account_role" "grant_role_to_user" {
  role_name = snowflake_account_role.role.name
  user_name = snowflake_user.user.name
}

# complete resource with external references
resource "snowflake_user_programmatic_access_token" "complete_with_external_references" {
  user                                      = snowflake_user.user.name
  name                                      = "TOKEN"
  role_restriction                          = snowflake_account_role.role.name
  days_to_expiry                            = 30
  mins_to_bypass_network_policy_requirement = 10
  disabled                                  = false
  comment                                   = "COMMENT"
}

# use the token returned from Snowflake and remember to mark it as sensitive
output "token" {
  value     = snowflake_user_programmatic_access_token.complete.token
  sensitive = true
}

# rotate the token regularly using the keepers field and time_rotating resource
resource "snowflake_user_programmatic_access_token" "rotating" {
  user = "USER"
  name = "TOKEN"
  keepers = {
    rotation_schedule = time_rotating.rotation_schedule.rotation_rfc3339
  }
}

# Note that the fields of this resource are updated only when Terraform is run.
# This means that the schedule may not be respected if Terraform is not run regularly.
resource "time_rotating" "rotation_schedule" {
  rotation_days = 30
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specifies the name for the programmatic access token; must be unique for the user. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `user` (String) The name of the user that the token is associated with. A user cannot use another user's programmatic access token to authenticate. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `comment` (String) Descriptive comment about the programmatic access token.
- `days_to_expiry` (Number) The number of days that the programmatic access token can be used for authentication. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `disabled` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Disables or enables the programmatic access token. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `expire_rotated_token_after_hours` (Number) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`-1`)) Sets the expiration time of the existing token secret to expire after the specified number of hours. You can set this to a value of 0 to expire the current token secret immediately. This field is only used when the token is rotated with `keepers` field.
- `keepers` (Map of String) Arbitrary map of values that, when changed, will trigger a new key to be generated.
- `mins_to_bypass_network_policy_requirement` (Number) The number of minutes during which a user can use this token to access Snowflake without being subject to an active network policy. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `role_restriction` (String) The name of the role used for privilege evaluation and object creation. This must be one of the roles that has already been granted to the user. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `rotated_token_name` (String) Name of the token that represents the prior secret. This field is updated only when the token is rotated.
- `show_output` (List of Object) Outputs the result of `SHOW USER PROGRAMMATIC ACCESS TOKENS` for the given user programmatic access token. (see [below for nested schema](#nestedatt--show_output))
- `token` (String, Sensitive) The token itself. Use this to authenticate to an endpoint. The data in this field is updated only when the token is created or rotated.

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
- `created_by` (String)
- `created_on` (String)
- `expires_at` (String)
- `mins_to_bypass_network_policy_requirement` (Number)
- `name` (String)
- `role_restriction` (String)
- `rotated_to` (String)
- `status` (String)
- `user_name` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_user_programmatic_access_token.example '"<user_name>"|"<token_name>"'
```
