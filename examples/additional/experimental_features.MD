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

#### TAGS_ALLOW_EMPTY_ALLOWED_VALUES
Enables behavior changes for the `allowed_values` field in the `snowflake_tag` resource.

When enabled, the three possible states in Snowflake for allowed values will be supported: `nil` (any value is allowed; whenever `allowed_values` are empty), `empty` (no value is allowed; handled by the `no_allowed_values` field), and `set` (all values defined in `allowed_values` are allowed).

Otherwise, the `no_allowed_values` field will be ignored (explicit changes will cause updates, but without any effect) and the `allowed_values` field will follow the old behavior: `nil` (any value is allowed; only available whenever tag resource is created without `allowed_values`), `empty` (no value is allowed; always set when updating from filled `allowed_values` set to empty one or completely removed from config), `set` (all values defined in `allowed_values` are allowed).

#### IMPORT_BOOLEAN_DEFAULT
Changes import behavior for boolean fields using the special `"default"` value.

When enabled, boolean fields using the special `"default"` value are set to `"default"` during import instead of the actual Snowflake value (e.g., `"false"`). This prevents unavoidable diffs on every plan after import.

Note: this is supported only on all stage resources (`snowflake_stage_external_s3`, `snowflake_stage_external_azure`, `snowflake_stage_external_gcs`, `snowflake_stage_external_s3_compatible`, and `snowflake_stage_internal`).

#### GRANTS_SAFE_DESTROY
When enabled, grant destroy operations silently succeed when the underlying Snowflake object (or its dependencies) no longer exists.

Currently supported by: `snowflake_grant_privileges_to_account_role`, `snowflake_grant_privileges_to_database_role`, `snowflake_grant_privileges_to_share`, `snowflake_grant_account_role`, `snowflake_grant_database_role`, `snowflake_grant_application_role`.

This prevents errors when, for example, a warehouse or role is deleted externally and the corresponding grant resource is later removed from the Terraform configuration.

Without this experiment, destroying such resources fails with `does not exist or not authorized`.

#### TAG_ASSOCIATION_SAFE_DESTROY
When enabled, tag association destroy operations silently succeed when the tagged object (or its parent hierarchy) no longer exists.

Currently supported by: `snowflake_tag_association`.

This prevents errors when, for example, a table or schema is deleted externally and the corresponding tag association resource is later removed from the Terraform configuration.

Without this experiment, destroying such resources fails with `does not exist or not authorized`.


