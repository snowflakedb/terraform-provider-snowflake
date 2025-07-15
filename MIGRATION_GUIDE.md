# Migration guide

This document is meant to help you migrate your Terraform config to the new newest version. In migration guides, we will only
describe deprecations or breaking changes and help you to change your configuration to keep the same (or similar) behavior
across different versions.

To keep your configuration up to date, we also recommend reading the [Snowflake BCR migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md)
for changes required after enabling given [Snowflake BCR Bundle](https://docs.snowflake.com/en/release-notes/behavior-changes).

> [!TIP]
> We highly recommend upgrading the versions one by one instead of bulk upgrades.
>
> To migrate particular resources, follow our [Resource Migration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration) guide for more details.
>
> In certain cases (like using the ancient provider versions), you can upgrade multiple versions at once. To do that:
> - read the target version documentation and all the intermediary migration guide entries;
> - focus on changes to authentication to make sure your provider is set up correctly in the newest version;
> - check changes to resource schemas; if in doubt, you can always simplify the resource and let the terraform figure out the changes (you can use plan output to make the configuration appropriate);
> - reimport your infrastructure using the target provider version, preferably in smaller chunks (or experiment with 1-2 resources of each type first).

> [!TIP]
> If you're still using the `Snowflake-Labs/snowflake` source, see [Upgrading from Snowflake-Labs Provider](./SNOWFLAKEDB_MIGRATION.md) to upgrade to the snowflakedb namespace.

## v2.2.0 ➞ v2.3.0

### *(new feature)* New `PROGRAMMATIC_ACCESS_TOKEN` authenticator option

We added a new `PROGRAMMATIC_ACCESS_TOKEN` option to the `authenticator` field in the provider. This feature enables authentication with `PROGRAMMATIC_ACCESS_TOKEN` authenticator in the Go driver. Read more in our [Authentication methods](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods) guide.

See [Snowflake official documentation](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens) for more information on PAT authentication.

### *(bugfix)* Fix `snowflake_functions` and `snowflake_procedures` data sources with 2025_03 Bundle enabled

> [!IMPORTANT]
> This behavior change in Snowflake was originally in the [2025_03 Bundle](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03_bundle) and intended to become enabled by default in the 2025_04 bundle. However, it has been [postponed](https://docs.snowflake.com/release-notes/bcr-bundles/un-bundled/bcr-1944) and a new release date has not been determined.
> After adjusting parsing data types in the provider, it handles arguments with and without attributes.

Check for more details and action steps needed in [Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands](./SNOWFLAKE_BCR_MIGRATION_GUIDE.md#argument-output-changes-for-show-functions-and-show-procedures-commands).
This fix was also backported to version v1.2.3.

References: [#3822](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3822)

### *(bugfix)* Fix all function and procedure resources with 2025_03 Bundle enabled

> [!IMPORTANT]
> This behavior change in Snowflake was originally in the [2025_03 Bundle](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03_bundle) and intended to become enabled by default in the 2025_04 bundle. However, it has been [postponed](https://docs.snowflake.com/release-notes/bcr-bundles/un-bundled/bcr-1944) and a new release date has not been determined.
> After adjusting parsing data types in the provider, it handles arguments with and without attributes.

Check for more details and action steps needed in [Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands](./SNOWFLAKE_BCR_MIGRATION_GUIDE.md#argument-output-changes-for-show-functions-and-show-procedures-commands).
This fix was also backported to version v1.2.3.

References: [#3823](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3823)

## v2.1.x ➞ v2.2.0
<a id="v210--v220"></a>

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

This fix was also backported to versions v1.0.6, v1.1.1, v1.2.2, v2.0.1, and v2.1.1.

### *(bugfix)* Fix grant_ownership resource for serverless tasks

Previously, it wasn't possible to use the `snowflake_grant_ownership` resource to grant ownership of serverless tasks.
In this version, we fixed the issue, and now you can use the resource to grant ownership of serverless tasks.

No configuration changes are needed.

References: [#3750](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3750)

### *(bugfix)* Fix external volume creation error handling

Errors in [`snowflake_external_volume`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/external_volume) resource creation were not handled and propagated properly, resulting in provider errors similar to:
```
Warning: Failed to query external volume. Marking the resource as removed.
│
│   with snowflake_external_volume.s3_volume,
│   on main.tf line 62, in resource "snowflake_external_volume" "s3_volume":
│   62: resource "snowflake_external_volume" "s3_volume" {
│
│ External Volume: "MY_S3_EXTERNAL_VOLUME", Err: object does not exist
╵
╷
│ Error: Provider produced inconsistent result after apply
│
│ When applying changes to snowflake_external_volume.s3_volume, provider
│ "provider[\"registry.terraform.io/snowflakedb/snowflake\"]" produced an unexpected new value: Root object
│ was present, but now absent.
│
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
```

Starting with this version, creation errors in `snowflake_external_volume` will be handled and propagated properly to the user.

No configuration changes are needed.

### *(new feature)* New fields in snowflake_cortex_search_service resource

We added a new `embedding_model` field to the `snowflake_cortex_search_service`. This field specifies the embedding model to use in the Cortex Search Service.
We updated the examples of using the resource with this field.
Additionally, we added a new `describe_output` field to handle this field properly (read more in our [design considerations](v1-preparations/CHANGES_BEFORE_V1.md#default-values)).

### *(new feature)* New consumption_billing_entity field in snowflake_account resource

The `snowflake_account` resource now has a new `consumption_billing_entity` field, which allows you to set the consumption billing entity for the account.
It is useful in case you have multiple billing entities in your account and want to set a specific one for the account.
You can find more details in [this](https://community.snowflake.com/s/article/ERROR-Multiple-suitable-billing-entities-exist-for-the-target-cloud) KB article.

No configuration changes are needed.

### The ORGADMIN checks removed from snowflake_account resource

Previously, the `snowflake_account` resource required the ORGADMIN role for operations to be executed.
In recent Snowflake changes that introduced organization accounts, more roles can now manage the account.
Because of that, to enable the resource to be used in more scenarios, we removed the ORGADMIN checks from the resource.

### *(new feature)* New tracking level

Every resource that is capable of setting tracing level (`database`, `shared_database`, `secondary_database`, `schema`) now supports the new `PROPAGATE` value.

### *(bugfix)* Fix how snowflake_user_authentication_policy_attachment resource handles missing objects it depends on

Previously, the `snowflake_user_authentication_policy_attachment` resource was not able to handle missing objects it depends on.
This means, if a user or authentication policy was removed manually outside Terraform, the provider would produce plans with errors like:
```
User 'XYZ' does not exist or not authorized
```
and only manual state management would help you to remove the resource from the state.

Now, the removal of the resource is handled properly and the resource is removed from the state automatically with the following warning:
```
Failed to find user authentication policy. Marking the resource as removed.
### or ###
Failed to get user policies. Marking the resource as removed.
```

If you are encountering this issue,
either bump the provider version to at least `v2.2.0` or [remove the resource from the state manually](https://developer.hashicorp.com/terraform/cli/commands/state/rm).

No configuration changes are needed.

References: [#3672](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3672)

### *(new feature)* snowflake_current_account resource
Added a new preview resource for managing an account that the provider is currently connected to. It's capable of managing attached parameters, resource_monitors, and more. See reference docs for [ALTERING ACCOUNT](https://docs.snowflake.com/en/sql-reference/sql/alter-account). You can read about the resources' limitations in the documentation in the registry.
This resource is intended to replace the `snowflake_account_parameter` resource, which will be deprecated in the future,
but some of the supported parameters in `snowflake_account_parameter` aren't supported in `snowflake_current_account`. Those parameters are:
- ENABLE_CONSOLE_OUTPUT
- ENABLE_PERSONAL_DATABASE
- PREVENT_LOAD_FROM_INLINE_URL

They are not supported, because they are not in the [official parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters).
Once they are publicly documented, they will be added to the `snowflake_current_account_resource` resource.

The `snowflake_current_account_resource` resource shouldn't be used with `snowflake_object_parameter` (with `on_account` field set) and `snowflake_account_parameter` resources in the same configuration, as it may lead to unexpected behavior. Unless they're used to manage the above parameters that are not supported.
The resource shouldn't be also used with `snowflake_account_password_policy_attachment`, `snowflake_network_policy_attachment`, `snowflake_account_authentication_policy_attachment` resources in the same configuration to manage policies on the current account, as it may lead to unexpected behavior.

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_current_account_resource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_service and snowflake_job_service resources
Added new preview resources for managing services and job services. See reference docs for [services](https://docs.snowflake.com/en/sql-reference/sql/create-service) and [job services](https://docs.snowflake.com/en/sql-reference/sql/execute-job-service). You can read about the resources' limitations in the documentation in the registry.

These features will be marked as stable in future releases. Breaking changes are expected, even without bumping the major version. To use these features, add `snowflake_service_resource` or `snowflake_job_service_resource` to `preview_features_enabled` field in the provider configuration, respectively.

### *(new feature)* snowflake_git_repository resource
Added a new preview resource for managing git repositories. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-git-repository). Note that `snowflake_api_integration` currently does not support `git_https_api` type. It will be added during the resource rework. Instead, you can use [execute](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/execute) resource.

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_git_repository_resource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_compute_pool resource
Added a new preview resource for managing compute pools. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool). The limitation of this resource is that identifiers with special or lower-case characters are not supported. This limitation in the provider follows the limitation in Snowflake (see the linked docs).

Managing compute pool state is limited. It is handled only by `initially_suspended`, `auto_suspend_secs`, and `auto_resume` fields. See the resource documentation for more details.

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_compute_pool_resource` to `preview_features_enabled` field in the provider configuration.

### *(bugfix)* Fix the behavior for empty privileges list in snowflake_grant_privileges_to_account_role, snowflake_grant_privileges_to_database_role, and snowflake_grant_privileges_to_share resources

Previously, it was possible to create `snowflake_grant_privileges_to_X` resources with an empty privilege list, which led to the following error:
```
│ Error: Failed to parse internal identifier
│ Error: [grant_privileges_to_database_role_identifier.go:79] invalid Privileges value: , should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges
```
After that, an identifier stored in state would be corrupted and only manual state manipulation would fix it.
We added validation to prevent this from happening. Now, if you try to create or update a resource with an empty privilege list, you will get the following error:

```
| Error: Not enough list items
|
|   with snowflake_grant_privileges_to_database_role.test,
|   on test.tf line 3, in resource "snowflake_grant_privileges_to_database_role" "test":
|    3:   privileges         = []
|
| Attribute privileges requires 1 item minimum, but config has only 0 declared.
```

and the validation error will prevent the state file from changing, which means you will be able to normally adjust the resource and reapply the configuration.

If you are experiencing this at the moment, you can fix it by running removing `snowflake_grant_privileges_to_database_role` from the state by running:
```shell
terraform state rm snowflake_grant_privileges_to_database_role.test # Replace `test` with the actual resource name
```
and apply it with the correct `privileges` list. If you don't want to apply the privileges again, make sure they are
revoked in Snowflake by running the corresponding [SHOW GRANTS](https://docs.snowflake.com/en/sql-reference/sql/show-grants) command
and then corresponding [REVOKE <privileges>](https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege) to remove unwanted privileges.

Other than that, no changes to the configurations are necessary.

Reference: [#3690](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3690).

### *(new feature)* snowflake_image_repository resource
Added a new preview resource for managing image repositories. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-image-repository). The limitation of this resource is that quoted names for special characters or case-sensitive names are not supported. Please use only characters compatible with [unquoted identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax#label-unquoted-identifier). The same constraint also applies to database and schema names where you create an image repository. This limitation in the provider follows the limitation in Snowflake (see the linked docs).

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_image_repository_resource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_git_repositories data source
Added a new preview data source for git repositories. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories).

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_git_repositories_datasource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_compute_pools data source
Added a new preview data source for compute pools. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools).

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_compute_pools_datasource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_image_repositories data source
Added a new preview data source for image repositories. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/show-image-repositories).

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_image_repositories_datasource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* snowflake_services data source
Added a new preview data source for services. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/show-services).

This feature will be marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add `snowflake_services_datasource` to `preview_features_enabled` field in the provider configuration.

### *(new feature)* Managing tags for image repositories, compute pools, services, and git repositories
The [snowflake_tag_association](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/tag_association) can now be used for managing tags in [image repositories](https://docs.snowflake.com/en/sql-reference/sql/create-image-repository), [compute pools](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool), [services](https://docs.snowflake.com/en/sql-reference/sql/create-service) and [git repositories](https://docs.snowflake.com/en/sql-reference/sql/create-git-repository).

### *(bugfix)* Fixed handling users' grants

In v2.1.0, we introduced a fix in handling users' grants ([migration guide](#bugfix-fixed-snowflake_grant_database_role-resource)), which addressed changes in the `2025_02` bundle. The username was parsed incorrectly if it had a prefix formed of `U`, `S`, `E`, and `R` characters. The username returned from `SHOW GRANTS` was incorrect in this case. Now, such names should be handled correctly.
No configuration changes are necessary.

### *(new feature)* Granting privileges on future cortex search services

As this is now available on Snowflake, we allow to grant privileges on future cortex search services both in `snowflake_grant_privileges_on_account_role` and `snowflake_grant_privileges_on_database_role`.

## v2.1.0 ➞ v2.1.1
<a id="v210---v211"></a>

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

This fix was also backported to versions v1.0.6, v1.1.1, v1.2.2, and v2.0.1.

## v2.0.x ➞ v2.1.0
<a id="v200--v210"></a>

### *(bugfix)* Fixed `snowflake_tag_association` resource

The `snowflake_tag_association` resource was crashing when performing the update operation (e.g., because the `tag_value` was changed)
for objects that are created on schema level. This was fixed, and now you can create tag associations for objects that are created on schema level.

Reference: [#3622](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3622).

### *(bugfix)* Added missing `DISABLE_USER_PRIVILEGE_GRANTS` account parameter

As part of the [2025_02 Bundle](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_02_bundle), support for User Based Access Control (UBAC) will be added ([BCR-1924](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_02/bcr-1924)).
It can be disabled by setting the `DISABLE_USER_PRIVILEGE_GRANTS` parameter to `true`.
This version adds the support for this parameter in the `snowflake_account_parameter` resource.

Reference: [#3639](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3639).

### *(bugfix)* Imports propagation to Snowflake for all `snowflake_procedure_*` resources

The `snowflake_procedure_python`, `snowflake_procedure_scala`, and `snowflake_procedure_java` resources were not propagating changes to `imports` set to Snowflake.
There is no `ALTER` to update the imports post-creation, so changes require dropping and recreating the given procedure.

No action is needed.

Reference: [#3401](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3401).

### *(bugfix)* Fixed permadiff issue with the `network_policy` attribute in all user resources

Using `snowflake_network_policy.my_policy.fully_qualified_name` directly as `network_policy` input for `snowflake_user`, `snowflake_service_user`, and `snowflake_legacy_service_user` resources could result in a permadiff (like ` ~ network_policy = "NETWORK_POLICY_ID" -> "\"NETWORK_POLICY_ID\""`).
This version adds appropriate validation and diff suppression to `network_policy` attribute, so such permadiffs are avoided.

No action is needed.

Reference: [#3655](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3655).

### *(bugfix)* Fixed snowflake_grant_database_role resource

The `2025_02` Snowflake BCR enables granting database roles directly to users.
This caused issues in the provider, leading to `Provider produced inconsistent result after apply` errors
when a database role was granted to a user. This version resolves the issue.
No configuration changes are necessary.

References: [#3629](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3629)

## v2.0.0 ➞ v2.0.1
<a id="v200---v201"></a>

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

This fix was also backported to versions v1.0.6, v1.1.1, and v1.2.2.

## v1.2.x ➞ v2.0.0
<a id="v121--v200"></a>

### Supported architectures

We have compiled a list to clarify which binaries are officially supported and which are provided additionally but not officially supported.
The lists are based on what the underlying [gosnowflake driver](https://github.com/snowflakedb/gosnowflake) supports and what [HashiCorp recommends for Terraform providers](https://developer.hashicorp.com/terraform/registry/providers/os-arch).

The provider officially supports the binaries built for the following OSes and architectures:
- Windows: amd64
- Linux: amd64 and arm64
- Darwin: amd64 and arm64

Currently, we also provide the binaries for the following OSes and architectures, but they are not officially supported, and we do not prioritize fixes for them:
- Windows: arm64 and 386
- Linux: 386
- Darwin: 386
- Freebsd: any architecture

### *(breaking change)* Changes in sensitive values
To ensure better security of users' data, we adjusted the fields containing sensitive information to be sensitive in the provider.
Some fields had to be removed due to Terraform SDK limitations (more on that in the [removal of sensitive fields](#removal-of-sensitive-fields) section).
This means these values will not be printed by Terraform during planning, etc. Note that the users are still responsible for storing the state securely.
Read more about sensitive values in the [Terraform documentation](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

Fields changed to sensitive:
- provider configuration: `passcode` field
- `snowflake_system_generate_scim_access_token` data source: `access_token` field
- `snowflake_api_authentication_integration_with_authorization_code_grant` resource: `oauth_client_id` and `oauth_client_secret` fields,
- `snowflake_api_authentication_integration_with_client_credentials` resource: `oauth_client_id` and `oauth_client_secret` fields,
- `snowflake_api_authentication_integration_with_jwt_bearer` resource: `oauth_client_id` and `oauth_client_secret` fields,
- `snowflake_saml2_integration` resource: `saml2_x509_cert` field
- `snowflake_storage_integration` resource: `azure_consent_url` field

If you reference one of these fields in an output or a variable block, then it needs to be marked as `sensitive = true` in the Terraform configuration. Read [Output documentation](https://developer.hashicorp.com/terraform/language/values/outputs#sensitive-suppressing-values-in-cli-output) and [Variable documentation](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output) for more details. In other case, you will get an error like this:
```
Planning failed. Terraform encountered an error while generating this plan.

╷
│ Error: Output refers to sensitive values
│
│   on 3565.tf line 84:
│   84: output "sensitive_output" {
│
│ To reduce the risk of accidentally exporting sensitive data that was intended to be only internal, Terraform requires that any root module output containing sensitive data be explicitly marked
│ as sensitive, to confirm your intent.
│
│ If you do intend to export this data, annotate the output value as sensitive by adding the following argument:
│     sensitive = true
╵
```

Some fields, like secure function definitions, can also contain sensitive values. However, because of [SDK v2](https://developer.hashicorp.com/terraform/plugin/sdkv2) limitations:
- There is no possibility to mark sensitive values conditionally ([reference](https://github.com/hashicorp/terraform-plugin-sdk/issues/736)). This means it is not possible to mark sensitive values based on other fields, like marking `body` based on the value of `secure` field in views, functions, and procedures. As a result, this field is not marked as sensitive. For such cases, we add disclaimers in the resource documentation.
- There is no possibility to mark sensitive values in nested fields ([reference](https://github.com/hashicorp/terraform-plugin-sdk/issues/201)). This means the nested fields, like these in `show_output` and `describe_output` cannot be sensitive. fields, like in `show_output` and `describe_output`, cannot be marked as sensitive.

Instead, we added notes in the documentation of the related resources. The full list includes:
- `snowflake_execute` resource: `execute`, `revert`, `query` and `query_results` fields,
- `snowflake_external_function` resource: `context_headers` and `header` fields,
- `snowflake_function_java` resource: `function_definition` and `show_output.arguments_raw` fields,
- `snowflake_function_javascript` resource: `function_definition` and `show_output.arguments_raw` fields,
- `snowflake_function_python` resource: `function_definition` and `show_output.arguments_raw` fields,
- `snowflake_function_scala` resource: `function_definition` and `show_output.arguments_raw` fields,
- `snowflake_function_sql` resource: `function_definition` and `show_output.arguments_raw` fields,
- `snowflake_legacy_service_user` resource: `display_name`, `show_output.display_name`, `show_output.email`, `show_output.login_name`, `show_output.first_name` and `show_output.last_name` fields,
- `snowflake_masking_policy` resource: `body` and `describe_output.body` fields,
- `snowflake_masking_policies` data source: `describe_output.body` field,
- `snowflake_materialized_view` resource: `statement` field,
- `snowflake_oauth_integration_for_custom_clients` resource: `oauth_redirect_uri` and `describe_output.oauth_redirect_uri` fields,
- `snowflake_oauth_integration_for_partner_applications` resource: `oauth_redirect_uri` and `describe_output.oauth_redirect_uri` fields,
- `snowflake_materialized_view` resource: `statement` field,
- `snowflake_procedure_java` resource: `procedure_definition` and `show_output.arguments_raw` fields,
- `snowflake_procedure_javascript` resource: `procedure_definition` and `show_output.arguments_raw` fields,
- `snowflake_procedure_python` resource: `procedure_definition` and `show_output.arguments_raw` fields,
- `snowflake_procedure_scala` resource: `procedure_definition` and `show_output.arguments_raw` fields,
- `snowflake_procedure_sql` resource: `procedure_definition` and `show_output.arguments_raw` fields,
- `snowflake_row_access_policy` resource: `body` and `describe_output.body` fields,
- `snowflake_row_access_policies` data source: `describe_output.body` field,
- `snowflake_security_integrations` data source: `describe_output.redirect_uri` field,
- `snowflake_service_user` resource: `display_name`, `show_output.display_name`, `show_output.email`, `show_output.login_name`, `show_output.first_name`, `show_output.middle_name` and `show_output.last_name` fields,
- `snowflake_task` resource: `config`, `show_output.config` and `show_output.definition` fields,
- `snowflake_tasks` data source: `show_output.config` and `show_output.definition` fields,
- `snowflake_user` resource: `display_name`, `show_output.display_name`, `show_output.email`, `show_output.login_name`, `show_output.first_name`, `show_output.middle_name` and `show_output.last_name` fields,
- `snowflake_users` data source: `display_name`, `email`, `login_name`, `first_name`, `middle_name` and `last_name` fields nested in `show_output` and `describe_output`,
- `snowflake_view` resource: `statement` and `show_output.text` fields,
- `snowflake_views` data source: `show_output.text` field,

#### Removal of sensitive fields

The following table represents fields removed from resources. They were removed because of the Terraform SDK limitations
on marking data as sensitive in objects or collections ([Terraform issue reference](https://github.com/hashicorp/terraform/issues/28222)). Removal of computed output fields may have an impact on detecting
external changes (on the Snowflake side) for (usually) top-level fields they were referring to (e.g. `describe_output.oauth_client_id` -> `oauth_client_id`).

> Note: We may bring those fields back after exploring a better approaches (e.g., by using the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)), as currently, our options with the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2) are limited in that regard.

| Resource name                                                            | Removed fields                                                                                                             |
|--------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
| `snowflake_api_authentication_integration_with_authorization_code_grant` | `describe_output.oauth_client_id`                                                                                          |
| `snowflake_api_authentication_integration_with_client_credentials`       | `describe_output.oauth_client_id`                                                                                          |
| `snowflake_api_authentication_integration_with_jwt_bearer`               | `describe_output.oauth_client_id`                                                                                          |
| `snowflake_oauth_integration_for_partner_applications`                   | `describe_output.oauth_client_id`, `describe_output.oauth_redirect_uri`                                                    |
| `snowflake_oauth_integration_for_custom_clients`                         | `describe_output.oauth_client_id`, `describe_output.oauth_redirect_uri`                                                    |
| `snowflake_saml2_integration`                                            | `describe_output.saml2_snowflake_x509_cert`, `describe_output.saml2_x509_cert`                                             |
| `snowflake_security_integrations` (data source)                          | `security_integrations.describe_output.saml2_snowflake_x509_cert`, `security_integrations.describe_output.saml2_x509_cert` |
| `snowflake_users` (data source)                                          | `users.describe_output.password`                                                                                           |

### *(breaking change)* Changes in default TOML format
As we have announced in [an earlier entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#new-toml-file-schema), now the provider uses the new TOML format by default (`use_legacy_toml_file` is `false` by default). This means that when you try running the v2 provider with the same provider configuration which worked before, you can get a following error: `Error: 260000: account is empty` error with non-empty `account` configuration after upgrading to v2.

Please adjust your TOML format, basing on our [example](https://registry.terraform.io/providers/snowflakedb/snowflake/2.0.0/docs#examples).

This is a breaking change because it requires adjustments on the user's side.

Read more details in the mentioned [migration guide entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#new-toml-file-schema).

Alternatively, specify `use_legacy_toml_file=true` in your configuration, but this is not recommended. The legacy format is deprecated and will be removed in the next major release (v3).

### *(breaking change)* Changes in TOML configuration file requirements
As we have announced in [an earlier entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#changes-in-toml-configuration-file-requirements), now file permissions are verified by default (`skip_toml_file_permission_verification` is `false` by default). This means that on non-Windows systems, when you run the provider, you can get a following error:
```
could not load config file: config file /Users/user/.snowflake/config has unsafe permissions - 0755
```
Please adjust your file permissions, e.g. `chmod 0600 ~/.snowflake/config`.

This is a breaking change because it requires adjustments on the user's side.

Read more details in the mentioned [migration guide entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#changes-in-toml-configuration-file-requirements).

Alternatively, specify `skip_toml_file_permission_verification=false` in your provider configuration and use the unchanged TOML file, but this is less secure and not recommended.

### *(breaking change)* Improved data type handling for snowflake_masking_policy and snowflake_row_access_policy

The provider bases its logic on data returned from Snowflake. For data types, the responses are not always full, e.g.:
- `NUMBER(20, 4)` can be returned as `NUMBER`;
- `NUMBER` can be returned as `NUMBER`.

When you create an object with data type without specifying its arguments (like `NUMBER` without specified scale and precision), Snowflake fill in the defaults based on [SQL data types reference](https://docs.snowflake.com/en/sql-reference-data-types).
To be able to detect changes in config properly and to react to some external changes, we updated the way how we handle the data types:
- We use the Snowflake defaults for data types on the provider side; this is the exception to our [common approach of not hardcoding the Snowflake defaults](v1-preparations/CHANGES_BEFORE_V1.md#default-values) in the provider; we decided that the data type default are far less likely to change.
- We save the full data type in the state; specifying `NUMBER` will result in storing `NUMBER(38,0)` in state; the same value will be sent to Snowflake.
- We react to changes in Snowflake only in certain changes; e.g. `NUMBER` -> `VARCHAR`, we can't react on the external change if Snowflake does not return the full data type definition (as above).
- We will gradually add this logic to all the resources. For now, it was only added to the stable resources: [`snowflake_masking_policy`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.0.0/docs/resources/masking_policy) and [`snowflake_row_access_policy`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.0.0/docs/resources/row_access_policy). We added state upgraders that should handle the state saving changes. There are no changes required, however, you may encounter non-empty plans in these two resources after bumping. Be careful and verify the plan thoroughly as these resources can handle updates only in a destructive manner (this is the limitation of Snowflake SQL syntax for [`ALTER MASKING POLICY`](https://docs.snowflake.com/en/sql-reference/sql/alter-masking-policy) and [`ALTER ROW ACCESS POLICY`](https://docs.snowflake.com/en/sql-reference/sql/alter-row-access-policy)).

### *(bugfix)* Fix CSV_TIMESTAMP_FORMAT handling in snowflake_account_parameters

[`CSV_TIMESTAMP_FORMAT`](https://docs.snowflake.com/en/sql-reference/parameters#csv-timestamp-format) lacked the single quotes in the constructed SQL query. No changes are required.

References: [#3580](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3580)

## v1.2.2 ➞ v1.2.3

### *(bugfix)* Fix `snowflake_functions` and `snowflake_procedures` data sources with 2025_03 Bundle enabled

Check for more details and action steps needed in [Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands](./SNOWFLAKE_BCR_MIGRATION_GUIDE.md#argument-output-changes-for-show-functions-and-show-procedures-commands).

References: [#3822](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3822)

### *(bugfix)* Fix all function and procedure resources with 2025_03 Bundle enabled

Check for more details and action steps needed in [Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands](./SNOWFLAKE_BCR_MIGRATION_GUIDE.md#argument-output-changes-for-show-functions-and-show-procedures-commands).

References: [#3823](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3823)

## v1.2.1 -> v1.2.2

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

This fix was also backported to versions v1.0.6 and v1.1.1.

## v1.2.0 ➞ v1.2.1
No migration needed.

## v1.1.x ➞ v1.2.0
<a id="v110--v120"></a>

### New behavior for Read and Delete operations when removing high-hierarchy objects
Some objects in Snowflake are created in hierarchy, for example, tables (database → schema → table).
When the user wants to remove the higher-hierarchy object (like a database), the lower-hierarchy objects should be removed beforehand.
Otherwise, Terraform would fail to remove the lower-hierarchy objects from the state,
and without manual state management it wouldn't be possible to remove this object, ending up in broken state (reference issue: [#1243](https://github.com/snowflakedb/terraform-provider-snowflake/issues/1243)).
This may only happen in particular cases, for example, if part of the hierarchy is managed outside Terraform or the configuration is missing dependencies between resources.
This behavior was described more in detail in [our documentation](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/object_renaming_research_summary#renaming-higher-hierarchy-objects).

For improved usability, we adjusted the Read and Delete operation implementations for all resources,
so that now, they're able to know the higher-hierarchy object is missing, and they can safely remove themselves from the state.

To demonstrate this behavior, let's take the following configuration:
```terraform
resource "snowflake_table" "test" {
  database = "TEST_DATABASE"
  schema   = "PUBLIC"
  name     = "TEMP_TABLE"
  column {
    name = "ID"
    type = "NUMBER"
  }
}
```
> Note: The `TEST_DATABASE` is created manually through Snowflake and the table configuration is already applied through Terraform.

When you remove the database by running `DROP DATABASE TEST_DATABASE` in Snowflake, and then run `terraform apply`,
previously, you would end up in the infinite loop of errors and only manual removal from state (`terraform state rm snowflake_table.test`)
would help you to remove the table from the state (and then from the configuration).

In the future, we are planning to do the same with object attachments, like grants, policies, etc. (To address cases like: [#3412](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3412))

### New TOML file schema
The TOML file schema before v1.2.0 was not consistent with the configuration keys in the provider. The main differences were:
- The keys in the provider contain an underscore (`_`) as a separator, but the TOML schema has fields without any separator.
- The field `driver_tracing` in the provider is related to `tracing` in the TOML schema.

These differences caused some confusion for the users. This is why we decided to introduce a new TOML schema addressing these flaws.
The new schema uses underscore (`_`) as a separator, and changes `tracing` to `driver_tracing` to be consistent with the provider schema.
You can see an example in [our registry](https://registry.terraform.io/providers/snowflakedb/snowflake/1.2.0/docs#order-precedence).

The default behavior is the same as before v1.2.0. You can enable the new behavior by setting `use_legacy_toml_file = false` in the provider, or by setting `SNOWFLAKE_USE_LEGACY_TOML_FILE=false` environmental variable.
Please note that we will change the default behavior in v2: the new TOML schema will be read by default, but there will still be a possibility to read the old format. However, we encourage you to use our new schema now and give us feedback.

NOTE: With this change, we are not bringing back the `account` field yet. This means that `account_name` and `organization_name` fields must still be used. We will discuss about this field after v2.

References: [#3553](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3553)

### Fixes in handling references to computed fields in context of `show_output`
The issue could arise in almost any object using show_output when the following steps happen:
1. Object is created.
2. Object's attribute is updated to computed value of other object added in this run.

In such situation the final plan could result in error like this one:
```
| Error: Provider produced inconsistent final plan
|
| When expanding the plan for snowflake_legacy_service_user.one to include new
| values learned so far during apply, provider
| "registry.terraform.io/hashicorp/snowflake" produced an invalid new value for
| .show_output: was known, but now unknown.
|
| This is a bug in the provider, which should be reported in the provider's own
| issue tracker.
```

This version fixes this behavior. No action should be required on user side.

References: [#3522](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3522)

## v1.1.0 ➞ v1.1.1
<a id="v110---v111"></a>

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

This fix was also backported to version v1.0.6.

## v1.0.x ➞ v1.1.0
<a id="v105--v110"></a>

### Timeouts in resources
By default, resource operation timeout after 20 minutes ([reference](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts#default-timeouts-and-deadline-exceeded-errors)). This caused some long running operations to timeout.

We already added configurable timeouts to [execute](https://registry.terraform.io/providers/snowflakedb/snowflake/1.0.4/docs/resources/execute#nested-schema-for-timeouts), [tag_association](https://registry.terraform.io/providers/snowflakedb/snowflake/1.0.4/docs/resources/tag_association#nested-schema-for-timeouts) and [cortex_search_service](https://registry.terraform.io/providers/snowflakedb/snowflake/1.0.4/docs/resources/cortex_search_service#nested-schema-for-timeouts) before. Now, we also allow setting them on all other resources.
Data sources will be supported in the future.

Read more about resource timeouts in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts).

References: [#3355](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3355)

### Fixes in `grant_privileges_to_account_role` resource
Using `grant_privileges_to_account_role` with `all_privileges=true` and `on_account = true` option started to fail recently due to newly introduced privileges in Snowflake:
```
003011 (42501): Grant partially executed: privileges [MANAGE LISTING AUTO FULFILLMENT, MANAGE ORGANIZATION SUPPORT CASES,
│ MANAGE POLARIS CONNECTIONS] not granted
```

Instead of failing the whole action, we return a warning instead and the operation execution continues, which aligns with the behavior in Snowsight. Note that for `all_privileges=true` the privileges list in the state is not populated, like before. If you want to detect differences in the privileges, use `privileges` list instead. If you want to make sure that the maximum privileges are granted, enable `always_apply`.

References: [#3507](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3507)

## v1.0.5 ➞ v1.0.6
<a id="v105---v106"></a>

### *(bugfix)* Fix `ENABLE_INTERNAL_STAGES_PRIVATELINK` mapping in `snowflake_account_parameter` resource

Due to incorrect mapping in setting account parameter logic in [`snowflake_account_parameter`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account_parameter), the [`ENABLE_INTERNAL_STAGES_PRIVATELINK`](https://docs.snowflake.com/en/sql-reference/parameters#enable-internal-stages-privatelink) could not be set. Setting it results in setting the [`ALLOW_ID_TOKEN`](https://docs.snowflake.com/en/sql-reference/parameters#allow-id-token) parameter instead. This version introduces the corrected mapping.

No configuration changes are needed. However, the provider won't set back the `ALLOW_ID_TOKEN` parameter value as we can't detect if setting its value was intentional (manually or through `snowflake_account_parameter`). Because of that, please verify your `ALLOW_ID_TOKEN` parameter and set it to the desired value.

## v1.0.4 ➞ v1.0.5

### Changes in TOML configuration file requirements
Before this version, it was possible to abuse the provider by providing a huge TOML config file which was read every time. To mitigate this, we set a limit of the supported file size to 10MB. For a larger TOML configuration file, the provider will fail.

We encourage you to make your TOML configuration file more restricted. Any privileges for a UNIX group or others should not be set (the maximum recommended privilege is `700`). You can set the expected privileges like `chmod 0600 ~/.snowflake/config`. This check is not enabled by default, but we introduced a new `skip_toml_file_permission_verification` boolean field to the provider's configuration with a default `true` value to enable this behavior.

If you want to check the TOML config file privileges, please specify `skip_toml_file_permission_verification=false` in your TF configuration or set `SKIP_TOML_FILE_PERMISSION_VERIFICATION=FALSE` environment variable. For a TOML configuration file with too broad permissions, the provider will fail.
This requirement can be checked only on non-Windows platforms. If you are using the provider on Windows, please make sure that your configuration file has not too permissive privileges.

This setting can be changed to `false` in the future, meaning verifying file permissions by default, so the preferred action is to set the proper permissions now and to disable skipping permission verification.

### Tracking external changes for oauth_redirect_uri in the snowflake_oauth_integration_for_partner_applications resource
From this version, the snowflake_oauth_integration_for_partner_applications resource is able to
detect changes on the Snowflake side and apply appropriate action from the provider level. This may produce
changes after running `terraform plan`, as before the configuration could contain different value than on the Snowflake side.

### Removal of instrumentation library
We decided to remove the instrumentation around the [Go Snowflake driver](https://github.com/snowflakedb/gosnowflake). It does not introduce any functional changes, however, it changes the way the Snowflake communication logs are turned on and how they are printed. Check [this section](FAQ.md#how-can-i-turn-on-logs) for more details.

`SF_TF_NO_INSTRUMENTED_SQL`, used to turn the instrumentation off, was removed because it is no longer needed.

These changes should not affect any existing workflows (unless you have custom logic based on the old logs output).

### Removal of additional debug logs for the `snowflake_grant_privileges_to_role` resource

The environment variable `SF_TF_ADDITIONAL_DEBUG_LOGGING` was used to turn on the additional logging in the `snowflake_grant_privileges_to_role` resource. The additional logger was later used in multiple other places. We are currently removing it completely; however, we plan to address the logging topic globally in the provider.

These changes should not affect any existing workflows (unless you have custom logic based on the additional logs output - `sf-tf-additional-debug` prefix).

## v1.0.3 ➞ v1.0.4

### Fixed external_function VARCHAR return_type
VARCHAR external_function return_type did not work correctly before ([#3392](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3392)) but was fixed in this version.

### New Go version and conflicts with Suricata-based firewalls (like AWS Network Firewall)
In this version we bumped our underlying Go version to v1.23.6.
Based on issue [#3421](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3421)
it seems it introduces changes to the standard library that may not be supported by other third party software.
The issue presents one of those changes that seem to be introduced in Golang's `crypto/tls` package.
One thing that is valuable in such cases is to check the [GODEBUG](https://go.dev/doc/godebug)
documentation page (especially [history section](https://go.dev/doc/godebug#history)).
It specifies a set of parameters which can be turned on/off depending on
what features of Go would you like to use or resign from. The solution for this issue was to set
the GODEBUG environment variable to `GODEBUG=tlskyber=0`.

## v1.0.2 ➞ v1.0.3

### Fixed METRIC_LEVEL parameter
METRIC_LEVEL account parameter did not work correctly before ([#3375](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3375)), but was fixed in this version.

### Fixed ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES parameter
ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES account parameter did not work correctly before ([#3344]). This parameter was of incorrect type, and the constructed queries did not provide the parameter's value during altering accounts. It has been fixed in this version.

### Changed documentation structure
We added `Preview` and `Stable` categories to the resources and data sources documentation, which clearly separates the preview and stable features in the documentation feature list.
We moved our technical guides to `guides` directory. This means that all such guides are available natively in the registry, similarly to [Unassigning policies](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/unassigning_policies) guide.
We also updated the links to point to the docs inside the registry. Note that our [Roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md) and [Migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md) are available in Github only.
This is a part of our effort to improve the provider documentation. We are open for your feedback and suggestions.

## v1.0.1 ➞ v1.0.2

### Fixed migration of account resource
Previously, during upgrading the provider from v0.99.0, when account fields `must_change_password` or `is_org_admin` were not set in state, the provider panicked. It has been fixed in this version.

### Add missing resource monitor in `snowflake_grant_ownership` resource
Resource monitor in not currently listed as option in `GRANT OWNERSHIP` documentation ([here](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters)) but this is a valid option. `snowflake_grant_ownership` was updated to support resource monitors.

References: [#3318](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3318)

### Timeouts in `snowflake_execute`
By default, resource operation timeouts after 20 minutes ([reference](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts#default-timeouts-and-deadline-exceeded-errors)). Because of generic nature of `snowflake_execute`, we decided to bump its default timeouts to 60 minutes; We also allowed setting them on the resource config level (following [official documentation](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts)).

References: [#3334](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3334)

## v1.0.0 ➞ v1.0.1

### Fixes in account parameters
As a follow-up of reworked `snowflake_account_parameter`, this version has several improvements regarding handling parameters.

#### Add missing parameters based on the docs and output of SHOW PARAMETERS IN ACCOUNT
Based on [parameters docs](https://docs.snowflake.com/en/sql-reference/parameters) and `SHOW PARAMETERS IN ACCOUNT`, we established a list of supported parameters. New supported or fixed parameters in `snowflake_account_parameter`:
- `ACTIVE_PYTHON_PROFILER`
- `CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS`
- `CORTEX_ENABLED_CROSS_REGION`
- `CSV_TIMESTAMP_FORMAT`
- `ENABLE_PERSONAL_DATABASE`
- `ENABLE_UNHANDLED_EXCEPTIONS_REPORTING`
- `ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES`
- `HYBRID_TABLE_LOCK_TIMEOUT`
- `JS_TREAT_INTEGER_AS_BIGINT`
- `PREVENT_UNLOAD_TO_INLINE_URL`
- `PREVENT_UNLOAD_TO_INTERNAL_STAGES`
- `PYTHON_PROFILER_MODULES`
- `PYTHON_PROFILER_TARGET_STAGE`
- `STORAGE_SERIALIZATION_POLICY`
- `TASK_AUTO_RETRY_ATTEMPTS`

#### Adjusted validations
Validations for number parameters are now relaxed. This is because a few of the value limits are soft limits in Snowflake, and can be changed externally.
We decided to keep validations for non-negative values. Affected parameters:
- `QUERY_TAG`
- `TWO_DIGIT_CENTURY_START`
- `WEEK_OF_YEAR_POLICY`
- `WEEK_START`
- `USER_TASK_TIMEOUT_MS`

We added non-negative validations for the following parameters:
- `CLIENT_PREFETCH_THREADS`
- `CLIENT_RESULT_CHUNK_SIZE`
- `CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY`
- `HYBRID_TABLE_LOCK_TIMEOUT`
- `JSON_INDENT`
- `STATEMENT_QUEUED_TIMEOUT_IN_SECONDS`
- `STATEMENT_TIMEOUT_IN_SECONDS`
- `TASK_AUTO_RETRY_ATTEMPTS`
- `USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS`

Note that enum parameters are still not validated by the provider - they are only validated in Snowflake. We will handle this during a small rework of the parameters in the future.

### Add missing preview features to config

Values:
- `snowflake_functions_datasource`
- `snowflake_procedures_datasource`
- `snowflake_tables_datasource`
  were missing in the `preview_features_enabled` attribute in the provider's config. They were added.

References: #3302

### functions and procedures docs updated

Argument names are automatically wrapped in double quotes, so:
- uppercase names should be used or
- argument name should be quoted in the procedure/function definition.

Updated the docs and the previous migration guide entry.

References: #3298

### python procedure docs updated

Importing python procedure is currently limited to procedures with snowflake-snowpark-python version explicitly set in Snowflake. Docs were updated.

References: #3303

## v0.100.0 ➞ v1.0.0

### Preview features flag
All of the preview features objects are now disabled by default. This includes:
- Resources
  - `snowflake_account_password_policy_attachment`
  - `snowflake_alert`
  - `snowflake_api_integration`
  - `snowflake_cortex_search_service`
  - `snowflake_dynamic_table`
  - `snowflake_external_function`
  - `snowflake_external_table`
  - `snowflake_external_volume`
  - `snowflake_failover_group`
  - `snowflake_file_format`
  - `snowflake_function_java`
  - `snowflake_function_javascript`
  - `snowflake_function_python`
  - `snowflake_function_scala`
  - `snowflake_function_sql`
  - `snowflake_managed_account`
  - `snowflake_materialized_view`
  - `snowflake_network_policy_attachment`
  - `snowflake_network_rule`
  - `snowflake_email_notification_integration`
  - `snowflake_notification_integration`
  - `snowflake_object_parameter`
  - `snowflake_password_policy`
  - `snowflake_pipe`
  - `snowflake_procedure_java`
  - `snowflake_procedure_javascript`
  - `snowflake_procedure_python`
  - `snowflake_procedure_scala`
  - `snowflake_procedure_sql`
  - `snowflake_sequence`
  - `snowflake_share`
  - `snowflake_stage`
  - `snowflake_storage_integration`
  - `snowflake_table`
  - `snowflake_table_column_masking_policy_application`
  - `snowflake_table_constraint`
  - `snowflake_user_public_keys`
  - `snowflake_user_password_policy_attachment`
- Data sources
  - `snowflake_current_account`
  - `snowflake_alerts`
  - `snowflake_cortex_search_services`
  - `snowflake_database`
  - `snowflake_database_role`
  - `snowflake_dynamic_tables`
  - `snowflake_external_functions`
  - `snowflake_external_tables`
  - `snowflake_failover_groups`
  - `snowflake_file_formats`
  - `snowflake_functions`
  - `snowflake_materialized_views`
  - `snowflake_pipes`
  - `snowflake_procedures`
  - `snowflake_current_role`
  - `snowflake_sequences`
  - `snowflake_shares`
  - `snowflake_parameters`
  - `snowflake_stages`
  - `snowflake_storage_integrations`
  - `snowflake_system_generate_scim_access_token`
  - `snowflake_system_get_aws_sns_iam_policy`
  - `snowflake_system_get_privatelink_config`
  - `snowflake_system_get_snowflake_platform_info`
  - `snowflake_tables`

If you want to have them enabled, add the feature name to the provider configuration (with `_datasource` or `_resource` suffix), like this:
```terraform
provider "snowflake" {
	preview_features_enabled = ["snowflake_current_account_datasource", "snowflake_alert_resource"]
}
```

Do not forget to add this line to all provider configurations using these features, including [provider aliases](https://developer.hashicorp.com/terraform/language/providers/configuration#alias-multiple-provider-configurations).

### Removed deprecated objects
All of the deprecated objects are removed from v1 release. This includes:
- Resources
  - `snowflake_database_old` - see [migration guide](#new-feature-new-database-resources)
  - `snowflake_role` - see [migration guide](#new-feature-new-snowflake_account_role-resource)
  - `snowflake_oauth_integration` - see [migration guide](#new-feature-snowflake_oauth_integration_for_custom_clients-and-snowflake_oauth_integration_for_partner_applications-resources)
  - `snowflake_saml_integration` - see [migration guide](#new-feature-snowflake_saml2_integration-resource)
  - `snowflake_session_parameter`
  - `snowflake_stream` - see [migration guide](#new-feature-snowflake_stream_on_directory_table-and-snowflake_stream_on_view-resource)
  - `snowflake_tag_masking_policy_association` - see [migration guide](#snowflake_tag_masking_policy_association-deprecation)
  - `snowflake_function`
  - `snowflake_procedure`
  - `snowflake_unsafe_execute` - see [migration guide](#unsafe_execute-resource-deprecation--new-execute-resource)
- Data sources
  - `snowflake_role` - see [migration guide](#snowflake_role-data-source-deprecation)
  - `snowflake_roles` - see [migration guide](#new-feature-account-role-data-source)
- Fields in the provider configuration:
  - `account` - see [migration guide](#behavior-change-deprecated-fields)
  - OAuth related fields - see [migration guide](#structural-change-oauth-api):
    - `oauth_access_token`
    - `oauth_client_id`
    - `oauth_client_secret`
    - `oauth_endpoint`
    - `oauth_redirect_url`
    - `oauth_refresh_token`
    - `browser_auth`
  - `private_key_path` - see [migration guide](#private_key_path-deprecation)
  - `region` - see [migration guide](#remove-redundant-information-region)
  - `session_params` - see [migration guide](#rename-session_params--params)
  - `username` - see [migration guide](#rename-username--user)
- Fields in `tag` resource:
  - `object_name`

Additionally, `JWT` value is no longer available for `authenticator` field in the provider configuration.

## v0.99.0 ➞ v0.100.0

### *(preview feature/deprecation)* Function and procedure resources

`snowflake_function` is now deprecated in favor of 5 new preview resources:

- `snowflake_function_java`
- `snowflake_function_javascript`
- `snowflake_function_python`
- `snowflake_function_scala`
- `snowflake_function_sql`

It will be removed with the v1 release. Please check the docs for the new resources and adjust your configuration files.
For no downtime migration, follow our [guide](./docs/guides/resource_migration.md).

The new resources are more aligned with current features like:
- external access integrations support
- secrets support
- argument default values

**Note**: argument names are now quoted automatically by the provider so remember about this while writing the function definition (argument name should be quoted or uppercase should be used for the argument name).

`snowflake_procedure` is now deprecated in favor of 5 new preview resources:

- `snowflake_procedure_java`
- `snowflake_procedure_javascript`
- `snowflake_procedure_python`
- `snowflake_procedure_scala`
- `snowflake_procedure_sql`

It will be removed with the v1 release. Please check the docs for the new resources and adjust your configuration files.
For no downtime migration, follow our [guide](./docs/guides/resource_migration.md).

The new resources are more aligned with current features like:
- external access integrations support
- secrets support
- argument default values

**Note**: argument names are now quoted automatically by the provider so remember about this while writing the procedure definition (argument name should be quoted or uppercase should be used for the argument name).

### *(new feature)* Account role data source
Added a new `snowflake_account_roles` data source for account roles. Now it reflects It's based on `snowflake_roles` data source.
`account_roles` field now organizes output of show under `show_output` field.

Before:
```terraform
output "simple_output" {
  value = data.snowflake_roles.test.roles[0].show_output[0].name
}
```
After:
```terraform
output "simple_output" {
  value = data.snowflake_account_roles.test.account_roles[0].show_output[0].name
}
```

### snowflake_roles data source deprecation
`snowflake_roles` is now deprecated in favor of `snowflake_account_roles` with a similar schema and behavior. It will be removed with the v1 release. Please adjust your configuration files.

### snowflake_account_parameter resource changes

#### *(behavior change)* resource deletion
During resource deleting, provider now uses `UNSET` instead of `SET` with the default value.

#### *(behavior change)* changes in `key` field
The value of `key` field is now case-insensitive and is validated. The list of supported values is available in the resource documentation.

### unsafe_execute resource deprecation / new execute resource

The `snowflake_unsafe_execute` gets deprecated in favor of the new resource `snowflake_execute`.
The `snowflake_execute` was build on top of `snowflake_unsafe_execute` with a few improvements.
The unsafe version will be removed with the v1 release, so please migrate to the `snowflake_execute` resource.

For no downtime migration, follow our [guide](./docs/guides/resource_migration.md).
When importing, remember that the given resource id has to be unique (using UUIDs is recommended).
Also, because of the nature of the resource, first apply after importing is necessary to "copy" values from the configuration to the state.

### snowflake_oauth_integration_for_partner_applications and snowflake_oauth_integration_for_custom_clients resource changes
#### *(behavior change)* `blocked_roles_list` field is no longer required

Previously, `blocked_roles_list` field was required to handle default account roles like `ACCOUNTADMIN`, `ORGADMIN`, and `SECURITYADMIN`.

Now, it is optional, because of using the value of `OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST` parameter (read more below).

No changes in the configuration are necessary.

#### *(behavior change)* new field `related_parameters`

To handle `blocked_roles_list` field properly in both of the resources, we introduce `related_parameters` field. This field is a list of parameters related to OAuth integrations. It is a computed-only field containing value of `OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST` account parameter (see [docs](https://docs.snowflake.com/en/sql-reference/parameters#oauth-add-privileged-roles-to-blocked-list)).

### snowflake_account resource changes

Changes:
- `admin_user_type` is now supported. No action required during the migration.
- `grace_period_in_days` is now required. The field should be explicitly set in the following versions.
- Account renaming is now supported.
- `is_org_admin` is a settable field (previously it was read-only field). Changing its value is also supported.
- `must_change_password` and `is_org_admin` type was changed from `bool` to bool-string (more on that [here](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#empty-values)). No action required during the migration.
- The underlying resource identifier was changed from `<account_locator>` to `<organization_name>.<account_name>`. Migration will be done automatically. Notice this introduces changes in how `snowflake_account` resource is imported.
- New `show_output` field was added (see [raw Snowflake output](./v1-preparations/CHANGES_BEFORE_V1.md#raw-snowflake-output)).

### snowflake_accounts data source changes
New filtering options:
- `with_history`

New output fields
- `show_output`

Breaking changes:
- `pattern` renamed to `like`
- `accounts` field now organizes output of show under `show_output` field and the output of show parameters under `parameters` field.

Before:
```terraform
output "simple_output" {
  value = data.snowflake_accounts.test.accounts[0].account_name
}
```
After:
```terraform
output "simple_output" {
  value = data.snowflake_accounts.test.accounts[0].show_output[0].account_name
}
```

### snowflake_tag_association resource changes
#### *(behavior change)* new id format
To provide more functionality for tagging objects, we have changed the resource id from `"TAG_DATABASE"."TAG_SCHEMA"."TAG_NAME"` to `"TAG_DATABASE"."TAG_SCHEMA"."TAG_NAME"|TAG_VALUE|OBJECT_TYPE`. This allows to group tags associations per tag ID, tag value and object type in one resource.
```
resource "snowflake_tag_association" "gold_warehouses" {
  object_identifiers = [snowflake_warehouse.w1.fully_qualified_name, snowflake_warehouse.w2.fully_qualified_name]
  object_type = "WAREHOUSE"
  tag_id      = snowflake_tag.tier.fully_qualified_name
  tag_value   = "gold"
}
resource "snowflake_tag_association" "silver_warehouses" {
  object_identifiers = [snowflake_warehouse.w3.fully_qualified_name]
  object_type = "WAREHOUSE"
  tag_id      = snowflake_tag.tier.fully_qualified_name
  tag_value   = "silver"
}
resource "snowflake_tag_association" "silver_databases" {
  object_identifiers = [snowflake_database.d1.fully_qualified_name]
  object_type = "DATABASE"
  tag_id      = snowflake_tag.tier.fully_qualified_name
  tag_value   = "silver"
}
```

Note that if you want to promote silver instances to gold, you can not simply change `tag_value` in `silver_warehouses`. Instead, you should first remove `object_identifiers` from `silver_warehouses`, run `terraform apply`, and then add the relevant `object_identifiers` in `gold_warehouses`, like this (note that `silver_warehouses` resource was deleted):
```
resource "snowflake_tag_association" "gold_warehouses" {
  object_identifiers = [snowflake_warehouse.w1.fully_qualified_name, snowflake_warehouse.w2.fully_qualified_name, snowflake_warehouse.w3.fully_qualified_name]
  object_type = "WAREHOUSE"
  tag_id      = snowflake_tag.tier.fully_qualified_name
  tag_value   = "gold"
}
```
and run `terraform apply` again.

Note that the order of operations is not deterministic in this case, and if you do these operations in one step, it is possible that the tag value will be changed first, and unset later because of removing the resource with old value.

The state is migrated automatically. There is no need to adjust configuration files, unless you use resource id `snowflake_tag_association.example.id` as a reference in other resources.

#### *(behavior change)* changed fields
Behavior of some fields was changed:
- `object_identifier` was renamed to `object_identifiers` and it is now a set of fully qualified names. Change your configurations from
```
resource "snowflake_tag_association" "table_association" {
  object_identifier {
    name     = snowflake_table.test.name
    database = snowflake_database.test.name
    schema   = snowflake_schema.test.name
  }
  object_type = "TABLE"
  tag_id      = snowflake_tag.test.fully_qualified_name
  tag_value   = "engineering"
}
```
to
```
resource "snowflake_tag_association" "table_association" {
  object_identifiers = [snowflake_table.test.fully_qualified_name]
  object_type = "TABLE"
  tag_id      = snowflake_tag.test.fully_qualified_name
  tag_value   = "engineering"
}
```
- `tag_id`  has now suppressed identifier quoting to prevent issues with Terraform showing permament differences, like [this one](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2982)
- `object_type` and `tag_id` are now marked as ForceNew

The state is migrated automatically. Please adjust your configuration files.

### Data type changes

As part of reworking functions, procedures, and any other resource utilizing Snowflake data types, we adjusted the parsing of data types to be more aligned with Snowflake (according to [docs](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types)).

Affected resources:
- `snowflake_function`
- `snowflake_procedure`
- `snowflake_table`
- `snowflake_external_function`
- `snowflake_masking_policy`
- `snowflake_row_access_policy`
- `snowflake_dynamic_table`
You may encounter non-empty plans in these resources after bumping.

Changes to the previous implementation/limitations:
- `BOOL` is no longer supported; use `BOOLEAN` instead.
- Following the change described [here](#bugfix-handle-data-type-diff-suppression-better-for-text-and-number), comparing and suppressing changes of data types was extended for all other data types with the following rules:
  - `CHARACTER`, `CHAR`, `NCHAR` now have the default size set to 1 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-text#char-character-nchar))
  - `BINARY` has default size set to 8388608 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-text#binary))
  - `TIME` has default precision set to 9 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-datetime#time))
  - `TIMESTAMP_LTZ` has default precision set to 9 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp)); supported aliases: `TIMESTAMPLTZ`, `TIMESTAMP WITH LOCAL TIME ZONE`.
  - `TIMESTAMP_NTZ` has default precision set to 9 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp)); supported aliases: `TIMESTAMPNTZ`, `TIMESTAMP WITHOUT TIME ZONE`, `DATETIME`.
  - `TIMESTAMP_TZ` has default precision set to 9 if not provided (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp)); supported aliases: `TIMESTAMPTZ`, `TIMESTAMP WITH TIME ZONE`.
- The session-settable `TIMESTAMP` is NOT supported ([docs](https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp))
- `VECTOR` type still is limited and will be addressed soon (probably before the release so it will be edited)

## v0.98.0 ➞ v0.99.0

### snowflake_tasks data source changes

New filtering options:
- `with_parameters`
- `like`
- `in`
- `starts_with`
- `root_only`
- `limit`

New output fields
- `show_output`
- `parameters`

Breaking changes:
- `database` and `schema` are right now under `in` field

Before:
```terraform
data "snowflake_tasks" "old_tasks" {
  database = "<database_name>"
  schema = "<schema_name>"
}
```
After:
```terraform
data "snowflake_tasks" "new_tasks" {
  in {
    # for IN SCHEMA specify:
    schema = "<database_name>.<schema_name>"

    # for IN DATABASE specify:
    database = "<database_name>"
  }
}
```
- `tasks` field now organizes output of show under `show_output` field and the output of show parameters under `parameters` field.

Before:
```terraform
output "simple_output" {
  value = data.snowflake_tasks.test.tasks[0].name
}
```
After:
```terraform
output "simple_output" {
  value = data.snowflake_tasks.test.tasks[0].show_output[0].name
}
```

### snowflake_task resource changes
New fields:
- `config` - enables to specify JSON-formatted metadata that can be retrieved in the `sql_statement` by using [SYSTEM$GET_TASK_GRAPH_CONFIG](https://docs.snowflake.com/en/sql-reference/functions/system_get_task_graph_config).
- `show_output` and `parameters` fields added for holding SHOW and SHOW PARAMETERS output (see [raw Snowflake output](./v1-preparations/CHANGES_BEFORE_V1.md#raw-snowflake-output)).
- Added support for finalizer tasks with `finalize` field. It conflicts with `after` and `schedule` (see [finalizer tasks](https://docs.snowflake.com/en/user-guide/tasks-graphs#release-and-cleanup-of-task-graphs)).

Changes:
- `enabled` field changed to `started` and type changed to string with only boolean values available (see ["empty" values](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values)). It is also now required field, so make sure it's explicitly set (previously it was optional with the default value set to `false`).
- `allow_overlapping_execution` type was changed to string with only boolean values available (see ["empty" values](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values)). Previously, it had the default set to `false` which will be migrated. If nothing will be set the provider will plan the change to `default` value. If you want to make sure it's turned off, set it explicitly to `false`.

Before:
```terraform
resource "snowflake_task" "example" {
  # ...
  enabled = true
  # ...
}
```
After:
```terraform
resource "snowflake_task" "example" {
  # ...
  started = true
  # ...
}
```
- `schedule` field changed from single value to a nested object that allows for specifying either minutes or cron

Before:
```terraform
resource "snowflake_task" "example" {
  # ...
  schedule = "5 MINUTES"
  # or
  schedule = "USING CRON * * * * * UTC"
  # ...
}
```
After:
```terraform
resource "snowflake_task" "example" {
  # ...
  schedule {
    minutes = 5
    # or
    using_cron = "* * * * * UTC"
  }
  # ...
}
```
- All task parameters defined in [the Snowflake documentation](https://docs.snowflake.com/en/sql-reference/parameters) added into the top-level schema and removed `session_parameters` map.

Before:
```terraform
resource "snowflake_task" "example" {
  # ...
  session_parameters = {
    QUERY_TAG = "<query_tag>"
  }
  # ...
}
```
After:
```terraform
resource "snowflake_task" "example" {
  # ...
  query_tag = "<query_tag>"
  # ...
}
```

- `after` field type was changed from `list` to `set` and the values were changed from names to fully qualified names.

Before:
```terraform
resource "snowflake_task" "example" {
  # ...
  after = ["<task_name>", snowflake_task.some_task.name]
  # ...
}
```
After:
```terraform
resource "snowflake_task" "example" {
  # ...
  after = ["<database_name>.<schema_name>.<task_name>", snowflake_task.some_task.fully_qualified_name]
  # ...
}
```

### *(new feature)* snowflake_tags datasource
Added a new datasource enabling querying and filtering tags. Notes:
- all results are stored in `tags` field.
- `like` field enables tags filtering by name.
- `in` field enables tags filtering by `account`, `database`, `schema`, `application` and `application_package`.
- `SHOW TAGS` output is enclosed in `show_output` field inside `tags`.

### snowflake_tag_masking_policy_association deprecation
`snowflake_tag_masking_policy_association` is now deprecated in favor of `snowflake_tag` with a new `masking_policy` field. It will be removed with the v1 release. Please adjust your configuration files.

### snowflake_tag resource changes
New fields:
  - `masking_policies` field that holds the associated masking policies.
  - `show_output` field that holds the response from SHOW TAGS.

#### *(breaking change)* Changed fields in snowflake_masking_policy resource
Changed fields:
  - `name` is now not marked as ForceNew. When this value is changed, the resource is renamed with `ALTER TAG`, instead of being recreated.
  - `allowed_values` type was changed from list to set. This causes different ordering to be ignored.
State will be migrated automatically.

#### *(breaking change)* Identifiers related changes
During [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`). Importing resources also needs to be adjusted (see [example](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/tag#import)).

Also, we added diff suppress function that prevents Terraform from showing differences, when only quoting is different.

No change is required, the state will be migrated automatically.

#### *(breaking change)* Required warehouse
For this resource, the provider now uses [tag references](https://docs.snowflake.com/en/sql-reference/functions/tag_references) to get information about masking policies attached to tags. This function requires a warehouse in the connection. Please, make sure you have either set a `DEFAULT_WAREHOUSE` for the user, or specified a warehouse in the provider configuration.

## v0.97.0 ➞ v0.98.0

### *(new feature)* snowflake_connections datasource
Added a new datasource enabling querying and filtering connections. Notes:
- all results are stored in `connections` field.
- `like` field enables connections filtering.
- SHOW CONNECTIONS output is enclosed in `show_output` field inside `connections`.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.


### *(new feature)* connection resources

Added a new resources for managing connections. We decided to split connection into two separate resources based on whether the connection is a primary or replicated (secondary). i.e.:

- `snowflake_primary_connection` is used to manage primary connection, with ability to enable failover to other accounts.
- `snowflake_secondary_connection` is used to manage replicated (secondary) connection.

To promote `snowflake_secondary_connection` to `snowflake_primary_connection`, resources need to be removed from the state, altered manually using:
```
ALTER CONNECTION <name> PRIMARY;
```
and then imported again, now as `snowflake_primary_connection`.

To demote `snowflake_primary_connection` back to `snowflake_secondary_connection`, resources need to be removed from the state, re-created manually using:
```
CREATE CONNECTION <name> AS REPLICA OF <organization_name>.<account_name>.<connection_name>;
```
and then imported as `snowflake_secondary_connection`.

For guidance on removing and importing resources into the state check [resource migration](./docs/guides/resource_migration.md).

See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-connection).

### snowflake_streams data source changes
New filtering options:
- `like`
- `in`
- `starts_with`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `streams` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### *(behavior change)* Provider configuration rework
On our road to v1, we have decided to rework configuration to address the most common issues (see a [roadmap entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#providers-configuration-rework)). We have created a list of topics we wanted to address before v1. We will prepare an announcement soon. The following subsections describe the things addressed in the v0.98.0.

#### *(behavior change)* new fields
We have added new fields to match the ones in [the driver](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#Config) and to simplify setting account name. Specifically:
- `include_retry_reason`, `max_retry_count`, `driver_tracing`, `tmp_directory_path` and `disable_console_login` are the new fields that are supported in the driver
- `disable_saml_url_check` will be added to the provider after upgrading the driver
- `account_name` and `organization_name` were added to improve handling account names. Execute `SELECT CURRENT_ORGANIZATION_NAME(), CURRENT_ACCOUNT_NAME();` to get the required values. Read more in [docs](https://docs.snowflake.com/en/user-guide/admin-account-identifier#using-an-account-name-as-an-identifier).

#### *(behavior change)* changed configuration of driver log level
To be more consistent with other configuration options, we have decided to add `driver_tracing` to the configuration schema. This value can also be configured by `SNOWFLAKE_DRIVER_TRACING` environmental variable and by `drivertracing` field in the TOML file. The previous `SF_TF_GOSNOWFLAKE_LOG_LEVEL` environmental variable is not supported now, and was removed from the provider.

#### *(behavior change)* deprecated fields
Because of new fields `account_name` and `organization_name`, `account` is now deprecated. It will be removed with the v1 release.
If you use Terraform configuration file, adjust it from
```terraform
provider "snowflake" {
	account = "ORGANIZATION-ACCOUNT"
}
```

to
```terraform
provider "snowflake" {
	organization_name = "ORGANIZATION"
	account_name    = "ACCOUNT"
}
```

If you use TOML configuration file, adjust it from
```toml
[default]
	account = "ORGANIZATION-ACCOUNT"
```

to
```toml
[default]
	organizationname = "ORGANIZATION"
	accountname    = "ACCOUNT"
```

If you use environmental variables, adjust them from
```bash
SNOWFLAKE_ACCOUNT = "ORGANIZATION-ACCOUNT"
```

```bash
SNOWFLAKE_ORGANIZATION_NAME = "ORGANIZATION"
SNOWFLAKE_ACCOUNT_NAME = "ACCOUNT"
```

This change may cause the connection host URL to change. If you get errors like
```
Error: open snowflake connection: Post "https://ORGANIZATION-ACCOUNT.snowflakecomputing.com:443/session/v1/login-request?requestId=[guid]&request_guid=[guid]&roleName=myrole": EOF
```
make sure that the host `ORGANIZATION-ACCOUNT.snowflakecomputing.com` is allowed to be reached from your network (i.e. not blocked by a firewall).

#### *(behavior change)* changed behavior of some fields
For the fields that are not deprecated, we focused on improving validations and documentation. Also, we adjusted some fields to match our [driver's](https://github.com/snowflakedb/gosnowflake) defaults. Specifically:
- Relaxed validations for enum fields like `protocol` and `authenticator`. Now, the case on such fields is ignored.
- `user`, `warehouse`, `role` - added a validation for an account object identifier
- `validate_default_parameters`, `client_request_mfa_token`, `client_store_temporary_credential`, `ocsp_fail_open`,  - to easily handle three-value logic (true, false, unknown) in provider's config, type of these fields was changed from boolean to string. For more details about default values, please refer to the [changes before v1](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.
- `client_ip` - added a validation for an IP address
- `port` - added a validation for a port number
- `okta_url`, `token_accessor.token_endpoint`, `client_store_temporary_credential` - added a validation for a URL address
- `login_timeout`, `request_timeout`, `jwt_expire_timeout`, `client_timeout`, `jwt_client_timeout`, `external_browser_timeout` - added a validation for setting this value to at least `0`
- `authenticator` - added a possibility to configure JWT flow with `SNOWFLAKE_JWT` (formerly, this was supported with `JWT`); the previous value `JWT` was left for compatibility, but will be removed before v1

### *(behavior change)* handling copy_grants
Currently, resources like `snowflake_view`, `snowflake_stream_on_table`, `snowflake_stream_on_external_table` and `snowflake_stream_on_directory_table`  support `copy_grants` field corresponding with `COPY GRANTS` during `CREATE`. The current behavior is that, when a change leading for recreation is detected (meaning a change that can not be handled by ALTER, but only by `CREATE OR REPLACE`), `COPY GRANTS` are used during recreation when `copy_grants` is set to `true`. Changing this field without changes in other field results in a noop because in this case there is no need to recreate a resource.

### *(new feature)* recovering stale streams
Starting from this version, the provider detects stale streams for `snowflake_stream_on_table`, `snowflake_stream_on_external_table` and `snowflake_stream_on_directory_table` and recreates them (optionally with `copy_grants`) to recover them. To handle this correctly, a new computed-only field `stale` has been added to these resource, indicating whether a stream is stale.

### *(new feature)* snowflake_stream_on_directory_table and snowflake_stream_on_view resource
Continuing changes made in [v0.97](#v0960--v0970), the new resource `snowflake_stream_on_directory_table` and `snowflake_stream_on_view` have been introduced to replace the previous `snowflake_stream` for streams on directory tables and streams on views.

To use the new `stream_on_directory_table`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_stage    = snowflake_stage.stage.fully_qualified_name

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_directory_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  stage             = snowflake_stage.stage.fully_qualified_name

  comment = "A stream."
}
```

To use the new `stream_on_view`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_view    = snowflake_view.view.fully_qualified_name

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_view" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  view             = snowflake_view.view.fully_qualified_name

  comment = "A stream."
}
```

Then, follow our [Resource migration guide](./docs/guides/resource_migration.md).

### *(new feature)* Secret resources
Added a new secrets resources for managing secrets.
We decided to split each secret flow into individual resources.
This segregation was based on the secret flows in CREATE SECRET. i.e.:
- `snowflake_secret_with_client_credentials`
- `snowflake_secret_with_authorization_code_grant`
- `snowflake_secret_with_basic_authentication`
- `snowflake_secret_with_generic_string`


See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-secret).

### *(bugfix)* Handle BCR Bundle 2024_08 in snowflake_user resource

[bcr 2024_08](https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_08/bcr-1798) changed the "empty" response in the `SHOW USERS` query. This provider version adapts to the new result types; it should be used if you want to have 2024_08 Bundle enabled on your account.

Note: Because [bcr 2024_07](https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_07/bcr-1692) changes the way how the `default_secondary_roles` attribute behaves, drift may be reported when enabling 2024_08 Bundle. Check [Handling default secondary roles](#breaking-change-handling-default-secondary-roles) for more context.

Connected issues: [#3125](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3125)

### *(bugfix)* Handle user import correctly

#### Context before the change

Password is empty after the `snowflake_user` import; we can't read it from the config or from Snowflake.
During the next terraform plan+apply it's updated to the "same" value.
It results in an error on Snowflake side: `New password rejected by current password policy. Reason: 'PRIOR_USE'.`

#### After the change

The error will be ignored on the provider side (after all, it means that the password in state is the same as on Snowflake side). Still, plan+apply is needed after importing user.

## v0.96.0 ➞ v0.97.0

### *(new feature)* snowflake_stream_on_table, snowflake_stream_on_external_table resource

To enhance clarity and functionality, the new resources `snowflake_stream_on_table` and `snowflake_stream_on_external_table` have been introduced to replace the previous `snowflake_stream`. Recognizing that the old resource carried multiple responsibilities within a single entity, we opted to divide it into more specialized resources.
The newly introduced resources are aligned with the latest Snowflake documentation at the time of implementation, and adhere to our [new conventions](#general-changes).
This segregation was based on the object on which the stream is created. The mapping between SQL statements and the resources is the following:
- `ON TABLE <table_name>` -> `snowflake_stream_on_table`
- `ON EXTERNAL TABLE <external_table_name>` -> `snowflake_stream_on_external_table` (this was previously not supported)

The resources for streams on directory tables and streams on views will be implemented in the future releases.

To use the new `stream_on_table`, change the old `stream` from
```terraform
resource "snowflake_stream" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  on_table    = snowflake_table.table.fully_qualified_name
  append_only = true

  comment = "A stream."
}
```

to

```terraform
resource "snowflake_stream_on_table" "stream" {
  name     = "stream"
  schema   = "schema"
  database = "database"

  table             = snowflake_table.table.fully_qualified_name
  append_only       = "true"

  comment = "A stream."
}
```


Then, follow our [Resource migration guide](./docs/guides/resource_migration.md).

### *(new feature)* new snowflake_service_user and snowflake_legacy_service_user resources

Release v0.95.0 introduced reworked `snowflake_user` resource. As [noted](#note-user-types), the new `SERVICE` and `LEGACY_SERVICE` user types were not supported.

This release introduces two new resources to handle these new user types: `snowflake_service_user` and `snowflake_legacy_service_user`.

Both resources have schemas almost identical to the `snowflake_user` resource with the following exceptions:
- `snowflake_service_user` does not contain the following fields (because they are not supported for the user of type `SERVICE` in Snowflake):
  - `password`
  - `first_name`
  - `middle_name`
  - `last_name`
  - `must_change_password`
  - `mins_to_bypass_mfa`
  - `disable_mfa`
- `snowflake_legacy_service_user` does not contain the following fields (because they are not supported for the user of type `LEGACY_SERVICE` in Snowflake):
  - `first_name`
  - `middle_name`
  - `last_name`
  - `mins_to_bypass_mfa`
  - `disable_mfa`

`snowflake_users` datasource was adjusted to handle different user types and `type` field was added to the `describe_output`.

If you used to manage service or legacy service users through `snowflake_user` resource (e.g. using `lifecycle.ignore_changes`) or `snowflake_unsafe_execute`, please migrate to the new resources following [our guidelines on resource migration](docs/guides/resource_migration.md).

E.g. change the old config from:

```terraform
resource "snowflake_user" "service_user" {
  lifecycle {
    ignore_changes = [user_type]
  }

  name         = "Snowflake Service User"
  login_name   = "service_user"
  email        = "service_user@snowflake.example"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}
```

to

```
resource "snowflake_service_user" "service_user" {
  name         = "Snowflake Service User"
  login_name   = "service_user"
  email        = "service_user@snowflake.example"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."
}

```

Then, follow our [resource migration guide](./docs/guides/resource_migration.md).

Connected issues: [#2951](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2951)

## v0.95.0 ➞ v0.96.0

### snowflake_masking_policies data source changes
New filtering options:
- `in`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `masking_policies` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### snowflake_masking_policy resource changes
New fields:
  - `show_output` field that holds the response from SHOW MASKING POLICIES.
  - `describe_output` field that holds the response from DESCRIBE MASKING POLICY.

#### *(breaking change)* Renamed fields in snowflake_masking_policy resource
Renamed fields:
  - `masking_expression` to `body`
Please rename these fields in your configuration files. State will be migrated automatically.

#### *(breaking change)* Removed fields from snowflake_masking_policy resource
Removed fields:
- `or_replace`
- `if_not_exists`
The value of these field will be removed from the state automatically.

#### *(breaking change)* Adjusted schema of arguments/signature
The field `signature` is renamed to `arguments` to be consistent with other resources.
Now, arguments are stored without nested `column` field. Please adjust that in your configs, like in the example below. State is migrated automatically.

The old configuration looks like this:
```
  signature {
    column {
      name = "val"
      type = "VARCHAR"
    }
  }
```

The new configuration looks like this:
```
  argument {
    name = "val"
    type = "VARCHAR"
  }
```

#### *(breaking change)* Identifiers related changes
During [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`). Importing resources also needs to be adjusted (see [example](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/row_access_policy#import)).

Also, we added diff suppress function that prevents Terraform from showing differences, when only quoting is different.

No change is required, the state will be migrated automatically.

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `exempt_other_policies` was changed from boolean to string.

For more details about default values, please refer to the [changes before v1](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.

### *(breaking change)* resource_monitor resource
Removed fields:
- `set_for_account` (will be settable on account resource, right now, the preferred way is to set it through unsafe_execute resource)
- `warehouses` (can be set on warehouse resource, optionally through unsafe_execute resource only if the warehouse is not managed by Terraform)
- `suspend_triggers` (now, `suspend_trigger` should be used)
- `suspend_immediate_triggers` (now, `suspend_immediate_trigger` should be used)

### *(breaking change)* resource_monitor data source
Changes:
- New filtering option `like`
- Now, the output of `SHOW RESOURCE MONITORS` is now inside `resource_monitors.*.show_output`. Here's the list of currently available fields:
    - `name`
    - `credit_quota`
    - `used_credits`
    - `remaining_credits`
    - `level`
    - `frequency`
    - `start_time`
    - `end_time`
    - `suspend_at`
    - `suspend_immediate_at`
    - `created_on`
    - `owner`
    - `comment`

### snowflake_row_access_policies data source changes
New filtering options:
- `in`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `row_access_policies` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

Please adjust your Terraform configuration files.

### snowflake_row_access_policy resource changes
New fields:
  - `show_output` field that holds the response from SHOW ROW ACCESS POLICIES.
  - `describe_output` field that holds the response from DESCRIBE ROW ACCESS POLICY.

#### *(breaking change)* Renamed fields in snowflake_row_access_policy resource
Renamed fields:
  - `row_access_expression` to `body`
Please rename these fields in your configuration files. State will be migrated automatically.

#### *(breaking change)* Adjusted schema of arguments/signature
The field `signature` is renamed to `arguments` to be consistent with other resources.
Now, arguments are stored as a list, instead of a map. Please adjust that in your configs. State is migrated automatically. Also, this means that order of the items matters and may be adjusted.


The old configuration looks like this:
```
  signature = {
    A = "VARCHAR",
    B = "VARCHAR"
  }
```

The new configuration looks like this:
```
  argument {
    name = "A"
    type = "VARCHAR"
  }
  argument {
    name = "B"
    type = "VARCHAR"
  }
```

Argument names are now case sensitive. All policies created previously in the provider have upper case argument names. If you used lower case before, please adjust your configs. Values in the state will be migrated to uppercase automatically.

#### *(breaking change)* Adjusted behavior on changing name
Previously, after changing `name` field, the resource was recreated. Now, the object is renamed with `RENAME TO`.

#### *(breaking change)* Mitigating permadiff on `body`
Previously, `body` of a policy was compared as a raw string. This led to permament diff because of leading newlines (see https://github.com/snowflakedb/terraform-provider-snowflake/issues/2053).

Now, similarly to handling statements in other resources, we replace blank characters with a space. The provider can cause false positives in cases where a change in case or run of whitespace is semantically significant.

#### *(breaking change)* Identifiers related changes
During [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`). Importing resources also needs to be adjusted (see [example](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/row_access_policy#import)).

Also, we added diff suppress function that prevents Terraform from showing differences, when only quoting is different.

No change is required, the state will be migrated automatically.

## v0.94.x ➞ v0.95.0

### *(breaking change)* database roles data source; field rename, schema structure changes, and adding missing filtering options

- `database` renamed to `in_database`
- Added `like` and `limit` filtering options
- `SHOW DATABASE ROLES` output is now put inside `database_roles.*.show_output`. Here's the list of currently available fields:
    - `created_on`
    - `name`
    - `is_default`
    - `is_current`
    - `is_inherited`
    - `granted_to_roles`
    - `granted_to_database_roles`
    - `granted_database_roles`
    - `owner`
    - `comment`
    - `owner_role_type`

### snowflake_views data source changes
New filtering options:
- `in`
- `like`
- `starts_with`
- `limit`
- `with_describe`

New output fields
- `show_output`
- `describe_output`

Breaking changes:
- `database` and `schema` are right now under `in` field
- `views` field now organizes output of show under `show_output` field and the output of describe under `describe_output` field.

### snowflake_view resource changes
New fields:
  - `row_access_policy`
  - `aggregation_policy`
  - `change_tracking`
  - `is_recursive`
  - `is_temporary`
  - `data_metric_schedule`
  - `data_metric_function`
  - `column`
- added `show_output` field that holds the response from SHOW VIEWS.
- added `describe_output` field that holds the response from DESCRIBE VIEW. Note that one needs to grant sufficient privileges e.g. with [grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership) on the tables used in this view. Otherwise, this field is not filled.

#### *(breaking change)* Removed fields from snowflake_view resource
Removed fields:
- `or_replace` - `OR REPLACE` is added by the provider automatically when `copy_grants` is set to `"true"`
- `tag` - Please, use [tag_association](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/tag_association) instead.
The value of these field will be removed from the state automatically.

#### *(breaking change)* Required warehouse
For this resource, the provider now uses [policy references](https://docs.snowflake.com/en/sql-reference/functions/policy_references) which requires a warehouse in the connection. Please, make sure you have either set a `DEFAULT_WAREHOUSE` for the user, or specified a warehouse in the provider configuration.

### Identifier changes

#### *(breaking change)* resource identifiers for schema and streamlit
During [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework) we decided to
migrate resource ids from pipe-separated to regular Snowflake identifiers (e.g. `<database_name>|<schema_name>` -> `"<database_name>"."<schema_name>"`).
Exception to that rule will be identifiers that consist of multiple parts (like in the case of [grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role#import)'s resource id).
The change was applied to already refactored resources (only in the case of `snowflake_schema` and `snowflake_streamlit` this will be a breaking change, because the rest of the objects are single part identifiers in the format of `<name>`):
- `snowflake_api_authentication_integration_with_authorization_code_grant`
- `snowflake_api_authentication_integration_with_client_credentials`
- `snowflake_api_authentication_integration_with_jwt_bearer`
- `snowflake_oauth_integration_for_custom_clients`
- `snowflake_oauth_integration_for_partner_applications`
- `snowflake_external_oauth_integration`
- `snowflake_saml2_integration`
- `snowflake_scim_integration`
- `snowflake_database`
- `snowflake_shared_database`
- `snowflake_secondary_database`
- `snowflake_account_role`
- `snowflake_network_policy`
- `snowflake_warehouse`

No change is required, the state will be migrated automatically.
The rest of the objects will be changed when working on them during [v1 object preparations](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1).

#### *(breaking change)* diff suppress for identifier quoting
(The same set of resources listed above was adjusted)
To prevent issues like [this one](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2982), we added diff suppress function that prevents Terraform from showing differences,
when only quoting is different. In some cases, Snowflake output (mostly from SHOW commands) was dictating which field should be additionally quoted and which shouldn't, but that should no longer be the case.
Like in the change above, the rest of the objects will be changed when working on them during [v1 object preparations](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1).

### New `fully_qualified_name` field in the resources.
We added a new `fully_qualified_name` to snowflake resources. This should help with referencing other resources in fields that expect a fully qualified name. For example, instead of
writing

```object_name = “\”${snowflake_table.database}\”.\”${snowflake_table.schema}\”.\”${snowflake_table.name}\””```

 now we can write

```object_name = snowflake_table.fully_qualified_name```

See more details in [identifiers guide](./docs/guides/identifiers.md#new-computed-fully-qualified-name-field-in-resources).

See [example usage](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role).

Some of the resources are excluded from this change:
- deprecated resources
  - `snowflake_database_old`
  - `snowflake_oauth_integration`
  - `snowflake_saml_integration`
- resources for which fully qualified name is not appropriate
  - `snowflake_account_parameter`
  - `snowflake_account_password_policy_attachment`
  - `snowflake_network_policy_attachment`
  - `snowflake_session_parameter`
  - `snowflake_table_constraint`
  - `snowflake_table_column_masking_policy_application`
  - `snowflake_tag_masking_policy_association`
  - `snowflake_tag_association`
  - `snowflake_user_password_policy_attachment`
  - `snowflake_user_public_keys`
  - grant resources

#### *(breaking change)* removed `qualified_name` from `snowflake_masking_policy`, `snowflake_network_rule`, `snowflake_password_policy` and `snowflake_table`
Because of introducing a new `fully_qualified_name` field for all of the resources, `qualified_name` was removed from `snowflake_masking_policy`, `snowflake_network_rule`,  `snowflake_password_policy` and `snowflake_table`. Please adjust your configurations. State is automatically migrated.

### snowflake_stage resource changes

#### *(bugfix)* Correctly handle renamed/deleted stage

Correctly handle the situation when stage was rename/deleted externally (earlier it resulted in a permanent loop). No action is required on the user's side.

Connected issues: [#2972](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2972)

### snowflake_table resource changes

#### *(bugfix)* Handle data type diff suppression better for text and number

Data types are not entirely correctly handled inside the provider (read more e.g. in [#2735](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2735)). It will be still improved with the upcoming function, procedure, and table rework. Currently, diff suppression was fixed for text and number data types in the table resource with the following assumptions/limitations:
- for numbers the default precision is 38 and the default scale is 0 (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-numeric#number))
- for number types the following types are treated as synonyms: `NUMBER`, `DECIMAL`, `NUMERIC`, `INT`, `INTEGER`, `BIGINT`, `SMALLINT`, `TINYINT`, `BYTEINT`
- for text the default length is 16777216 (following the [docs](https://docs.snowflake.com/en/sql-reference/data-types-text#varchar))
- for text types the following types are treated as synonyms: `VARCHAR`, `CHAR`, `CHARACTER`, `STRING`, `TEXT`
- whitespace and casing is ignored
- if the type arguments cannot be parsed the defaults are used and therefore diff may be suppressed unexpectedly (please report such cases)

No action is required on the user's side.

Connected issues: [#3007](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3007)

### snowflake_user resource changes

Because of the multiple changes in the resource, the easiest migration way is to follow our [migration guide](./docs/guides/resource_migration.md) to perform zero downtime migration. Alternatively, it is possible to follow some pointers below. Either way, familiarize yourself with the resource changes before version bumping. Also, check the [design decisions](./v1-preparations/CHANGES_BEFORE_V1.md).

#### *(breaking change)* user parameters added to snowflake_user resource

On our road to V1 we changed the approach to Snowflake parameters on the object level; now, we add them directly to the resource. This is a **breaking change** because now:
- Leaving the config empty does not set the default value on the object level but uses the one from hierarchy on Snowflake level instead (so after version bump, the diff running `UNSET` statements is expected).
- This change is not compatible with `snowflake_object_parameter` - you have to set the parameter inside `snowflake_user` resource **IF** you manage users through terraform **AND** you want to set the parameter on the user level.

For more details, check the [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters).

The following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added to the `snowflake_user` resource:
 - [ABORT_DETACHED_QUERY](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query)
 - [AUTOCOMMIT](https://docs.snowflake.com/en/sql-reference/parameters#autocommit)
 - [BINARY_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#binary-input-format)
 - [BINARY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#binary-output-format)
 - [CLIENT_MEMORY_LIMIT](https://docs.snowflake.com/en/sql-reference/parameters#client-memory-limit)
 - [CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX](https://docs.snowflake.com/en/sql-reference/parameters#client-metadata-request-use-connection-ctx)
 - [CLIENT_PREFETCH_THREADS](https://docs.snowflake.com/en/sql-reference/parameters#client-prefetch-threads)
 - [CLIENT_RESULT_CHUNK_SIZE](https://docs.snowflake.com/en/sql-reference/parameters#client-result-chunk-size)
 - [CLIENT_RESULT_COLUMN_CASE_INSENSITIVE](https://docs.snowflake.com/en/sql-reference/parameters#client-result-column-case-insensitive)
 - [CLIENT_SESSION_KEEP_ALIVE](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive)
 - [CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY](https://docs.snowflake.com/en/sql-reference/parameters#client-session-keep-alive-heartbeat-frequency)
 - [CLIENT_TIMESTAMP_TYPE_MAPPING](https://docs.snowflake.com/en/sql-reference/parameters#client-timestamp-type-mapping)
 - [DATE_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#date-input-format)
 - [DATE_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#date-output-format)
 - [ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION](https://docs.snowflake.com/en/sql-reference/parameters#enable-unload-physical-type-optimization)
 - [ERROR_ON_NONDETERMINISTIC_MERGE](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-merge)
 - [ERROR_ON_NONDETERMINISTIC_UPDATE](https://docs.snowflake.com/en/sql-reference/parameters#error-on-nondeterministic-update)
 - [GEOGRAPHY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#geography-output-format)
 - [GEOMETRY_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#geometry-output-format)
 - [JDBC_TREAT_DECIMAL_AS_INT](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-decimal-as-int)
 - [JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-timestamp-ntz-as-utc)
 - [JDBC_USE_SESSION_TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-use-session-timezone)
 - [JSON_INDENT](https://docs.snowflake.com/en/sql-reference/parameters#json-indent)
 - [LOCK_TIMEOUT](https://docs.snowflake.com/en/sql-reference/parameters#lock-timeout)
 - [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#log-level)
 - [MULTI_STATEMENT_COUNT](https://docs.snowflake.com/en/sql-reference/parameters#multi-statement-count)
 - [NOORDER_SEQUENCE_AS_DEFAULT](https://docs.snowflake.com/en/sql-reference/parameters#noorder-sequence-as-default)
 - [ODBC_TREAT_DECIMAL_AS_INT](https://docs.snowflake.com/en/sql-reference/parameters#odbc-treat-decimal-as-int)
 - [QUERY_TAG](https://docs.snowflake.com/en/sql-reference/parameters#query-tag)
 - [QUOTED_IDENTIFIERS_IGNORE_CASE](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case)
 - [ROWS_PER_RESULTSET](https://docs.snowflake.com/en/sql-reference/parameters#rows-per-resultset)
 - [S3_STAGE_VPCE_DNS_NAME](https://docs.snowflake.com/en/sql-reference/parameters#s3-stage-vpce-dns-name)
 - [SEARCH_PATH](https://docs.snowflake.com/en/sql-reference/parameters#search-path)
 - [SIMULATED_DATA_SHARING_CONSUMER](https://docs.snowflake.com/en/sql-reference/parameters#simulated-data-sharing-consumer)
 - [STATEMENT_QUEUED_TIMEOUT_IN_SECONDS](https://docs.snowflake.com/en/sql-reference/parameters#statement-queued-timeout-in-seconds)
 - [STATEMENT_TIMEOUT_IN_SECONDS](https://docs.snowflake.com/en/sql-reference/parameters#statement-timeout-in-seconds)
 - [STRICT_JSON_OUTPUT](https://docs.snowflake.com/en/sql-reference/parameters#strict-json-output)
 - [TIMESTAMP_DAY_IS_ALWAYS_24H](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-day-is-always-24h)
 - [TIMESTAMP_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-input-format)
 - [TIMESTAMP_LTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ltz-output-format)
 - [TIMESTAMP_NTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-ntz-output-format)
 - [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-output-format)
 - [TIMESTAMP_TYPE_MAPPING](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-type-mapping)
 - [TIMESTAMP_TZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#timestamp-tz-output-format)
 - [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#timezone)
 - [TIME_INPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#time-input-format)
 - [TIME_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#time-output-format)
 - [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#trace-level)
 - [TRANSACTION_ABORT_ON_ERROR](https://docs.snowflake.com/en/sql-reference/parameters#transaction-abort-on-error)
 - [TRANSACTION_DEFAULT_ISOLATION_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#transaction-default-isolation-level)
 - [TWO_DIGIT_CENTURY_START](https://docs.snowflake.com/en/sql-reference/parameters#two-digit-century-start)
 - [UNSUPPORTED_DDL_ACTION](https://docs.snowflake.com/en/sql-reference/parameters#unsupported-ddl-action)
 - [USE_CACHED_RESULT](https://docs.snowflake.com/en/sql-reference/parameters#use-cached-result)
 - [WEEK_OF_YEAR_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#week-of-year-policy)
 - [WEEK_START](https://docs.snowflake.com/en/sql-reference/parameters#week-start)
 - [ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR](https://docs.snowflake.com/en/sql-reference/parameters#enable-unredacted-query-syntax-error)
 - [NETWORK_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#network-policy)
 - [PREVENT_UNLOAD_TO_INTERNAL_STAGES](https://docs.snowflake.com/en/sql-reference/parameters#prevent-unload-to-internal-stages)

Connected issues: [#2938](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2938)

#### *(breaking change)* Changes in sensitiveness of name, login_name, and display_name

According to https://docs.snowflake.com/en/sql-reference/functions/all_user_names#usage-notes, `NAME`s are not considered sensitive data and `LOGIN_NAME`s are. Previous versions of the provider had this the other way around. In this version, `name` attribute was unmarked as sensitive, whereas `login_name` was marked as sensitive. This may break your configuration if you were using `login_name`s before e.g. in a `for_each` loop.

The `display_name` attribute was marked as sensitive. It defaults to `name` if not provided on Snowflake side. Because `name` is no longer sensitive, we also change the setting for the `display_name`.

Connected issues: [#2662](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2662), [#2668](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2668).

#### *(bugfix)* Correctly handle `default_warehouse`, `default_namespace`, and `default_role`

During the [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework), we generalized how we compute the differences correctly for the identifier fields (read more in [this document](./docs/guides/identifiers_rework_design_decisions.md)). Proper suppressor was applied to `default_warehouse`, `default_namespace`, and `default_role`. Also, all these three attributes were corrected (e.g. handling spaces/hyphens in names).

Connected issues: [#2836](https://github.com/snowflakedb/terraform-provider-snowflake/pull/2836), [#2942](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2942)

#### *(bugfix)* Correctly handle failed update

Not every attribute can be updated in the state during read (like `password` in the `snowflake_user` resource). In situations where update fails, we may end up with an incorrect state (read more in https://github.com/hashicorp/terraform-plugin-sdk/issues/476). We use a deprecated method from the plugin SDK, and now, for partially failed updates, we preserve the resource's previous state. It fixed this kind of situations for `snowflake_user` resource.

Connected issues: [#2970](https://github.com/snowflakedb/terraform-provider-snowflake/pull/2970)

#### *(breaking change)* Handling default secondary roles

Old field `default_secondary_roles` was removed in favour of the new, easier, `default_secondary_roles_option` because the only possible options that can be currently set are `('ALL')` and `()`.  The logic to handle set element changes was convoluted and error-prone. Additionally, [bcr 2024_07](https://docs.snowflake.com/en/release-notes/bcr-bundles/2024_07/bcr-1692) complicated the matter even more.

Now:
- the default value is `DEFAULT` - it falls back to Snowflake default (so `()` before and `('ALL')` after the BCR)
- to explicitly set to `('ALL')` use `ALL`
- to explicitly set to `()` use `NONE`

While migrating, the old `default_secondary_roles` will be removed from the state automatically and `default_secondary_roles_option` will be constructed based on the previous value (in some cases apply may be necessary).

Connected issues: [#3038](https://github.com/snowflakedb/terraform-provider-snowflake/pull/3038)

#### *(breaking change)* Attributes changes

Attributes that are no longer computed:
- `login_name`
- `display_name`
- `disabled`
- `default_role`

New fields:
- `middle_name`
- `days_to_expiry`
- `mins_to_unlock`
- `mins_to_bypass_mfa`
- `disable_mfa`
- `default_secondary_roles_option`
- `show_output` - holds the response from `SHOW USERS`. Remember that the field will be only recomputed if one of the user attributes is changed.
- `parameters` - holds the response from `SHOW PARAMETERS IN USER`.

Removed fields:
- `has_rsa_public_key`
- `default_secondary_roles` - replaced with `default_secondary_roles_option`

Default changes:
- `must_change_password`
- `disabled`

Type changes:
- `must_change_password`: bool -> string (To easily handle three-value logic (true, false, unknown) in provider's configs, read more in https://github.com/snowflakedb/terraform-provider-snowflake/blob/751239b7d2fee4757471db6c03b952d4728ee099/v1-preparations/CHANGES_BEFORE_V1.md?plain=1#L24)
- `disabled`: bool -> string (To easily handle three-value logic (true, false, unknown) in provider's configs, read more in https://github.com/snowflakedb/terraform-provider-snowflake/blob/751239b7d2fee4757471db6c03b952d4728ee099/v1-preparations/CHANGES_BEFORE_V1.md?plain=1#L24)

#### *(breaking change)* refactored snowflake_users datasource
> **IMPORTANT NOTE:** when querying users you don't have permissions to, the querying options are limited.
You won't get almost any field in `show_output` (only empty or default values), the DESCRIBE command will return error when called, so you have to set `with_describe = false`; the SHOW PARAMETERS command will return error if called too, so you have to set `with_parameters = false`.

Changes:
- account checking logic was entirely removed
- `pattern` renamed to `like`
- `like`, `starts_with`, and `limit` filters added
- `SHOW USERS` output is enclosed in `show_output` field inside `users` (all the previous fields in `users` map were removed)
- Added outputs from **DESC USER** and **SHOW PARAMETERS IN USER** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC USER** (with `with_describe` turned on) and **SHOW PARAMETERS IN USER** (with `with_parameters` turned on) **per user** returned by **SHOW USERS**.
  The outputs of both commands are held in `users` entry, where **DESC USER** is saved in the `describe_output` field, and **SHOW PARAMETERS IN USER** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

Connected issues: [#2902](https://github.com/snowflakedb/terraform-provider-snowflake/pull/2902)

#### *(breaking change)* snowflake_user_public_keys usage with snowflake_user

`snowflake_user_public_keys` is a resource allowing to set keys for the given user. Before this version, it was possible to have `snowflake_user` and `snowflake_user_public_keys` used next to each other.
Because the logic handling the keys in `snowflake_user` was fixed, it is advised to use `snowflake_user_public_keys` only when user is not managed through terraform. Having both resources configured for the same user will result in improper behavior.

To migrate, in case of having two resources:
- copy the keys to `rsa_public_key` and `rsa_public_key2` in `snowflake_user`
- remove `snowflake_user_public_keys` from state (following [Resource migration guide](./docs/guides/resource_migration.md#resource-migration))
- remove `snowflake_user_public_keys` from config

#### *(breaking change)* snowflake_network_policy_attachment usage with snowflake_user

`snowflake_network_policy_attachment` changes are similar to the changes to `snowflake_user_public_keys` above. It is advised to use `snowflake_network_policy_attachment` only when user is not managed through terraform. Having both resources configured for the same user will result in improper behavior.

To migrate, in case of having two resources:
- copy network policy to [network_policy](https://registry.terraform.io/providers/snowflakedb/snowflake/0.95.0/docs/resources/user#network_policy) attribute in the `snowflake_user` resource
- remove `snowflake_network_policy_attachment` from state (following [Resource migration guide](./docs/guides/resource_migration.md#resource-migration))
- remove `snowflake_network_policy_attachment` from config

References: [#3048](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/3048), [#3058](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3058)

#### *(note)* snowflake_user_password_policy_attachment and other user policies

`snowflake_user_password_policy_attachment` is not addressed in the current version.
Attaching other user policies is not addressed in the current version.

Both topics will be addressed in the following versions.

#### *(note)* user types

`service` and `legacy_service` user types are currently not supported. They will be supported in the following versions as separate resources (namely `snowflake_service_user` and `snowflake_legacy_service_user`).

If you used the existing `snowflake_user` and altered its type externally (manually or through `snowflake_unsafe_execute`), then after migrating to v0.95.0 the provider will try to recreate it as a `person` type.

Because `snowflake_service_user` and `snowflake_legacy_service_user` resources are available in v0.97.0 version, you can temporarily suppress these changes to allow version-by-version migration. To do that, use [`ignore_changes`](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) meta-attribute:

```hcl
resource "snowflake_user" "example_user" {
  # ...
  lifecycle {
    ignore_changes = [ user_type ]
  }
}
```

## v0.94.0 ➞ v0.94.1
### changes in snowflake_schema

In order to avoid dropping `PUBLIC` schemas, we have decided to use `ALTER` instead of `OR REPLACE` during creation. In the future we are planning to use `CREATE OR ALTER` when it becomes available for schems.

## v0.93.0 ➞ v0.94.0
### *(breaking change)* changes in snowflake_scim_integration

In order to fix issues in v0.93.0, when a resource has Azure scim client, `sync_password` field is now set to `default` value in the state. State will be migrated automatically.

### *(breaking change)* refactored snowflake_schema resource

Renamed fields:
- renamed `is_managed` to `with_managed_access`
- renamed `data_retention_days` to `data_retention_time_in_days`

Please rename these fields in your configuration files. State will be migrated automatically.

Removed fields:
- `tag`
The value of this field will be removed from the state automatically. Please, use [tag_association](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/tag_association) instead.

New fields:
- the following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added:
    - `max_data_extension_time_in_days`
    - `external_volume`
    - `catalog`
    - `replace_invalid_characters`
    - `default_ddl_collation`
    - `storage_serialization_policy`
    - `log_level`
    - `trace_level`
    - `suspend_task_after_num_failures`
    - `task_auto_retry_attempts`
    - `user_task_managed_initial_warehouse_size`
    - `user_task_timeout_ms`
    - `user_task_minimum_trigger_interval_in_seconds`
    - `quoted_identifiers_ignore_case`
    - `enable_console_output`
    - `pipe_execution_paused`
- added `show_output` field that holds the response from SHOW SCHEMAS.
- added `describe_output` field that holds the response from DESCRIBE SCHEMA. Note that one needs to grant sufficient privileges e.g. with [grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership) on all objects in the schema. Otherwise, this field is not filled.
- added `parameters` field that holds the response from SHOW PARAMETERS IN SCHEMA.

We allow creating and managing `PUBLIC` schemas now. When the name of the schema is `PUBLIC`, it's created with `OR_REPLACE`. Please be careful with this operation, because you may experience data loss. `OR_REPLACE` does `DROP` before `CREATE`, so all objects in the schema will be dropped and this is not visible in Terraform plan. To restore data-related objects that might have been accidentally or intentionally deleted, pleas read about [Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel). The alternative is to import `PUBLIC` schema manually and then manage it with Terraform. We've decided this based on [#2826](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2826).

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `is_transient` and `with_managed_access` was changed from boolean to string.

Terraform should recreate resources for configs lacking `is_transient` (`DROP` and then `CREATE` will be run underneath). To prevent this behavior, please set the `is_transient` field to the desired value (`"true"` for transient schemas, `"false"` for non-transient ones).
For more details about default values, please refer to the [changes before v1](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#default-values) document.

Terraform should perform an action for configs lacking `with_managed_access` (`ALTER SCHEMA DISABLE MANAGED ACCESS` will be run underneath which should not affect the Snowflake object, because `MANAGED ACCESS` is not set by default)
### *(breaking change)* refactored snowflake_schemas datasource
Changes:
- `database` is removed and can be specified inside `in` field.
- `like`, `in`, `starts_with`, and `limit` fields enable filtering.
- SHOW SCHEMAS output is enclosed in `show_output` field inside `schemas`.
- Added outputs from **DESC SCHEMA** and **SHOW PARAMETERS IN SCHEMA** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC SCHEMA** (with `with_describe` turned on) and **SHOW PARAMETERS IN SCHEMA** (with `with_parameters` turned on) **per schema** returned by **SHOW SCHEMAS**.
  The outputs of both commands are held in `schemas` entry, where **DESC SCHEMA** is saved in the `describe_output` field, and **SHOW PARAMETERS IN SCHEMA** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(new feature)* new snowflake_account_role resource

Already existing `snowflake_role` was deprecated in favor of the new `snowflake_account_role`. The old resource got upgraded to
have the same features as the new one. The only difference is the deprecation message on the old resource.

New fields:
- added `show_output` field that holds the response from SHOW ROLES. Remember that the field will be only recomputed if one of the fields (`name` or `comment`) are changed.

### *(breaking change)* refactored snowflake_roles data source

Changes:
- New `in_class` filtering option to filter out roles by class name, e.g. `in_class = "SNOWFLAKE.CORE.BUDGET"`
- `pattern` was renamed to `like`
- output of SHOW is enclosed in `show_output`, so before, e.g. `roles.0.comment` is now `roles.0.show_output.0.comment`

### *(new feature)* snowflake_streamlit resource
Added a new resource for managing streamlits. See reference [docs](https://docs.snowflake.com/en/sql-reference/sql/create-streamlit). In this resource, we decided to split `ROOT_LOCATION` in Snowflake to two fields: `stage` representing stage fully qualified name and `directory_location` containing a path within this stage to root location.

### *(new feature)* snowflake_streamlits datasource
Added a new datasource enabling querying and filtering stremlits. Notes:
- all results are stored in `streamlits` field.
- `like`, `in`, and `limit` fields enable streamlits filtering.
- SHOW STREAMLITS output is enclosed in `show_output` field inside `streamlits`.
- Output from **DESC STREAMLIT** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `streamlits`.
  The additional parameters call **DESC STREAMLIT** (with `with_describe` turned on) **per streamlit** returned by **SHOW STREAMLITS**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(new feature)* refactored snowflake_network_policy resource

No migration required.

New behavior:
- `name` is no longer marked as ForceNew parameter. When changed, now it will perform ALTER RENAME operation, instead of re-creating with the new name.
- Additional validation was added to `blocked_ip_list` to inform about specifying `0.0.0.0/0` ip. More details in the [official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-network-policy#usage-notes).

New fields:
- `show_output` and `describe_output` added to hold the results returned by `SHOW` and `DESCRIBE` commands. Those fields will only be recomputed when specified fields change

### *(new feature)* snowflake_network_policies datasource

Added a new datasource enabling querying and filtering network policies. Notes:
- all results are stored in `network_policies` field.
- `like` field enables filtering.
- SHOW NETWORK POLICIES output is enclosed in `show_output` field inside `network_policies`.
- Output from **DESC NETWORK POLICY** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `network_policies`.
  The additional parameters call **DESC NETWORK POLICY** (with `with_describe` turned on) **per network policy** returned by **SHOW NETWORK POLICIES**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### *(fix)* snowflake_warehouse resource

Because of the issue [#2948](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2948), we are relaxing the validations for the Snowflake parameter values. Read more in [CHANGES_BEFORE_V1.md](v1-preparations/CHANGES_BEFORE_V1.md#validations).

## v0.92.0 ➞ v0.93.0

### general changes

With this change we introduce the first resources redesigned for the V1. We have made a few design choices that will be reflected in these and in the further reworked resources. This includes:
- Handling the [default values](./v1-preparations/CHANGES_BEFORE_V1.md#default-values).
- Handling the ["empty" values](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values).
- Handling the [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters).
- Saving the [config values in the state](./v1-preparations/CHANGES_BEFORE_V1.md#config-values-in-the-state).
- Providing a ["raw Snowflake output"](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values) for the managed resources.

They are all described in short in the [changes before v1 doc](./v1-preparations/CHANGES_BEFORE_V1.md). Please familiarize yourself with these changes before the upgrade.

### old grant resources removal
Following the [announcement](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2736) we have removed the old grant resources.
The two resources [snowflake_role_ownership_grant](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/role_ownership_grant) and
[snowflake_user_ownership_grant](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/user_ownership_grant) were not listed in the announcement,
but they were also marked as deprecated ones. We are removing them too to conclude the grants redesign saga.

#### Grant resource mappings
As previous resources had multiple responsibilities within a single entity, we opted to divide them into more specialized resources.
Because of that, they (mostly) cannot be mapped one to one. To migrate the old grant resources to the new ones, use the following mapping rules:
- If you are using `shares` field, use the [snowflake_grant_privileges_to_share](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_share) resource.
- If you are using `roles` field, use the [snowflake_grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role) resource.
- For OWNERSHIP privilege manipulation, use the dedicated [snowflake_grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership) resource.
- New grant resources expect fully qualified identifiers (e.g., [grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role#object_name-2) with `object_name`), instead of parts of identifier split between few fields (e.g., [function_grant](https://registry.terraform.io/providers/snowflakedb/snowflake/0.90.0/docs/resources/function_grant#database_name-8) with `database_name`, `schema_name`, `function_name`, and `argument_data_types`). Use `fully_qualified_name` field whenever possible to avoid identifier issues (look [here](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources)).
- The `enable_masking_grants` field that could be found, for example, in [masking_policy_grant](https://registry.terraform.io/providers/snowflakedb/snowflake/0.90.0/docs/resources/masking_policy_grant) resource is now "enabled by default" for all new grant resources. In the new resources, there's no way to disable this setting. We left it as a [future topic](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#future-topics) if such need arises.
- The `on_all` and `on_future` fields were transformed from [boolean type](https://registry.terraform.io/providers/snowflakedb/snowflake/0.90.0/docs/resources/schema_grant#on_all-9) to [nested object](https://registry.terraform.io/providers/snowflakedb/snowflake/0.90.0/docs/resources/schema_grant#on_all-9). Now, it reflects more the [GRANT PRIVILEGE](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege) documentation.
- The `revert_ownership_to_role_name` option that was available in the [user_ownership_grant](https://registry.terraform.io/providers/snowflakedb/snowflake/0.90.0/docs/resources/user_ownership_grant#revert_ownership_to_role_name-24) is not supported in the new [grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership) resource, but [we have it in our plans](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grant_ownership_resource_overview#future-plans) as an improvement.

> Because new resources are more general and promote single responsibility, you may end up with higher resource count than before (at least "under the hood," because new resources are easier to work with `for_each` can be easily generated for different configurations).
> We plan to improve this topic in the future, as it is a highly requested feature (mostly because resource count is usually a measure used for billing).

### *(new feature)* Api authentication resources
Added new api authentication resources, i.e.:
- `snowflake_api_authentication_integration_with_authorization_code_grant`
- `snowflake_api_authentication_integration_with_client_credentials`
- `snowflake_api_authentication_integration_with_jwt_bearer`

See reference [doc](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth).

### *(new feature)* snowflake_oauth_integration_for_custom_clients and snowflake_oauth_integration_for_partner_applications resources

To enhance clarity and functionality, the new resources `snowflake_oauth_integration_for_custom_clients` and `snowflake_oauth_integration_for_partner_applications` have been introduced
to replace the previous `snowflake_oauth_integration`. Recognizing that the old resource carried multiple responsibilities within a single entity, we opted to divide it into two more specialized resources.
The newly introduced resources are aligned with the latest Snowflake documentation at the time of implementation, and adhere to our [new conventions](#general-changes).
This segregation was based on the `oauth_client` attribute, where `CUSTOM` corresponds to `snowflake_oauth_integration_for_custom_clients`,
while other attributes align with `snowflake_oauth_integration_for_partner_applications`.

### *(new feature)* snowflake_security_integrations datasource
Added a new datasource enabling querying and filtering all types of security integrations. Notes:
- all results are stored in `security_integrations` field.
- `like` field enables security integrations filtering.
- SHOW SECURITY INTEGRATIONS output is enclosed in `show_output` field inside `security_integrations`.
- Output from **DESC SECURITY INTEGRATION** (which can be turned off by declaring `with_describe = false`, **it's turned on by default**) is enclosed in `describe_output` field inside `security_integrations`.
  **DESC SECURITY INTEGRATION** returns different properties based on the integration type. Consult the documentation to check which ones will be filled for which integration.
  The additional parameters call **DESC SECURITY INTEGRATION** (with `with_describe` turned on) **per security integration** returned by **SHOW SECURITY INTEGRATIONS**.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

### snowflake_external_oauth_integration resource changes

#### *(behavior change)* Renamed fields
Renamed fields:
- `type` to `external_oauth_type`
- `issuer` to `external_oauth_issuer`
- `token_user_mapping_claims` to `external_oauth_token_user_mapping_claim`
- `snowflake_user_mapping_attribute` to `external_oauth_snowflake_user_mapping_attribute`
- `scope_mapping_attribute` to `external_oauth_scope_mapping_attribute`
- `jws_keys_urls` to `external_oauth_jws_keys_url`
- `rsa_public_key` to `external_oauth_rsa_public_key`
- `rsa_public_key_2` to `external_oauth_rsa_public_key_2`
- `blocked_roles` to `external_oauth_blocked_roles_list`
- `allowed_roles` to `external_oauth_allowed_roles_list`
- `audience_urls` to `external_oauth_audience_list`
- `any_role_mode` to `external_oauth_any_role_mode`
- `scope_delimiter` to `external_oauth_scope_delimiter`
to align with Snowflake docs. Please rename this field in your configuration files. State will be migrated automatically.

#### *(behavior change)* Force new for multiple attributes after removing from config
Conditional force new was added for the following attributes when they are removed from config. There are no alter statements supporting UNSET on these fields.
- `external_oauth_rsa_public_key`
- `external_oauth_rsa_public_key_2`
- `external_oauth_scope_mapping_attribute`
- `external_oauth_jws_keys_url`
- `external_oauth_token_user_mapping_claim`

#### *(behavior change)* Conflicting fields
Fields listed below can not be set at the same time in Snowflake. They are marked as conflicting fields.
- `external_oauth_jws_keys_url` <-> `external_oauth_rsa_public_key`
- `external_oauth_jws_keys_url` <-> `external_oauth_rsa_public_key_2`
- `external_oauth_allowed_roles_list` <-> `external_oauth_blocked_roles_list`

#### *(behavior change)* Changed diff suppress for some fields
The fields listed below had diff suppress which removed '-' from strings. Now, this behavior is removed, so if you had '-' in these strings, please remove them. Note that '-' in these values is not allowed by Snowflake.
- `external_oauth_snowflake_user_mapping_attribute`
- `external_oauth_type`
- `external_oauth_any_role_mode`

### *(new feature)* snowflake_saml2_integration resource

The new `snowflake_saml2_integration` is introduced and deprecates `snowflake_saml_integration`. It contains new fields
and follows our new conventions making it more stable. The old SAML integration wasn't changed, so no migration needed,
but we recommend to eventually migrate to the newer counterpart.

### snowflake_scim_integration resource changes
#### *(behavior change)* Changed behavior of `sync_password`

Now, the `sync_password` field will set the state value to `default` whenever the value is not set in the config. This indicates that the value on the Snowflake side is set to the Snowflake default.

> [!WARNING]
> This change causes issues for Azure scim client (see [#2946](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2946)). The workaround is to remove the resource from the state with `terraform state rm`, add `sync_password = true` to the config, and import with `terraform import "snowflake_scim_integration.test" "aad_provisioning"`. After these steps, there should be no errors and no diff on this field. This behavior is fixed in v0.94 with state upgrader.


#### *(behavior change)* Renamed fields

Renamed field `provisioner_role` to `run_as_role` to align with Snowflake docs. Please rename this field in your configuration files. State will be migrated automatically.

#### *(new feature)* New fields
Fields added to the resource:
- `enabled`
- `sync_password`
- `comment`

#### *(behavior change)* Changed behavior of `enabled`
New field `enabled` is required. Previously the default value during create in Snowflake was `true`. If you created a resource with Terraform, please add `enabled = true` to have the same value.

#### *(behavior change)* Force new for multiple attributes
ForceNew was added for the following attributes (because there are no usable SQL alter statements for them):
- `scim_client`
- `run_as_role`

### snowflake_warehouse resource changes

Because of the multiple changes in the resource, the easiest migration way is to follow our [migration guide](./docs/guides/resource_migration.md) to perform zero downtime migration. Alternatively, it is possible to follow some pointers below. Either way, familiarize yourself with the resource changes before version bumping. Also, check the [design decisions](./v1-preparations/CHANGES_BEFORE_V1.md).

#### *(potential behavior change)* Default values removed
As part of the [redesign](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1) we are removing the default values for attributes having their defaults on Snowflake side to reduce coupling with the provider (read more in [default values](./v1-preparations/CHANGES_BEFORE_V1.md#default-values)). Because of that the following defaults were removed:
- `comment` (previously `""`)
- `enable_query_acceleration` (previously `false`)
- `query_acceleration_max_scale_factor` (previously `8`)
- `warehouse_type` (previously `"STANDARD"`)
- `max_concurrency_level` (previously `8`)
- `statement_queued_timeout_in_seconds` (previously `0`)
- `statement_timeout_in_seconds` (previously `172800`)

**Beware!** For attributes being Snowflake parameters (in case of warehouse: `max_concurrency_level`, `statement_queued_timeout_in_seconds`, and `statement_timeout_in_seconds`), this is a breaking change (read more in [Snowflake parameters](./v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters)). Previously, not setting a value for them was treated as a fallback to values hardcoded on the provider side. This caused warehouse creation with these parameters set on the warehouse level (and not using the Snowflake default from hierarchy; read more in the [parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters)). To keep the previous values, fill in your configs to the default values listed above.

All previous defaults were aligned with the current Snowflake ones, however it's not possible to distinguish between filled out value and no value in the automatic state upgrader. Therefore, if the given attribute is not filled out in your configuration, terraform will try to perform update after the change (to UNSET the given attribute to the Snowflake default); it should result in no changes on Snowflake object side, but it is required to make Terraform state aligned with your config. **All** other optional fields that were not set inside the config at all (because of the change in handling state logic on our provider side) will follow the same logic. To avoid the need for the changes, fill out the default fields in your config. Alternatively, run `terraform apply`; no further changes should be shown as a part of the plan.

#### *(note)* Automatic state migrations
There are three migrations that should happen automatically with the version bump:
- incorrect `2XLARGE`, `3XLARGE`, `4XLARGE`, `5XLARGE`, `6XLARGE` values for warehouse size are changed to the proper ones
- deprecated `wait_for_provisioning` attribute is removed from the state
- old empty resource monitor attribute is cleaned (earlier it was set to `"null"` string)

#### *(fix)* Warehouse size UNSET

Before the changes, removing warehouse size from the config was not handled properly. Because UNSET is not supported for warehouse size (check the [docs](https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse#properties-parameters) - usage notes for unset) and there are multiple defaults possible, removing the size from config will result in the resource recreation.

#### *(behavior change)* Validation changes
As part of the [redesign](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1) we are adjusting validations or removing them to reduce coupling between Snowflake and the provider. Because of that the following validations were removed/adjusted/added:
- `max_cluster_count` - adjusted: added higher bound (10) according to Snowflake docs
- `min_cluster_count` - adjusted: added higher bound (10) according to Snowflake docs
- `auto_suspend` - adjusted: added `0` as valid value
- `warehouse_size` - adjusted: removed incorrect `2XLARGE`, `3XLARGE`, `4XLARGE`, `5XLARGE`, `6XLARGE` values
- `resource_monitor` - added: validation for a valid identifier (still subject to change during [identifiers rework](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework))
- `max_concurrency_level` - added: validation according to MAX_CONCURRENCY_LEVEL parameter docs
- `statement_queued_timeout_in_seconds` - added: validation according to STATEMENT_QUEUED_TIMEOUT_IN_SECONDS parameter docs
- `statement_timeout_in_seconds` - added: validation according to STATEMENT_TIMEOUT_IN_SECONDS parameter docs

#### *(behavior change)* Deprecated `wait_for_provisioning` field removed
`wait_for_provisioning` field was deprecated a long time ago. It's high time it was removed from the schema.

#### *(behavior change)* `query_acceleration_max_scale_factor` conditional logic removed
Previously, the `query_acceleration_max_scale_factor` was depending on `enable_query_acceleration` parameter, but it is not required on Snowflake side. After migration, `terraform plan` should suggest changes if `enable_query_acceleration` was earlier set to false (manually or from default) and if `query_acceleration_max_scale_factor` was set in config.

#### *(behavior change)* `initially_suspended` forceNew removed
Previously, the `initially_suspended` attribute change caused the resource recreation. This attribute is used only during creation (to create suspended warehouse). There is no reason to recreate the whole object just to have initial state changed.

#### *(behavior change)* Boolean type changes
To easily handle three-value logic (true, false, unknown) in provider's configs, type of `auto_resume` and `enable_query_acceleration` was changed from boolean to string. This should not require updating existing configs (boolean/int value should be accepted and state will be migrated to string automatically), however we recommend changing config values to strings. Terraform should perform an action for configs lacking `auto_resume` or `enable_query_acceleration` (`ALTER WAREHOUSE UNSET AUTO_RESUME` and/or `ALTER WAREHOUSE UNSET ENABLE_QUERY_ACCELERATION` will be run underneath which should not affect the Snowflake object, because `auto_resume` and `enable_query_acceleration` are false by default).

#### *(note)* `resource_monitor` validation and diff suppression
`resource_monitor` is an identifier and handling logic may be still slightly changed as part of https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework. It should be handled automatically (without needed manual actions on user side), though, but it is not guaranteed.

#### *(behavior change)* snowflake_warehouses datasource
- Added `like` field to enable warehouse filtering
- Added missing fields returned by SHOW WAREHOUSES and enclosed its output in `show_output` field.
- Added outputs from **DESC WAREHOUSE** and **SHOW PARAMETERS IN WAREHOUSE** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
  The additional parameters call **DESC WAREHOUSE** (with `with_describe` turned on) and **SHOW PARAMETERS IN WAREHOUSE** (with `with_parameters` turned on) **per warehouse** returned by **SHOW WAREHOUSES**.
  The outputs of both commands are held in `warehouses` entry, where **DESC WAREHOUSE** is saved in the `describe_output` field, and **SHOW PARAMETERS IN WAREHOUSE** in the `parameters` field.
  It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

You can read more in ["raw Snowflake output"](./v1-preparations/CHANGES_BEFORE_V1.md#empty-values).

### *(new feature)* new database resources
As part of the [preparation for v1](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#preparing-essential-ga-objects-for-the-provider-v1), we split up the database resource into multiple ones:
- Standard database - can be used as `snowflake_database` (replaces the old one and is used to create databases with optional ability to become a primary database ready for replication)
- Shared database - can be used as `snowflake_shared_database` (used to create databases from externally defined shares)
- Secondary database - can be used as `snowflake_secondary_database` (used to create replicas of databases from external sources)

All the field changes in comparison to the previous database resource are:
- `is_transient`
    - in `snowflake_shared_database`
        - removed: the field is removed from `snowflake_shared_database` as it doesn't have any effect on shared databases.
- `from_database` - database cloning was entirely removed and is not possible by any of the new database resources.
- `from_share` - the parameter was moved to the dedicated resource for databases created from shares `snowflake_shared_database`. Right now, it's a text field instead of a map. Additionally, instead of legacy account identifier format we're expecting the new one that with share looks like this: `<organization_name>.<account_name>.<share_name>`. For more information on account identifiers, visit the [official documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier).
- `from_replication` - the parameter was moved to the dedicated resource for databases created from primary databases `snowflake_secondary_database`
- `replication_configuration` - renamed: was renamed to `configuration` and is only available in the `snowflake_database`. Its internal schema changed that instead of list of accounts, we expect a list of nested objects with accounts for which replication (and optionally failover) should be enabled. More information about converting between both versions [here](#resource-renamed-snowflake_database---snowflake_database_old). Additionally, instead of legacy account identifier format we're expecting the new one that looks like this: `<organization_name>.<account_name>` (it will be automatically migrated to the recommended format by the state upgrader). For more information on account identifiers, visit the [official documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier).
- `data_retention_time_in_days`
  - in `snowflake_shared_database`
      - removed: the field is removed from `snowflake_shared_database` as it doesn't have any effect on shared databases.
  - in `snowflake_database` and `snowflake_secondary_database`
    - adjusted: now, it uses different approach that won't set it to -1 as a default value, but rather fills the field with the current value from Snowflake (this still can change).
- added: The following set of [parameters](https://docs.snowflake.com/en/sql-reference/parameters) was added to every database type:
    - `max_data_extension_time_in_days`
    - `external_volume`
    - `catalog`
    - `replace_invalid_characters`
    - `default_ddl_collation`
    - `storage_serialization_policy`
    - `log_level`
    - `trace_level`
    - `suspend_task_after_num_failures`
    - `task_auto_retry_attempts`
    - `user_task_managed_initial_warehouse_size`
    - `user_task_timeout_ms`
    - `user_task_minimum_trigger_interval_in_seconds`
    - `quoted_identifiers_ignore_case`
    - `enable_console_output`

The split was done (and will be done for several objects during the refactor) to simplify the resource on maintainability and usage level.
Its purpose was also to divide the resources by their specific purpose rather than cramping every use case of an object into one resource.

### *(behavior change)* Resource renamed snowflake_database -> snowflake_database_old
We made a decision to use the existing `snowflake_database` resource for redesigning it into a standard database.
The previous `snowflake_database` was renamed to `snowflake_database_old` and the current `snowflake_database`
contains completely new implementation that follows our guidelines we set for V1.
When upgrading to the 0.93.0 version, the automatic state upgrader should cover the migration for databases that didn't have the following fields set:
- `from_share` (now, the new `snowflake_shared_database` should be used instead)
- `from_replica` (now, the new `snowflake_secondary_database` should be used instead)
- `replication_configuration`

For configurations containing `replication_configuraiton` like this one:
```terraform
resource "snowflake_database" "test" {
  name = "<name>"
  replication_configuration {
    accounts = ["<account_locator>", "<account_locator_2>"]
    ignore_edition_check = true
  }
}
```

You have to transform the configuration into the following format (notice the change from account locator into the new account identifier format):
```terraform
resource "snowflake_database" "test" {
  name = "%s"
  replication {
    enable_to_account {
      account_identifier = "<organization_name>.<account_name>"
      with_failover      = false
    }
    enable_to_account {
      account_identifier = "<organization_name_2>.<account_name_2>"
      with_failover      = false
    }
  }
  ignore_edition_check = true
}
```

If you had `from_database` set, you should follow our [resource migration guide](./docs/guides/resource_migration.md) to remove
the database from state to later import it in the newer version of the provider.
Otherwise, it may cause issues when migrating to v0.93.0.
For now, we're dropping the possibility to create a clone database from other databases.
The only way will be to clone a database manually and import it as `snowflake_database`, but if
cloned databases diverge in behavior from standard databases, it may cause issues.

For databases with one of the fields mentioned above, manual migration will be needed.
Please refer to our [migration guide](./docs/guides/resource_migration.md) to perform zero downtime migration.

If you would like to upgrade to the latest version and postpone the upgrade, you still have to perform the manual migration
to the `snowflake_database_old` resource by following the [zero downtime migrations document](./docs/guides/resource_migration.md).
The only difference would be that instead of writing/generating new configurations you have to just rename the existing ones to contain `_old` suffix.

### *(behavior change)* snowflake_databases datasource
- `terse` and `history` fields were removed.
- `replication_configuration` field was removed from `databases`.
- `pattern` was replaced by `like` field.
- Additional filtering options added (`limit`).
- Added missing fields returned by SHOW DATABASES and enclosed its output in `show_output` field.
- Added outputs from **DESC DATABASE** and **SHOW PARAMETERS IN DATABASE** (they can be turned off by declaring `with_describe = false` and `with_parameters = false`, **they're turned on by default**).
The additional parameters call **DESC DATABASE** (with `with_describe` turned on) and **SHOW PARAMETERS IN DATABASE** (with `with_parameters` turned on) **per database** returned by **SHOW DATABASES**.
The outputs of both commands are held in `databases` entry, where **DESC DATABASE** is saved in the `describe_output` field, and **SHOW PARAMETERS IN DATABASE** in the `parameters` field.
It's important to limit the records and calls to Snowflake to the minimum. That's why we recommend assessing which information you need from the data source and then providing strong filters and turning off additional fields for better plan performance.

## v0.89.0 ➞ v0.90.0
### snowflake_table resource changes
#### *(behavior change)* Validation to column type added
While solving issue [#2733](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2733) we have introduced diff suppression for `column.type`. To make it work correctly we have also added a validation to it. It should not cause any problems, but it's worth noting in case of any data types used that the provider is not aware of.

### snowflake_procedure resource changes
#### *(behavior change)* Validation to arguments type added
Diff suppression for `arguments.type` is needed for the same reason as above for `snowflake_table` resource.

### tag_masking_policy_association resource changes
Now the `tag_masking_policy_association` resource will only accept fully qualified names separated by dot `.` instead of pipe `|`.

Before
```terraform
resource "snowflake_tag_masking_policy_association" "name" {
    tag_id            = snowflake_tag.this.id
    masking_policy_id = snowflake_masking_policy.example_masking_policy.id
}
```

After
```terraform
resource "snowflake_tag_masking_policy_association" "name" {
    tag_id            = "\"${snowflake_tag.this.database}\".\"${snowflake_tag.this.schema}\".\"${snowflake_tag.this.name}\""
    masking_policy_id = "\"${snowflake_masking_policy.example_masking_policy.database}\".\"${snowflake_masking_policy.example_masking_policy.schema}\".\"${snowflake_masking_policy.example_masking_policy.name}\""
}
```

It's more verbose now, but after identifier rework it should be similar to the previous form.

## v0.88.0 ➞ v0.89.0
#### *(behavior change)* ForceNew removed
The `ForceNew` field was removed in favor of in-place Update for `name` parameter in:
- `snowflake_file_format`
- `snowflake_masking_policy`
So from now, these objects won't be re-created when the `name` changes, but instead only the name will be updated with `ALTER .. RENAME TO` statements.

## v0.87.0 ➞ v0.88.0

### snowflake_role data source deprecation

Already existing `snowflake_role` was deprecated in favor of the new `snowflake_roles`. You can have a similar behavior like before by specifying `pattern` field. Please adjust your Terraform configurations.

### snowflake_procedure resource changes
#### *(behavior change)* Execute as validation added
From now on, the `snowflake_procedure`'s `execute_as` parameter allows only two values: OWNER and CALLER (case-insensitive). Setting other values earlier resulted in falling back to the Snowflake default (currently OWNER) and creating a permadiff.

### snowflake_grants datasource changes
`snowflake_grants` datasource was refreshed as part of the ongoing [Grants Redesign](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#redesigning-grants).

#### *(behavior change)* role fields renames
To be aligned with the convention in other grant resources, `role` was renamed to `account_role` for the following fields:
- `grants_to.role`
- `grants_of.role`
- `future_grants_to.role`.

To migrate simply change `role` to `account_role` in the aforementioned fields.

#### *(behavior change)* grants_to.share type change
`grants_to.share` was a text field. Because Snowflake introduced new syntax `SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name>` (check more in the [docs](https://docs.snowflake.com/en/sql-reference/sql/show-grants#variants)) the type was changed to object. To migrate simply change:
```terraform
data "snowflake_grants" "example_to_share" {
  grants_to {
    share = "some_share"
  }
}
```
to
```terraform
data "snowflake_grants" "example_to_share" {
  grants_to {
    share {
      share_name = "some_share"
    }
  }
}
```
Note: `in_application_package` is not yet supported.

#### *(behavior change)* future_grants_in.schema type change
`future_grants_in.schema` was an object field allowing to set required `schema_name` and optional `database_name`. Our strategy is to be explicit, so the schema field was changed to string and fully qualified name is expected. To migrate change:
```terraform
data "snowflake_grants" "example_future_in_schema" {
  future_grants_in {
    schema {
      database_name = "some_database"
      schema_name   = "some_schema"
    }
  }
}
```
to
```terraform
data "snowflake_grants" "example_future_in_schema" {
  future_grants_in {
    schema = "\"some_database\".\"some_schema\""
  }
}
```
#### *(new feature)* grants_to new options
`grants_to` was enriched with three new options:
- `application`
- `application_role`
- `database_role`

No migration work is needed here.

#### *(new feature)* grants_of new options
`grants_to` was enriched with two new options:
- `database_role`
- `application_role`

No migration work is needed here.

#### *(new feature)* future_grants_to new options
`future_grants_to` was enriched with one new option:
- `database_role`

No migration work is needed here.

#### *(documentation)* improvements
Descriptions of attributes were altered. More examples were added (both for old and new features).

## v0.86.0 ➞ v0.87.0
### snowflake_database resource changes
#### *(behavior change)* External object identifier changes

Previously, in `snowflake_database` when creating a database form share, it was possible to provide `from_share.provider`
in the format of `<org_name>.<account_name>`. It worked even though we expected account locator because our "external" identifier wasn't quoting its string representation.
To be consistent with other identifier types, we quoted the output of "external" identifiers which makes such configurations break
(previously, they were working "by accident"). To fix it, the previous format of `<org_name>.<account_name>` has to be changed
to account locator format `<account_locator>` (mind that it's now case-sensitive). The account locator can be retrieved by calling `select current_account();` on the sharing account.
In the future we would like to eventually come back to the `<org_name>.<account_name>` format as it's recommended by Snowflake.

### Provider configuration changes

#### **IMPORTANT** *(bug fix)* Configuration hierarchy
There were several issues reported about the configuration hierarchy, e.g. [#2294](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2294) and [#2242](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2242).
In fact, the order of precedence described in the docs was not followed. This have led to the incorrect behavior.

After migrating to this version, the hierarchy from the docs should be followed:
```text
The Snowflake provider will use the following order of precedence when determining which credentials to use:
1) Provider Configuration
2) Environment Variables
3) Config File
```

**BEWARE**: your configurations will be affected with that change because they may have been leveraging the incorrect configurations precedence. Please be sure to check all the configurations before running terraform.

### snowflake_failover_group resource changes
#### *(bug fix)* ACCOUNT PARAMETERS is returned as PARAMETERS from SHOW FAILOVER GROUPS
Longer context in [#2517](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2517).
After this change, one apply may be required to update the state correctly for failover group resources using `ACCOUNT PARAMETERS`.

### snowflake_database, snowflake_schema, and snowflake_table resource changes
#### *(behavior change)* Database `data_retention_time_in_days` + Schema `data_retention_days` + Table `data_retention_time_in_days`
For context [#2356](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2356).
To make data retention fields truly optional (previously they were producing plan every time when no value was set),
we added `-1` as a possible value, and it is set as default. That got rid of the unexpected plans when no value is set and added possibility to use default value assigned by Snowflake (see [the data retention period](https://docs.snowflake.com/en/user-guide/data-time-travel#data-retention-period)).

### snowflake_table resource changes
#### *(behavior change)* Table `data_retention_days` field removed in favor of `data_retention_time_in_days`
For context [#2356](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2356).
To define data retention days for table `data_retention_time_in_days` should be used as deprecated `data_retention_days` field is being removed.

## v0.85.0 ➞ v0.86.0
### snowflake_table_constraint resource changes

#### *(behavior change)* NOT NULL removed from possible types
The `type` of the constraint was limited back to `UNIQUE`, `PRIMARY KEY`, and `FOREIGN KEY`.
The reason for that is, that syntax for Out-of-Line constraint ([docs](https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint#out-of-line-unique-primary-foreign-key)) does not contain `NOT NULL`.
It is noted as a behavior change but in some way it is not; with the previous implementation it did not work at all with `type` set to `NOT NULL` because the generated statement was not a valid Snowflake statement.

We will consider adding `NOT NULL` back because it can be set by `ALTER COLUMN columnX SET NOT NULL`, but first we want to revisit the whole resource design.

#### *(behavior change)* table_id reference
The docs were inconsistent. Example prior to 0.86.0 version showed using the `table.id` as the `table_id` reference. The description of the `table_id` parameter never allowed such a value (`table.id` is a `|`-delimited identifier representation and only the `.`-separated values were listed in the docs: https://registry.terraform.io/providers/snowflakedb/snowflake/0.85.0/docs/resources/table_constraint#required. The misuse of `table.id` parameter will result in error after migrating to 0.86.0. To make the config work, please remove and reimport the constraint resource from the state as described in [resource migration doc](./docs/guides/resource_migration.md).

After discussions in [#2535](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2535) we decided to provide a temporary workaround in 0.87.0 version, so that the manual migration is not necessary. It allows skipping the migration and jumping straight to 0.87.0 version. However, the temporary workaround will be gone in one of the future versions. Please adjust to the newly suggested reference with the new resources you create.

### snowflake_external_function resource changes

#### *(behavior change)* return_null_allowed default is now true
The `return_null_allowed` attribute default value is now `true`. This is a behavior change because it was `false` before. The reason it was changed is to match the expected default value in the [documentation](https://docs.snowflake.com/en/sql-reference/sql/create-external-function#optional-parameters) `Default: The default is NULL (i.e. the function can return NULL values).`

#### *(behavior change)* comment is no longer required
The `comment` attribute is now optional. It was required before, but it is not required in Snowflake API.

### snowflake_external_functions data source changes

#### *(behavior change)* schema is now required with database
The `schema` attribute is now required with `database` attribute to match old implementation `SHOW EXTERNAL FUNCTIONS IN SCHEMA "<database>"."<schema>"`. In the future this may change to make schema optional.

## vX.XX.X ➞ v0.85.0
<a id="vxxxx---v0850"></a>

### Migration from old (grant) resources to new ones

In recent changes, we introduced a new grant resources to replace the old ones.
To aid with the migration, we wrote a guide to show one of the possible ways to migrate deprecated resources to their new counter-parts.
As the guide is more general and applies to every version (and provider), we moved it [here](./docs/guides/resource_migration.md).

### snowflake_procedure resource changes
#### *(deprecation)* return_behavior
`return_behavior` parameter is deprecated because it is also deprecated in the Snowflake API.

### snowflake_function resource changes
#### *(behavior change)* return_type
`return_type` has become force new because there is no way to alter it without dropping and recreating the function.

## v0.84.0 ➞ v0.85.0

### snowflake_stage resource changes

#### *(behavior change/regression)* copy_options
Setting `copy_options` to `ON_ERROR = 'CONTINUE'` would result in a permadiff. Use `ON_ERROR = CONTINUE` (without single quotes) or bump to v0.89.0 in which the behavior was fixed.

### snowflake_notification_integration resource changes
#### *(behavior change)* notification_provider
`notification_provider` becomes required and has three possible values `AZURE_STORAGE_QUEUE`, `AWS_SNS`, and `GCP_PUBSUB`.
It is still possible to set it to `AWS_SQS` but because there is no underlying SQL, so it will result in an error.
Attributes `aws_sqs_arn` and `aws_sqs_role_arn` will be ignored.
Computed attributes `aws_sqs_external_id` and `aws_sqs_iam_user_arn` won't be updated.

#### *(behavior change)* force new for multiple attributes
Force new was added for the following attributes (because no usable SQL alter statements for them):
- `azure_storage_queue_primary_uri`
- `azure_tenant_id`
- `gcp_pubsub_subscription_name`
- `gcp_pubsub_topic_name`

#### *(deprecation)* direction
`direction` parameter is deprecated because it is added automatically on the SDK level.

#### *(deprecation)* type
`type` parameter is deprecated because it is added automatically on the SDK level (and basically it's always `QUEUE`).

## v0.73.0 ➞ v0.74.0
### Provider configuration changes

In this change we have done a provider refactor to make it more complete and customizable by supporting more options that
were already available in Golang Snowflake driver. This lead to several attributes being added and a few deprecated.
We will focus on the deprecated ones and show you how to adapt your current configuration to the new changes.

#### *(rename)* username ➞ user
Provider field `username` were renamed to `user`. Adjust your provider configuration like below:
```terraform
provider "snowflake" {
  # before
  username = "username"

  # after
  user = "username"
}
```

#### *(structural change)* OAuth API
Provider fields regarding Oauth were renamed and nested. Adjust your provider configuration like below:

```terraform
provider "snowflake" {
  # before
  browser_auth        = false
  oauth_access_token  = "<access_token>"
  oauth_refresh_token = "<refresh_token>"
  oauth_client_id     = "<client_id>"
  oauth_client_secret = "<client_secret>"
  oauth_endpoint      = "<endpoint>"
  oauth_redirect_url  = "<redirect_uri>"

  # after
  authenticator = "ExternalBrowser"
  token         = "<access_token>"
  token_accessor {
    refresh_token   = "<refresh_token>"
    client_id       = "<client_id>"
    client_secret   = "<client_secret>"
    token_endpoint  = "<endpoint>"
    redirect_uri    = "<redirect_uri>"
  }
}
```

#### *(remove redundant information)* region

Specifying a region is a legacy thing and according to https://docs.snowflake.com/en/user-guide/admin-account-identifier
you can specify a region as a part of account parameter. Specifying account parameter with the region is also considered legacy,
but with this approach it will be easier to convert only your account identifier to the new preferred way of specifying account identifier.

```terraform
provider "snowflake" {
  # before
  region = "<cloud_region_id>"

  # after
  account = "<account_locator>.<cloud_region_id>"
}
```

#### private_key_path deprecation
Provider field `private_key_path` is now deprecated in favor of `private_key` and `file` Terraform function (see [docs](https://developer.hashicorp.com/terraform/language/functions/file)). Adjust your provider configuration like below:

```terraform
provider "snowflake" {
  # before
  private_key_path = "<filepath>"

  # after
  private_key = file("<filepath>")
}
```

#### *(rename)* session_params ➞ params
Provider field `session_params` were renamed to `params`. Adjust your provider configuration like below:
```terraform
provider "snowflake" {
  # before
  session_params = {}

  # after
  params = {}
}
```

#### *(behavior change)* authenticator (JWT)

Before the change `authenticator` parameter did not have to be set for private key authentication and was deduced by the provider. The change is a result of the introduced configuration alignment with an underlying [gosnowflake driver](https://github.com/snowflakedb/gosnowflake). The authentication type is required there, and it defaults to user+password one. From this version, set `authenticator` to `JWT` explicitly.
