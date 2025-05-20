locals {
  application_role_identifier = "\"${var.application_name}\".\"${var.application_role_name}\""
}

resource "snowflake_grant_application_role" "g" {
  application_role_name    = local.application_role_identifier
  parent_account_role_name = "\"${var.parent_account_role_name}\""
}
