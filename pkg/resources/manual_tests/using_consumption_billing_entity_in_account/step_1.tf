# Commands to run
# - terraform apply
# - terraform plan
# - terraform apply -destroy

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "2.4.0"
    }
  }
}

provider "snowflake" {
  role = "GLOBALORGADMIN"
}

resource "snowflake_account" "test" {
  name                       = "<random_name>"
  admin_name                 = "<random_admin_name>"
  admin_rsa_public_key       = "<random_admin_rsa_public_key>"
  email                      = "<random_admin_email>"
  edition                    = "ENTERPRISE"
  grace_period_in_days       = 3
  consumption_billing_entity = "<test_consumption_billing_entity>"
}
