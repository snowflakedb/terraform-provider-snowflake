---
page_title: "Strict Privilege Management"
subcategory: ""
description: |-

---

# Strict privilege management

Some time ago, during the identifier redesign, we introduced a new set of privilege-granting resources ([more details](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions)).
With this change, we removed the possibility to specify the `enable_multiple_grants` flag ([more details](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#implicit-enable_multiple_grants-enabled-by-default)).
We wanted to bring back this option to allow you to specify whether resource should revoke any privileges granted externally.
This functionality was added as `strict_privilege_management` flag that, for now, only exists in the [grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role) resource.

It's not enabled by default and to use it, you have to enable this feature on the provider level
by adding `GRANTS_STRICT_PRIVILEGE_MANAGEMENT` value to the [`experimental_features_enabled`](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#experimental_features_enabled-1) provider field.
It's similar to the existing [`preview_features_enabled`](https://registry.terraform.io/providers/snowflakedb/snowflake/2.10.0/docs#preview_features_enabled-1),
but instead of enabling the use of the whole resources, it's meant to slightly alter the provider's behavior.
**It's still considered a preview feature, even when applied to the stable resources.**

Related feature request: [#3973](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3973)

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

## Limitations

The new flag has several limitations originating from the current resource design.
They can't be resolved without introducing the breaking changes.
We will consider design changes for future major releases.

### Conflicting fields

The newly introduced flag is conflicting with the following set of fields:
- `all_privileges`
- `on_schema.0.all_schemas_in_database`
- `on_schema_object.0.all`

### Grant option

When using `strict_privilege_management`, all privileges must use the same `with_grant_option`.
That's because you cannot create two resources with `strict_privilege_management` working on the same object and role, as they will conflict with each other.
External grants are detected regardless of their grant option setting.

When two resources configured with `strict_privilege_management` attempt to manage different privileges on the same object for the same role.
The outcome is non-deterministic, because these resources are not dependent on each other, but more or less they would
grant the privilege defined in its configuration while simultaneously revoking the privilege granted by the other resource.

> Note: One resource with the flag enabled would also conflict, but in a different way, because the flag is trying to get rid of privileges not defined by a given resource where it's enabled.

Below example shows such conflict:

```terraform
provider "snowflake" {
  experimental_features_enabled = [ "GRANTS_STRICT_PRIVILEGE_MANAGEMENT" ]
}

resource "snowflake_grant_privileges_to_account_role" "conflicting_1" {
  account_role_name = "TEST_ROLE"
  privileges = [ "MODIFY" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "conflicting_2" {
  account_role_name = "TEST_ROLE"
  privileges = [ "MONITOR" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = true
}
```

The solution is then to combine two resources into one, central privilege management point for a given object and role.
One big change is that you have to decide which `with_grant_option` is right in your case.
Merged privileges into one resource with `with_grant_option = false`:

```terraform
# provider block ...

resource "snowflake_grant_privileges_to_account_role" "conflicting_1" {
  account_role_name = "TEST_ROLE"
  privileges = [ "MODIFY", "MONITOR" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = false
}
```

Another solution involves creating two separate roles, with the second role dependent on the first.
The grants for the first role would be managed by one privilege-granting resource **without** the `with_grant_option`.
The grants for the second role would be managed by a different privilege-granting resource **with** the `with_grant_option` enabled.
You can then use a parent role that includes the privileges from both of these roles.
Below example shows this could be represented within Terraform configuration:

```terraform
provider "snowflake" {
  experimental_features_enabled = [ "GRANTS_STRICT_PRIVILEGE_MANAGEMENT" ]
}

resource "snowflake_account_role" "parent" {
  name = "TEST_ROLE_PARENT" 
}

resource "snowflake_account_role" "child" {
  name = "TEST_ROLE_CHILD"
}

resource "snowflake_grant_account_role" "create_role_dependency" {
  role_name = snowflake_account_role.child.name
  parent_role_name = snowflake_account_role.parent.name
}

resource "snowflake_grant_privileges_to_account_role" "non_conflicting_1" {
  account_role_name = snowflake_account_role.parent.name
  privileges = [ "MODIFY" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "non_conflicting_2" {
  account_role_name = snowflake_account_role.child.name
  privileges = [ "MONITOR" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = true
}
```

## Future grants

Future grants are allowed to be specified with the newly added `strict_privilege_management`.
We were facing the limitation of not being able to show the difference between regular and future grants.
Due to the difference on the Snowflake side for both grants, we decided to keep to current grants discovery logic.
This means, that the privilege-granting resources with `strict_privilege_management` flag enabled will detect external changes
only for privilege type ("regular" or future) used within the resource. This is also motivated by the difference within the SQL statements,
but also, how logically the grants are stored:
- The regular grants are stored `on` objects (e.g. `table`)
- The future grants stored are `in` objects (either `database` or `schema`) and they describe grants on lower level (e.g. `in database` describe grants on `schema`, or even lower `table` level)

This distinction shows, that both privileges shouldn't be managed within one resource instance.
It also shows, that there are two levels of future grants (`database` and `schema` levels), that should be considered
when using `strict_privilege_management` flag for full future grants control. Defining future grants for both levels, would
ensure, that no additional external privilege is granted on the Snowflake side manually (as `schema` future grants
are not concerned by `database` future grants, and vice versa).

The below example shows that despite referencing the same database and role, the resources don't conflict,
as one works for privileges `on` database and second one for privileges `in` database for future tables:

```terraform
resource "snowflake_grant_privileges_to_account_role" "on_database" {
  account_role_name = "TEST_ROLE"
  privileges = [ "MONITOR" ]
  on_account_object {
    object_type = "DATABASE"
    object_name = "TEST_DATABASE"
  }
  strict_privilege_management = true
  with_grant_option = true
}

resource "snowflake_grant_privileges_to_account_role" "in_database" {
  privileges        = [ "SELECT", "INSERT" ]
  account_role_name = "TEST_ROLE"
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = "TEST_DATABASE"
    }
  }
}
```

### Delayed application

This is a Terraform limitation, as are limited with what we can show in the plan.
Because of this, we decided to delay the revoking action in certain scenarios, so no revokes will be called without showing proper plan beforehand.
Otherwise, the revokes could happen "under the hood", making you unable to make conscious decision whether a given set of privileges should be revoked or not.
The delay will occur whenever you update already existing [grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role) resource
to have the flag set to `true`. The flag setting takes place in the first `terraform apply`,
and the proper `terraform plan` with potential privilege revokes can be applied with the second `terraform apply`.

