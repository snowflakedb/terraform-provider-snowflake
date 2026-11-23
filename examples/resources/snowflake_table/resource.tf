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
# mutually exclusive. column.default.expression is a raw SQL expression, so string
# literals keep their quotes and a sequence is referenced with NEXTVAL. An empty
# expression ("") is an empty string default, the same as the literal "''".
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
    not_null = false
    comment  = "sequence-backed id"

    default {
      expression = "${snowflake_sequence.sequence.fully_qualified_name}.NEXTVAL"
    }
  }

  column {
    name     = "identity"
    type     = "NUMBER(38,0)"
    not_null = false

    identity {
      start_num = 1
      step_num  = 3
    }
  }

  column {
    name     = "data"
    type     = "text"
    not_null = true
    collate  = "en-ci"
    comment  = "collated text"

    default {
      expression = "'pending'"
    }

    masking_policy {
      policy_name = "DATABASE.SCHEMA.MASKING_POLICY"
      using       = ["data", "extra"]
    }
  }

  column {
    name = "DATE"
    type = "TIMESTAMP_NTZ(9)"

    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  column {
    name = "notes"
    type = "VARCHAR(256)"

    # An empty string default. "''" is equivalent.
    default {
      expression = ""
    }
  }

  column {
    name    = "extra"
    type    = "VARIANT"
    comment = "extra data"
  }

  # A virtual column is computed from the expression at query time and is not stored.
  column {
    name = "double_id"
    type = "NUMBER(38,0)"
    as   = "\"id\" * 2"
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
