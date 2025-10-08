variable "extra_ips" {
  type = list(string)
  # This variable is required on purpose. It should be set to other extra IPs we have in the network policies.
}

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "2.8.0"
    }
  }
}

provider "snowflake" {
  profile = "default"
}

provider "snowflake" {
  profile = "secondary_test_account"
  alias   = "secondary"
}

provider "snowflake" {
  profile = "third_test_account"
  alias   = "third"
}

provider "snowflake" {
  profile = "fourth_test_account"
  alias   = "fourth"
}

locals {
  all_github_ips  = compact(split("\n", file("github_ipv4_subnets.txt")))
  extra_ips       = var.extra_ips
  all_allowed_ips = concat(local.extra_ips, local.all_github_ips)
  comment         = "Allows for connections only comming from GitHub Actions or behind VPN"
}

resource "snowflake_network_policy" "restricted_access_primary" {
  name            = "RESTRICTED_ACCESS"
  allowed_ip_list = local.all_allowed_ips
  comment         = local.comment
}

# resource "snowflake_network_policy" "restricted_access_secondary" {
#   provider        = snowflake.secondary
#   name            = "RESTRICTED_ACCESS"
#   allowed_ip_list = local.all_allowed_ips
#   comment         = local.comment
# }

# resource "snowflake_network_policy" "restricted_access_third" {
#   provider        = snowflake.third
#   name            = "RESTRICTED_ACCESS"
#   allowed_ip_list = local.all_allowed_ips
#   comment         = local.comment
# }

# resource "snowflake_network_policy" "restricted_access_fourth" {
#   provider        = snowflake.fourth
#   name            = "RESTRICTED_ACCESS"
#   allowed_ip_list = local.all_allowed_ips
#   comment         = local.comment
# }
