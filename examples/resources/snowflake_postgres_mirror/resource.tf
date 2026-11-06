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
