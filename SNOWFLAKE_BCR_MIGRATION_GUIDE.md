# Snowflake BCR migration guide

This document is meant to help you migrate your Terraform config and maintain compatibility after enabling given [Snowflake BCR Bundle](https://docs.snowflake.com/en/release-notes/behavior-changes).
Some of the breaking changes on Snowflake side may be not compatible with the current version of the Terraform provider, so you may need to update your Terraform config to adapt to the new behavior.
As some changes may require work on the provider side, we advise you to always use the latest version of the provider ([new features and fixes policy](https://docs.snowflake.com/en/user-guide/terraform#new-features-and-fixes)).
To avoid any issues and follow [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md) when migrating to newer versions.
According to the [Bundle Lifecycle](https://docs.snowflake.com/en/release-notes/intro-bcr-releases#bundle-lifecycle), changes are eventually enabled by default without the possibility to disable them, so it's important to know what is going to be introduced beforehand.
If you would like to test the new behavior before it is enabled by default, you can use the [SYSTEM\$ENABLE_BEHAVIOR_CHANGE_BUNDLE](https://docs.snowflake.com/en/sql-reference/functions/system_enable_behavior_change_bundle)
command to enable the bundle manually, and then the [SYSTEM\$DISABLE_BEHAVIOR_CHANGE_BUNDLE](https://docs.snowflake.com/en/sql-reference/functions/system_disable_behavior_change_bundle) command to disable it.

Remember that only changes that affect the provider are listed here, to get the full list of changes, please refer to the [Snowflake BCR Bundle documentation](https://docs.snowflake.com/en/release-notes/behavior-changes).
The `snowflake_execute` resource won't be listed here, as it is users' responsibility to check the SQL commands executed and adapt them to the new behavior.

## [Unbundled changes](https://docs.snowflake.com/en/release-notes/bcr-bundles/un-bundled/unbundled-behavior-changes)

### Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands

> [!IMPORTANT]
> This change has been rolled back from the BCR 2025_03.

Changed format in `Arguments` column from `SHOW FUNCTIONS/PROCEDURES` output is not compatible with the provider parsing function. It leads to:
- [`snowflake_functions`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.2.0/docs/data-sources/functions) and [`snowflake_procedures`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.2.0/docs/data-sources/procedures) being inoperable. Check: [#3822](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3822).
- All function and all procedure resources failing to read their state from Snowflake, which leads to removing them from terraform state (if `terraform apply` or `terraform plan --refresh-only` is run). Check: [#3823](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3823).

The parsing was improved and is available starting with the [2.3.0](https://registry.terraform.io/providers/snowflakedb/snowflake/2.3.0/docs/) version of the provider. This fix was also backported to the [1.2.3](https://github.com/snowflakedb/terraform-provider-snowflake/releases/tag/v1.2.3) version.

To use the provider with the bundles containing this change:
1. Bump the provider to 2.3.0 version (or 1.2.3 version).
2. Affected data sources should work without any further actions after bumping.
3. If your function/procedure resources were removed from terraform state (you can check it by running `terraform state list`), you need to reimport them (follow our [resource migration guide](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration)).
4. If your function/procedure resources are still in the terraform state, they should work any further actions after bumping.

Reference: [BCR-1944](https://docs.snowflake.com/release-notes/bcr-bundles/un-bundled/bcr-1944)

## [Bundle 2025_06](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06_bundle)

### Changes in authentication policies
<!-- TODO(SNOW-2187814): Update this entry. -->

> [!IMPORTANT]
> The [BCR-2086](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06/bcr-2086) change has been rolled back from the BCR 2025_04 and was moved to 2025_06.

> [!IMPORTANT]
> These change has not been addressed in the provider yet. They will be addressed in the next versions of the provider.
> As a workaround, please use the [execute](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/execute) resource.

The `MFA_AUTHENTICATION_METHODS` property is deprecated. Setting the `MFA_AUTHENTICATION_METHODS` property returns an error. If you use the [authentication_policy](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/authentication_policy) resource with `mfa_authentication_methods` field
and have this bundle enabled, the provider will return an error.
The new way of handling authentication methods is `ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION` which will be handled in this resource in the next versions.

Additionally, the allowed values for `MFA_ENROLLMENT` are changed: `OPTIONAL` is removed and `REQUIRED_PASSWORD_ONLY` and `REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY` are added.


Reference: [BCR-2086](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06/bcr-2086), [BCR-2097](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06/bcr-2097)

### Snowflake OAuth authentication: Change in the network policy used for a request from client to Snowflake

This change modifies the behavior of authentication with active network policies. Please verify that your network policy configuration allows connection by the provider after activating this change.

Additionally, this change adds the possibility to assign network policies to External Oauth integrations.

Setting the `network_policy` field in `external_oauth_integration` resource is not yet supported in the provider, and it will be handled in the future. As a workaround, please use the [execute](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/execute) resource.

Reference: [BCR-2094](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06/bcr-2094)

### Snowpark Container Services job service: Retention-time increase

In the provider, the job_service resource forces setting the `ASYNC` option.

Before the change, Snowflake automatically deletes the job service 7 days after completion.

After the change, Snowflake retains job services for 14 days after completion.

In most cases, this change in Snowflake should have no effect on the provider. Optionally, you can manually drop the completed jobs.

Reference: [BCR-2093](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_06/bcr-2093)

## [Bundle 2025_05](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_05_bundle)

### Key-pair authentication for Google Cloud accounts in the us-central1 region

Previously, when you used key-pair authentication from a Snowflake account in the Google Cloud us-central1 region, specifying the account by using an account locator with additional segments was supported.

Now, when you use key-pair authentication across all cloud platforms and regions, you must specify the account by using only the account locator without additional segments. See our [Authentication methods](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods) guide for authentication overview in the provider.

Reference: [BCR-2055](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_05/bcr-2055)

### File formats and stages: Enforce dependency checks

You can't drop or recreate a file format or stage that has dependent external tables. You also can't alter the location of a stage with dependent external tables. To perform these operations, first drop the dependent external tables manually.

Reference: [BCR-1989](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_05/bcr-1989)

## [Bundle 2025_04](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_04_bundle)

### `MFA_AUTHENTICATION_METHODS` in authentication policy now only includes `PASSWORD` by default

Previously, the created authentication policies with default `MFA_AUTHENTICATION_METHODS` had both `[PASSWORD, SAML]` values.
In this BCR, the default value is changed to only `PASSWORD`. This can cause a permadiff on the optional `mfa_authentication_methods` field in `authentication_policy` resource.
To address this, you can either specify this attribute in the resource configuration, or use the [ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments#lifecycle) meta argument.

This resource is still in preview, and we are planning to rework it in the near future. Handling of default values will be improved.

Reference: [BCR-1971](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_04/bcr-1971)

### Primary role requires stage access during `CREATE EXTERNAL TABLE` command

Creating an external table succeeds only if a user’s primary role has the `USAGE` privilege on the stage referenced in the `snowflake_external_table` resource. If you manage external tables in the provider, please grant the `USAGE` privilege on the relevant stages to the connection role.

Reference: [BCR-1993](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_04/bcr-1993)

## [Bundle 2025_03](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03_bundle)

### The `CREATE DATA EXCHANGE LISTING` privilege rename

The `CREATE DATA EXCHANGE LISTING` that is granted on account was changed to just `CREATE LISTING`.
If you are using any of the privilege-granting resources, such as [snowflake_grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role)
to perform no downtime migration, you may want to follow our [resource migration guide](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration).
Basically the steps are:
- Remove the resource from the state
- Adjust it to use the new privilege name, i.e. `CREATE LISTING`
- Re-import the resource into the state (with correct privilege name in the imported identifier)

Reference: [BCR-1926](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1926)

### New maximum size limits for database objects

Max sizes for the few data types were increased.

There are no immediate impacts found on the provider execution.
However, as explained in the [Data type changes](./MIGRATION_GUIDE.md#data-type-changes) section of our migration guide, the provider fills out the data type attributes (like size) if they are not provided by the user.
Sizes of `VARCHAR` and `BINARY` data types (when no size is specified) will continue to use the old defaults in the provider (16MB and 8MB respectively).
If you want to use bigger sizes after enabling the Bundle, please specify them explicitly.

These default values may be changed in the future versions of the provider.

Reference: [BCR-1942](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1942)

### Python UDFs and stored procedures: Stop implicit auto-injection of the psutil package

The `psutil` package is no longer implicitly injected into Python UDFs and stored procedures.
Adjust your configuration to use the `psutil` package explicitly in your Python UDFs and stored procedures, like so:
```terraform
resource "snowflake_procedure_python" "test" {
  packages = ["psutil==5.9.0"]
  # other arguments...
}
```

Reference: [BCR-1948](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1948)
