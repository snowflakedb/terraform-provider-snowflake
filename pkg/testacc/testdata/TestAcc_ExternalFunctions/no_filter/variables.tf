variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "name" {
  type = string
}

variable "api_allowed_prefixes" {
  type = list(string)
}

variable "url_of_proxy_and_resource" {
  type = string
}

variable "comment" {
  type = string
}
