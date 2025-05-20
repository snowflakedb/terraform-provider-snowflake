variable "name_1" {
  type = string
}

variable "name_2" {
  type = string
}

variable "name_3" {
  type = string
}

variable "like" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "arguments" {
  type = set(map(string))
}

variable "body" {
  type = string
}
