### Active experiments

The following experiments are currently active. Depending on the feedback, we may decide to include them as default behavior/stable feature of the provider in the future.

To share feedback please reach out to us through your Snowflake account manager.

#### WAREHOUSE_SHOW_IMPROVED_PERFORMANCE
It's meant to improve the performance for accounts with many warehouses.

When enabled, it uses a slightly different SHOW query to read warehouse details (`SHOW WAREHOUSES LIKE '<identifier>' STARTS WITH '<identifier>' LIMIT 1`).

**Important**: to benefit from this improvement, you need to have it enabled also on your Snowflake account. To do this, please reach out to us through your Snowflake Account Manager.

#### GRANTS_STRICT_PRIVILEGE_MANAGEMENT
The new `strict_privilege_management` flag was added to the `snowflake_grant_privileges_to_account_role` resource.

It has similar behavior to the `enable_multiple_grants` flag present in the old grant resources, and it makes the resource able to detect external changes for privileges other than those present in the configuration which can make the `snowflake_grant_privileges_to_account_role` resource a central point of knowledge privilege management for a given object and role.

Read more in our [strict privilege management](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/strict_privilege_management) guide.

This feature works independently of the `GRANTS_IMPORT_VALIDATION` flag.

#### PARAMETERS_IGNORE_VALUE_CHANGES_IF_NOT_ON_OBJECT_LEVEL
Currently, not setting the parameter value on the object level can unnecessarily react to external changes to this parameter's value on the higher levels (e.g. not setting `data_retention_time_in_days` on `snowflake_schema` can result in non-empty plan when the parameter value changes on the database/account level).

When enabled, the provider ignores changes to the parameter value happening on the higher hierarchy levels.

#### PARAMETERS_REDUCED_OUTPUT
Currently, the `parameters` field in various resources contains a verbatim output for the `SHOW PARAMETERS IN <object>` command. One of the fields contained in the output is the `description`. It does not change and is repeated for all objects containing the given parameter. It leads to an excessive output (check e.g., [#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118)).

To mitigate the problem, we are adding this option to reduce the output to only `value` and `level` fields, which should significantly reduce the state size. **Note**: it's also affecting the `parameters` output for data sources.

We considered the option to remove the `parameters` output completely, however, we plan to change the external change logic detection to use it (to make it consistent with other attributes using `show_output` and because we won't be able to implement the current logic when switching to the Terraform Plugin Framework) and it still allows referencing the parameter value/level from other parts of the configuration.

#### USER_ENABLE_DEFAULT_WORKLOAD_IDENTITY
The new `default_workload_identity_federation` field was added to the `snowflake_legacy_service_user` and `snowflake_service_user` resources. This field allows for managing WIFs. Due to feature complexity, it requires enabling this experiment.

Read more in our [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-feature-workload-identity-federation-support-for-service-users).

#### GRANTS_IMPORT_VALIDATION
Enables import validation for the `snowflake_grant_privileges_to_account_role` resource.

When enabled, importing a grant resource with a fixed set of privileges (`privileges` field) will validate that the specified privileges actually exist in Snowflake with the correct `with_grant_option` setting, and error immediately if they don't match.

This feature works independently of the `GRANTS_STRICT_PRIVILEGE_MANAGEMENT` flag.


