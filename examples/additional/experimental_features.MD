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

#### PARAMETERS_IGNORE_VALUE_CHANGES_IF_NOT_ON_OBJECT_LEVEL
Currently, not setting the parameter value on the object level can unnecessarily react to external changes to this parameter's value on the higher levels (e.g. not setting `data_retention_time_in_days` on `snowflake_schema` can result in non-empty plan when the parameter value changes on the database/account level).

When enabled, the provider ignores changes to the parameter value happening on the higher hierarchy levels.

#### USER_ENABLE_DEFAULT_WORKLOAD_IDENTITY
TODO


