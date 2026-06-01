# https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table#examples
resource "snowflake_dynamic_table" "dt" {
  name     = "product"
  database = "mydb"
  schema   = "myschema"
  target_lag {
    maximum_duration = "20 minutes"
  }
  warehouse = "mywh"
  query     = "SELECT product_id, product_name FROM \"mydb\".\"myschema\".\"staging_table\""
  comment   = "example comment"
}

# Optional: attach column-level constraints (e.g. NOT NULL) at CREATE time.
# When `column` blocks are provided, the list must enumerate every output
# column of `query`, in order, to match Snowflake's CREATE DYNAMIC TABLE
# (<col_list>) ... AS <query> syntax. Changes to `column` blocks force the
# dynamic table to be recreated.
resource "snowflake_dynamic_table" "dt_with_columns" {
  name     = "product_required"
  database = "mydb"
  schema   = "myschema"
  target_lag {
    downstream = true
  }
  warehouse = "mywh"
  query     = "SELECT product_id, product_name FROM \"mydb\".\"myschema\".\"staging_table\""

  column {
    name = "product_id"
  }
  column {
    name     = "product_name"
    not_null = true
  }
}
