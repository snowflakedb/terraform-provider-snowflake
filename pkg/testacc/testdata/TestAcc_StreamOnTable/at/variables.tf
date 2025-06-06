variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "table" {
  type = string
}

variable "copy_grants" {
  type = bool
}

variable "show_initial_rows" {
  type = string
}

variable "append_only" {
  type = string
}

variable "at" {
  type = map(string)
}

variable "comment" {
  type = string
}
