resource "snowflake_database" "test" {
  name = var.name

  replication {
    enable_to_account {
      account_identifier = var.account_identifier
      with_failover      = var.with_failover
    }

    ignore_edition_check = var.ignore_edition_check
  }
}
