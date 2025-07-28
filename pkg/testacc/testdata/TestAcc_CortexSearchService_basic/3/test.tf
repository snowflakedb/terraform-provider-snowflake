resource "snowflake_warehouse" "t" {
  name           = var.warehouse
  warehouse_size = "XSMALL"
}

resource "snowflake_table" "t" {
  database        = var.database
  schema          = var.schema
  name            = var.table_name
  change_tracking = true

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }

  column {
    name = "SOME_TEXT"
    type = "VARCHAR"
  }

  column {
    name = "SOME_OTHER_TEXT"
    type = "VARCHAR(32)"
  }
}

resource "snowflake_cortex_search_service" "css" {
  depends_on = [snowflake_table.t, snowflake_warehouse.t]
  on         = var.on
  attributes = var.attributes
  name       = var.name
  database   = var.database
  schema     = var.schema
  target_lag = "2 minutes"
  warehouse  = var.warehouse
  query      = var.query
  comment    = var.comment
}
