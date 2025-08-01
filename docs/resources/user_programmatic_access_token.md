---
page_title: "snowflake_user_programmatic_access_token Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage user programmatic access tokens. For more information, check user programmatic access tokens documentation https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token. A programmatic access token is a token that can be used to authenticate to an endpoint. See Using programmatic access tokens for authentication https://docs.snowflake.com/en/user-guide/programmatic-access-tokens user guide for more details.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Read more about PAT support in the provider in our [Authentication Methods guide](../guides/authentication_methods#managing-pats).

-> **Note** External changes to `mins_to_bypass_network_policy_requirement` are not handled by the provider because the value changes continuously on Snowflake side after setting it.

-> **Note** External changes to `days_to_expiry` are not handled by the provider because Snowflake returns `expires_at` which is the token expiration date. Also, the provider does not handle expired tokens automatically. Please change the value of `days_to_expiry` to force a new expiration date.

-> **Note** External changes to `token` are not handled by the provider because the data in this field can be updated only when the token is created or rotated.

-> **Note** Rotating a token can be done by changing the value of `keeper` field. See an example below.

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

# Use the token returned from Snowflake and remember to mark it as sensitive.
output "token" {
  value     = snowflake_user_programmatic_access_token.complete.token
  sensitive = true
}

# Token Rotation

# Rotate the token regularly using the keeper field and time_rotating resource.
resource "snowflake_user_programmatic_access_token" "rotating" {
  user = "USER"
  name = "TOKEN"

  # Use the keeper field to force token rotation. If, and only if, the value changes
  # from a non-empty to a different non-empty value (or is known after apply), the token will be rotated.
  # When you add this key or remove this key from the config, the token will not be rotated.
  # When the token is rotated, the `token` and `rotated_token_name` fields are marked as computed.
  keeper = time_rotating.rotation_schedule.rotation_rfc3339
}

# Note that the fields of this resource are updated only when Terraform is run.
# This means that the schedule may not be respected if Terraform is not run regularly.
resource "time_rotating" "rotation_schedule" {
  rotation_days = 30
}
```

## Token rotation clarifications
The token value returned from Snowflake is stored in the Terraform state and can stay constant until the `keeper` field is changed from a non-empty to a different non-empty value (or is known after apply).
Then, the token is rotated and the new `token` value is stored in the state. In this case, the `token` and `rotated_token_name` fields are marked as computed.
You can use the `keeper` argument in this resource to store an arbitrary string. Fill it with a value that should persist unless you want a new token.
The key gets is not rotated when the `keeper` field is added to or removed from the configuration.
Keep in mind that `keeper` isn't treated as sensitive data, so any values you use for this field will appear as plain text in the Terraform outputs.

In the example above, when you first add a `keeper` field, you will see a plan output similar to:
```
time_rotating.rotation_schedule: Refreshing state... [id=2025-07-22T07:59:20Z]
snowflake_user_programmatic_access_token.complete_with_external_references: Refreshing state... [id="PAT"|"TOKEN"]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # snowflake_user_programmatic_access_token.complete_with_external_references will be updated in-place
  ~ resource "snowflake_user_programmatic_access_token" "complete_with_external_references" {
        id                                        = "\"PAT\"|\"TOKEN\""
      + keeper                                    = "2025-07-22T07:59:20Z"
        name                                      = "TOKEN"
        # (9 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```
This plan means that the value of `keeper` is saved in the state, but the token will not be rotated yet.
After 30 days pass, the `rotation_schedule` resource will return a new timestamp. So, the value of the `keeper` field is changed automatically.
You will see a plan output similar to:
```
time_rotating.rotation_schedule: Refreshing state... [id=2025-07-22T08:06:51Z]
snowflake_user_programmatic_access_token.complete_with_external_references: Refreshing state... [id="PAT"|"TOKEN"]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # time_rotating.rotation_schedule has been deleted
  - resource "time_rotating" "rotation_schedule" {
        id               = "2025-07-22T08:06:51Z"
      - rfc3339          = "2025-07-22T08:06:51Z" -> null
        # (9 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to
these changes.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
  ~ update in-place

Terraform will perform the following actions:

  # snowflake_user_programmatic_access_token.complete_with_external_references will be updated in-place
  ~ resource "snowflake_user_programmatic_access_token" "complete_with_external_references" {
        id                                        = "\"PAT\"|\"TOKEN\""
      ~ keeper                                    = "2025-07-22T08:06:51Z" -> (known after apply)
        name                                      = "TOKEN"
      + rotated_token_name                        = (known after apply)
      ~ token                                     = (sensitive value)
        # (8 unchanged attributes hidden)
    }

  # time_rotating.rotation_schedule will be created
  + resource "time_rotating" "rotation_schedule" {
      + day              = 30
      + hour             = (known after apply)
      + id               = (known after apply)
      + minute           = (known after apply)
      + month            = (known after apply)
      + rfc3339          = (known after apply)
      + rotation_minutes = (known after apply)
      + rotation_rfc3339 = (known after apply)
      + second           = (known after apply)
      + unix             = (known after apply)
      + year             = (known after apply)
    }

Plan: 1 to add, 1 to change, 0 to destroy.
```
After you apply this plan, the token will be rotated and the new `token` value will be stored in the state each time the value of `time_rotating.rotation_schedule.rotation_rfc3339` is changed.
When you remove the `keeper` field from the configuration, the token will not be rotated, see the plan output below:
```
time_rotating.rotation_schedule: Refreshing state... [id=2025-07-22T08:06:51Z]
snowflake_user_programmatic_access_token.complete_with_external_references: Refreshing state... [id="PAT"|"TOKEN"]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
  ~ update in-place

Terraform will perform the following actions:

  # snowflake_user_programmatic_access_token.complete_with_external_references will be updated in-place
  ~ resource "snowflake_user_programmatic_access_token" "complete_with_external_references" {
        id                                        = "\"PAT\"|\"TOKEN\""
      - keeper                                    = "2025-07-22T08:06:51Z" -> null
        name                                      = "TOKEN"
        # (9 unchanged attributes hidden)
    }

  # time_rotating.rotation_schedule will be created
  + resource "time_rotating" "rotation_schedule" {
      + day              = 30
      + hour             = (known after apply)
      + id               = (known after apply)
      + minute           = (known after apply)
      + month            = (known after apply)
      + rfc3339          = (known after apply)
      + rotation_minutes = (known after apply)
      + rotation_rfc3339 = (known after apply)
      + second           = (known after apply)
      + unix             = (known after apply)
      + year             = (known after apply)
    }

Plan: 1 to add, 1 to change, 0 to destroy.
```

This way you can cancel the rotation schedule without rotating the token (the `token` and `rotated_token_name` fields are not marked as computed).

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specifies the name for the programmatic access token; must be unique for the user. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `user` (String) The name of the user that the token is associated with. A user cannot use another user's programmatic access token to authenticate. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.

### Optional

- `comment` (String) Descriptive comment about the programmatic access token.
- `days_to_expiry` (Number) The number of days that the programmatic access token can be used for authentication. This field cannot be altered after the token is created. Instead, you must rotate the token with the `keeper` field. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `disabled` (String) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`default`)) Disables or enables the programmatic access token. Available options are: "true" or "false". When the value is not set in the configuration the provider will put "default" there which means to use the Snowflake default for this value.
- `expire_rotated_token_after_hours` (Number) (Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`-1`)) This field is only used when the token is rotated by changing the `keeper` field. Sets the expiration time of the existing token secret to expire after the specified number of hours. You can set this to a value of 0 to expire the current token secret immediately.
- `keeper` (String) Arbitrary string that, if and only if, changed from a non-empty to a different non-empty value (or known after apply), will trigger a key to be rotated. When you add this field to the configuration, or remove it from the configuration, the rotation is not triggered. When the token is rotated, the `token` and `rotated_token_name` fields are marked as computed.
- `mins_to_bypass_network_policy_requirement` (Number) The number of minutes during which a user can use this token to access Snowflake without being subject to an active network policy. External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".
- `role_restriction` (String) The name of the role used for privilege evaluation and object creation. This must be one of the roles that has already been granted to the user. Due to technical limitations (read more [here](../guides/identifiers_rework_design_decisions#known-limitations-and-identifier-recommendations)), avoid using the following characters: `|`, `.`, `"`.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `rotated_token_name` (String) Name of the token that represents the prior secret. This field is updated only when the token is rotated. In this case, the field is marked as computed.
- `show_output` (List of Object) Outputs the result of `SHOW USER PROGRAMMATIC ACCESS TOKENS` for the given user programmatic access token. (see [below for nested schema](#nestedatt--show_output))
- `token` (String, Sensitive) The token itself. Use this to authenticate to an endpoint. The data in this field is updated only when the token is created or rotated. In this case, the field is marked as computed.

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
