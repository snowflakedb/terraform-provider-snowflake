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

### Argument output changes for SHOW FUNCTIONS and SHOW PROCEDURES commands

Changed format in `Arguments` column from `SHOW FUNCTIONS/PROCEDURES` output is not compatible with the provider parsing function. It leads to:
- [`snowflake_functions`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.2.0/docs/data-sources/functions) and [`snowflake_procedures`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.2.0/docs/data-sources/procedures) being inoperable. Check: [#3822](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3822).
- All function and all procedure resources failing to read their state from Snowflake, which leads to removing them from terraform state (if `terraform apply` or `terraform plan --refresh-only` is run). Check: [#3823](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3823).

The parsing was improved and is available starting with the [2.3.0](https://registry.terraform.io/providers/snowflakedb/snowflake/2.3.0/docs/) version of the provider. This fix was also backported to the [1.2.3](https://github.com/snowflakedb/terraform-provider-snowflake/releases/tag/v1.2.3) version.

To use the provider with the 2025_03 Bundle:
1. Bump the provider to 2.3.0 version (or 1.2.3 version).
2. Affected data sources should work without any further actions after bumping.
3. If your function/procedure resources were removed from terraform state (you can check it by running `terraform state list`), you need to reimport them (follow our [resource migration guide](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration)).
4. If your function/procedure resources are still in the terraform state, they should work any further actions after bumping.

Reference: [BCR-1944](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1944)

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
