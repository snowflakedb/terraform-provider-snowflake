resource "snowflake_storage_lifecycle_policy" "slp" {
  database = "prod"
  schema   = "security"
  name     = "default_storage_lifecycle_policy"
  argument {
    name = "VAL"
    type = "VARCHAR"
  }
  body = "LENGTH(VAL) > 0"
}

resource "snowflake_table" "table" {
  database = "prod"
  schema   = "security"
  name     = "my_table"
  column {
    name = "VAL"
    type = "VARCHAR"
  }
}

resource "snowflake_dynamic_table" "dynamic_table" {
  database = "prod"
  schema   = "security"
  name     = "my_dynamic_table"
  target_lag {
    maximum_duration = "20 minutes"
  }
  warehouse = "my_warehouse"
  query     = "SELECT VAL FROM \"prod\".\"security\".\"my_table\""
}

# Attaching a storage lifecycle policy to a table
resource "snowflake_table_storage_lifecycle_policy_attachment" "table_attachment" {
  table_name                    = snowflake_table.table.fully_qualified_name
  table_type                    = "TABLE"
  storage_lifecycle_policy_name = snowflake_storage_lifecycle_policy.slp.fully_qualified_name
  on                            = ["VAL"]
}

# Attaching a storage lifecycle policy to a dynamic table
resource "snowflake_table_storage_lifecycle_policy_attachment" "dynamic_table_attachment" {
  table_name                    = snowflake_dynamic_table.dynamic_table.fully_qualified_name
  table_type                    = "DYNAMIC_TABLE"
  storage_lifecycle_policy_name = snowflake_storage_lifecycle_policy.slp.fully_qualified_name
  on                            = ["VAL"]
}
