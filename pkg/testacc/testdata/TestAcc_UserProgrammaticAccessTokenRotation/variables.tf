variable "name" {
  type = string
}

variable "user" {
  type = string
}

variable "keepers" {
  type    = map(string)
  default = null
}

variable "expire_rotated_token_after_hours" {
  type    = number
  default = null
}
