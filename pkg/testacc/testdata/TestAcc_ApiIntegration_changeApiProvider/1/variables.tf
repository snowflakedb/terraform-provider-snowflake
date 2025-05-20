variable "name" {
  type = string
}

variable "api_provider" {
  type = string
}

variable "api_aws_role_arn" {
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

