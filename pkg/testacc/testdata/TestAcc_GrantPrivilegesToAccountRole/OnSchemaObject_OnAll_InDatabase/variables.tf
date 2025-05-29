variable "name" {
  type = string
}

variable "privileges" {
  type = list(string)
}

variable "database" {
  type = string
}

variable "object_type_plural" {
  type = string
}

variable "with_grant_option" {
  type = bool
}
