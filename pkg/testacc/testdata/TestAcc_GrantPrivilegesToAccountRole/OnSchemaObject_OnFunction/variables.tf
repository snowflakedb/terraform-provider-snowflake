variable "name" {
  type = string
}

variable "function_name" {
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

variable "with_grant_option" {
  type = bool
}

variable "argument_type" {
  type = string
}
