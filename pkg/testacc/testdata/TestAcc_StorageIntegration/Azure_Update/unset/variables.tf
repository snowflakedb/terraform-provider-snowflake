variable "name" {
  type = string
}

variable "allowed_locations" {
  type = set(string)
}

variable "azure_tenant_id" {
  type = string
}
