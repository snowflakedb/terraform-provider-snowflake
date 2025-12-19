---
page_title: "Strict Privilege Management"
subcategory: ""
description: |-

---

# Strict privilege management

Some time ago, during the identifier redesign, we introduced a new set of privilege-granting resources ([more details](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions)).
With this change, we removed the possibility to specify the `enable_multiple_grants` flag ([more details](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#implicit-enable_multiple_grants-enabled-by-default)).
We wanted to bring back this option to allow you to specify whether resource should revoke any privileges granted externally.
This functionality was added as `strict_privilege_management` flag that, for now, only exists on the `snowflake_grant_privileges_to_account_role` resource.

It's not enabled by default and to use it, you have to enable this feature on the provider level by adding `GRANTS_STRICT_PRIVILEGE_MANAGEMENT` to `experimental_features_enabled`.
It's similar to the existing [`preview_features_enabled`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.10.0/docs#preview_features_enabled-1).
Instead of enabling the use of the whole resources, it's meant to slightly alter the provider's behavior.
**It's still considered a preview feature, even when applied to the stable resources.**

## Usage example

To quickly show how the flag works, let's say we have a database `TEST_DATABASE` and a role `TEST_ROLE`.
We want to ensure the role has a list of configured privileges set on the database to the role (e.g., `MODIFY` and `MONITOR`), and
no other privilege should be available for this role on this database. Now, we could fulfill such requirement with the following configuration:

```terraform
provider "snowflake" {
  experimental_features_enabled = [ "GRANTS_STRICT_PRIVILEGE_MANAGEMENT" ]
}

resource "snowflake_grant_privileges_to_account_role" "test" {
  account_role_name = "TEST_ROLE"
  privileges = [ "MODIFY", "MONITOR" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
}
```

Now, let's say someone grants a new privilege `USAGE` outside of the Terraform configuration (e.g., in SnowSight) on `TEST_DATABASE` to `TEST_ROLE` as follows:

```sql
GRANT USAGE ON DATABASE "TEST_DATABASE" TO ROLE "TEST_ROLE";
```

Because we have the new `strict_privilege_management` flag set, the resource during the next `terraform plan` will propose the following change:

```log
~ privileges = [ "MODIFY", "MONITOR", "USAGE" ] -> [ "MODIFY", "MONITOR" ]
```

This means the resource can now detect when privileges are granted externally (outside of Terraform) and will revoke them to match the configuration.

## Behavior with future grants

[//]: # (TODO: Explain)
- how it can be used with future grants
- how regular and future grants coexists

## Limitations

The new flag has several limitations due to the current resource implementation.
We plan to address these in future updates.

### Conflicting fields

The newly introduced flag is conflicting with the following set of fields:
- `all_privileges`
- `on_schema.0.all_schemas_in_database`
- `on_schema_object.0.all`

[//]: # (TODO: Explain why)

### Grants option

The grant option is tracked at the resource level, not per privilege.
When using `strict_privilege_management`, all privileges must use the same grant option (either all with or all without).
Using two resources with different grant options will cause conflicts.
External grants are detected regardless of the grant option setting.

### Delayed application

This, on the other hand, is a Terraform limitation. We are limited with what we can show in the plan.
Because of this, we decided to delay the revoking action in certain scenarios, so no revokes will be called without showing proper plan.
Otherwise, the revokes could happen "under the hood", making you unable to make conscious decision whether given set of privileges should be revoked or not.
The delays will occur whenever you:
- Create a new `snowflake_grant_privileges_to_account_role` resource with the flag set to `true`
- or update already existing `snowflake_grant_privileges_to_account_role` to have the flag set to `true`

The flag setting takes place in the first `terraform apply`, and the proper plan with revokes can be applied in the second `terraform apply`.

Related: [#3973](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3973)

