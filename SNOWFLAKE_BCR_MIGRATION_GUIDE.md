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

(will be filled in soon)

Reference: [BCR-1944](https://docs.snowflake.com/en/release-notes/bcr-bundles/2025_03/bcr-1944)

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
