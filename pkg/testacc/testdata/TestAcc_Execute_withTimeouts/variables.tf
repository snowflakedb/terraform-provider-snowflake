variable "execute" {
  type = string
}

variable "revert" {
  type = string
}

variable "query" {
  type = string
}

variable "create_timeout" {
  type    = string
  default = "21m"
}

variable "read_timeout" {
  type    = string
  default = "22m"
}

variable "update_timeout" {
  type    = string
  default = "23m"
}

variable "delete_timeout" {
  type    = string
  default = "24m"
}
