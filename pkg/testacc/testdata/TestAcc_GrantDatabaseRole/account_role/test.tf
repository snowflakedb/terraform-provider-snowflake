variable "database_role_name" {
  type = string
}

variable "parent_role_name" {
  type = string
}

variable "database" {
  type = string
}

resource "snowflake_database_role" "database_role" {
  database = var.database
  name     = var.database_role_name
}

resource "snowflake_account_role" "parent_role" {
  name = var.parent_role_name
}

resource "snowflake_grant_database_role" "g" {
  database_role_name = snowflake_database_role.database_role.fully_qualified_name
  parent_role_name   = snowflake_account_role.parent_role.name
}
