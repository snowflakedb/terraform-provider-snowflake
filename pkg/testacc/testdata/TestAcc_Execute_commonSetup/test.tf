resource "snowflake_execute" "test" {
  execute = var.execute
  revert  = var.revert
}
