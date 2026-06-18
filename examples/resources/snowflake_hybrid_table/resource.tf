# basic resource
resource "snowflake_hybrid_table" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "HYBRID_TABLE"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }

  primary_key {
    keys = ["ID"]
  }
}

# complete resource
resource "snowflake_hybrid_table" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "HYBRID_TABLE"
  comment  = "A hybrid table for HTAP workloads"

  data_retention_time_in_days     = 7
  max_data_extension_time_in_days = 14

  column {
    name     = "ID"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "NAME"
    type     = "VARCHAR(256)"
    nullable = true
    collate  = "en-ci"
    comment  = "Name column"
  }

  column {
    name     = "CREATED_AT"
    type     = "TIMESTAMP_NTZ"
    nullable = false
    default {
      expression = "CURRENT_TIMESTAMP()"
    }
  }

  primary_key {
    name = "pk_hybrid_table"
    keys = ["ID"]
  }

  index {
    name    = "idx_name"
    columns = ["NAME"]
  }

  index {
    name            = "idx_name_created_at"
    columns         = ["NAME"]
    include_columns = ["CREATED_AT"]
  }
}
