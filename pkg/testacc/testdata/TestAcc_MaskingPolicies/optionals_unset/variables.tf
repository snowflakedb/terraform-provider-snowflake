variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "argument" {
  type = set(map(string))
}

variable "body" {
  type = string
}

variable "comment" {
  type = string
}

variable "return_data_type" {
  type = string
}
