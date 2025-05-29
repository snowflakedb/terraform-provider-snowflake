resource "snowflake_share" "test" {
  name = var.share
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = var.database
}

data "snowflake_grants" "test" {
  depends_on = [snowflake_grant_privileges_to_share.test]

  grants_to {
    share {
      share_name = snowflake_share.test.name
    }
  }
}
