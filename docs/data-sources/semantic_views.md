---
page_title: "snowflake_semantic_views Data Source - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Data source used to get details of filtered semantic views. Filtering is aligned with the current possibilities for SHOW SEMANTIC VIEWS https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views query. The results of SHOW and DESCRIBE are encapsulated in one output collection semantic_views.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

# snowflake_semantic_views (Data source)

Data source used to get details of filtered semantic views. Filtering is aligned with the current possibilities for [SHOW SEMANTIC VIEWS](https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `semantic_views`.

## Example Usage

```terraform
# Simple usage
data "snowflake_semantic_views" "simple" {
}

output "simple_output" {
  value = data.snowflake_semantic_views.simple.semantic_views
}

# Filtering (like)
data "snowflake_semantic_views" "like" {
  like = "semantic-view-name"
}

output "like_output" {
  value = data.snowflake_semantic_views.like.semantic_views
}

# Filtering by prefix (like)
data "snowflake_semantic_views" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_semantic_views.like_prefix.semantic_views
}

# Filtering (starts_with)
data "snowflake_semantic_views" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_semantic_views.starts_with.semantic_views
}

# Filtering (in)
data "snowflake_semantic_views" "in_account" {
  in {
    account = true
  }
}

data "snowflake_semantic_views" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_semantic_views" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_semantic_views.in_account.semantic_views,
    "database" : data.snowflake_semantic_views.in_database.semantic_views,
    "schema" : data.snowflake_semantic_views.in_schema.semantic_views,
  }
}

# Filtering (limit)
data "snowflake_semantic_views" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_semantic_views.limit.semantic_views
}

# Without additional data (to limit the number of calls made for every found semantic view)
data "snowflake_semantic_views" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SEMANTIC VIEW for every semantic view found and attaches its output to semantic_views.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_semantic_views.only_show.semantic_views
}

# Ensure the number of semantic views is equal to at least one element (with the use of postcondition)
data "snowflake_semantic_views" "assert_with_postcondition" {
  like = "semantic-view-name%"
  lifecycle {
    postcondition {
      condition     = length(self.semantic_views) > 0
      error_message = "there should be at least one semantic view"
    }
  }
}

# Ensure the number of semantic views is equal to exactly one element (with the use of check block)
check "semantic_view_check" {
  data "snowflake_semantic_views" "assert_with_check_block" {
    like = "semantic-view-name"
  }

  assert {
    condition     = length(data.snowflake_semantic_views.assert_with_check_block.semantic_views) == 1
    error_message = "semantic views filtered by '${data.snowflake_semantic_views.assert_with_check_block.like}' returned ${length(data.snowflake_semantic_views.assert_with_check_block.semantic_views)} semantic views where one was expected"
  }
}
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

