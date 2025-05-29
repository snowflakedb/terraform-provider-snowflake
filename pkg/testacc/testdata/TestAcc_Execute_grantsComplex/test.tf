resource "snowflake_execute" "test" {
  for_each = { for index, db_grant in var.database_grants : index => db_grant }
  execute  = "GRANT ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} TO ROLE ${each.value.role_id}"
  revert   = "REVOKE ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} FROM ROLE ${each.value.role_id}"
}
