# Hybrid tables with foreign key relationships
#
# This example demonstrates referential integrity between hybrid tables
# using both inline and out-of-line foreign key constraints.
#
# Note: Hybrid tables require preview feature enablement in the provider:
# provider "snowflake" {
#   preview_features_enabled = ["snowflake_hybrid_table_resource"]
# }

# Parent table: Products
resource "snowflake_hybrid_table" "products" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "PRODUCTS"
  comment  = "Product catalog table (parent)"

  column {
    name        = "product_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
    comment     = "Unique product identifier"
  }

  column {
    name     = "product_name"
    type     = "VARCHAR(100)"
    nullable = false
  }

  column {
    name = "price"
    type = "NUMBER(10,2)"
  }

  column {
    name     = "category"
    type     = "VARCHAR(50)"
    nullable = false
  }

  # Index for FK lookups
  index {
    name    = "idx_product_category"
    columns = ["category"]
  }

  data_retention_time_in_days = 7
}

# Child table: Order Items with inline foreign key
resource "snowflake_hybrid_table" "order_items" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "ORDER_ITEMS"
  comment  = "Order line items with product references (child)"

  # Depends on products table being created first
  depends_on = [snowflake_hybrid_table.products]

  column {
    name        = "order_item_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name     = "order_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  # Column with inline foreign key constraint
  column {
    name     = "product_id"
    type     = "NUMBER(38,0)"
    nullable = false
    comment  = "References products.product_id"

    foreign_key {
      table_name  = snowflake_hybrid_table.products.fully_qualified_name
      column_name = "product_id"
    }
  }

  column {
    name = "quantity"
    type = "NUMBER(10,0)"
  }

  column {
    name = "unit_price"
    type = "NUMBER(10,2)"
  }

  data_retention_time_in_days = 7
}

# Another child table: Reviews with out-of-line foreign key
resource "snowflake_hybrid_table" "reviews" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "PRODUCT_REVIEWS"
  comment  = "Product reviews with foreign key constraint"

  depends_on = [snowflake_hybrid_table.products]

  column {
    name        = "review_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name     = "product_id"
    type     = "NUMBER(38,0)"
    nullable = false
    comment  = "Product being reviewed"
  }

  column {
    name = "rating"
    type = "NUMBER(1,0)"
  }

  column {
    name = "review_text"
    type = "VARCHAR(1000)"
  }

  column {
    name = "review_date"
    type = "DATE"
  }

  # Out-of-line foreign key constraint
  foreign_key {
    name               = "fk_review_product"
    columns            = ["product_id"]
    references_table   = snowflake_hybrid_table.products.fully_qualified_name
    references_columns = ["product_id"]
  }

  # Index for FK column
  index {
    name    = "idx_review_product"
    columns = ["product_id"]
  }

  data_retention_time_in_days = 7
}

# Composite foreign key example
resource "snowflake_hybrid_table" "order_shipments" {
  database = "MY_DATABASE"
  schema   = "MY_SCHEMA"
  name     = "ORDER_SHIPMENTS"
  comment  = "Shipment tracking for orders"

  depends_on = [snowflake_hybrid_table.order_items]

  column {
    name        = "shipment_id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name     = "order_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name     = "order_item_id"
    type     = "NUMBER(38,0)"
    nullable = false
  }

  column {
    name = "tracking_number"
    type = "VARCHAR(50)"
  }

  column {
    name = "ship_date"
    type = "DATE"
  }

  # Composite foreign key referencing order_items
  foreign_key {
    name               = "fk_shipment_order_item"
    columns            = ["order_item_id"]
    references_table   = snowflake_hybrid_table.order_items.fully_qualified_name
    references_columns = ["order_item_id"]
  }

  data_retention_time_in_days = 7
}

# Outputs
output "products_fqn" {
  value = snowflake_hybrid_table.products.fully_qualified_name
}

output "order_items_fqn" {
  value = snowflake_hybrid_table.order_items.fully_qualified_name
}

output "reviews_fqn" {
  value = snowflake_hybrid_table.reviews.fully_qualified_name
}
