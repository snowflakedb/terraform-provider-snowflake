# Basic - only required fields
resource "snowflake_table" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

# Complete - every configurable field set. column.default and column.identity are
# mutually exclusive, and the three default kinds (constant, expression, sequence)
# are mutually exclusive, so each is shown on a different column.
resource "snowflake_table" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "TABLE"
  comment  = "A table."

  is_transient                    = "true"
  cluster_by                      = ["to_date(DATE)"]
  change_tracking                 = true
  enable_schema_evolution         = "true"
  data_retention_time_in_days     = 1
  max_data_extension_time_in_days = 30
  default_ddl_collation           = "en-ci"
  error_logging                   = "true"
  row_timestamp                   = "true"

  column {
    name     = "id"
    type     = "NUMBER(38,0)"
    nullable = true
    comment  = "sequence-backed id"

    default {
      sequence = snowflake_sequence.sequence.fully_qualified_name
    }
  }

  column {
    name     = "identity"
    type     = "NUMBER(38,0)"
    nullable = true

    identity {
      start_num = 1
      step_num  = 3
    }
  }

  column {
    name     = "data"
    type     = "text"
    nullable = false
    collate  = "en-ci"
    comment  = "collated text"

    default {
      constant = "pending"
    }

    masking_policy = "DATABASE.SCHEMA.MASKING_POLICY"
  }

  column {
    name = "DATE"
    type = "TIMESTAMP_NTZ(9)"

    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  column {
    name    = "extra"
    type    = "VARIANT"
    comment = "extra data"
  }

  row_access_policy {
    policy_name = "DATABASE.SCHEMA.ROW_ACCESS_POLICY"
    on          = ["id"]
  }

  aggregation_policy {
    policy_name = "DATABASE.SCHEMA.AGGREGATION_POLICY"
    entity_key  = ["id"]
  }

  join_policy {
    policy_name       = "DATABASE.SCHEMA.JOIN_POLICY"
    allowed_join_keys = ["id"]
  }
}

resource "snowflake_sequence" "sequence" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "sequence"
}
