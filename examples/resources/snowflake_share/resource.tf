# basic resource
resource "snowflake_share" "basic" {
  name = "SHARE_NAME"
}

# complete resource
resource "snowflake_share" "complete" {
  name                = "SHARE_NAME"
  comment             = "A share."
  accounts            = ["organizationName.accountName"]
  secure_objects_only = "false"
}

resource "snowflake_database" "example" {
  # remember to define dependency between objects on a share, because shared objects have to be dropped before dropping share
  depends_on = [snowflake_share.complete]
  name       = "test"
}
