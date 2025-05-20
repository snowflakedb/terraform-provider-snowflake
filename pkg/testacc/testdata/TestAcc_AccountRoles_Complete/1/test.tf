resource "snowflake_account_role" "test1" {
  name    = var.account_role_name_1
  comment = var.comment
}

resource "snowflake_account_role" "test2" {
  name    = var.account_role_name_2
  comment = var.comment
}

resource "snowflake_account_role" "test3" {
  name    = var.account_role_name_3
  comment = var.comment
}

data "snowflake_account_roles" "test" {
  depends_on = [
    snowflake_account_role.test1,
    snowflake_account_role.test2,
    snowflake_account_role.test3,
  ]
  like = var.like
}
