#basic resource
resource "snowflake_semantic_view" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "SEMANTIC_VIEW"
  tables {
    table_name = "TABLE_NAME"
  }
  metrics {
    semantic_expression {
      qualified_expression_name = "TABLE_NAME.METRIC_NAME"
      sql_expression            = "SQL_EXPRESSION"
    }
  }
}

# complete resource
resource "snowflake_semantic_view" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "SEMANTIC_VIEW"
  comment  = "comment"
  tables {
    comment     = "comment"
    primary_key = ["COL1"]
    unique {
      values = ["COL2"]
    }
    table_alias = "TABLE_ALIAS"
    table_name  = "TABLE_NAME"
    synonym     = ["synonym", "synonym"]
  }
  tables {
    comment     = "comment"
    primary_key = ["COL1"]
    unique {
      values = ["COL2"]
    }
    table_alias = "TABLE_ALIAS"
    table_name  = "TABLE_NAME"
    synonym     = ["synonym"]
  }
  relationships {
    relationship_identifier = "RELATIONSHIP_NAME"
    table_name_or_alias {
      table_alias = "TABLE_ALIAS"
    }
    relationship_columns = ["COL1", "COL2"]
    referenced_table_name_or_alias {
      table_alias = "TABLE_ALIAS"
    }
    referenced_relationship_columns = ["COL1", "COL2"]
  }
  facts {
    comment                   = "comment"
    qualified_expression_name = "TABLE_ALIAS.FACT_NAME"
    sql_expression            = "SQL_EXPRESSION"
    synonym                   = ["synonym"]
  }
  dimensions {
    comment                   = "comment"
    qualified_expression_name = "TABLE_ALIAS.DIMENSION_NAME"
    sql_expression            = "SQL_EXPRESSION"
    synonym                   = ["synonym"]
  }
  metrics {
    semantic_expression {
      comment                   = "comment"
      qualified_expression_name = "TABLE_ALIAS.METRIC_NAME"
      sql_expression            = "SQL_EXPRESSION"
      synonym                   = ["synonym"]
    }
  }
  metrics {
    window_function {
      metric = "METRIC_EXPRESSION"
      over_clause {
        partition_by        = "PARTITION_CLAUSE"
        order_by            = "ORDER_BY_CLAUSE"
        window_frame_clause = "WINDOW_FRAME_CLAUSE"
      }
      window_function = "ALIAS.WINDOW_FUNCTION_NAME"
    }
  }
}