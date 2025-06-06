resource "snowflake_row_access_policy" "test_1" {
  name     = var.name_1
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body = var.body
}

resource "snowflake_row_access_policy" "test_2" {
  name     = var.name_2
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body = var.body
}

resource "snowflake_row_access_policy" "test_3" {
  name     = var.name_3
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body = var.body
}


data "snowflake_row_access_policies" "test" {
  depends_on = [snowflake_row_access_policy.test_1, snowflake_row_access_policy.test_2, snowflake_row_access_policy.test_3]

  like = var.like
}
