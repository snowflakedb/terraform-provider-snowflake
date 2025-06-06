variable "name" {
  type = string
}

variable "azure_tenant_id" {
  type = string
}

variable "azure_ad_application_id" {
  type = string
}

variable "api_allowed_prefixes" {
  type = list(string)
}

variable "api_blocked_prefixes" {
  type = list(string)
}

variable "api_key" {
  type = string
}

variable "comment" {
  type = string
}

variable "enabled" {
  type = bool
}
