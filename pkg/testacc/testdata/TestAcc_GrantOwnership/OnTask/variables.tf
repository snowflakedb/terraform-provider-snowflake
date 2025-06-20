variable "account_role_name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "task" {
  type = string
}

variable "warehouse" {
  type    = string
  default = null
}

variable "warehouse_init_size" {
  type    = string
  default = null
}
