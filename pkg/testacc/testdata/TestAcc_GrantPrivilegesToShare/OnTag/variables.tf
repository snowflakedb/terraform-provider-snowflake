variable "to_share" {
  type = string
}

variable "privileges" {
  type = list(string)
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "on_tag" {
  type = string
}
