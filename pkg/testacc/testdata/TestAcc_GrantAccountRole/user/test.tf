variable "role_name" {
  type = string
}

variable "user_name" {
  type = string
}

resource "snowflake_account_role" "role" {
  name = var.role_name
}

resource "snowflake_user" "user" {
  name = var.user_name
}

resource "snowflake_grant_account_role" "g" {
  role_name = snowflake_account_role.role.name
  user_name = snowflake_user.user.name
}
