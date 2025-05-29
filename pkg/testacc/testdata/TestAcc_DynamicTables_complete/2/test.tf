resource "snowflake_table" "t" {
  database        = var.database
  schema          = var.schema
  name            = var.table_name
  change_tracking = true
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_dynamic_table" "dt" {
  depends_on = [snowflake_table.t]
  name       = var.name
  database   = var.database
  schema     = var.schema
  target_lag {
    maximum_duration = "2 minutes"
  }
  warehouse = var.warehouse
  query     = var.query
  comment   = var.comment
}

data "snowflake_dynamic_tables" "dts" {
  like {
    pattern = snowflake_dynamic_table.dt.name
  }
  in {
    schema = "\"${snowflake_dynamic_table.dt.database}\".\"${snowflake_dynamic_table.dt.schema}\""
  }
  starts_with = snowflake_dynamic_table.dt.name
  limit {
    rows = 1
  }
}
