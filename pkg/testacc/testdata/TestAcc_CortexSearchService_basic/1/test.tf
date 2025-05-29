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
}

resource "snowflake_cortex_search_service" "css" {
  depends_on = [snowflake_table.t]
  name       = var.name
  on         = var.on
  database   = var.database
  schema     = var.schema
  target_lag = "2 minutes"
  warehouse  = var.warehouse
  query      = var.query
  comment    = var.comment
}
