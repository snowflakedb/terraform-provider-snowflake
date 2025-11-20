# basic resource
resource "snowflake_semantic_view" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "SEMANTIC_VIEW"

  metrics {
    semantic_expression {
      qualified_expression_name = "\"lt1\".\"m1\""
      sql_expression            = "SUM(\"lt1\".\"a1\")"
    }
  }

  tables {
    table_alias = "lt1"
    table_name  = snowflake_table.test.fully_qualified_name
  }
}

# complete resource
resource "snowflake_semantic_view" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "SEMANTIC_VIEW"
  comment  = "comment"

  dimensions {
    comment                   = "dimension comment"
    qualified_expression_name = "\"lt1\".\"d2\""
    sql_expression            = "\"lt1\".\"a2\""
    synonym = ["dim2"]
  }

  facts {
    comment                   = "fact comment"
    qualified_expression_name = "\"lt1\".\"f2\""
    sql_expression            = "\"lt1\".\"a1\""
    synonym = ["fact2"]
  }

  metrics {
    semantic_expression {
      comment                   = "semantic expression comment"
      qualified_expression_name = "\"lt1\".\"m1\""
      sql_expression            = "SUM(\"lt1\".\"a1\")"
      synonym = ["sem1", "baseSem"]
    }
  }

  metrics {
    window_function {
      over_clause {
        partition_by = "\"lt1\".\"d2\""
      }
      qualified_expression_name = "\"lt1\".\"wf1\""
      sql_expression            = "SUM(\"lt1\".\"m1\")"
    }
  }

  relationships {
    referenced_relationship_columns = ["a1", "a2"]
    referenced_table_name_or_alias {
      table_alias = "lt1"
    }
    relationship_columns = ["a1", "a2"]
    relationship_identifier = "r2"
    table_name_or_alias {
      table_alias = "lt2"
    }
  }

  tables {
    comment     = "logical table 1 comment"
    primary_key = ["a1"]
    synonym = ["orders", "sales"]
    table_alias = "lt1"
    table_name  = snowflake_table.test.fully_qualified_name
    unique {
      values = ["a2"]
    }
    unique {
      values = ["a3", "a4"]
    }
  }
  tables {
    comment     = ""
    primary_key = ["a1"]
    table_alias = "lt2"
    table_name  = snowflake_table.test2.fully_qualified_name
  }
}
