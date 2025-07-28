variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "statement" {
  type = string
}

variable "copy_grants" {
  type = bool
}

variable "is_secure" {
  type = bool
}

variable "column" {
  type = set(map(string))
}
