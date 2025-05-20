variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "authentication_methods" {
  type = set(string)
}

variable "mfa_authentication_methods" {
  type = set(string)
}

variable "mfa_enrollment" {
  type = string
}

variable "client_types" {
  type = set(string)
}

variable "security_integrations" {
  type = set(string)
}

variable "comment" {
  type = string
}
