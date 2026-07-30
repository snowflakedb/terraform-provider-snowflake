# Basic - only required fields
resource "snowflake_iceberg_table" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }
  column {
    name = "NAME"
    type = "VARCHAR(16777216)"
  }
}

# Complete - every field set (except cluster_by, which conflicts with partition_by - see below)
resource "snowflake_iceberg_table" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"
  comment  = "COMMENT"

  # Snowflake-managed parameters
  external_volume                 = "EXTERNAL_VOLUME"
  catalog                         = "SNOWFLAKE"
  catalog_sync                    = "CATALOG_INTEGRATION"
  target_file_size                = "64MB"
  storage_serialization_policy    = "OPTIMIZED"
  data_retention_time_in_days     = 5
  max_data_extension_time_in_days = 10
  enable_data_compaction          = true
  enable_iceberg_merge_on_read    = true
  base_location                   = "iceberg_table"
  path_layout                     = "FLAT"
  change_tracking                 = "true"
  iceberg_version                 = 2
  error_logging                   = "true"

  column {
    name     = "ID"
    type     = "NUMBER(38,0)"
    not_null = "true"
    comment  = "Primary identifier"
  }
  column {
    name    = "NAME"
    type    = "VARCHAR(16777216)"
    comment = "Name of the entity"
    masking_policy {
      policy_name = "MASKING_POLICY"
      using       = ["NAME"]
    }
  }
  column {
    name = "REGION"
    type = "VARCHAR(16777216)"
    projection_policy {
      policy_name = "PROJECTION_POLICY"
    }
  }
  column {
    name = "STATUS"
    type = "VARCHAR(16777216)"
  }
  column {
    name = "CATEGORY"
    type = "VARCHAR(16777216)"
    masking_policy {
      policy_name = "CONDITIONAL_MASKING_POLICY"
      using       = ["CATEGORY", "STATUS"]
    }
  }
  column {
    name = "CREATED_AT"
    type = "TIMESTAMP_NTZ(9)"
    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }
  column {
    name = "REF_ID"
    type = "NUMBER(38,0)"
    default {
      expression = "2"
    }
  }

  primary_key_constraint {
    name               = "PK"
    column             = ["ID"]
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Primary key constraint"
  }

  unique_constraint {
    name               = "NAME_UQ"
    column             = ["NAME"]
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Unique constraint on name"
  }

  foreign_key_constraint {
    name               = "FK"
    column             = ["REF_ID"]
    table_name         = "OTHER_DATABASE.OTHER_SCHEMA.OTHER_TABLE"
    ref_column         = ["ID"]
    match              = "SIMPLE"
    on_update          = "CASCADE"
    on_delete          = "SET NULL"
    enforced           = "false"
    deferrable         = "true"
    initially_deferred = "true"
    enable             = "true"
    validate           = "true"
    rely               = "true"
    comment            = "Foreign key constraint"
  }

  check_constraint {
    name       = "CHK"
    expression = "ID > 0"
    validate   = "true"
  }

  row_access_policy {
    policy_name = "ROW_ACCESS_POLICY"
    on          = ["ID"]
  }

  aggregation_policy {
    policy_name = "AGGREGATION_POLICY"
    entity_key  = ["ID"]
  }

  partition_by {
    identity = "REGION"
  }
  partition_by {
    bucket {
      num_buckets = 4
      column      = "ID"
    }
  }
  partition_by {
    truncate {
      width  = 10
      column = "NAME"
    }
  }
  partition_by {
    year = "CREATED_AT"
  }
  partition_by {
    month = "CREATED_AT"
  }
  partition_by {
    day = "CREATED_AT"
  }
  partition_by {
    hour = "CREATED_AT"
  }
}

# cluster_by conflicts with partition_by, so it is shown on a separate resource.
resource "snowflake_iceberg_table" "complete_with_cluster_by" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }
  column {
    name = "NAME"
    type = "VARCHAR(16777216)"
  }

  cluster_by = ["ID", "NAME"]
}
