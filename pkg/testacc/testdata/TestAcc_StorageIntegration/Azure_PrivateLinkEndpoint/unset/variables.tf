variable "name" {
  description = "The name of the storage integration"
  type        = string
}

variable "allowed_locations" {
  description = "Explicitly limits external stages that use the integration to reference one or more storage locations"
  type        = list(string)
}

variable "azure_tenant_id" {
  description = "The ID for your Office 365 tenant that the Azure service principal belongs to"
  type        = string
}

variable "use_private_link_endpoint" {
  description = "For an Azure storage integration, whether to use private links for connections to Azure"
  type        = bool
  default     = null
}
