variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "privileges" {
  type    = list(string)
  default = ["USAGE"]
}

variable "with_grant_option" {
  type = bool
}
