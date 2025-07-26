resource "snowflake_table" "base_table" {
  database        = var.database
  schema          = var.schema
  name            = "${var.on_table}_base"
  change_tracking = true
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_dynamic_table" "test_dynamic_table" {
  name       = var.on_table
  database   = var.database
  schema     = var.schema
  target_lag {
    maximum_duration = "2 minutes"
  }
  warehouse = var.warehouse
  query = <<-EOT
    with temp as (
      select "id" from ${snowflake_table.base_table.fully_qualified_name}
    )
    select * from temp
  EOT
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = var.to_share
  privileges  = ["USAGE"]
  on_database = var.database
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = var.to_share
  privileges = var.privileges
  on_table   = snowflake_dynamic_table.test_dynamic_table.fully_qualified_name
  depends_on = [snowflake_grant_privileges_to_share.test_setup]
}
