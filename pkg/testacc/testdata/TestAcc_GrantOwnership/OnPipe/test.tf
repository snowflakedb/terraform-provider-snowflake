resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_stage" "test" {
  database = var.database
  schema   = var.schema
  name     = var.stage
}

resource "snowflake_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.table

  column {
    type = "NUMBER(38,0)"
    name = "ID"
  }
}

resource "snowflake_pipe" "test" {
  database       = var.database
  schema         = var.schema
  name           = var.pipe
  copy_statement = "copy into \"${snowflake_table.test.database}\".\"${snowflake_table.test.schema}\".\"${snowflake_table.test.name}\"(ID) from @\"${snowflake_stage.test.database}\".\"${snowflake_stage.test.schema}\".\"${snowflake_stage.test.name}\""
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name

  on {
    object_type = "PIPE"
    object_name = snowflake_pipe.test.fully_qualified_name
  }
}
