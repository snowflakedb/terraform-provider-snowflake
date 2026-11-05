---
page_title: "snowflake_postgres_mirror Resource - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Resource used to manage Postgres data mirrors, which replicate Postgres tables into a Snowflake database. For more information, check data mirroring documentation https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

-> **Note** Before you create a mirror, grant the `SNOWFLAKE.POSTGRES_MIRROR_ADMIN` application role to the role that manages it, and grant `USAGE` on the Postgres instance to the `SNOWFLAKE` application. If the mirror uses a `warehouse` or an `external_volume`, the `SNOWFLAKE` application needs `USAGE` on it too. The provider has no resource for granting privileges to an application yet, so the example uses [`snowflake_execute`](./execute) for those grants. Do not revoke these grants from the `SNOWFLAKE` application while mirrors use them, because that breaks the mirrors. For details, see [Create a mirror](https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring-create).

-> **Note** `terraform apply` returns once the mirror is created, before the initial copy of the tables finishes. To follow the initial copy, call `SNOWFLAKE.POSTGRES.LIST_MIRRORED_TABLES('<mirror_name>')`: tables show `SNAPSHOTTING` during the initial copy and `REPLICATING` once they caught up.

-> **Note** Changes to `postgres_tables`, `postgres_schemas`, `refresh_interval`, and `warehouse` are applied in place with `ALTER_MIRROR`, without copying the existing tables again. Changes to `name`, `postgres_instance`, `postgres_database`, `target_database`, `target_table_type`, and `external_volume` recreate the mirror, which copies all tables again.

-> **Note** If you manage a mirror today with `snowflake_execute` calling `CREATE_MIRROR`, remove that resource from the state with `terraform state rm` (destroying it would run its `revert` statement), remove its block from the configuration, and import the mirror by name into `snowflake_postgres_mirror`.

-> **Note** Snowflake reports `target_table_type = SNOWFLAKE` when the configuration omits it. The attribute forces a new resource, so importing in enforcing mode plans a replacement. Permissive mode leaves the attribute unset, and a configuration that omits it does not replace the mirror. See [Importing resources](../guides/importing_resources#import-modes) for how the two modes work.

-> **Known limitation** Data mirroring is available only on AWS and Azure, and not on the BURSTABLE instance tier. When the source Postgres instance cannot be reached, `DESCRIBE_MIRROR` returns no tables or schemas. The provider then keeps the whole previous state, including `describe_output`, and reports a warning. Changes made outside of Terraform are not detected until the tables and schemas can be read again, and import fails until then.

# snowflake_postgres_mirror (Resource)

Resource used to manage Postgres data mirrors, which replicate Postgres tables into a Snowflake database. For more information, check [data mirroring documentation](https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring).

## Example Usage

```terraform
## Prerequisites
# The role that manages the mirror needs the postgres_mirror_admin application role,
# and the SNOWFLAKE application needs USAGE on the Postgres instance.
resource "snowflake_grant_application_role" "mirror_admin" {
  application_role_name    = "\"SNOWFLAKE\".\"POSTGRES_MIRROR_ADMIN\""
  parent_account_role_name = "my_role"
}

resource "snowflake_execute" "instance_usage" {
  execute = "GRANT USAGE ON POSTGRES INSTANCE \"my_instance\" TO APPLICATION SNOWFLAKE"
  # Revoking USAGE from the SNOWFLAKE application breaks active mirrors, so destroy keeps the grant.
  revert = "SELECT 1"
}

# Only when the mirror uses a warehouse.
resource "snowflake_execute" "warehouse_usage" {
  execute = "GRANT USAGE ON WAREHOUSE \"MY_WAREHOUSE\" TO APPLICATION SNOWFLAKE"
  revert  = "SELECT 1"
}

# Only when the mirror writes Iceberg tables to an external volume.
resource "snowflake_execute" "external_volume_usage" {
  execute = "GRANT USAGE ON EXTERNAL VOLUME \"MY_EXTERNAL_VOLUME\" TO APPLICATION SNOWFLAKE"
  revert  = "SELECT 1"
}

## Minimal
resource "snowflake_postgres_mirror" "basic" {
  name              = "my_mirror"
  postgres_instance = "my_instance"
  postgres_database = "postgres"
  target_database   = "MY_MIRROR_DB"
  postgres_tables   = ["public.orders", "public.customers"]

  depends_on = [snowflake_grant_application_role.mirror_admin, snowflake_execute.instance_usage]
}

## Complete
resource "snowflake_postgres_mirror" "complete" {
  name              = "my_complete_mirror"
  postgres_instance = "my_instance"
  postgres_database = "postgres"
  target_database   = "MY_COMPLETE_MIRROR_DB"
  postgres_tables   = ["public.orders"]
  postgres_schemas  = ["sales"]
  refresh_interval {
    minutes = 1
  }
  warehouse                       = "MY_WAREHOUSE"
  drop_target_database_on_destroy = true

  depends_on = [snowflake_grant_application_role.mirror_admin, snowflake_execute.instance_usage, snowflake_execute.warehouse_usage]
}

## Iceberg tables on an external volume
resource "snowflake_postgres_mirror" "iceberg" {
  name              = "my_iceberg_mirror"
  postgres_instance = "my_instance"
  postgres_database = "postgres"
  target_database   = "MY_ICEBERG_MIRROR_DB"
  postgres_tables   = ["public.orders"]
  target_table_type = "ICEBERG"
  external_volume   = "MY_EXTERNAL_VOLUME"

  depends_on = [snowflake_grant_application_role.mirror_admin, snowflake_execute.instance_usage, snowflake_execute.external_volume_usage]
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specifies the name of the mirror. It must be unique across all Postgres instances in the account. Snowflake stores the name lowercased, and the lowercased name is the resource ID.
- `postgres_database` (String) Specifies the source database on the Postgres instance.
- `postgres_instance` (String) Specifies the source Postgres instance (case-sensitive). The role used by the provider must own the instance.
- `target_database` (String) Specifies the Snowflake database the mirror creates and keeps up to date. The database must not exist yet, and Snowflake uppercases its name. The mirror protects the database, so do not manage it with `snowflake_database`.

### Optional

- `drop_target_database_on_destroy` (Boolean) (Default: `true`) Specifies whether destroying the mirror also drops its target database. A kept target database can never be the target of a mirror again, so with `false`, a replacement that keeps the same `target_database` fails. A kept target database stays owned by the `SNOWFLAKE` application. To manage or drop it, first transfer its ownership with `GRANT OWNERSHIP ON DATABASE`, which needs the `MANAGE GRANTS` privilege.
- `external_volume` (String) Specifies the external volume for Iceberg target tables. Valid only with `target_table_type = "ICEBERG"`. When not set, the tables use Snowflake-managed storage. Grant `USAGE` on the volume to `APPLICATION SNOWFLAKE` first.
- `postgres_schemas` (Set of String) Specifies the source schemas to mirror. Every current and future table in them is included. Schema names are literal and case-sensitive, written without Postgres quoting, such as `MixedCase` rather than `"MixedCase"`. Schemas are added and removed in place.
- `postgres_tables` (Set of String) Specifies the source tables to mirror, as schema-qualified names (`schema.table`) written the way `DESCRIBE_MIRROR` returns them: unquoted and lowercase, except for parts Postgres requires to be quoted, such as `public."MixedCase"`. Other spellings of the same table show as a change on every plan. Tables are added and removed in place, without copying the other tables again.
- `refresh_interval` (Block List, Max: 1) Specifies how often changes are merged into the target tables. Set exactly one of `hours`, `minutes`, or `seconds`. When not set, Snowflake uses its default interval. Removing the block keeps the current interval, because `ALTER_MIRROR` cannot reset it. (see [below for nested schema](#nestedblock--refresh_interval))
- `target_table_type` (String) Specifies the format of the target tables. When not set, Snowflake uses `SNOWFLAKE`. Valid values are (case-insensitive): `SNOWFLAKE` | `ICEBERG`.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `warehouse` (String) Specifies the warehouse the mirror's apply task runs on. When not set, the task runs on serverless compute. Snowflake uppercases the name, so the warehouse itself must have an uppercase name; the case used in the configuration does not matter. Grant `USAGE` on the warehouse to `APPLICATION SNOWFLAKE` first.

### Read-Only

- `describe_output` (List of Object) Outputs the result of `DESCRIBE_MIRROR` for the given Postgres mirror. (see [below for nested schema](#nestedatt--describe_output))
- `fully_qualified_name` (String) Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).
- `id` (String) The ID of this resource.

<a id="nestedblock--refresh_interval"></a>
### Nested Schema for `refresh_interval`

Optional:

- `hours` (Number) Specifies the refresh interval in hours. (conflicts with `minutes` and `seconds`)
- `minutes` (Number) Specifies the refresh interval in minutes. (conflicts with `hours` and `seconds`)
- `seconds` (Number) Specifies the refresh interval in seconds. (conflicts with `hours` and `minutes`)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--describe_output"></a>
### Nested Schema for `describe_output`

Read-Only:

- `apply_lag_bytes` (Number)
- `apply_lag_operations` (Number)
- `apply_lag_seconds` (Number)
- `catalog_integration_name` (String)
- `error_message` (String)
- `external_volume` (String)
- `feed_only` (Boolean)
- `last_applied_commit_lsn` (String)
- `last_applied_commit_time` (String)
- `last_apply_time` (String)
- `last_operation_id` (Number)
- `last_operation_time` (String)
- `name` (String)
- `postgres_database` (String)
- `postgres_instance` (String)
- `postgres_schemas` (List of String)
- `postgres_status` (String)
- `postgres_tables` (List of String)
- `recent_query_durations` (String)
- `refresh_interval` (String)
- `status` (String)
- `target_database` (String)
- `target_table_type` (String)
- `warehouse` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_postgres_mirror.example '<mirror_name>'
```
