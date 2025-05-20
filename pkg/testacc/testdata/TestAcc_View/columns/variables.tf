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

variable "projection_name" {
  type = string
}

variable "masking_name" {
  type = string
}

variable "masking_using" {
  type    = list(string)
  default = null
}
