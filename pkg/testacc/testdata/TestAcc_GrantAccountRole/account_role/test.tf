variable "role_name" {
  type = string
}

variable "parent_role_name" {
  type = string
}

resource "snowflake_account_role" "role" {
  name = var.role_name
}

resource "snowflake_account_role" "parent_role" {
  name = var.parent_role_name
}

resource "snowflake_grant_account_role" "g" {
  role_name        = snowflake_account_role.role.name
  parent_role_name = snowflake_account_role.parent_role.name
}
