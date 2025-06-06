resource "snowflake_masking_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  dynamic "argument" {
    for_each = var.argument
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
  body                  = var.body
  return_data_type      = var.return_data_type
  exempt_other_policies = var.exempt_other_policies
  comment               = var.comment
}
