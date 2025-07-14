## Minimal
resource "snowflake_organization_account" "minimal" {
  name                 = "ACCOUNT_NAME"
  admin_name           = var.admin_name
  admin_password       = var.admin_password
  email                = var.email
  edition              = "ENTERPRISE"
  grace_period_in_days = 3
}

## Complete
resource "snowflake_organization_account" "complete" {
  name                 = "ACCOUNT_NAME"
  admin_name           = var.admin_name
  admin_rsa_public_key = var.admin_rsa_public_key
  first_name           = var.first_name
  last_name            = var.last_name
  email                = var.email
  edition              = "ENTERPRISE"
  region_group         = "PUBLIC"
  region               = "AWS_US_WEST_2"
  comment              = "some comment"
  grace_period_in_days = 3
}

variable "admin_name" {
  type      = string
  sensitive = true
}

variable "email" {
  type      = string
  sensitive = true
}

variable "admin_password" {
  type      = string
  sensitive = true
}

variable "admin_rsa_public_key" {
  type      = string
  sensitive = true
}

variable "first_name" {
  type      = string
  sensitive = true
}

variable "last_name" {
  type      = string
  sensitive = true
}
