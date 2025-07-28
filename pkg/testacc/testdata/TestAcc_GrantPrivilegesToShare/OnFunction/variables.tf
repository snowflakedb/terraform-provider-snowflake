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

variable "argument_type" {
  type = string
}
