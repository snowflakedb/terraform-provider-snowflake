resource "snowflake_execute" "test" {
  execute = var.execute
  revert  = var.revert
  query   = var.query

  timeouts {
    create = var.create_timeout
    read   = var.read_timeout
    update = var.update_timeout
    delete = var.delete_timeout
  }
}
