variable "name" {}
variable "azure_tenant_id" {}
variable "allowed_locations" { type = set(string) }
variable "use_private_link_endpoint" { type = bool }
